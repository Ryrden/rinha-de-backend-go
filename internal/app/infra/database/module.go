package database

import (
	"github.com/ryrden/rinha-de-backend-go/internal/app/infra/database/clientdb"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(NewPostgresDatabase),
	clientdb.Module,
)
