package ssh

import (
	"context"
	"fmt"
	"io"
	"strings"
)

type RemoteExec struct {
	client *client
	envs   *Envs
	output io.Writer
}

func (e *RemoteExec) Close() error {
	e.client.Done()
	return nil
}

func (e *RemoteExec) WithEnvs(envs *Envs) Command {
	e.envs = envs
	return e
}

func (e *RemoteExec) RunCtx(ctx context.Context, cmd string) error {
	sess, err := e.client.client.NewSession()
	if err != nil {
		return err
	}
	closed := false
	defer func() {
		if !closed {
			_ = sess.Close()
			closed = true
		}
	}()
	if e.envs != nil && !e.envs.Empty() {
		cmd = fmt.Sprintf("%s && %s", strings.Join(e.envs.SliceKV(), " "), cmd)
	}

	if e.output != nil {
		sess.Stdout = e.output
		sess.Stderr = e.output
	}
	if ctx != nil {
		closeCh := make(chan struct{})
		defer close(closeCh)
		go func() {
			select {
			case <-ctx.Done():
				if !closed {
					_ = sess.Close()
					closed = true
				}
				return
			case <-closeCh:
				return
			}
		}()
	}
	return sess.Run(cmd)
}

func (e *RemoteExec) Run(cmd string) error {
	return e.RunCtx(nil, cmd)
}
