package execution

import (
	logging "github.com/puppetlabs/insights-logging"
	"github.com/puppetlabs/nebula/pkg/io"
)

type ExecutorRuntime interface {
	IO() *io.IO
	Logger() logging.Logger
}
