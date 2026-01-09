package repository

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/ent/schema/schematype"
)

type RefundOrderTestSuite struct {
	RepositoryTestSuite
	repo      *RefundOrderRepository
	orderRepo *OrderRepository
	ctx       context.Context
}

func TestRefundOrderTestSuite(t *testing.T) {
	suite.Run(t, new(RefundOrderTestSuite))
}

func (s *RefundOrderTestSuite) SetupTest() {
	s.RepositoryTestSuite.SetupTest()
	s.repo = &RefundOrderRepository{Client: s.client}
	s.orderRepo = &OrderRepository{Client: s.client}
	s.ctx = context.Background()
}

func (s *RefundOrderTestSuite) newTestOrder(merchantID, storeID uuid.UUID, orderNo string) *domain.Order {
	return &domain.Order{
		ID:            uuid.New(),
		MerchantID:    merchantID,
		StoreID:       storeID,
		BusinessDate:  "2026-01-08",
		OrderNo:       orderNo,
		DiningMode:    domain.DiningModeDineIn,
		OrderStatus:   domain.OrderStatusCompleted,
		PaymentStatus: domain.PaymentStatusPaid,
		Channel:       domain.ChannelPOS,
		Store:         domain.OrderStore{ID: storeID, MerchantID: merchantID, StoreName: "测试门店"},
		Pos:           domain.OrderPOS{ID: uuid.New(), Name: "POS001"},
		Cashier:       domain.OrderCashier{CashierID: uuid.New(), CashierName: "收银员"},
		Amount:        domain.OrderAmount{AmountPaid: decimal.NewFromInt(100)},
		PaidAt:        time.Now(),
	}
}

func (s *RefundOrderTestSuite) newTestRefundOrder(merchantID, storeID, originOrderID uuid.UUID, refundNo string) *domain.RefundOrder {
	return &domain.RefundOrder{
		ID:               uuid.New(),
		MerchantID:       merchantID,
		StoreID:          storeID,
		BusinessDate:     "2026-01-08",
		RefundNo:         refundNo,
		OriginOrderID:    originOrderID,
		OriginOrderNo:    "ORD-001",
		OriginPaidAt:     time.Now(),
		OriginAmountPaid: decimal.NewFromInt(100),
		RefundType:       domain.RefundTypeFull,
		RefundStatus:     domain.RefundStatusPending,
		RefundReasonCode: domain.RefundReasonCustomerRequest,
		RefundReason:     "顾客要求退款",
		Store:            domain.OrderStore{ID: storeID, MerchantID: merchantID, StoreName: "测试门店"},
		Channel:          domain.ChannelPOS,
		Pos:              domain.OrderPOS{ID: uuid.New(), Name: "POS001"},
		Cashier:          domain.OrderCashier{CashierID: uuid.New(), CashierName: "收银员"},
		RefundAmount: domain.RefundAmount{
			ItemsSubtotal: decimal.NewFromInt(100),
			RefundTotal:   decimal.NewFromInt(100),
		},
		Remark: "测试退款",
	}
}

func (s *RefundOrderTestSuite) TestRefundOrder_Create() {
	merchantID := uuid.New()
	storeID := uuid.New()
	originOrderID := uuid.New()

	s.T().Run("创建成功", func(t *testing.T) {
		ro := s.newTestRefundOrder(merchantID, storeID, originOrderID, "RF-001")

		err := s.repo.Create(s.ctx, ro)
		require.NoError(t, err)
		require.NotEqual(t, uuid.Nil, ro.ID)

		dbRO := s.client.RefundOrder.GetX(s.ctx, ro.ID)
		require.Equal(t, ro.MerchantID, dbRO.MerchantID)
		require.Equal(t, ro.StoreID, dbRO.StoreID)
		require.Equal(t, ro.BusinessDate, dbRO.BusinessDate)
		require.Equal(t, ro.RefundNo, dbRO.RefundNo)
		require.Equal(t, ro.OriginOrderID, dbRO.OriginOrderID)
		require.Equal(t, domain.RefundTypeFull, dbRO.RefundType)
		require.Equal(t, domain.RefundStatusPending, dbRO.RefundStatus)
		require.Equal(t, "测试退款", dbRO.Remark)
	})

	s.T().Run("唯一键冲突返回 Conflict", func(t *testing.T) {
		ro1 := s.newTestRefundOrder(merchantID, storeID, originOrderID, "RF-DUP")
		ro2 := s.newTestRefundOrder(merchantID, storeID, originOrderID, "RF-DUP")

		require.NoError(t, s.repo.Create(s.ctx, ro1))
		err := s.repo.Create(s.ctx, ro2)
		require.Error(t, err)
		require.True(t, domain.IsConflict(err))
	})
}

func (s *RefundOrderTestSuite) TestRefundOrder_CreateWithProducts() {
	merchantID := uuid.New()
	storeID := uuid.New()
	originOrderID := uuid.New()

	s.T().Run("创建退款单含商品明细", func(t *testing.T) {
		ro := s.newTestRefundOrder(merchantID, storeID, originOrderID, "RF-PROD-001")
		ro.RefundProducts = []domain.RefundOrderProduct{
			{
				OriginOrderProductID: uuid.New(),
				OriginOrderItemID:    "ITEM-001",
				ProductID:            uuid.New(),
				ProductName:          "宫保鸡丁",
				ProductType:          domain.ProductTypeNormal,
				OriginQty:            2,
				OriginPrice:          decimal.NewFromInt(38),
				OriginSubtotal:       decimal.NewFromInt(76),
				OriginTax:            decimal.NewFromFloat(4.56),
				OriginTotal:          decimal.NewFromFloat(80.56),
				RefundQty:            1,
				RefundSubtotal:       decimal.NewFromInt(38),
				RefundTax:            decimal.NewFromFloat(2.28),
				RefundTotal:          decimal.NewFromFloat(40.28),
				RefundReason:         "顾客不想要了",
			},
			{
				OriginOrderProductID: uuid.New(),
				OriginOrderItemID:    "ITEM-002",
				ProductID:            uuid.New(),
				ProductName:          "雪碧",
				ProductType:          domain.ProductTypeNormal,
				OriginQty:            1,
				OriginPrice:          decimal.NewFromInt(8),
				OriginSubtotal:       decimal.NewFromInt(8),
				OriginTax:            decimal.NewFromFloat(0.48),
				OriginTotal:          decimal.NewFromFloat(8.48),
				RefundQty:            1,
				RefundSubtotal:       decimal.NewFromInt(8),
				RefundTax:            decimal.NewFromFloat(0.48),
				RefundTotal:          decimal.NewFromFloat(8.48),
				RefundReason:         "饮料有问题",
			},
		}

		err := s.repo.Create(s.ctx, ro)
		require.NoError(t, err)
		require.NotEqual(t, uuid.Nil, ro.ID)

		// 验证商品明细已创建
		found, err := s.repo.FindByID(s.ctx, ro.ID)
		require.NoError(t, err)
		require.Len(t, found.RefundProducts, 2)
		require.Equal(t, "宫保鸡丁", found.RefundProducts[0].ProductName)
		require.Equal(t, "雪碧", found.RefundProducts[1].ProductName)
		require.Equal(t, 1, found.RefundProducts[0].RefundQty)
	})
}

func (s *RefundOrderTestSuite) TestRefundOrder_FindByID() {
	merchantID := uuid.New()
	storeID := uuid.New()
	originOrderID := uuid.New()

	s.T().Run("正常查询", func(t *testing.T) {
		ro := s.newTestRefundOrder(merchantID, storeID, originOrderID, "RF-FIND-001")
		require.NoError(t, s.repo.Create(s.ctx, ro))

		found, err := s.repo.FindByID(s.ctx, ro.ID)
		require.NoError(t, err)
		require.Equal(t, ro.ID, found.ID)
		require.Equal(t, ro.RefundNo, found.RefundNo)
		require.Equal(t, ro.MerchantID, found.MerchantID)
		require.Equal(t, ro.StoreID, found.StoreID)
		require.Equal(t, domain.RefundTypeFull, found.RefundType)
	})

	s.T().Run("不存在的ID返回NotFound", func(t *testing.T) {
		_, err := s.repo.FindByID(s.ctx, uuid.New())
		require.Error(t, err)
		require.True(t, domain.IsNotFound(err))
	})
}

func (s *RefundOrderTestSuite) TestRefundOrder_Update() {
	merchantID := uuid.New()
	storeID := uuid.New()
	originOrderID := uuid.New()

	s.T().Run("正常更新", func(t *testing.T) {
		ro := s.newTestRefundOrder(merchantID, storeID, originOrderID, "RF-UPD-001")
		require.NoError(t, s.repo.Create(s.ctx, ro))

		// 更新状态和备注
		ro.RefundStatus = domain.RefundStatusCompleted
		ro.Remark = "已完成退款"
		ro.RefundedAt = time.Now()
		ro.RefundedBy = uuid.New()
		ro.RefundedByName = "操作员A"

		err := s.repo.Update(s.ctx, ro)
		require.NoError(t, err)

		// 验证更新结果
		found, err := s.repo.FindByID(s.ctx, ro.ID)
		require.NoError(t, err)
		require.Equal(t, domain.RefundStatusCompleted, found.RefundStatus)
		require.Equal(t, "已完成退款", found.Remark)
		require.Equal(t, "操作员A", found.RefundedByName)
	})

	s.T().Run("更新不存在的ID返回NotFound", func(t *testing.T) {
		ro := &domain.RefundOrder{
			ID:           uuid.New(),
			RefundStatus: domain.RefundStatusCompleted,
		}
		err := s.repo.Update(s.ctx, ro)
		require.Error(t, err)
		require.True(t, domain.IsNotFound(err))
	})
}

func (s *RefundOrderTestSuite) TestRefundOrder_Delete() {
	merchantID := uuid.New()
	storeID := uuid.New()
	originOrderID := uuid.New()

	s.T().Run("正常软删除", func(t *testing.T) {
		ro := s.newTestRefundOrder(merchantID, storeID, originOrderID, "RF-DEL-001")
		require.NoError(t, s.repo.Create(s.ctx, ro))

		err := s.repo.Delete(s.ctx, ro.ID)
		require.NoError(t, err)

		// 正常查询应该找不到（软删除）
		_, err = s.repo.FindByID(s.ctx, ro.ID)
		require.Error(t, err)
		require.True(t, domain.IsNotFound(err))

		// 跳过软删除可以查到
		ctxSkip := schematype.SkipSoftDelete(s.ctx)
		dbRO := s.client.RefundOrder.GetX(ctxSkip, ro.ID)
		require.NotNil(t, dbRO.DeletedAt)
	})

	s.T().Run("删除不存在的ID返回NotFound", func(t *testing.T) {
		err := s.repo.Delete(s.ctx, uuid.New())
		require.Error(t, err)
		require.True(t, domain.IsNotFound(err))
	})
}

func (s *RefundOrderTestSuite) TestRefundOrder_List() {
	merchantID := uuid.New()
	storeID := uuid.New()
	originOrderID := uuid.New()

	// 创建测试数据
	for i := 1; i <= 5; i++ {
		ro := s.newTestRefundOrder(merchantID, storeID, originOrderID, fmt.Sprintf("RF-LIST-%03d", i))
		if i%2 == 0 {
			ro.RefundType = domain.RefundTypePartial
			ro.RefundStatus = domain.RefundStatusCompleted
		}
		require.NoError(s.T(), s.repo.Create(s.ctx, ro))
	}

	s.T().Run("分页查询", func(t *testing.T) {
		list, total, err := s.repo.List(s.ctx, domain.RefundOrderListParams{
			MerchantID: merchantID,
			Page:       1,
			Size:       2,
		})
		require.NoError(t, err)
		require.Equal(t, 5, total)
		require.Len(t, list, 2)
	})

	s.T().Run("按退款类型过滤", func(t *testing.T) {
		list, total, err := s.repo.List(s.ctx, domain.RefundOrderListParams{
			MerchantID: merchantID,
			RefundType: domain.RefundTypePartial,
		})
		require.NoError(t, err)
		require.Equal(t, 2, total)
		require.Len(t, list, 2)
		for _, ro := range list {
			require.Equal(t, domain.RefundTypePartial, ro.RefundType)
		}
	})

	s.T().Run("按退款状态过滤", func(t *testing.T) {
		list, total, err := s.repo.List(s.ctx, domain.RefundOrderListParams{
			MerchantID:   merchantID,
			RefundStatus: domain.RefundStatusCompleted,
		})
		require.NoError(t, err)
		require.Equal(t, 2, total)
		require.Len(t, list, 2)
		for _, ro := range list {
			require.Equal(t, domain.RefundStatusCompleted, ro.RefundStatus)
		}
	})

	s.T().Run("软删记录不出现", func(t *testing.T) {
		// 创建并删除一条记录
		ro := s.newTestRefundOrder(merchantID, storeID, originOrderID, "RF-LIST-DEL")
		require.NoError(t, s.repo.Create(s.ctx, ro))
		require.NoError(t, s.repo.Delete(s.ctx, ro.ID))

		// 列表查询不应包含已删除记录
		list, total, err := s.repo.List(s.ctx, domain.RefundOrderListParams{
			MerchantID: merchantID,
			RefundNo:   "RF-LIST-DEL",
		})
		require.NoError(t, err)
		require.Equal(t, 0, total)
		require.Len(t, list, 0)
	})
}

func (s *RefundOrderTestSuite) TestRefundOrder_FindByOriginOrderID() {
	merchantID := uuid.New()
	storeID := uuid.New()
	originOrderID := uuid.New()

	s.T().Run("按原订单ID查询", func(t *testing.T) {
		// 创建多条关联同一原订单的退款单
		for i := 1; i <= 3; i++ {
			ro := s.newTestRefundOrder(merchantID, storeID, originOrderID, fmt.Sprintf("RF-ORIGIN-%03d", i))
			require.NoError(t, s.repo.Create(s.ctx, ro))
		}

		list, err := s.repo.FindByOriginOrderID(s.ctx, originOrderID)
		require.NoError(t, err)
		require.Len(t, list, 3)
		for _, ro := range list {
			require.Equal(t, originOrderID, ro.OriginOrderID)
		}
	})

	s.T().Run("不存在的原订单ID返回空列表", func(t *testing.T) {
		list, err := s.repo.FindByOriginOrderID(s.ctx, uuid.New())
		require.NoError(t, err)
		require.Len(t, list, 0)
	})
}
