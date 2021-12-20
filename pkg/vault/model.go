package vault

const (
	DefaultVaultKubernetesHost = "https://kubernetes.default.svc"

	VaultRootToken = "root-token"
	VaultUnsealKey = "unseal-key"

	VaultKubernetesHost   = "kubernetes_host"
	VaultKubernetesCACert = "kubernetes_ca_cert"
	VaultTokenReviewerJWT = "token_reviewer_jwt"

	VaultAuthKubernetesConfig = "auth/kubernetes/config"
	VaultAuthKubernetesPath   = "kubernetes/"
	VaultAuthKubernetesType   = "kubernetes"

	VaultSysAuth = "sys/auth"
)

type VaultKeys struct {
	RootToken  string
	UnsealKeys []string
}
