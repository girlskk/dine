package domain

import "context"

//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/datastore.go -package=mock . DataStore
type DataStore interface {
	Atomic(ctx context.Context, fn func(ctx context.Context, ds DataStore) error) error
	IsTransactionActive() bool
	AddHook(hook func())
	GetFrontendUserRepository() FrontendUserRepository
	BackendUserRepo() BackendUserRepository
	ProductUnitRepo() ProductUnitRepository
	ProductAttrRepo() ProductAttrRepository
	GetAdminUserRepository() AdminUserRepository
	ProductCategoryRepo() CategoryRepository
	ProductRecipeRepo() ProductRecipeRepository
	ProductSpecRepo() ProductSpecRepository
	ProductRepo() ProductRepository
	ProductSpecRelRepo() ProductSpecRelRepository
	SetMealDetailRepo() SetMealDetailRepository
	TableAreaRepo() TableAreaRepository
	TableRepo() TableRepository
	OrderRepo() OrderRepository
	StoreRepo() StoreRepository
	PaymentRepo() PaymentRepository
	ReconciliationRecordRepo() ReconciliationRecordRepository
	PointSettlementRepo() PointSettlementRepository
	StoreAccountRepo() StoreAccountRepository
	DataExportRepo() DataExportRepository
	StoreWithdrawRepo() StoreWithdrawRepository
	CustomerRepo() CustomerRepository
	OrderCartRepo() OrderCartRepository
}

//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/datacache.go -package=mock . DataCache
type DataCache interface {
	SetFrontendUser(ctx context.Context, user *FrontendUser) error
	GetFrontendUser(ctx context.Context, id int) (*FrontendUser, error)
	SetStore(ctx context.Context, store *Store) error
	GetStore(ctx context.Context, id int) (*Store, error)
}
