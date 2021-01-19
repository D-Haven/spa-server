{{/*
Expand the name of the chart.
*/}}
{{- define "spa-server.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "spa-server.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "spa-server.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "spa-server.labels" -}}
helm.sh/chart: {{ include "spa-server.chart" . }}
{{ include "spa-server.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "spa-server.selectorLabels" -}}
app.kubernetes.io/name: {{ include "spa-server.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "spa-server.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "spa-server.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Retrieve requested TLS Secret
*/}}
{{- define "spa-server.tls.type" -}}
{{- if empty (lookup "v1" "Secret" .Release.Namespace (index .Values.ingress.tls 0).secretName) -}}
missing
{{- else }}
{{- (lookup "v1" "Secret" .Release.Namespace (index .Values.ingress.tls 0).secretName).type }}
{{- end }}
{{- end }}

{{/*
Validate TLS Secret if it is provided
*/}}
{{- define "spa-server.tls.useTLS" -}}
{{- if gt (len .Values.ingress.tls) 0 }}
{{- if not (eq "kubernetes.io/tls" (include "spa-server.tls.type" .)) }}
{{- cat "Secret" (index .Values.ingress.tls 0).secretName "is" (include "spa-server.tls.type" .) "-- expecting type kubernetes.io/tls" | fail }}
{{- end }}
true
{{- else }}
false
{{- end }}
{{- end }}