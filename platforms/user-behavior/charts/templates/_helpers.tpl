{{- define "shopee-user-behavior.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- define "shopee-user-behavior.fullname" -}}
{{- printf "%s-%s" .Release.Name .Chart.Name | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- define "shopee-user-behavior.labels" -}}
app: {{ include "shopee-user-behavior.name" . }}
version: v1
{{- end -}}
