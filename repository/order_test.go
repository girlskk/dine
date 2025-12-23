package repository

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/ent"
	entorder "gitlab.jiguang.dev/pos-dine/dine/ent/order"
	"gitlab.jiguang.dev/pos-dine/dine/ent/schema/schematype"
)

type OrderTestSuite struct {
	RepositoryTestSuite
	repo *OrderRepository
	ctx  context.Context
}

func TestOrderTestSuite(t *testing.T) {
	suite.Run(t, new(OrderTestSuite))
}

func (s *OrderTestSuite) SetupTest() {
	s.RepositoryTestSuite.SetupTest()
	s.repo = &OrderRepository{Client: s.client}
	s.ctx = context.Background()
}

func (s *OrderTestSuite) newTestOrder(storeID, orderNo string) *domain.Order {
	return &domain.Order{
		ID:           uuid.New(),
		MerchantID:   uuid.NewString(),
		StoreID:      storeID,
		BusinessDate: "2025-12-22",
		OrderNo:      orderNo,
		DiningMode:   "DINE_IN",
		Store:        json.RawMessage(`{"store_id":"` + storeID + `"}`),
		Cart:         json.RawMessage(`[]`),
		Products:     json.RawMessage(`[]`),
		Amount:       json.RawMessage(`{}`),
	}
}

func (s *OrderTestSuite) createEntOrder(storeID, orderNo string, createdAt time.Time, paymentStatus string) *ent.Order {
	builder := s.client.Order.Create().
		SetID(uuid.New()).
		SetMerchantID(uuid.NewString()).
		SetStoreID(storeID).
		SetBusinessDate("2025-12-22").
		SetOrderNo(orderNo).
		SetDiningMode(entorder.DiningMode("DINE_IN")).
		SetCreatedAt(createdAt)

	if paymentStatus != "" {
		builder = builder.SetPaymentStatus(entorder.PaymentStatus(paymentStatus))
	}

	return builder.SaveX(s.ctx)
}

func (s *OrderTestSuite) TestOrder_Create() {
	s.T().Run("创建成功", func(t *testing.T) {
		storeID := uuid.NewString()
		order := s.newTestOrder(storeID, "NO-001")

		err := s.repo.Create(s.ctx, order)
		require.NoError(t, err)
		require.NotEqual(t, uuid.Nil, order.ID)
		require.False(t, order.CreatedAt.IsZero())
		require.False(t, order.UpdatedAt.IsZero())

		dbOrder := s.client.Order.GetX(s.ctx, order.ID)
		require.Equal(t, order.MerchantID, dbOrder.MerchantID)
		require.Equal(t, order.StoreID, dbOrder.StoreID)
		require.Equal(t, order.BusinessDate, dbOrder.BusinessDate)
		require.Equal(t, order.OrderNo, dbOrder.OrderNo)
		require.Equal(t, entorder.DiningMode("DINE_IN"), dbOrder.DiningMode)
		require.Equal(t, []byte(order.Store), []byte(dbOrder.Store))
	})

	s.T().Run("非法枚举值返回 ParamsError", func(t *testing.T) {
		storeID := uuid.NewString()
		order := s.newTestOrder(storeID, "NO-ENUM")
		order.OrderStatus = "OPEN"

		err := s.repo.Create(s.ctx, order)
		require.Error(t, err)
		require.True(t, domain.IsParamsError(err))
	})

	s.T().Run("唯一键冲突返回 Conflict", func(t *testing.T) {
		storeID := uuid.NewString()

		o1 := s.newTestOrder(storeID, "NO-DUP")
		o2 := s.newTestOrder(storeID, "NO-DUP")

		require.NoError(t, s.repo.Create(s.ctx, o1))
		err := s.repo.Create(s.ctx, o2)
		require.Error(t, err)
		require.True(t, domain.IsConflict(err))
	})

	s.T().Run("软删后可复用相同 order_no", func(t *testing.T) {
		storeID := uuid.NewString()

		o1 := s.newTestOrder(storeID, "NO-REUSE")
		require.NoError(t, s.repo.Create(s.ctx, o1))
		require.NoError(t, s.repo.Delete(s.ctx, o1.ID))

		o2 := s.newTestOrder(storeID, "NO-REUSE")
		err := s.repo.Create(s.ctx, o2)
		require.NoError(t, err)
	})
}

func (s *OrderTestSuite) TestOrder_FindByID() {
	storeID := uuid.NewString()
	order := s.newTestOrder(storeID, "NO-GET")
	require.NoError(s.T(), s.repo.Create(s.ctx, order))

	s.T().Run("正常查询", func(t *testing.T) {
		found, err := s.repo.FindByID(s.ctx, order.ID)
		require.NoError(t, err)
		require.Equal(t, order.ID, found.ID)
		require.Equal(t, order.StoreID, found.StoreID)
		require.Equal(t, order.OrderNo, found.OrderNo)
		require.Equal(t, "DINE_IN", found.DiningMode)
	})

	s.T().Run("不存在的ID", func(t *testing.T) {
		_, err := s.repo.FindByID(s.ctx, uuid.New())
		require.Error(t, err)
		require.True(t, domain.IsNotFound(err))
	})
}

func (s *OrderTestSuite) TestOrder_Update() {
	storeID := uuid.NewString()
	order := s.newTestOrder(storeID, "NO-UPD")
	require.NoError(s.T(), s.repo.Create(s.ctx, order))

	s.T().Run("正常更新", func(t *testing.T) {
		newBusinessDate := "2025-12-23"
		newProducts := json.RawMessage(`[{"order_item_id":"1","product_name":"可乐","qty":1}]`)
		newAmount := json.RawMessage(`{"amount_due":100}`)

		upd := &domain.Order{
			ID:            order.ID,
			BusinessDate:  newBusinessDate,
			OrderStatus:   "PLACED",
			PaymentStatus: "PAID",
			Products:      newProducts,
			Amount:        newAmount,
		}

		err := s.repo.Update(s.ctx, upd)
		require.NoError(t, err)

		dbOrder := s.client.Order.GetX(s.ctx, order.ID)
		require.Equal(t, newBusinessDate, dbOrder.BusinessDate)
		require.Equal(t, entorder.OrderStatus("PLACED"), dbOrder.OrderStatus)
		require.Equal(t, entorder.PaymentStatus("PAID"), dbOrder.PaymentStatus)
		require.Equal(t, []byte(newProducts), []byte(dbOrder.Products))
		require.Equal(t, []byte(newAmount), []byte(dbOrder.Amount))
	})

	s.T().Run("更新不存在的ID", func(t *testing.T) {
		err := s.repo.Update(s.ctx, &domain.Order{ID: uuid.New(), BusinessDate: "2025-12-24"})
		require.Error(t, err)
		require.True(t, domain.IsNotFound(err))
	})

	s.T().Run("非法枚举值返回 ParamsError", func(t *testing.T) {
		upd := &domain.Order{
			ID:          order.ID,
			OrderStatus: "OPEN",
		}
		err := s.repo.Update(s.ctx, upd)
		require.Error(t, err)
		require.True(t, domain.IsParamsError(err))
	})
}

func (s *OrderTestSuite) TestOrder_Delete() {
	storeID := uuid.NewString()
	order := s.newTestOrder(storeID, "NO-DEL")
	require.NoError(s.T(), s.repo.Create(s.ctx, order))

	s.T().Run("正常软删", func(t *testing.T) {
		err := s.repo.Delete(s.ctx, order.ID)
		require.NoError(t, err)

		// 默认查询应被软删拦截为 not found
		_, err = s.client.Order.Get(s.ctx, order.ID)
		require.Error(t, err)
		require.True(t, ent.IsNotFound(err))

		// SkipSoftDelete 可以查到并验证 deleted_at 已设置
		ctx := schematype.SkipSoftDelete(s.ctx)
		dbOrder := s.client.Order.GetX(ctx, order.ID)
		require.Greater(t, dbOrder.DeletedAt, int64(0))
	})

	s.T().Run("删除不存在的ID", func(t *testing.T) {
		err := s.repo.Delete(s.ctx, uuid.New())
		require.Error(t, err)
		require.True(t, domain.IsNotFound(err))
	})
}

func (s *OrderTestSuite) TestOrder_List() {
	storeID := uuid.NewString()
	merchantID := uuid.NewString()

	// 手工创建三条，保证 created_at 有序，便于验证排序
	base := time.Date(2025, 12, 22, 10, 0, 0, 0, time.UTC)

	// 直接用 ent 创建以控制 created_at，但仍满足 repo.List 的过滤字段
	o1 := s.client.Order.Create().
		SetID(uuid.New()).
		SetMerchantID(merchantID).
		SetStoreID(storeID).
		SetBusinessDate("2025-12-22").
		SetOrderNo("NO-L1").
		SetDiningMode(entorder.DiningMode("DINE_IN")).
		SetCreatedAt(base.Add(1 * time.Second)).
		SaveX(s.ctx)
	o2 := s.client.Order.Create().
		SetID(uuid.New()).
		SetMerchantID(merchantID).
		SetStoreID(storeID).
		SetBusinessDate("2025-12-22").
		SetOrderNo("NO-L2").
		SetDiningMode(entorder.DiningMode("DINE_IN")).
		SetPaymentStatus(entorder.PaymentStatus("PAID")).
		SetCreatedAt(base.Add(2 * time.Second)).
		SaveX(s.ctx)
	o3 := s.client.Order.Create().
		SetID(uuid.New()).
		SetMerchantID(merchantID).
		SetStoreID(storeID).
		SetBusinessDate("2025-12-22").
		SetOrderNo("NO-L3").
		SetDiningMode(entorder.DiningMode("DINE_IN")).
		SetCreatedAt(base.Add(3 * time.Second)).
		SaveX(s.ctx)

	s.T().Run("分页 + created_at 倒序", func(t *testing.T) {
		items, total, err := s.repo.List(s.ctx, domain.OrderListParams{
			MerchantID: merchantID,
			StoreID:    storeID,
			Page:       1,
			Size:       2,
		})
		require.NoError(t, err)
		require.Equal(t, 3, total)
		require.Len(t, items, 2)
		require.Equal(t, o3.ID, items[0].ID)
		require.Equal(t, o2.ID, items[1].ID)

		items2, total2, err := s.repo.List(s.ctx, domain.OrderListParams{
			MerchantID: merchantID,
			StoreID:    storeID,
			Page:       2,
			Size:       2,
		})
		require.NoError(t, err)
		require.Equal(t, 3, total2)
		require.Len(t, items2, 1)
		require.Equal(t, o1.ID, items2[0].ID)
	})

	s.T().Run("按 order_no 过滤", func(t *testing.T) {
		items, total, err := s.repo.List(s.ctx, domain.OrderListParams{
			MerchantID: merchantID,
			StoreID:    storeID,
			OrderNo:    "NO-L2",
			Page:       1,
			Size:       10,
		})
		require.NoError(t, err)
		require.Equal(t, 1, total)
		require.Len(t, items, 1)
		require.Equal(t, o2.ID, items[0].ID)
	})

	s.T().Run("按 payment_status 过滤", func(t *testing.T) {
		items, total, err := s.repo.List(s.ctx, domain.OrderListParams{
			MerchantID:    merchantID,
			StoreID:       storeID,
			PaymentStatus: "PAID",
			Page:          1,
			Size:          10,
		})
		require.NoError(t, err)
		require.Equal(t, 1, total)
		require.Len(t, items, 1)
		require.Equal(t, o2.ID, items[0].ID)
	})

	s.T().Run("非法 order_status 过滤返回 ParamsError", func(t *testing.T) {
		_, _, err := s.repo.List(s.ctx, domain.OrderListParams{
			MerchantID:  merchantID,
			StoreID:     storeID,
			OrderStatus: "OPEN",
			Page:        1,
			Size:        10,
		})
		require.Error(t, err)
		require.True(t, domain.IsParamsError(err))
	})

	s.T().Run("非法 payment_status 过滤返回 ParamsError", func(t *testing.T) {
		_, _, err := s.repo.List(s.ctx, domain.OrderListParams{
			MerchantID:    merchantID,
			StoreID:       storeID,
			PaymentStatus: "PAYED",
			Page:          1,
			Size:          10,
		})
		require.Error(t, err)
		require.True(t, domain.IsParamsError(err))
	})

	s.T().Run("非法 order_type 过滤返回 ParamsError", func(t *testing.T) {
		_, _, err := s.repo.List(s.ctx, domain.OrderListParams{
			MerchantID: merchantID,
			StoreID:    storeID,
			OrderType:  "SALE_ORDER",
			Page:       1,
			Size:       10,
		})
		require.Error(t, err)
		require.True(t, domain.IsParamsError(err))
	})

	s.T().Run("软删记录不应出现在列表", func(t *testing.T) {
		require.NoError(t, s.repo.Delete(s.ctx, o2.ID))

		items, total, err := s.repo.List(s.ctx, domain.OrderListParams{
			MerchantID: merchantID,
			StoreID:    storeID,
			Page:       1,
			Size:       10,
		})
		require.NoError(t, err)
		require.Equal(t, 2, total)
		require.Len(t, items, 2)
		// created_at 倒序：o3 在前，o1 在后（o2 已软删）
		require.Equal(t, o3.ID, items[0].ID)
		require.Equal(t, o1.ID, items[1].ID)
	})
}
