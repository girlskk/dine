package repository

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/ent"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

var _ domain.CustomerRepository = (*CustomerRepository)(nil)

type CustomerRepository struct {
	Client *ent.Client
}

func NewCustomerRepository(client *ent.Client) *CustomerRepository {
	return &CustomerRepository{
		Client: client,
	}
}

func (repo *CustomerRepository) Find(ctx context.Context, id int) (u *domain.Customer, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CustomerRepository.Find")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	customer, err := repo.Client.Customer.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			err = domain.NotFoundError(err)
		}
		return nil, err
	}

	return convertCustomer(customer), nil
}

func (repo *CustomerRepository) FindOrCreate(ctx context.Context, customer *domain.Customer) (id int, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CustomerRepository.FindOrCreate")
	defer func() {
		util.SpanErrFinish(span, err)
	}()

	id, err = repo.Client.Customer.Create().
		SetNickname(customer.Nickname).
		SetPhone(customer.Phone).
		SetAvatar(customer.Avatar).
		SetGender(customer.Gender).
		OnConflict().
		UpdateUpdatedAt().
		ID(ctx)

	return id, err
}

func convertCustomer(eu *ent.Customer) *domain.Customer {
	if eu == nil {
		return nil
	}

	return &domain.Customer{
		ID:        eu.ID,
		Nickname:  eu.Nickname,
		Phone:     eu.Phone,
		Avatar:    eu.Avatar,
		Gender:    eu.Gender,
		CreatedAt: eu.CreatedAt,
		UpdatedAt: eu.UpdatedAt,
	}
}
