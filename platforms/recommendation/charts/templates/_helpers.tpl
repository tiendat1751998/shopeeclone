{{- define "shopee-recommendation.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- define "shopee-recommendation.fullname" -}}
{{- printf "%s-%s" .Release.Name .Chart.Name | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- define "shopee-recommendation.labels" -}}
app: {{ include "shopee-recommendation.name" . }}
version: v1
{{- end -}}
