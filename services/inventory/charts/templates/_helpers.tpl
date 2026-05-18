{{- define "shopee-inventory.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- define "shopee-inventory.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name (include "shopee-inventory.name" .) | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- define "shopee-inventory.labels" -}}
helm.sh/chart: {{ include "shopee-inventory.name" . }}-{{ .Chart.Version }}
app.kubernetes.io/name: {{ include "shopee-inventory.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}
