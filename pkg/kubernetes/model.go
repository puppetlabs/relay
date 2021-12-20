package kubernetes

const (
	KubernetesSecretDataCACertificate = "ca.crt"
	KubernetesSecretDataToken         = "token"
)

type KubernetesCertificateData struct {
	CACertificate string
	Token         string
}
