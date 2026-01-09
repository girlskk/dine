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
	"gitlab.jiguang.dev/pos-dine/dine/pkg/logging"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var _ domain.AdminUserInteractor = (*AdminUserInteractor)(nil)

type AdminUserInteractor struct {
	AuthConfig domain.AuthConfig
	DS         domain.DataStore
}

func NewAdminUserInteractor(authConfig domain.AuthConfig, ds domain.DataStore) *AdminUserInteractor {
	return &AdminUserInteractor{
		AuthConfig: authConfig,
		DS:         ds,
	}
}

func (interactor *AdminUserInteractor) Login(ctx context.Context, username, password string) (token string, expAt time.Time, err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "AdminUserInteractor.Login")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	user, err := interactor.DS.AdminUserRepo().FindByUsername(ctx, username)
	if err != nil {
		err = fmt.Errorf("failed to find user by username: %w", err)
		return
	}
	if err = user.CheckPassword(password); err != nil {
		return
	}
	if !user.Enabled {
		err = domain.ErrUserDisabled
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

func (interactor *AdminUserInteractor) Logout(ctx context.Context) error { _ = ctx; return nil }

func (interactor *AdminUserInteractor) Authenticate(ctx context.Context, token string) (user *domain.AdminUser, err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "AdminUserInteractor.Authenticate")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	logger := logging.FromContext(ctx).Named("AdminUserInteractor.Authenticate")

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

	user, err = interactor.DS.AdminUserRepo().Find(ctx, claims.ID)
	if err != nil {
		if domain.IsNotFound(err) {
			err = domain.ErrTokenInvalid
			return
		}
		err = fmt.Errorf("failed to find user[%d]: %w", claims.ID, err)
		return
	}

	return
}

func (interactor *AdminUserInteractor) Create(ctx context.Context, user *domain.AdminUser) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "AdminUserInteractor.Create")
	defer func() { util.SpanErrFinish(span, err) }()

	return interactor.DS.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		// check username exists
		exists, err := ds.AdminUserRepo().Exists(ctx, domain.AdminUserExistsParams{
			Username: user.Username,
		})
		if err != nil {
			return err
		}
		if exists {
			return domain.ErrUsernameExist
		}
		// check department assigned
		if user.DepartmentID == uuid.Nil {
			return domain.ErrUserDepartmentRequired
		}
		department, err := ds.DepartmentRepo().FindByID(ctx, user.DepartmentID)
		if err != nil {
			return err
		}
		if !department.Enable {
			return domain.ErrDepartmentDisabled
		}
		if string(department.DepartmentType) != string(domain.UserTypeAdmin) {
			return domain.ErrUserDepartmentTypeMismatch
		}
		// check role assigned
		if len(user.RoleIDs) == 0 {
			return domain.ErrUserRoleRequired
		}
		role, err := ds.RoleRepo().FindByID(ctx, user.RoleIDs[0])
		if err != nil {
			return err
		}
		if !role.Enable {
			return domain.ErrRoleDisabled
		}
		if string(role.RoleType) != string(domain.UserTypeAdmin) {
			return domain.ErrUserRoleTypeMismatch
		}
		user.ID = uuid.New()
		err = ds.AdminUserRepo().Create(ctx, user)
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

func (interactor *AdminUserInteractor) Update(ctx context.Context, user *domain.AdminUser) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "AdminUserInteractor.Update")
	defer func() { util.SpanErrFinish(span, err) }()

	if len(user.RoleIDs) == 0 {
		return domain.ErrUserRoleRequired
	}
	return interactor.DS.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		exists, err := ds.AdminUserRepo().Exists(ctx, domain.AdminUserExistsParams{
			Username:  user.Username,
			ExcludeID: user.ID,
		})
		if err != nil {
			return err
		}
		if exists {
			return domain.ErrUsernameExist
		}
		oldUser, err := ds.AdminUserRepo().Find(ctx, user.ID)
		if err != nil {
			return err
		}
		if oldUser.IsSuperAdmin {
			return domain.ErrSuperUserCannotUpdate
		}
		// check department assigned
		if user.DepartmentID == uuid.Nil {
			return domain.ErrUserDepartmentRequired
		}
		if user.DepartmentID != oldUser.DepartmentID {
			department, err := ds.DepartmentRepo().FindByID(ctx, user.DepartmentID)
			if err != nil {
				return err
			}
			if !department.Enable {
				return domain.ErrDepartmentDisabled
			}
			if string(department.DepartmentType) != string(domain.UserTypeAdmin) {
				return domain.ErrUserDepartmentTypeMismatch
			}
		}

		// check role assigned
		if len(user.RoleIDs) == 0 {
			return domain.ErrUserRoleRequired
		}
		role, err := ds.RoleRepo().FindByID(ctx, user.RoleIDs[0])
		if err != nil {
			return err
		}
		if !role.Enable {
			return domain.ErrRoleDisabled
		}
		if string(role.RoleType) != string(domain.UserTypeAdmin) {
			return domain.ErrUserRoleTypeMismatch
		}
		oldUser.Username = user.Username
		oldUser.Nickname = user.Nickname
		oldUser.DepartmentID = user.DepartmentID
		oldUser.RealName = user.RealName
		oldUser.Gender = user.Gender
		oldUser.Email = user.Email
		oldUser.PhoneNumber = user.PhoneNumber
		oldUser.Enabled = user.Enabled
		oldUser.RoleIDs = user.RoleIDs
		err = ds.AdminUserRepo().Update(ctx, oldUser)
		if err != nil {
			return err
		}

		userRole, err := ds.UserRoleRepo().FindOneByUser(ctx, user)
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

func (interactor *AdminUserInteractor) Delete(ctx context.Context, id uuid.UUID) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "AdminUserInteractor.Delete")
	defer func() { util.SpanErrFinish(span, err) }()

	return interactor.DS.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		user, err := ds.AdminUserRepo().Find(ctx, id)
		if err != nil {
			return err
		}
		if user.IsSuperAdmin {
			return domain.ErrSuperUserCannotDelete
		}
		err = ds.UserRoleRepo().DeleteByUsers(ctx, domain.UserTypeAdmin, id)
		if err != nil {
			return err
		}
		err = ds.AdminUserRepo().Delete(ctx, id)
		if err != nil {
			return err
		}
		return nil
	})
}

func (interactor *AdminUserInteractor) GetUser(ctx context.Context, id uuid.UUID) (user *domain.AdminUser, err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "AdminUserInteractor.GetUser")
	defer func() { util.SpanErrFinish(span, err) }()

	err = interactor.DS.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		user, err = ds.AdminUserRepo().Find(ctx, id)
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

func (interactor *AdminUserInteractor) GetUsers(ctx context.Context, pager *upagination.Pagination, filter *domain.AdminUserListFilter, orderBys ...domain.AdminUserOrderBy) (users []*domain.AdminUser, total int, err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "AdminUserInteractor.GetUsers")
	defer func() { util.SpanErrFinish(span, err) }()

	err = interactor.DS.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		var userRoles []*domain.UserRole
		var roleUserIDs []uuid.UUID
		var roleIDs []uuid.UUID
		if filter.RoleID != uuid.Nil {
			userRoles, err = ds.UserRoleRepo().GetByRoleIDs(ctx, domain.UserTypeAdmin, filter.RoleID)
			if err != nil {
				return err
			}
			roleUserIDs = lo.Map(userRoles, func(item *domain.UserRole, _ int) uuid.UUID { return item.UserID })
			if len(roleUserIDs) == 0 {
				return nil
			}
			filter.UserIDs = roleUserIDs

			roleIDs = lo.Map(userRoles, func(item *domain.UserRole, _ int) uuid.UUID { return item.RoleID })
		}
		users, total, err = ds.AdminUserRepo().GetUsers(ctx, pager, filter, orderBys...)
		if err != nil {
			return err
		}
		if len(users) == 0 {
			return nil
		}
		uIds := lo.Map(users, func(item *domain.AdminUser, _ int) uuid.UUID { return item.ID })
		if filter.RoleID == uuid.Nil {
			userRoles, err = ds.UserRoleRepo().GetByUserIDs(ctx, domain.UserTypeAdmin, uIds...)
			if err != nil {
				return err
			}
			roleIDs = lo.Map(userRoles, func(item *domain.UserRole, _ int) uuid.UUID { return item.RoleID })
		}
		if len(roleIDs) == 0 {
			return nil
		}
		roles, err := ds.RoleRepo().ListByIDs(ctx, roleIDs...)
		if err != nil {
			return err
		}
		if len(roles) == 0 {
			return nil
		}
		roleMap := lo.SliceToMap(roles, func(item *domain.Role) (uuid.UUID, *domain.Role) {
			return item.ID, item
		})
		if len(userRoles) == 0 {
			return nil
		}
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

// SimpleUpdate implements toggling simple fields for AdminUser (e.g., enabled)
func (interactor *AdminUserInteractor) SimpleUpdate(ctx context.Context, updateField domain.AdminUserSimpleUpdateField, params domain.AdminUserSimpleUpdateParams) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "AdminUserInteractor.SimpleUpdate")
	defer func() { util.SpanErrFinish(span, err) }()

	err = interactor.DS.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		user, err := ds.AdminUserRepo().Find(ctx, params.ID)
		if err != nil {
			if domain.IsNotFound(err) {
				return domain.ParamsError(domain.ErrUserNotExists)
			}
			return err
		}
		switch updateField {
		case domain.AdminUserSimpleUpdateFieldEnable:
			if user.Enabled == params.Enabled {
				return nil
			}
			if !params.Enabled {
				return domain.ErrSuperUserCannotDisable
			}
			user.Enabled = params.Enabled
		case domain.AdminUserSimpleUpdateFieldPassword:
			if params.Password == "" {
				return domain.ErrPasswordCannotBeEmpty
			}
			err = user.SetPassword(params.Password)
			if err != nil {
				return fmt.Errorf("failed to set password: %w", err)
			}
		default:
			return domain.ParamsError(fmt.Errorf("unsupported update field: %v", updateField))
		}
		err = interactor.DS.AdminUserRepo().Update(ctx, user)
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
