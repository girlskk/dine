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

var _ domain.BackendUserInteractor = (*BackendUserInteractor)(nil)

type BackendUserInteractor struct {
	AuthConfig domain.AuthConfig
	DS         domain.DataStore
}

func NewBackendUserInteractor(authConfig domain.AuthConfig, dataStore domain.DataStore) *BackendUserInteractor {
	return &BackendUserInteractor{
		AuthConfig: authConfig,
		DS:         dataStore,
	}
}

func (interactor *BackendUserInteractor) Login(ctx context.Context, username, password string) (token string, expAt time.Time, err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "BackendUserInteractor.Login")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	user, err := interactor.DS.BackendUserRepo().FindByUsername(ctx, username)
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

func (interactor *BackendUserInteractor) Logout(ctx context.Context) error { return nil }

func (interactor *BackendUserInteractor) Authenticate(ctx context.Context, token string) (user *domain.BackendUser, err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "BackendUserInteractor.Authenticate")
	defer func() {
		util.SpanErrFinish(span, err)
	}()
	logger := logging.FromContext(ctx).Named("BackendUserInteractor.Authenticate")

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

	user, err = interactor.DS.BackendUserRepo().Find(ctx, claims.ID)
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

func (interactor *BackendUserInteractor) Create(ctx context.Context, user *domain.BackendUser) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "BackendUserInteractor.Create")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	return interactor.DS.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		exists, err := ds.AdminUserRepo().Exists(ctx, domain.AdminUserExistsParams{
			Username: user.Username,
		})
		if err != nil {
			return err
		}
		if exists {
			return domain.ErrUsernameExist
		}
		user.ID = uuid.New()
		err = ds.BackendUserRepo().Create(ctx, user)
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

func (interactor *BackendUserInteractor) Update(ctx context.Context, user *domain.BackendUser) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "BackendUserInteractor.Update")
	defer func() {
		util.SpanErrFinish(span, err)
	}()
	if len(user.RoleIDs) == 0 {
		return domain.ErrUserRoleRequired
	}
	return interactor.DS.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		exists, err := ds.BackendUserRepo().Exists(ctx, domain.BackendUserExistsParams{
			Username:  user.Username,
			ExcludeID: user.ID,
		})
		if err != nil {
			return err
		}
		if exists {
			return domain.ErrUsernameExist
		}
		userRole, err := interactor.DS.UserRoleRepo().FindOneByUser(ctx, user)
		if err != nil {
			return err
		}
		err = ds.BackendUserRepo().Update(ctx, user)
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

func (interactor *BackendUserInteractor) Delete(ctx context.Context, id uuid.UUID) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "BackendUserInteractor.Delete")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	return interactor.DS.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		err := ds.UserRoleRepo().DeleteByUsers(ctx, domain.UserTypeBackend, id)
		if err != nil {
			return err
		}
		err = ds.BackendUserRepo().Delete(ctx, id)
		if err != nil {
			return err
		}
		return nil
	})
}

func (interactor *BackendUserInteractor) GetUser(ctx context.Context, id uuid.UUID) (user *domain.BackendUser, err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "BackendUserInteractor.GetUser")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	err = interactor.DS.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		user, err = ds.BackendUserRepo().Find(ctx, id)
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

func (interactor *BackendUserInteractor) GetUsers(ctx context.Context, pager *upagination.Pagination, filter *domain.BackendUserListFilter, orderBys ...domain.BackendUserOrderBy) (users []*domain.BackendUser, total int, err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "BackendUserInteractor.GetUsers")
	defer func() { util.SpanErrFinish(span, err) }()

	err = interactor.DS.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		users, total, err = ds.BackendUserRepo().GetUsers(ctx, pager, filter, orderBys...)
		if err != nil {
			return err
		}
		uIds := lo.Map(users, func(item *domain.BackendUser, _ int) uuid.UUID { return item.ID })
		userRoles, err := ds.UserRoleRepo().GetByUserIDs(ctx, domain.UserTypeBackend, uIds...)
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

// SimpleUpdate implements toggling simple fields for BackendUser (e.g., enabled)
func (interactor *BackendUserInteractor) SimpleUpdate(ctx context.Context, updateField domain.BackendUserSimpleUpdateField, params domain.BackendUserSimpleUpdateParams) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "BackendUserInteractor.SimpleUpdate")
	defer func() { util.SpanErrFinish(span, err) }()

	err = interactor.DS.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		user, err := ds.BackendUserRepo().Find(ctx, params.ID)
		if err != nil {
			if domain.IsNotFound(err) {
				return domain.ParamsError(domain.ErrUserNotExists)
			}
			return err
		}
		switch updateField {
		case domain.BackendUserSimpleUpdateFieldEnable:
			if user.Enabled == params.Enabled {
				return nil
			}
			user.Enabled = params.Enabled
		default:
			return domain.ParamsError(fmt.Errorf("unsupported update field: %v", updateField))
		}
		err = interactor.DS.BackendUserRepo().Update(ctx, user)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return
	}
	return
}
