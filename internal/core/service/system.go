package service

import "github.com/gin-gonic/gin"

type Pinger interface {
	Ping(ctx *gin.Context) error
}

type SystemService struct {
	Pinger
}

func NewSystemService(store Pinger) *SystemService {
	return &SystemService{store}
}
