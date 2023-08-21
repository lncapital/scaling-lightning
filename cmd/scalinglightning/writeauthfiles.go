package scalinglightning

import (
	"fmt"
	"path"

	sl "github.com/scaling-lightning/scaling-lightning/pkg/network"
	"github.com/spf13/cobra"
)

var writeAuthFilesCmd = &cobra.Command{
	Use:   "writeauthfiles",
	Short: "Output the auth files for a node or all nodes",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		processDebugFlag(cmd)
		nodeName := cmd.Flag("node").Value.String()
		authFilesDir := cmd.Flag("dir").Value.String()
		all, err := cmd.Flags().GetBool("all")
		if err != nil {
			fmt.Printf("Problem getting all flag: %v\n", err.Error())
			return
		}

		slnetwork, err := sl.DiscoverStartedNetwork(kubeConfigPath)
		if err != nil {
			fmt.Printf(
				"Problem with network discovery, is there a network running? Error: %v\n",
				err.Error(),
			)
			return
		}
		foundANode := false
		originalAuthFilesDir := authFilesDir
		for _, node := range slnetwork.LightningNodes {
			if all {
				authFilesDir = path.Join(originalAuthFilesDir, node.GetName())
			}
			if node.GetName() == nodeName || all {
				err := node.WriteAuthFilesToDirectory(authFilesDir)
				if err != nil {
					fmt.Printf("Problem writing auth files: %v\n", err.Error())
					return
				}
				foundANode = true
			}
		}
		if foundANode {
			fmt.Println("Files written")
			return
		}

		allNames := []string{}
		for _, node := range slnetwork.LightningNodes {
			allNames = append(allNames, node.GetName())
		}
		fmt.Printf(
			"Can't find node(s), here are the lightning nodes that are running: %v\n",
			allNames,
		)
	},
}

func init() {
	rootCmd.AddCommand(writeAuthFilesCmd)

	writeAuthFilesCmd.Flags().
		StringP("node", "n", "", "The name of the node to download the auth files for")

	writeAuthFilesCmd.Flags().
		BoolP("all", "a", false, "Download the auth files for all nodes")

	writeAuthFilesCmd.Flags().
		StringP("dir", "o", "", "The directory to write the auth files to")
	writeAuthFilesCmd.MarkFlagRequired("dir")

}
