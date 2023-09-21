{{/*
    Return the proper image name
*/}}
{{- define "ephemerup.image" -}}
{{ include "common.images.image" (dict "imageRoot" .Values.image "global" .Values.global) }}
{{- end -}}


