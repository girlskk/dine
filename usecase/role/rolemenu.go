package role

import (
	"context"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var _ domain.RoleMenuInteractor = (*RoleMenuInteractor)(nil)

type RoleMenuInteractor struct {
	DS domain.DataStore
}

func NewRoleMenuInteractor(ds domain.DataStore) *RoleMenuInteractor {
	return &RoleMenuInteractor{DS: ds}
}

func (interactor *RoleMenuInteractor) SetRoleMenu(ctx context.Context, roleID uuid.UUID, paths []string) (err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "RoleMenuInteractor.SetRoleMenu")
	defer func() { util.SpanErrFinish(span, err) }()

	err = interactor.DS.Atomic(ctx, func(ctx context.Context, ds domain.DataStore) error {
		role, err := ds.RoleRepo().FindByID(ctx, roleID)
		if err != nil {
			if domain.IsNotFound(err) {
				return domain.ErrRoleNotExists
			}
			return err
		}

		// fetch existing role-menu relations
		oldRelations, err := ds.RoleMenuRepo().GetByRoleID(ctx, roleID)
		if err != nil {
			return err
		}

		// build a slice of existing paths for quick lookup
		oldPaths := lo.Map(oldRelations, func(r *domain.RoleMenu, _ int) string { return r.Path })

		// deleteIDs: old relations whose path is not present in the new paths
		deleteIDs := lo.FilterMap(oldRelations, func(r *domain.RoleMenu, _ int) (uuid.UUID, bool) {
			if lo.Contains(paths, r.Path) {
				return uuid.Nil, false
			}
			return r.ID, true
		})

		// newPaths: paths that are not present in old relations
		newPaths := lo.Filter(paths, func(p string, _ int) bool {
			return !lo.Contains(oldPaths, p)
		})

		// dedupe newPaths just in case
		newPaths = lo.Uniq(newPaths)

		if len(newPaths) > 0 {
			if err := ds.RoleMenuRepo().CreateBulkByRoleIDPaths(ctx, role, newPaths); err != nil {
				return err
			}
		}

		if len(deleteIDs) > 0 {
			if err := ds.RoleMenuRepo().Deletes(ctx, deleteIDs); err != nil {
				return err
			}
		}

		return nil
	})
	return
}

func (interactor *RoleMenuInteractor) RoleMenuList(ctx context.Context, roleID uuid.UUID) (paths []string, err error) {
	span, ctx := util.StartSpan(ctx, "usecase", "RoleMenuInteractor.RoleMenuList")
	defer func() { util.SpanErrFinish(span, err) }()

	roleMenus, err := interactor.DS.RoleMenuRepo().GetByRoleID(ctx, roleID)
	if err != nil {
		return
	}

	paths = lo.Map(roleMenus, func(rm *domain.RoleMenu, _ int) string {
		return rm.Path
	})
	return
}
