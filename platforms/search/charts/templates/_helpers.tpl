{{- define "shopee-search.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- define "shopee-search.fullname" -}}
{{- printf "%s-%s" .Release.Name .Chart.Name | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- define "shopee-search.labels" -}}
app: {{ include "shopee-search.name" . }}
version: v1
{{- end -}}
