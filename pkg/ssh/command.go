package ssh

import "context"

type Command interface {
	WithEnvs(envs *Envs) Command
	RunCtx(ctx context.Context, cmd string) error
	Run(cmd string) error
	Close() error
}
