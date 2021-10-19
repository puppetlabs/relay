module github.com/puppetlabs/relay

go 1.16

require (
	github.com/cli/browser v1.0.0
	github.com/docker/go-connections v0.4.0
	github.com/eiannone/keyboard v0.0.0-20200508000154-caf4b762e807
	github.com/fatih/color v1.13.0
	github.com/go-swagger/go-swagger v0.23.0
	github.com/hashicorp/hcl2 v0.0.0-20191002203319-fb75b3253c80
	github.com/jedib0t/go-pretty v4.3.0+incompatible
	github.com/jetstack/cert-manager v0.16.0
	github.com/kballard/go-shellquote v0.0.0-20180428030007-95032a82bc51
	github.com/puppetlabs/errawr-gen v1.0.1
	github.com/puppetlabs/errawr-go/v2 v2.2.0
	github.com/puppetlabs/leg/encoding v0.1.0
	github.com/puppetlabs/leg/httputil v0.1.4
	github.com/puppetlabs/leg/k8sutil v0.6.1
	github.com/puppetlabs/leg/timeutil v0.4.2
	github.com/puppetlabs/leg/workdir v0.1.0
	github.com/puppetlabs/relay-client-go/client v0.4.2
	github.com/puppetlabs/relay-client-go/models v1.0.6
	github.com/puppetlabs/relay-core v0.0.0-20211019182941-e1deb885f736
	github.com/rancher/helm-controller v0.6.3
	github.com/rancher/k3d/v4 v4.4.7
	github.com/shurcooL/vfsgen v0.0.0-20200824052919-0d455de96546
	github.com/spf13/cobra v1.2.1
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.8.1
	github.com/stretchr/testify v1.7.0
	golang.org/x/crypto v0.0.0-20210513164829-c07d793c2f9a
	k8s.io/api v0.22.2
	k8s.io/apiextensions-apiserver v0.20.1
	k8s.io/apimachinery v0.22.2
	k8s.io/apiserver v0.20.6
	k8s.io/client-go v11.0.1-0.20190805182717-6502b5e7b1b5+incompatible
	k8s.io/klog/v2 v2.9.0 // indirect
	knative.dev/caching v0.0.0-20210215030244-1212288570f0
	sigs.k8s.io/controller-runtime v0.8.1
)

replace (
	k8s.io/api => k8s.io/api v0.19.7
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.19.7
	k8s.io/apimachinery => k8s.io/apimachinery v0.19.7
	k8s.io/client-go => k8s.io/client-go v0.19.7
)

exclude github.com/Sirupsen/logrus v1.4.0

exclude github.com/Sirupsen/logrus v1.3.0

exclude github.com/Sirupsen/logrus v1.2.0

exclude github.com/Sirupsen/logrus v1.1.1

exclude github.com/Sirupsen/logrus v1.1.0
