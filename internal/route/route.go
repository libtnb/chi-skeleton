package route

import "github.com/google/wire"

// ProviderSet is route providers.
var ProviderSet = wire.NewSet(NewHttp)
