package repository

import (
	"context"

	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/ent"
	"gitlab.jiguang.dev/pos-dine/dine/ent/schema/schematype"
	setmealdetail "gitlab.jiguang.dev/pos-dine/dine/ent/setmealdetail"
	setmealgroup "gitlab.jiguang.dev/pos-dine/dine/ent/setmealgroup"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var _ domain.SetMealGroupRepository = (*SetMealGroupRepository)(nil)

type SetMealGroupRepository struct {
	Client *ent.Client
}

func NewSetMealGroupRepository(client *ent.Client) *SetMealGroupRepository {
	return &SetMealGroupRepository{
		Client: client,
	}
}

// CreateGroups 批量创建套餐组
func (repo *SetMealGroupRepository) CreateGroups(ctx context.Context, groups []*domain.SetMealGroup) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "SetMealGroupRepository.CreateGroups")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	if len(groups) == 0 {
		return nil
	}

	builders := make([]*ent.SetMealGroupCreate, 0, len(groups))
	for _, group := range groups {
		builder := repo.Client.SetMealGroup.Create().
			SetID(group.ID).
			SetProductID(group.ProductID).
			SetName(group.Name).
			SetSelectionType(group.SelectionType)

		builders = append(builders, builder)
	}

	_, err = repo.Client.SetMealGroup.CreateBulk(builders...).Save(ctx)
	return err
}

// CreateDetails 批量创建套餐组详情
func (repo *SetMealGroupRepository) CreateDetails(ctx context.Context, details []*domain.SetMealDetail) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "SetMealGroupRepository.CreateDetails")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	if len(details) == 0 {
		return nil
	}

	builders := make([]*ent.SetMealDetailCreate, 0, len(details))
	for _, detail := range details {
		builder := repo.Client.SetMealDetail.Create().
			SetID(detail.ID).
			SetGroupID(detail.GroupID).
			SetProductID(detail.ProductID).
			SetQuantity(detail.Quantity).
			SetIsDefault(detail.IsDefault)

		if len(detail.OptionalProductIDs) > 0 {
			builder = builder.SetOptionalProductIds(detail.OptionalProductIDs)
		}

		builders = append(builders, builder)
	}

	_, err = repo.Client.SetMealDetail.CreateBulk(builders...).Save(ctx)
	return err
}

// DeleteByProductID 根据商品ID删除套餐组和套餐组详情（物理删除）
func (repo *SetMealGroupRepository) DeleteByProductID(ctx context.Context, productID uuid.UUID) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "SetMealGroupRepository.DeleteByProductID")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	// 使用 SkipSoftDelete context 确保物理删除（硬删除）
	skipSoftDeleteCtx := schematype.SkipSoftDelete(ctx)

	// 查询该商品下的所有套餐组ID
	groupIDs, err := repo.Client.SetMealGroup.Query().
		Where(setmealgroup.ProductIDEQ(productID)).
		IDs(skipSoftDeleteCtx)
	if err != nil {
		return err
	}

	// 先删除套餐组详情（因为套餐组详情有外键关联到套餐组）
	if len(groupIDs) > 0 {
		_, err = repo.Client.SetMealDetail.Delete().
			Where(setmealdetail.GroupIDIn(groupIDs...)).
			Exec(skipSoftDeleteCtx)
		if err != nil {
			return err
		}
	}

	// 再删除套餐组
	_, err = repo.Client.SetMealGroup.Delete().
		Where(setmealgroup.ProductIDEQ(productID)).
		Exec(skipSoftDeleteCtx)
	if err != nil {
		return err
	}

	return nil
}

// CheckProductBelongToSetMeal 检查商品是否属于套餐组
func (repo *SetMealGroupRepository) CheckProductBelongToSetMeal(ctx context.Context, productID uuid.UUID) (res bool, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "SetMealGroupRepository.CheckProductBelongToSetMeal")
	defer func() {
		util.SpanErrFinish(span, err)
	}()
	exist, err := repo.Client.SetMealDetail.Query().
		Where(setmealdetail.ProductIDEQ(productID)).
		Exist(ctx)
	if err != nil {
		return false, err
	}
	return exist, nil
}

// ============================================
// 转换函数
// ============================================

func convertSetMealGroupToDomain(eg *ent.SetMealGroup) *domain.SetMealGroup {
	if eg == nil {
		return nil
	}

	group := &domain.SetMealGroup{
		ID:            eg.ID,
		Name:          eg.Name,
		ProductID:     eg.ProductID,
		SelectionType: eg.SelectionType,
		CreatedAt:     eg.CreatedAt,
		UpdatedAt:     eg.UpdatedAt,
	}

	for _, detail := range eg.Edges.Details {
		group.Details = append(group.Details, convertSetMealDetailToDomain(detail))
	}

	return group
}

func convertSetMealDetailToDomain(ed *ent.SetMealDetail) *domain.SetMealDetail {
	if ed == nil {
		return nil
	}

	detail := &domain.SetMealDetail{
		ID:                 ed.ID,
		GroupID:            ed.GroupID,
		ProductID:          ed.ProductID,
		Quantity:           ed.Quantity,
		IsDefault:          ed.IsDefault,
		OptionalProductIDs: ed.OptionalProductIds,
		CreatedAt:          ed.CreatedAt,
		UpdatedAt:          ed.UpdatedAt,
	}

	return detail
}
