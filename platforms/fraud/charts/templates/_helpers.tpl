{{- define "shopee-fraud.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- define "shopee-fraud.fullname" -}}
{{- printf "%s-%s" .Release.Name .Chart.Name | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- define "shopee-fraud.labels" -}}
app: {{ include "shopee-fraud.name" . }}
version: v1
{{- end -}}
