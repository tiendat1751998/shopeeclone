{{- define "shopee-gateway.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{- define "shopee-gateway.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- include "shopee-gateway.name" . }}-{{ .Release.Name }}
{{- end }}
{{- end }}

{{- define "shopee-gateway.labels" -}}
app: {{ include "shopee-gateway.name" . }}
release: {{ .Release.Name }}
tier: platform
{{- end }}
