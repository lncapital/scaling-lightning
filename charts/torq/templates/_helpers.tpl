{{/*
 Expand the name of the chart.
 */}}
{{- define "torq.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
 Create a default fully qualified app name.
 We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
 If release name contains chart name it will be used as a full name.
 */}}
{{- define "torq.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "torq.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
 Common labels
 */}}
{{- define "torq.labels" -}}
helm.sh/chart: {{ include "torq.chart" . }}
{{ include "torq.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
 Selector labels
 */}}
{{- define "torq.selectorLabels" -}}
app.kubernetes.io/name: {{ include "torq.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}
