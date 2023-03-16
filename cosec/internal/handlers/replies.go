package handlers

import (
	"github.com/50HJ/Intelli-Mall/cosec/internal"
	"github.com/50HJ/Intelli-Mall/cosec/internal/models"
	"github.com/50HJ/Intelli-Mall/internal/am"
	"github.com/50HJ/Intelli-Mall/internal/registry"
	"github.com/50HJ/Intelli-Mall/internal/sec"
)

func NewReplyHandlers(reg registry.Registry, orchestrator sec.Orchestrator[*models.CreateOrderData], mws ...am.MessageHandlerMiddleware) am.MessageHandler {
	return am.NewReplyHandler(reg, orchestrator, mws...)
}

func RegisterReplyHandlers(subscriber am.MessageSubscriber, handlers am.MessageHandler) error {
	_, err := subscriber.Subscribe(internal.CreateOrderReplyChannel, handlers, am.GroupName("cosec-replies"))
	return err
}
