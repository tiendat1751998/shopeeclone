{{- define "shopee-advertising.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- define "shopee-advertising.fullname" -}}
{{- printf "%s-%s" .Release.Name .Chart.Name | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- define "shopee-advertising.labels" -}}
app: {{ include "shopee-advertising.name" . }}
version: v1
{{- end -}}
