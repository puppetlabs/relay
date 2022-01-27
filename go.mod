module github.com/puppetlabs/relay

go 1.16

require (
	github.com/cli/browser v1.0.0
	github.com/eiannone/keyboard v0.0.0-20200508000154-caf4b762e807
	github.com/fatih/color v1.13.0
	github.com/go-swagger/go-swagger v0.23.0
	github.com/google/wire v0.5.0
	github.com/jedib0t/go-pretty v4.3.0+incompatible
	github.com/kballard/go-shellquote v0.0.0-20180428030007-95032a82bc51
	github.com/puppetlabs/errawr-gen v1.0.1
	github.com/puppetlabs/errawr-go/v2 v2.2.0
	github.com/puppetlabs/leg/encoding v0.1.0
	github.com/puppetlabs/leg/httputil v0.1.4
	github.com/puppetlabs/leg/k8sutil v0.6.4
	github.com/puppetlabs/leg/timeutil v0.4.2
	github.com/puppetlabs/leg/workdir v0.1.0
	github.com/puppetlabs/relay-client-go/client v0.4.5
	github.com/puppetlabs/relay-client-go/models v1.0.7
	github.com/puppetlabs/relay-core v0.0.0-20220126182442-87b4675b8fdf
	github.com/rancher/helm-controller v0.6.3
	github.com/rancher/k3d/v5 v5.1.0
	github.com/shurcooL/vfsgen v0.0.0-20200824052919-0d455de96546
	github.com/spf13/cobra v1.2.1
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.8.1
	github.com/stretchr/testify v1.7.0
	golang.org/x/crypto v0.0.0-20210711020723-a769d52b0f97
	k8s.io/api v0.23.0
	k8s.io/apiextensions-apiserver v0.22.2
	k8s.io/apimachinery v0.23.0
	k8s.io/apiserver v0.20.6
	k8s.io/client-go v11.0.1-0.20190805182717-6502b5e7b1b5+incompatible
	knative.dev/caching v0.0.0-20210215030244-1212288570f0
	sigs.k8s.io/controller-runtime v0.8.3
)

replace (
	k8s.io/api => k8s.io/api v0.19.7
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.19.7
	k8s.io/apimachinery => k8s.io/apimachinery v0.19.7
	k8s.io/client-go => k8s.io/client-go v0.19.7
)

exclude github.com/golangci/golangci-lint v1.42.0

exclude github.com/Sirupsen/logrus v1.4.0

exclude github.com/Sirupsen/logrus v1.3.0

exclude github.com/Sirupsen/logrus v1.2.0

exclude github.com/Sirupsen/logrus v1.1.1

exclude github.com/Sirupsen/logrus v1.1.0
