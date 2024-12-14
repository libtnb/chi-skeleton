package bootstrap

import "github.com/google/wire"

// ProviderSet is bootstrap providers.
var ProviderSet = wire.NewSet(NewConf, NewRouter, NewHttp, NewDB, NewMigrate, NewValidator, NewTranslator)
