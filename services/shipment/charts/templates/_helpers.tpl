{{- define "shopee-shipment.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- define "shopee-shipment.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name (include "shopee-shipment.name" .) | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- define "shopee-shipment.labels" -}}
helm.sh/chart: {{ include "shopee-shipment.name" . }}-{{ .Chart.Version }}
app.kubernetes.io/name: {{ include "shopee-shipment.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}
