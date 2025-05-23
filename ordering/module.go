package ordering

import (
	"context"
	"database/sql"

	"github.com/rs/zerolog"

	"github.com/50HJ/Intelli-Mall/baskets/basketspb"
	"github.com/50HJ/Intelli-Mall/depot/depotpb"
	"github.com/50HJ/Intelli-Mall/internal/am"
	"github.com/50HJ/Intelli-Mall/internal/amotel"
	"github.com/50HJ/Intelli-Mall/internal/amprom"
	"github.com/50HJ/Intelli-Mall/internal/ddd"
	"github.com/50HJ/Intelli-Mall/internal/di"
	"github.com/50HJ/Intelli-Mall/internal/es"
	"github.com/50HJ/Intelli-Mall/internal/jetstream"
	pg "github.com/50HJ/Intelli-Mall/internal/postgres"
	"github.com/50HJ/Intelli-Mall/internal/postgresotel"
	"github.com/50HJ/Intelli-Mall/internal/registry"
	"github.com/50HJ/Intelli-Mall/internal/registry/serdes"
	"github.com/50HJ/Intelli-Mall/internal/system"
	"github.com/50HJ/Intelli-Mall/internal/tm"
	"github.com/50HJ/Intelli-Mall/ordering/internal/application"
	"github.com/50HJ/Intelli-Mall/ordering/internal/constants"
	"github.com/50HJ/Intelli-Mall/ordering/internal/domain"
	"github.com/50HJ/Intelli-Mall/ordering/internal/grpc"
	"github.com/50HJ/Intelli-Mall/ordering/internal/handlers"
	"github.com/50HJ/Intelli-Mall/ordering/internal/rest"
	"github.com/50HJ/Intelli-Mall/ordering/orderingpb"
)

type Module struct{}

func (Module) Startup(ctx context.Context, mono system.Service) (err error) {
	return Root(ctx, mono)
}

func Root(ctx context.Context, svc system.Service) (err error) {
	container := di.New()
	// setup Driven adapters
	container.AddSingleton(constants.RegistryKey, func(c di.Container) (any, error) {
		reg := registry.New()
		if err := registrations(reg); err != nil {
			return nil, err
		}
		if err := basketspb.Registrations(reg); err != nil {
			return nil, err
		}
		if err := orderingpb.Registrations(reg); err != nil {
			return nil, err
		}
		if err := depotpb.Registrations(reg); err != nil {
			return nil, err
		}
		return reg, nil
	})
	stream := jetstream.NewStream(svc.Config().Nats.Stream, svc.JS(), svc.Logger())
	container.AddSingleton(constants.DomainDispatcherKey, func(c di.Container) (any, error) {
		return ddd.NewEventDispatcher[ddd.Event](), nil
	})
	container.AddScoped(constants.DatabaseTransactionKey, func(c di.Container) (any, error) {
		return svc.DB().Begin()
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
	container.AddScoped(constants.OrdersRepoKey, func(c di.Container) (any, error) {
		tx := postgresotel.Trace(c.Get(constants.DatabaseTransactionKey).(*sql.Tx))
		reg := c.Get(constants.RegistryKey).(registry.Registry)
		return es.NewAggregateRepository[*domain.Order](
			domain.OrderAggregate,
			c.Get(constants.RegistryKey).(registry.Registry),
			es.AggregateStoreWithMiddleware(
				pg.NewEventStore(constants.EventsTableName, tx, reg),
				pg.NewSnapshotStore(constants.SnapshotsTableName, tx, reg),
			),
		), nil
	})

	// setup application
	container.AddScoped(constants.ApplicationKey, func(c di.Container) (any, error) {
		return application.New(
			c.Get(constants.OrdersRepoKey).(domain.OrderRepository),
			c.Get(constants.DomainDispatcherKey).(*ddd.EventDispatcher[ddd.Event]),
		), nil
	})
	container.AddScoped(constants.DomainEventHandlersKey, func(c di.Container) (any, error) {
		return handlers.NewDomainEventHandlers(c.Get(constants.EventPublisherKey).(am.EventPublisher)), nil
	})
	container.AddScoped(constants.IntegrationEventHandlersKey, func(c di.Container) (any, error) {
		return handlers.NewIntegrationEventHandlers(
			c.Get(constants.RegistryKey).(registry.Registry),
			c.Get(constants.ApplicationKey).(application.App),
			tm.InboxHandler(c.Get(constants.InboxStoreKey).(tm.InboxStore)),
		), nil
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
	if err = handlers.RegisterIntegrationEventHandlersTx(container); err != nil {
		return err
	}
	if err = handlers.RegisterCommandHandlersTx(container); err != nil {
		return err
	}
	startOutboxProcessor(ctx, outboxProcessor, svc.Logger())

	return nil
}

func registrations(reg registry.Registry) (err error) {
	serde := serdes.NewJsonSerde(reg)

	// Order
	if err = serde.Register(domain.Order{}, func(v any) error {
		order := v.(*domain.Order)
		order.Aggregate = es.NewAggregate("", domain.OrderAggregate)
		return nil
	}); err != nil {
		return err
	}
	// order events
	if err = serde.Register(domain.OrderCreated{}); err != nil {
		return err
	}
	if err = serde.Register(domain.OrderRejected{}); err != nil {
		return err
	}
	if err = serde.Register(domain.OrderApproved{}); err != nil {
		return err
	}
	if err = serde.Register(domain.OrderCanceled{}); err != nil {
		return err
	}
	if err = serde.Register(domain.OrderReadied{}); err != nil {
		return err
	}
	if err = serde.Register(domain.OrderCompleted{}); err != nil {
		return err
	}
	// order snapshots
	if err = serde.RegisterKey(domain.OrderV1{}.SnapshotName(), domain.OrderV1{}); err != nil {
		return err
	}

	return nil
}

func startOutboxProcessor(ctx context.Context, outboxProcessor tm.OutboxProcessor, logger zerolog.Logger) {
	go func() {
		err := outboxProcessor.Start(ctx)
		if err != nil {
			logger.Error().Err(err).Msg("ordering outbox processor encountered an error")
		}
	}()
}
