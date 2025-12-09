package repository

import (
	"context"
	"database/sql"
	"fmt"
	"sync"

	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/ent"
)

var _ domain.DataStore = (*Repository)(nil)

type Repository struct {
	transactionActive        bool
	hooks                    []func()
	mu                       sync.Mutex
	client                   *ent.Client
	frontedUserRepo          *FrontendUserRepository
	backendUserRepo          *BackendUserRepository
	productUnitRepo          *ProductUnitRepository
	productAttrRepo          *ProductAttrRepository
	adminUserRepo            *AdminUserRepository
	productCategoryRepo      *CategoryRepository
	productRecipeRepo        *ProductRecipeRepository
	productSpecRepo          *ProductSpecRepository
	productRepo              *ProductRepository
	productSpecRelRepo       *ProductSpecRelRepository
	setmealDetailRepo        *SetMealDetailRepository
	tableareaRepo            *TableAreaRepository
	tableRepo                *TableRepository
	orderRepo                *OrderRepository
	storeRepo                *StoreRepository
	paymentRepo              *PaymentRepository
	reconciliationRecordRepo *ReconciliationRecordRepository
	pointSettlementRepo      *PointSettlementRepository
	storeAccountRepo         *StoreAccountRepository
	dataExportRepo           *DataExportRepository
	storeWithdrawRepo        *StoreWithdrawRepository
	customerRepo             *CustomerRepository
	orderCartRepo            *OrderCartRepository
}

func (repo *Repository) IsTransactionActive() bool {
	return repo.transactionActive
}

func (repo *Repository) AddHook(hook func()) {
	if !repo.transactionActive {
		panic("transaction not started")
	}
	repo.mu.Lock()
	defer repo.mu.Unlock()
	repo.hooks = append(repo.hooks, hook)
}

func (repo *Repository) executeHooks() {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	for _, hook := range repo.hooks {
		hook()
	}
}

func (repo *Repository) TableAreaRepo() domain.TableAreaRepository {
	return repo.tableareaRepo
}

func (repo *Repository) ProductAttrRepo() domain.ProductAttrRepository {
	return repo.productAttrRepo
}

func (repo *Repository) ProductUnitRepo() domain.ProductUnitRepository {
	return repo.productUnitRepo
}

func (repo *Repository) ProductCategoryRepo() domain.CategoryRepository {
	return repo.productCategoryRepo
}

func (repo *Repository) ProductRecipeRepo() domain.ProductRecipeRepository {
	return repo.productRecipeRepo
}

func (repo *Repository) ProductSpecRepo() domain.ProductSpecRepository {
	return repo.productSpecRepo
}

func (repo *Repository) ProductRepo() domain.ProductRepository {
	return repo.productRepo
}

func (repo *Repository) ProductSpecRelRepo() domain.ProductSpecRelRepository {
	return repo.productSpecRelRepo
}

func (repo *Repository) SetMealDetailRepo() domain.SetMealDetailRepository {
	return repo.setmealDetailRepo
}

func (repo *Repository) TableRepo() domain.TableRepository {
	return repo.tableRepo
}

func (repo *Repository) OrderRepo() domain.OrderRepository {
	return repo.orderRepo
}

func (repo *Repository) BackendUserRepo() domain.BackendUserRepository {
	return repo.backendUserRepo
}

func (repo *Repository) StoreRepo() domain.StoreRepository {
	return repo.storeRepo
}

func (repo *Repository) PaymentRepo() domain.PaymentRepository {
	return repo.paymentRepo
}

func (repo *Repository) ReconciliationRecordRepo() domain.ReconciliationRecordRepository {
	return repo.reconciliationRecordRepo
}

func (repo *Repository) PointSettlementRepo() domain.PointSettlementRepository {
	return repo.pointSettlementRepo
}

func (repo *Repository) StoreAccountRepo() domain.StoreAccountRepository {
	return repo.storeAccountRepo
}

func (repo *Repository) DataExportRepo() domain.DataExportRepository {
	return repo.dataExportRepo
}

func (repo *Repository) StoreWithdrawRepo() domain.StoreWithdrawRepository {
	return repo.storeWithdrawRepo
}

func (repo *Repository) CustomerRepo() domain.CustomerRepository {
	return repo.customerRepo
}

func (repo *Repository) OrderCartRepo() domain.OrderCartRepository {
	return repo.orderCartRepo
}

func New(client *ent.Client) *Repository {
	return &Repository{
		client:                   client,
		frontedUserRepo:          NewFrontendUserRepository(client),
		productUnitRepo:          NewProductUnitRepository(client),
		productAttrRepo:          NewProductAttrRepository(client),
		adminUserRepo:            NewAdminUserRepository(client),
		productCategoryRepo:      NewCategoryRepository(client),
		productRecipeRepo:        NewProductRecipeRepository(client),
		productSpecRepo:          NewProductSpecRepository(client),
		productRepo:              NewProductRepository(client),
		productSpecRelRepo:       NewProductSpecRelRepository(client),
		setmealDetailRepo:        NewSetMealDetailRepository(client),
		tableareaRepo:            NewTableAreaRepository(client),
		tableRepo:                NewTableRepository(client),
		orderRepo:                NewOrderRepository(client),
		backendUserRepo:          NewBackendUserRepository(client),
		storeRepo:                NewStoreRepository(client),
		paymentRepo:              NewPaymentRepository(client),
		reconciliationRecordRepo: NewReconciliationRecordRepository(client),
		pointSettlementRepo:      NewPointSettlementRepository(client),
		storeAccountRepo:         NewStoreAccountRepository(client),
		dataExportRepo:           NewDataExportRepository(client),
		storeWithdrawRepo:        NewStoreWithdrawRepository(client),
		customerRepo:             NewCustomerRepository(client),
		orderCartRepo:            NewOrderCartRepository(client),
	}
}

func (repo *Repository) Atomic(ctx context.Context, fn func(ctx context.Context, ds domain.DataStore) error) error {
	if repo.transactionActive {
		panic("transaction already started")
	}

	var r *Repository
	defer func() {
		if r == nil || len(r.hooks) == 0 {
			return
		}

		r.executeHooks()
	}()

	return withTx(ctx, repo.client, func(tx *ent.Tx) error {
		r = New(tx.Client())
		r.transactionActive = true
		return fn(ctx, r)
	})
}

func (repo *Repository) GetFrontendUserRepository() domain.FrontendUserRepository {
	return repo.frontedUserRepo
}

func (repo *Repository) GetAdminUserRepository() domain.AdminUserRepository {
	return repo.adminUserRepo
}

func withTx(ctx context.Context, client *ent.Client, fn func(tx *ent.Tx) error) error {
	tx, err := client.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelRepeatableRead})
	if err != nil {
		return err
	}
	defer func() {
		if v := recover(); v != nil {
			_ = tx.Rollback()
			panic(v)
		}
	}()
	if err := fn(tx); err != nil {
		if rerr := tx.Rollback(); rerr != nil {
			err = fmt.Errorf("rolling back transaction: %w", rerr)
		}
		return err
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}
	return nil
}
