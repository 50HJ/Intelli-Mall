package commands

import (
	"context"

	"github.com/50HJ/Intelli-Mall/internal/ddd"
	"github.com/50HJ/Intelli-Mall/stores/internal/domain"
)

type DisableParticipation struct {
	ID string
}

type DisableParticipationHandler struct {
	stores    domain.StoreRepository
	publisher ddd.EventPublisher[ddd.Event]
}

func NewDisableParticipationHandler(stores domain.StoreRepository, publisher ddd.EventPublisher[ddd.Event]) DisableParticipationHandler {
	return DisableParticipationHandler{
		stores:    stores,
		publisher: publisher,
	}
}

func (h DisableParticipationHandler) DisableParticipation(ctx context.Context, cmd DisableParticipation) error {
	store, err := h.stores.Load(ctx, cmd.ID)
	if err != nil {
		return err
	}

	event, err := store.DisableParticipation()
	if err != nil {
		return err
	}

	err = h.stores.Save(ctx, store)
	if err != nil {
		return err
	}

	return h.publisher.Publish(ctx, event)
}
