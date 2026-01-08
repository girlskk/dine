package domain

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
)

var (
	ErrUserRoleNotExists         = errors.New("用户角色关系不存在")
	ErrUserRoleExists            = errors.New("用户角色关系已存在")
	ErrRoleAssignedCannotDisable = errors.New("角色已分配用户，无法禁用")
	ErrRoleAssignedCannotDelete  = errors.New("角色已分配用户，无法删除")
)

// UserRoleRepository 用户角色关系仓储接口
//
//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mock/user_role_repository.go -package=mock . UserRoleRepository
type UserRoleRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*UserRole, error)
	FindOneByUser(ctx context.Context, user User) (*UserRole, error)
	Create(ctx context.Context, userRole *UserRole) error
	CreateBulk(ctx context.Context, userRoles []*UserRole) error
	CreateBulkByRoleIDUsers(ctx context.Context, roleID uuid.UUID, users []User) error
	CreateBulkByUserIDRoles(ctx context.Context, user User, roles []uuid.UUID) error
	Update(ctx context.Context, userRole *UserRole) error
	Deletes(ctx context.Context, ids ...uuid.UUID) error
	DeleteByRoles(ctx context.Context, roleIDs ...uuid.UUID) error
	DeleteByUsers(ctx context.Context, userType UserType, userIDs ...uuid.UUID) error
	GetUserRoles(ctx context.Context, pager *upagination.Pagination, filter *UserRoleListFilter, orderBys ...UserRoleListOrderBy) ([]*UserRole, int, error)
	GetByRoleIDs(ctx context.Context, userType UserType, roleID ...uuid.UUID) ([]*UserRole, error)
	GetByUserIDs(ctx context.Context, userType UserType, userID ...uuid.UUID) ([]*UserRole, error)
}

type UserRoleListOrderByType int

const (
	_ UserRoleListOrderByType = iota
	UserRoleListOrderByID
	UserRoleListOrderByCreatedAt
)

type UserRoleListOrderBy struct {
	OrderBy UserRoleListOrderByType
	Desc    bool
}

func NewUserRoleListOrderByID(desc bool) UserRoleListOrderBy {
	return UserRoleListOrderBy{OrderBy: UserRoleListOrderByID, Desc: desc}
}

func NewUserRoleListOrderByCreatedAt(desc bool) UserRoleListOrderBy {
	return UserRoleListOrderBy{OrderBy: UserRoleListOrderByCreatedAt, Desc: desc}
}

type UserRole struct {
	ID         uuid.UUID `json:"id"`
	UserType   UserType  `json:"user_type"`
	UserID     uuid.UUID `json:"user_id"`
	RoleID     uuid.UUID `json:"role_id"`
	MerchantID uuid.UUID `json:"merchant_id"`
	StoreID    uuid.UUID `json:"store_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type UserRoleListFilter struct {
	UserType   UserType  `json:"user_type"`
	UserID     uuid.UUID `json:"user_id"`
	RoleID     uuid.UUID `json:"role_id"`
	MerchantID uuid.UUID `json:"merchant_id"`
	StoreID    uuid.UUID `json:"store_id"`
}
