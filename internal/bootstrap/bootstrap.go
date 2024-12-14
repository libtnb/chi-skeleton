package bootstrap

import "github.com/google/wire"

// ProviderSet is bootstrap providers.
var ProviderSet = wire.NewSet(NewConf, NewLog, NewCli, NewRouter, NewHttp, NewDB, NewMigrate, NewSession, NewCron, NewCrypter)
