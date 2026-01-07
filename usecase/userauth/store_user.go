package userauth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/ent"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/logging"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var _ domain.StoreUserInteractor = (*StoreUserInteractor)(nil)

type StoreUserInteractor struct {
	AuthConfig domain.AuthConfig
	DataStore  domain.DataStore
}

func NewStoreUserInteractor(authConfig domain.AuthConfig, dataStore domain.DataStore) *StoreUserInteractor {
	return &StoreUserInteractor{
		AuthConfig: authConfig,
		DataStore:  dataStore,
	}
}

func (interactor *StoreUserInteractor) Login(ctx context.Context, username, password string) (token string, expAt time.Time, err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "StoreUserInteractor.Login")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	user, err := interactor.DataStore.StoreUserRepo().FindByUsername(ctx, username)
	if err != nil {
		return
	}
	if err = user.CheckPassword(password); err != nil {
		return
	}

	expAt = time.Now().Add(time.Duration(interactor.AuthConfig.Expire) * time.Second)
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, AuthToken{
		ID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expAt),
		},
	})

	token, err = claims.SignedString([]byte(interactor.AuthConfig.Secret))
	if err != nil {
		err = fmt.Errorf("failed to sign token: %w", err)
		return
	}

	return
}

func (interactor *StoreUserInteractor) Logout(ctx context.Context) error { _ = ctx; return nil }

func (interactor *StoreUserInteractor) Authenticate(ctx context.Context, token string) (user *domain.StoreUser, err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "StoreUserInteractor.Authenticate")
	defer func() {
		util.SpanErrFinish(span, err)
	}()
	logger := logging.FromContext(ctx).Named("StoreUserInteractor.Authenticate")

	claims := &AuthToken{}
	tokenInfo, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(interactor.AuthConfig.Secret), nil
	})
	if err != nil {
		if !errors.Is(err, jwt.ErrTokenNotValidYet) && !errors.Is(err, jwt.ErrTokenExpired) {
			logger.Errorw("failed to parse token", "error", err)
		}
		err = domain.ErrTokenInvalid
		return
	}

	if !tokenInfo.Valid {
		err = domain.ErrTokenInvalid
		return
	}

	user, err = interactor.DataStore.StoreUserRepo().Find(ctx, claims.ID)
	if err != nil {
		if ent.IsNotFound(err) {
			err = domain.ErrTokenInvalid
			return
		}
		err = fmt.Errorf("failed to find user[%d]: %w", claims.ID, err)
		return
	}

	return
}

func (interactor *StoreUserInteractor) Create(ctx context.Context, user *domain.StoreUser) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "StoreUserInteractor.Create")
	defer func() { util.SpanErrFinish(span, err) }()

	return interactor.DataStore.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		exists, err := ds.StoreUserRepo().Exists(ctx, domain.StoreUserExistsParams{
			Username: user.Username,
		})
		if err != nil {
			return err
		}
		if exists {
			return fmt.Errorf("%w", domain.ErrBackendUserUsernameExist)
		}
		user.ID = uuid.New()
		err = ds.StoreUserRepo().Create(ctx, user)
		if err != nil {
			return err
		}
		if len(user.RoleIDs) > 0 {
			err = ds.UserRoleRepo().CreateBulkByUserIDRoles(ctx, user, user.RoleIDs)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (interactor *StoreUserInteractor) Update(ctx context.Context, user *domain.StoreUser) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "StoreUserInteractor.Update")
	defer func() { util.SpanErrFinish(span, err) }()

	if len(user.RoleIDs) == 0 {
		return domain.ErrBackendUserRoleRequired
	}
	return interactor.DataStore.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		exists, err := ds.StoreUserRepo().Exists(ctx, domain.StoreUserExistsParams{
			Username:  user.Username,
			ExcludeID: user.ID,
		})
		if err != nil {
			return err
		}
		if exists {
			return domain.ErrBackendUserUsernameExist
		}

		userRole, err := ds.UserRoleRepo().FindOneByUser(ctx, user)
		if err != nil {
			return err
		}
		err = ds.StoreUserRepo().Update(ctx, user)
		if err != nil {
			return err
		}

		if userRole.RoleID != user.RoleIDs[0] {
			userRole.RoleID = user.RoleIDs[0]
			err = ds.UserRoleRepo().Update(ctx, userRole)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (interactor *StoreUserInteractor) Delete(ctx context.Context, id uuid.UUID) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "StoreUserInteractor.Delete")
	defer func() { util.SpanErrFinish(span, err) }()

	return interactor.DataStore.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		err := ds.UserRoleRepo().DeleteByUsers(ctx, domain.UserTypeStore, id)
		if err != nil {
			return err
		}
		err = ds.StoreUserRepo().Delete(ctx, id)
		if err != nil {
			return err
		}
		return nil
	})
}

func (interactor *StoreUserInteractor) GetUser(ctx context.Context, id uuid.UUID) (user *domain.StoreUser, err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "StoreUserInteractor.GetUser")
	defer func() { util.SpanErrFinish(span, err) }()

	err = interactor.DataStore.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		user, err = ds.StoreUserRepo().Find(ctx, id)
		if err != nil {
			return err
		}

		userRole, err := ds.UserRoleRepo().FindOneByUser(ctx, user)
		if err != nil {
			return err
		}
		user.RoleIDs = []uuid.UUID{userRole.RoleID}
		role, err := ds.RoleRepo().FindByID(ctx, userRole.RoleID)
		if err != nil {
			return err
		}
		user.RoleList = []*domain.Role{role}
		return nil
	})

	return
}

func (interactor *StoreUserInteractor) GetUsers(ctx context.Context, pager *upagination.Pagination, filter *domain.StoreUserListFilter, orderBys ...domain.StoreUserOrderBy) (users []*domain.StoreUser, total int, err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "StoreUserInteractor.GetUsers")
	defer func() { util.SpanErrFinish(span, err) }()

	err = interactor.DataStore.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		users, total, err = ds.StoreUserRepo().GetUsers(ctx, pager, filter, orderBys...)
		if err != nil {
			return err
		}
		uIds := lo.Map(users, func(item *domain.StoreUser, _ int) uuid.UUID { return item.ID })
		userRoles, err := ds.UserRoleRepo().GetByUserIDs(ctx, domain.UserTypeStore, uIds...)
		if err != nil {
			return err
		}
		roleIDs := lo.Map(userRoles, func(item *domain.UserRole, _ int) uuid.UUID { return item.RoleID })
		roles, err := ds.RoleRepo().ListByIDs(ctx, roleIDs...)
		if err != nil {
			return err
		}
		roleMap := lo.SliceToMap(roles, func(item *domain.Role) (uuid.UUID, *domain.Role) {
			return item.ID, item
		})
		userRoleMap := lo.SliceToMap(userRoles, func(item *domain.UserRole) (uuid.UUID, *domain.UserRole) {
			return item.UserID, item
		})
		for _, user := range users {
			if ur, ok := userRoleMap[user.ID]; ok {
				user.RoleIDs = []uuid.UUID{ur.RoleID}
				if role, ok := roleMap[ur.RoleID]; ok {
					user.RoleList = []*domain.Role{role}
				}
			}
		}
		return nil
	})

	return
}
