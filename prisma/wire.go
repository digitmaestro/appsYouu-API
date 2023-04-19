package prisma

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewClient,
	NewConfig,
)
