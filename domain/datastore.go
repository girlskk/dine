package domain

import "context"

//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/datastore.go -package=mock . DataStore
type DataStore interface {
	Atomic(ctx context.Context, fn func(ctx context.Context, ds DataStore) error) error
	IsTransactionActive() bool
	AddHook(hook func())
	AdminUserRepo() AdminUserRepository
	CategoryRepo() CategoryRepository
	BackendUserRepo() BackendUserRepository
	ProductUnitRepo() ProductUnitRepository
}

//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/datacache.go -package=mock . DataCache
type DataCache interface {
}
