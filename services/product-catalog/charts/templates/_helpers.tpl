{{- define "shopee-product-catalog.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- define "shopee-product-catalog.fullname" -}}
{{- printf "%s-%s" .Release.Name .Chart.Name | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- define "shopee-product-catalog.labels" -}}
app: {{ include "shopee-product-catalog.name" . }}
version: v1
{{- end -}}
