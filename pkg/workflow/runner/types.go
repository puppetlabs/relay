package runner

type RunnerKind string

const (
	RunnerKindGKEClusterProvisioner RunnerKind = "gke-cluster-provisioner"
	RunnerKindShell                 RunnerKind = "shell"
	RunnerKindWorkflow              RunnerKind = "workflow"
	RunnerKindHelmDeploy            RunnerKind = "helm-deploy"
)
