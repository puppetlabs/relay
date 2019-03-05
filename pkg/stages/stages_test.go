package stages

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestStageExecution(t *testing.T) {
	driver := NewShellDriver()

	cmds := []Command{
		NewShellCommand(driver, "ls"),
		NewShellCommand(driver, "pwd"),
	}

	stage := NewCommandStage(driver, cmds)

	err := stage.Execute()

	spew.Dump(err)
}
