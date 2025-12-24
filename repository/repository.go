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
	adminUserRepo            *AdminUserRepository
	categoryRepo             *CategoryRepository
	backendUserRepo          *BackendUserRepository
	merchantRepo             *MerchantRepository
	storeRepo                *StoreRepository
	merchantRenewalRepo      *MerchantRenewalRepository
	merchantBusinessTypeRepo *MerchantBusinessTypeRepository
	remarkRepo               *RemarkRepository
	remarkCategoryRepo       *RemarkCategoryRepository
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

func New(client *ent.Client) *Repository {
	return &Repository{
		client: client,
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

func (repo *Repository) AdminUserRepo() domain.AdminUserRepository {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	if repo.adminUserRepo == nil {
		repo.adminUserRepo = NewAdminUserRepository(repo.client)
	}
	return repo.adminUserRepo
}

func (repo *Repository) CategoryRepo() domain.CategoryRepository {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	if repo.categoryRepo == nil {
		repo.categoryRepo = NewCategoryRepository(repo.client)
	}
	return repo.categoryRepo
}

func (repo *Repository) BackendUserRepo() domain.BackendUserRepository {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	if repo.backendUserRepo == nil {
		repo.backendUserRepo = NewBackendUserRepository(repo.client)
	}
	return repo.backendUserRepo
}

func (repo *Repository) MerchantRepo() domain.MerchantRepository {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	if repo.merchantRepo == nil {
		repo.merchantRepo = NewMerchantRepository(repo.client)
	}
	return repo.merchantRepo
}

func (repo *Repository) StoreRepo() domain.StoreRepository {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	if repo.storeRepo == nil {
		repo.storeRepo = NewStoreRepository(repo.client)
	}
	return repo.storeRepo
}

func (repo *Repository) MerchantRenewalRepo() domain.MerchantRenewalRepository {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	if repo.merchantRenewalRepo == nil {
		repo.merchantRenewalRepo = NewMerchantRenewalRepository(repo.client)
	}
	return repo.merchantRenewalRepo
}

func (repo *Repository) MerchantBusinessTypeRepo() domain.MerchantBusinessTypeRepository {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	if repo.merchantBusinessTypeRepo == nil {
		repo.merchantBusinessTypeRepo = NewMerchantBusinessTypeRepository(repo.client)
	}
	return repo.merchantBusinessTypeRepo
}

func (repo *Repository) RemarkRepo() domain.RemarkRepository {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	if repo.remarkRepo == nil {
		repo.remarkRepo = NewRemarkRepository(repo.client)
	}
	return repo.remarkRepo
}

func (repo *Repository) RemarkCategoryRepo() domain.RemarkCategoryRepository {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	if repo.remarkCategoryRepo == nil {
		repo.remarkCategoryRepo = NewRemarkCategoryRepository(repo.client)
	}
	return repo.remarkCategoryRepo
}
