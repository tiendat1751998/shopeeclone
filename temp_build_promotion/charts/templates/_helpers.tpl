{{- define "shopee-promotion.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- define "shopee-promotion.fullname" -}}
{{- printf "%s-%s" .Release.Name .Chart.Name | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- define "shopee-promotion.labels" -}}
app: {{ include "shopee-promotion.name" . }}
version: v1
{{- end -}}
