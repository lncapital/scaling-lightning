package kube

import (
	"encoding/json"
	"os/exec"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/rs/zerolog/log"
)

const (
	mainNamespace    = "sl"
	traefikNamespace = "sl-traefik"
)

func KubeCP(kubeconfig string, source string, destination string) error {

	//on windows remove the C: from the path as kubeCP doesn't support a path with ":"
	if strings.HasPrefix(strings.ToLower(source), "c:") {
		source = source[2:]
	}
	if strings.HasPrefix(strings.ToLower(destination), "c:") {
		destination = destination[2:]
	}

	//TODO: sanitise inputs here
	kubectlCmd := exec.Command( //nolint:gosec
		"kubectl",
		"--kubeconfig",
		kubeconfig,
		"-n",
		mainNamespace,
		"cp",
		source,
		destination,
	)
	kubectlOut, err := kubectlCmd.CombinedOutput()
	log.Debug().Msgf("kubectl output was: %v", string(kubectlOut))
	if err != nil {
		log.Error().Err(err).Msgf("kubectl output was: %v", string(kubectlOut))
		return errors.Wrap(err, "Running kubectl cp command")
	}
	return nil
}

type ingressRouteTCPData struct {
	Spec struct {
		EntryPoints []string `json:"entryPoints"`
	} `json:"spec"`
}

type endpointsData struct {
	Subsets []struct {
		Ports []struct {
			Name     string `json:"name"`
			Port     uint16 `json:"port"`
			Protocol string `json:"protocol"`
		} `json:"ports"`
	} `json:"subsets"`
}

type endpointMode int

const (
	ModeHTTP endpointMode = iota
	ModeTCP
)

func GetEndpointForNode(kubeconfig string, ingressName string, mode endpointMode) (uint16, error) {
	crd := "ingressroutetcps.traefik.containo.us"
	if mode == ModeHTTP {
		crd = "ingressroutes.traefik.io"
	}

	// TODO: sanitise inputs here
	kubectlCmd := exec.Command( //nolint:gosec
		"kubectl",
		"--kubeconfig",
		kubeconfig,
		"-n",
		mainNamespace,
		"get",
		crd,
		ingressName,
		"-o",
		"json",
	)
	kubectlOut, err := kubectlCmd.Output()
	log.Debug().Msgf("kubectl output was: %v", string(kubectlOut))
	if err != nil {
		log.Error().Err(err).Msgf("kubectl output was: %v", string(kubectlOut))
		return 0, errors.Wrap(err, "Running kubectl get endpoints command")
	}
	ingressRouteTCPData := ingressRouteTCPData{}
	err = json.Unmarshal(kubectlOut, &ingressRouteTCPData)
	if err != nil {
		return 0, errors.Wrap(err, "Unmarshalling ingressRouteTCPData")
	}
	if len(ingressRouteTCPData.Spec.EntryPoints) != 1 {
		return 0, errors.New("Expected 1 entrypoint")
	}
	entrypoint := ingressRouteTCPData.Spec.EntryPoints[0]

	// TODO: sanitise inputs here
	kubectlCmd = exec.Command( //nolint:gosec
		"kubectl",
		"--kubeconfig",
		kubeconfig,
		"-n",
		traefikNamespace,
		"get",
		"endpoints",
		"traefik",
		"-o",
		"json",
	)

	kubectlOut, err = kubectlCmd.Output()
	log.Debug().Msgf("kubectl output was: %v", string(kubectlOut))
	if err != nil {
		log.Error().Err(err).Msgf("kubectl output was: %v", string(kubectlOut))
		return 0, errors.Wrap(err, "Running kubectl get endpoints command")
	}

	endpointsData := endpointsData{}
	err = json.Unmarshal(kubectlOut, &endpointsData)
	if err != nil {
		return 0, errors.Wrap(err, "Unmarshalling endpointsData")
	}
	for _, subset := range endpointsData.Subsets {
		for _, port := range subset.Ports {
			if port.Name == entrypoint {
				return port.Port, nil
			}
		}
	}
	return 0, errors.New("Couldn't find port")
}
