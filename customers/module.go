package customers

import (
	"context"
	"database/sql"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/rs/zerolog"

	"github.com/50HJ/Intelli-Mall/customers/customerspb"
	"github.com/50HJ/Intelli-Mall/customers/internal/application"
	"github.com/50HJ/Intelli-Mall/customers/internal/constants"
	"github.com/50HJ/Intelli-Mall/customers/internal/domain"
	"github.com/50HJ/Intelli-Mall/customers/internal/grpc"
	"github.com/50HJ/Intelli-Mall/customers/internal/handlers"
	"github.com/50HJ/Intelli-Mall/customers/internal/postgres"
	"github.com/50HJ/Intelli-Mall/customers/internal/rest"
	"github.com/50HJ/Intelli-Mall/internal/am"
	"github.com/50HJ/Intelli-Mall/internal/amotel"
	"github.com/50HJ/Intelli-Mall/internal/amprom"
	"github.com/50HJ/Intelli-Mall/internal/ddd"
	"github.com/50HJ/Intelli-Mall/internal/di"
	"github.com/50HJ/Intelli-Mall/internal/jetstream"
	pg "github.com/50HJ/Intelli-Mall/internal/postgres"
	"github.com/50HJ/Intelli-Mall/internal/postgresotel"
	"github.com/50HJ/Intelli-Mall/internal/registry"
	"github.com/50HJ/Intelli-Mall/internal/system"
	"github.com/50HJ/Intelli-Mall/internal/tm"
)

type Module struct{}

func (m Module) Startup(ctx context.Context, mono system.Service) (err error) {
	return Root(ctx, mono)
}

func Root(ctx context.Context, svc system.Service) (err error) {
	container := di.New()
	// setup Driven adapters
	container.AddSingleton(constants.RegistryKey, func(c di.Container) (any, error) {
		reg := registry.New()
		if err := customerspb.Registrations(reg); err != nil {
			return nil, err
		}
		return reg, nil
	})
	stream := jetstream.NewStream(svc.Config().Nats.Stream, svc.JS(), svc.Logger())
	container.AddSingleton(constants.DomainDispatcherKey, func(c di.Container) (any, error) {
		return ddd.NewEventDispatcher[ddd.AggregateEvent](), nil
	})
	container.AddScoped(constants.DatabaseTransactionKey, func(c di.Container) (any, error) {
		return svc.DB().Begin()
	})
	container.AddScoped(constants.CustomersRepoKey, func(c di.Container) (any, error) {
		return postgres.NewCustomerRepository(
			constants.CustomersTableName,
			postgresotel.Trace(c.Get(constants.DatabaseTransactionKey).(*sql.Tx)),
		), nil
	})
	sentCounter := amprom.SentMessagesCounter(constants.ServiceName)
	container.AddScoped(constants.MessagePublisherKey, func(c di.Container) (any, error) {
		tx := postgresotel.Trace(c.Get(constants.DatabaseTransactionKey).(*sql.Tx))
		outboxStore := pg.NewOutboxStore(constants.OutboxTableName, tx)
		return am.NewMessagePublisher(
			stream,
			amotel.OtelMessageContextInjector(),
			sentCounter,
			tm.OutboxPublisher(outboxStore),
		), nil
	})
	container.AddSingleton(constants.MessageSubscriberKey, func(c di.Container) (any, error) {
		return am.NewMessageSubscriber(
			stream,
			amotel.OtelMessageContextExtractor(),
			amprom.ReceivedMessagesCounter(constants.ServiceName),
		), nil
	})
	container.AddScoped(constants.EventPublisherKey, func(c di.Container) (any, error) {
		return am.NewEventPublisher(
			c.Get(constants.RegistryKey).(registry.Registry),
			c.Get(constants.MessagePublisherKey).(am.MessagePublisher),
		), nil
	})
	container.AddScoped(constants.ReplyPublisherKey, func(c di.Container) (any, error) {
		return am.NewReplyPublisher(
			c.Get(constants.RegistryKey).(registry.Registry),
			c.Get(constants.MessagePublisherKey).(am.MessagePublisher),
		), nil
	})
	container.AddScoped(constants.InboxStoreKey, func(c di.Container) (any, error) {
		tx := postgresotel.Trace(c.Get(constants.DatabaseTransactionKey).(*sql.Tx))
		return pg.NewInboxStore(constants.InboxTableName, tx), nil
	})
	// Prometheus counters
	customersRegistered := promauto.NewCounter(prometheus.CounterOpts{
		Name: constants.CustomersRegisteredCount,
	})

	// setup application
	container.AddScoped(constants.ApplicationKey, func(c di.Container) (any, error) {
		return application.NewInstrumentedApp(application.New(
			c.Get(constants.CustomersRepoKey).(domain.CustomerRepository),
			c.Get(constants.DomainDispatcherKey).(*ddd.EventDispatcher[ddd.AggregateEvent]),
		), customersRegistered), nil
	})
	container.AddScoped(constants.DomainEventHandlersKey, func(c di.Container) (any, error) {
		return handlers.NewDomainEventHandlers(c.Get(constants.EventPublisherKey).(am.EventPublisher)), nil
	})
	container.AddScoped(constants.CommandHandlersKey, func(c di.Container) (any, error) {
		return handlers.NewCommandHandlers(
			c.Get(constants.RegistryKey).(registry.Registry),
			c.Get(constants.ApplicationKey).(application.App),
			c.Get(constants.ReplyPublisherKey).(am.ReplyPublisher),
			tm.InboxHandler(c.Get(constants.InboxStoreKey).(tm.InboxStore)),
		), nil
	})
	outboxProcessor := tm.NewOutboxProcessor(
		stream,
		pg.NewOutboxStore(constants.OutboxTableName, svc.DB()),
	)

	// setup Driver adapters
	if err = grpc.RegisterServerTx(container, svc.RPC()); err != nil {
		return err
	}
	if err = rest.RegisterGateway(ctx, svc.Mux(), svc.Config().Rpc.Address()); err != nil {
		return err
	}
	if err = rest.RegisterSwagger(svc.Mux()); err != nil {
		return err
	}
	handlers.RegisterDomainEventHandlersTx(container)
	if err = handlers.RegisterCommandHandlersTx(container); err != nil {
		return err
	}
	startOutboxProcessor(ctx, outboxProcessor, svc.Logger())

	return nil
}

func startOutboxProcessor(ctx context.Context, outboxProcessor tm.OutboxProcessor, logger zerolog.Logger) {
	go func() {
		err := outboxProcessor.Start(ctx)
		if err != nil {
			logger.Error().Err(err).Msg("customers outbox processor encountered an error")
		}
	}()
}
