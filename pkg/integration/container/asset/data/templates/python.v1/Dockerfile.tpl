{{- $PackageFileName := printf "nebula_sdk-%s-py3-none-any.whl" (.SDKVersion | trimPrefix "v") -}}
{{- $PackageRepoURL := printf "https://packages.nebula.puppet.net/sdk/support/python/%s" .SDKVersion -}}

{{- $FilePath := printf "/nebula/step-%s.py" .Name -}}
FROM {{ .Settings.Image }}
RUN apk --no-cache add bash ca-certificates curl git jq openssh && update-ca-certificates
{{- if .Settings.AdditionalAlpinePackages }}
RUN apk --no-cache add{{ range .Settings.AdditionalAlpinePackages }} {{ . }}{{ end }}
{{- end }}
RUN pip --no-cache-dir install "{{ printf "%s/%s" $PackageRepoURL $PackageFileName }}"
{{- if .Settings.AdditionalPythonPackages }}
RUN pip --no-cache-dir install{{ range .Settings.AdditionalPythonPackages }} {{ . }}{{ end }}
{{- end }}
{{- range .Settings.AdditionalCommands }}
RUN ["/bin/bash", "-c", {{ . | mustToJson }}]
{{- end }}
COPY "./{{ .Settings.CommandPath }}" "{{ $FilePath }}"
ENTRYPOINT []
CMD ["python3", "{{ $FilePath }}"]
