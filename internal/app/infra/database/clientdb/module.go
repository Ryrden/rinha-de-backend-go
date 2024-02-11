package clientdb

import (
	"github.com/ryrden/rinha-de-backend-go/internal/app/domain/client"
	"go.uber.org/fx"
)

var Module = fx.Provide(
	// TODO: Add the NewDispatcher function to the list of providers
	//NewDispatcher,
	fx.Annotate(
		NewClientRepository,
		fx.As(new(client.Repository)),
	),
)
