{{- define "shopee-payment.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- define "shopee-payment.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name (include "shopee-payment.name" .) | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- define "shopee-payment.labels" -}}
helm.sh/chart: {{ include "shopee-payment.name" . }}-{{ .Chart.Version }}
app.kubernetes.io/name: {{ include "shopee-payment.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}
