//go:build !windows

package daemon

import (
	"os/exec"
	"syscall"
)

type platformAttrs struct{}

func (p platformAttrs) setSysProcAttr(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setsid: true, // Create a new session
	}
}
