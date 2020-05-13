{{- $PackageFileName := printf "ni-%s-linux-amd64.tar.xz" .SDKVersion -}}
{{- $PackageSHA256FileName := printf "%s.sha256" $PackageFileName -}}
{{- $PackageRepoURL := printf "https://packages.nebula.puppet.net/sdk/ni/%s" .SDKVersion -}}
FROM {{ .Images.base.Ref }}
RUN set -eux ; \
    mkdir -p /tmp/ni && \
    cd /tmp/ni && \
    wget {{ printf "%s/%s" $PackageRepoURL $PackageFileName }} && \
    wget {{ printf "%s/%s" $PackageRepoURL $PackageSHA256FileName }} && \
    echo "$( cat {{ $PackageSHA256FileName }} )  {{ $PackageFileName }}" | sha256sum -c - && \
    tar -xvJf {{ $PackageFileName }} && \
    mv ni-{{ .SDKVersion }}*-linux-amd64 /usr/local/bin/ni && \
    cd - && \
    rm -fr /tmp/ni
