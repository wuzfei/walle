package server

import (
	"context"
	"go.uber.org/zap"
)

type Job struct {
	log *zap.Logger
}

func NewJob(
	log *zap.Logger,
) *Job {
	return &Job{
		log: log,
	}
}
func (j *Job) Start(ctx context.Context) error {
	// eg: kafka consumer
	return nil
}
func (j *Job) Stop(ctx context.Context) error {
	return nil
}
