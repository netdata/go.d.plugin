package nvme

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
)

func (n *NVMe) validateConfig() error {
	if n.BinaryPath == "" {
		return errors.New("'binary_path' can not be empty")
	}

	return nil
}

func (n *NVMe) initNVMeCLIExec() (nvmeCLI, error) {
	nvmePath, err := exec.LookPath(n.BinaryPath)
	if err != nil {
		return nil, err
	}

	var sudoPath string
	if os.Getuid() != 0 {
		sudoPath, err = exec.LookPath("sudo")
		if err != nil {
			return nil, err
		}
	}

	if sudoPath != "" {
		ctx, cancel := context.WithTimeout(context.Background(), n.Timeout.Duration)
		defer cancel()

		_, err := exec.CommandContext(ctx, sudoPath, "-n", "-l", nvmePath).Output()
		if err != nil {
			return nil, fmt.Errorf("can not execute '%s' with sudo: %v", n.BinaryPath, err)
		}
	}

	return &nvmeCLIExec{
		sudoPath: sudoPath,
		nvmePath: nvmePath,
		timeout:  n.Timeout.Duration,
	}, nil
}
