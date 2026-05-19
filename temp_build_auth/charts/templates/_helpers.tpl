{{- define "shopee-auth.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{- define "shopee-auth.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- include "shopee-auth.name" . }}-{{ .Release.Name }}
{{- end }}
{{- end }}

{{- define "shopee-auth.labels" -}}
app: {{ include "shopee-auth.name" . }}
release: {{ .Release.Name }}
tier: services
{{- end }}
