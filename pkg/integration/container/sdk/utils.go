package sdk

const repoContentURL = "https://api.github.com/repos/relay-integrations/container-definitions/contents"

func url(path string) string {
	return repoContentURL + withSlashPrefix(path)
}

func withSlashPrefix(name string) string {
	if len(name) == 0 {
		return "/"
	}

	if name[0] != '/' {
		return "/" + name
	}

	return name
}
