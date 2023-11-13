package service

import (
	"go.uber.org/zap"
	"yema.dev/internal/repository"
	"yema.dev/pkg/helper/sid"
	"yema.dev/pkg/jwt"
)

type Service struct {
	log *zap.Logger
	sid *sid.Sid
	jwt *jwt.JWT
	tm  repository.Transaction
}

func NewService(tm repository.Transaction, log *zap.Logger, sid *sid.Sid, jwt *jwt.JWT) *Service {
	return &Service{
		log: log,
		sid: sid,
		jwt: jwt,
		tm:  tm,
	}
}
