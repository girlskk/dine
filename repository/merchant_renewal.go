package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/ent"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var _ domain.MerchantRenewalRepository = (*MerchantRenewalRepository)(nil)

type MerchantRenewalRepository struct {
	Client *ent.Client
}

func NewMerchantRenewalRepository(client *ent.Client) *MerchantRenewalRepository {
	return &MerchantRenewalRepository{
		Client: client,
	}
}

func (repo *MerchantRenewalRepository) GetByMerchant(ctx context.Context, merchantId uuid.UUID) (renewals []*domain.MerchantRenewal, err error) {
	//TODO implement me
	return
}

func (repo *MerchantRenewalRepository) Create(ctx context.Context, merchantRenewal *domain.MerchantRenewal) (err error) {
	span, ctx := util.StartSpan(ctx, "repository", "MerchantRenewalRepository.Create")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	if merchantRenewal == nil {
		err = fmt.Errorf("merchantRenewal is nil")
		return
	}

	_, err = repo.Client.MerchantRenewal.Create().SetID(merchantRenewal.ID).
		SetMerchantID(merchantRenewal.MerchantID).
		SetPurchaseDuration(merchantRenewal.PurchaseDuration).
		SetPurchaseDurationUnit(merchantRenewal.PurchaseDurationUnit).
		SetOperatorName(merchantRenewal.OperatorName).
		SetOperatorAccount(merchantRenewal.OperatorAccount).
		Save(ctx)
	if err != nil {
		err = fmt.Errorf("failed to create merchantRenewal: %w", err)
		return
	}
	return
}
