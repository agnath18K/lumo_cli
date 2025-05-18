//go:build windows

package daemon

import (
	"os/exec"
	"syscall"
)

type platformAttrs struct{}

func (p platformAttrs) setSysProcAttr(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP,
	}
}
