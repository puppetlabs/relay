module github.com/puppetlabs/relay

go 1.14

require (
	github.com/Microsoft/hcsshim/test v0.0.0-20201109231005-3cc00bc91358 // indirect
	github.com/containerd/cgroups v0.0.0-20200817152742-7a3c009711fb // indirect
	github.com/containerd/containerd v1.4.0 // indirect
	github.com/containerd/continuity v0.0.0-20200710164510-efbc4488d8fe // indirect
	github.com/containerd/fifo v0.0.0-20200410184934-f15a3290365b // indirect
	github.com/containerd/ttrpc v1.0.1 // indirect
	github.com/containerd/typeurl v1.0.1 // indirect
	github.com/fairwindsops/rbac-manager v0.9.4
	github.com/fatih/color v1.10.0
	github.com/go-swagger/go-swagger v0.23.0
	github.com/gobuffalo/flect v0.2.2 // indirect
	github.com/gogo/googleapis v1.4.0 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/googleapis/gnostic v0.5.1 // indirect
	github.com/hashicorp/hcl2 v0.0.0-20191002203319-fb75b3253c80
	github.com/imdario/mergo v0.3.11 // indirect
	github.com/inconshreveable/log15 v0.0.0-20180818164646-67afb5ed74ec
	github.com/jedib0t/go-pretty v4.3.0+incompatible
	github.com/jetstack/cert-manager v0.16.0
	github.com/kballard/go-shellquote v0.0.0-20180428030007-95032a82bc51
	github.com/moby/sys/mountinfo v0.1.3 // indirect
	github.com/onsi/ginkgo v1.14.0 // indirect
	github.com/opencontainers/selinux v1.6.0 // indirect
	github.com/prometheus/client_golang v1.7.1 // indirect
	github.com/puppetlabs/errawr-gen v1.0.1
	github.com/puppetlabs/errawr-go/v2 v2.2.0
	github.com/puppetlabs/horsehead/v2 v2.16.0
	github.com/puppetlabs/relay-core v0.0.0-20201211131443-4ba4da9d77b0
	github.com/rancher/helm-controller v0.6.3
	github.com/rancher/k3d/v3 v3.0.1
	github.com/shurcooL/vfsgen v0.0.0-20181202132449-6a9ea43bcacd
	github.com/sirupsen/logrus v1.6.0
	github.com/spf13/cobra v1.1.1
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.7.0
	github.com/stretchr/testify v1.6.1
	github.com/syndtr/gocapability v0.0.0-20200815063812-42c35b437635 // indirect
	github.com/willf/bitset v1.1.11 // indirect
	golang.org/x/crypto v0.0.0-20200820211705-5c72a883971a
	golang.org/x/mod v0.4.0 // indirect
	golang.org/x/net v0.0.0-20201209123823-ac852fbbde11 // indirect
	golang.org/x/sys v0.0.0-20201211090839-8ad439b19e0f // indirect
	golang.org/x/text v0.3.4 // indirect
	golang.org/x/tools v0.0.0-20201211025543-abf6a1d87e11 // indirect
	google.golang.org/grpc v1.34.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776
	k8s.io/api v0.20.0
	k8s.io/apiextensions-apiserver v0.20.0
	k8s.io/apimachinery v0.20.0
	k8s.io/apiserver v0.18.5
	k8s.io/client-go v11.0.1-0.20190805182717-6502b5e7b1b5+incompatible
	k8s.io/component-base v0.18.5
	k8s.io/klog/v2 v2.1.0 // indirect
	k8s.io/kubernetes v1.17.2
	k8s.io/utils v0.0.0-20201110183641-67b214c5f920 // indirect
	knative.dev/caching v0.0.0-20200630172829-a78409990d76
	sigs.k8s.io/controller-runtime v0.5.11
	sigs.k8s.io/controller-tools v0.4.1 // indirect
)

replace (
	github.com/docker/distribution => github.com/docker/distribution v0.0.0-20200319173657-742aab907b54
	github.com/googleapis/gnostic => github.com/googleapis/gnostic v0.3.1
	k8s.io/api => k8s.io/api v0.17.12
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.17.12
	k8s.io/apimachinery => k8s.io/apimachinery v0.17.12
	k8s.io/apiserver => k8s.io/apiserver v0.17.12
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.17.12
	k8s.io/client-go => k8s.io/client-go v0.17.12
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.17.12
	k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.17.12
	k8s.io/code-generator => k8s.io/code-generator v0.17.12
	k8s.io/component-base => k8s.io/component-base v0.17.12
	k8s.io/cri-api => k8s.io/cri-api v0.17.12
	k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.17.12
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.17.12
	k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.17.12
	k8s.io/kube-proxy => k8s.io/kube-proxy v0.17.12
	k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.17.12
	k8s.io/kubectl => k8s.io/kubectl v0.17.12
	k8s.io/kubelet => k8s.io/kubelet v0.17.12
	k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.17.12
	k8s.io/metrics => k8s.io/metrics v0.17.12
	k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.17.12
)

exclude github.com/Sirupsen/logrus v1.4.0

exclude github.com/Sirupsen/logrus v1.3.0

exclude github.com/Sirupsen/logrus v1.2.0

exclude github.com/Sirupsen/logrus v1.1.1

exclude github.com/Sirupsen/logrus v1.1.0
