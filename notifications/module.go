package notifications

import (
	"context"

	"github.com/50HJ/Intelli-Mall/customers/customerspb"
	"github.com/50HJ/Intelli-Mall/internal/am"
	"github.com/50HJ/Intelli-Mall/internal/amotel"
	"github.com/50HJ/Intelli-Mall/internal/amprom"
	"github.com/50HJ/Intelli-Mall/internal/jetstream"
	pg "github.com/50HJ/Intelli-Mall/internal/postgres"
	"github.com/50HJ/Intelli-Mall/internal/postgresotel"
	"github.com/50HJ/Intelli-Mall/internal/registry"
	"github.com/50HJ/Intelli-Mall/internal/system"
	"github.com/50HJ/Intelli-Mall/internal/tm"
	"github.com/50HJ/Intelli-Mall/notifications/internal/application"
	"github.com/50HJ/Intelli-Mall/notifications/internal/constants"
	"github.com/50HJ/Intelli-Mall/notifications/internal/grpc"
	"github.com/50HJ/Intelli-Mall/notifications/internal/handlers"
	"github.com/50HJ/Intelli-Mall/notifications/internal/postgres"
	"github.com/50HJ/Intelli-Mall/ordering/orderingpb"
)

type Module struct{}

func (m Module) Startup(ctx context.Context, mono system.Service) (err error) {
	return Root(ctx, mono)
}

func Root(ctx context.Context, svc system.Service) (err error) {
	// setup Driven adapters
	reg := registry.New()
	if err = customerspb.Registrations(reg); err != nil {
		return err
	}
	if err = orderingpb.Registrations(reg); err != nil {
		return err
	}
	inboxStore := pg.NewInboxStore(constants.InboxTableName, svc.DB())
	messageSubscriber := am.NewMessageSubscriber(
		jetstream.NewStream(svc.Config().Nats.Stream, svc.JS(), svc.Logger()),
		amotel.OtelMessageContextExtractor(),
		amprom.ReceivedMessagesCounter(constants.ServiceName),
	)
	customers := postgres.NewCustomerCacheRepository(
		constants.CustomersCacheTableName,
		postgresotel.Trace(svc.DB()),
		grpc.NewCustomerRepository(svc.Config().Rpc.Service(constants.CustomersServiceName)),
	)

	// setup application
	app := application.New(customers)
	integrationEventHandlers := handlers.NewIntegrationEventHandlers(
		reg, app, customers,
		tm.InboxHandler(inboxStore),
	)

	// setup Driver adapters
	if err := grpc.RegisterServer(ctx, app, svc.RPC()); err != nil {
		return err
	}
	if err = handlers.RegisterIntegrationEventHandlers(messageSubscriber, integrationEventHandlers); err != nil {
		return err
	}

	return nil
}
