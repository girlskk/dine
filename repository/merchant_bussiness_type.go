package repository

import (
	"context"

	"github.com/samber/lo"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/ent"
	"gitlab.jiguang.dev/pos-dine/dine/ent/merchantbusinesstype"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var _ domain.MerchantBusinessTypeRepository = (*MerchantBusinessTypeRepository)(nil)

type MerchantBusinessTypeRepository struct {
	Client *ent.Client
}

func NewMerchantBusinessTypeRepository(client *ent.Client) *MerchantBusinessTypeRepository {
	return &MerchantBusinessTypeRepository{
		Client: client,
	}
}

func (repo MerchantBusinessTypeRepository) FindById(ctx context.Context, id int) (businessType *domain.MerchantBusinessType, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "MerchantRepository.FindByID")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	em, err := repo.Client.MerchantBusinessType.Query().
		Where(merchantbusinesstype.ID(id)).
		Only(ctx)
	if err != nil {
		return
	}
	return convertMerchantBusinessType(em), nil
}

func (repo MerchantBusinessTypeRepository) GetAll(ctx context.Context) (ts []*domain.MerchantBusinessType, err error) {
	span, ctx := util.StartSpan(ctx, "repository", "MerchantRepository.GetAll")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	ems, err := repo.Client.MerchantBusinessType.Query().All(ctx)
	if err != nil {
		return
	}
	ts = lo.Map(ems, func(em *ent.MerchantBusinessType, _ int) *domain.MerchantBusinessType {
		return convertMerchantBusinessType(em)
	})
	return
}

func convertMerchantBusinessType(em *ent.MerchantBusinessType) *domain.MerchantBusinessType {
	return &domain.MerchantBusinessType{
		ID:       em.ID,
		TypeCode: em.TypeCode,
		TypeName: em.TypeName,
	}
}
