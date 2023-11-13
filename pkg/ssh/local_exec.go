package ssh

import (
	"context"
	"io"
	"os/exec"
)

type LocalExec struct {
	envs   *Envs
	output io.Writer
}

func NewLocalExec(output io.Writer) *LocalExec {
	return &LocalExec{
		output: output,
		envs:   NewEnvs(),
	}
}

func (e *LocalExec) Close() error {
	return nil
}

func (e *LocalExec) WithEnvs(envs *Envs) Command {
	e.envs = envs
	return e
}

func (e *LocalExec) RunCtx(ctx context.Context, cmd string) error {
	var command *exec.Cmd
	if ctx == nil {
		command = exec.Command("bash", "-c", cmd)
	} else {
		command = exec.CommandContext(ctx, "bash", "-c", cmd)
	}
	if e.envs != nil && !e.envs.Empty() {
		command.Env = e.envs.SliceKV()
	}
	if e.output != nil {
		command.Stderr = e.output
		command.Stdout = e.output
	}
	return command.Run()
}

func (e *LocalExec) Run(cmd string) error {
	return e.RunCtx(nil, cmd)
}
