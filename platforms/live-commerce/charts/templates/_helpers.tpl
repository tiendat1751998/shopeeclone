{{- define "shopee-live-commerce.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- define "shopee-live-commerce.fullname" -}}
{{- printf "%s-%s" .Release.Name .Chart.Name | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- define "shopee-live-commerce.labels" -}}
app: {{ include "shopee-live-commerce.name" . }}
version: v1
{{- end -}}
