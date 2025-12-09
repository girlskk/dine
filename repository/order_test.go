package repository

import (
	"context"

	"github.com/shopspring/decimal"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/ent"
	"gitlab.jiguang.dev/pos-dine/dine/ent/order"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/upagination"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/util"
)

func (suite *RepositoryTestSuite) TestOrder_Create() {
	ctx := context.Background()
	store := suite.createTestStore(ctx)
	area := suite.createTestTableArea(ctx, store)
	table := suite.createTestTable(ctx, area)

	suite.Run("创建基本订单", func() {
		dorder := &domain.Order{
			No:              util.RandomString(10),
			Type:            domain.OrderTypeDineIn,
			Source:          domain.OrderSourceOffline,
			Status:          domain.OrderStatusUnpaid,
			TotalPrice:      decimal.NewFromFloat(100),
			Discount:        decimal.Zero,
			RealPrice:       decimal.NewFromFloat(100),
			PointsAvailable: decimal.Zero,
			StoreID:         store.ID,
			StoreName:       store.Name,
			TableID:         table.ID,
			TableName:       table.Name,
			PeopleNumber:    4,
			CreatorID:       1,
			CreatorName:     "测试用户",
			Items: []*domain.OrderItem{
				{
					ProductID:     1,
					Name:          "测试商品",
					Type:          domain.ProductTypeSingle,
					AllowPointPay: true,
					Quantity:      decimal.NewFromFloat(2),
					Price:         decimal.NewFromFloat(50),
					Amount:        decimal.NewFromFloat(100),
					ProductSnapshot: domain.OrderProductInfoSnapshot{
						Price:    decimal.NewFromFloat(50),
						SpecName: "默认规格",
						UnitName: "份",
					},
				},
			},
		}

		var newOrder *domain.Order
		var err error
		err = withTx(ctx, suite.client, func(tx *ent.Tx) error {
			repo := NewOrderRepository(tx.Client())
			newOrder, err = repo.Create(ctx, dorder)
			return err
		})
		suite.Require().NoError(err)

		suite.Positive(newOrder.ID)
		suite.Equal(dorder.No, newOrder.No)
		suite.NotZero(newOrder.CreatedAt)
		suite.NotZero(newOrder.UpdatedAt)

		// 验证订单是否创建成功
		savedOrder, err := suite.client.Order.Get(ctx, newOrder.ID)
		suite.NoError(err)
		suite.Equal(dorder.No, savedOrder.No)
		suite.Equal(dorder.TotalPrice.String(), savedOrder.TotalPrice.String())
		suite.Equal(dorder.PeopleNumber, savedOrder.PeopleNumber)

		// 验证订单商品是否创建成功
		items, err := savedOrder.QueryItems().All(ctx)
		suite.NoError(err)
		suite.Len(items, 1)
		suite.Equal(dorder.Items[0].ProductID, items[0].ProductID)
		suite.Equal(dorder.Items[0].Amount.String(), items[0].Amount.String())
	})

	suite.Run("创建带套餐的订单", func() {
		dorder := &domain.Order{
			No:              util.RandomString(10),
			Type:            domain.OrderTypeDineIn,
			Source:          domain.OrderSourceOffline,
			Status:          domain.OrderStatusUnpaid,
			TotalPrice:      decimal.NewFromFloat(200),
			Discount:        decimal.Zero,
			RealPrice:       decimal.NewFromFloat(200),
			PointsAvailable: decimal.Zero,
			StoreID:         store.ID,
			StoreName:       store.Name,
			TableID:         table.ID,
			TableName:       table.Name,
			PeopleNumber:    2,
			CreatorID:       1,
			CreatorName:     "测试用户",
			Items: []*domain.OrderItem{
				{
					ProductID:     2,
					Name:          "测试套餐",
					Type:          domain.ProductTypeSetMeal,
					AllowPointPay: true,
					Quantity:      decimal.NewFromFloat(1),
					Price:         decimal.NewFromFloat(200),
					Amount:        decimal.NewFromFloat(200),
					ProductSnapshot: domain.OrderProductInfoSnapshot{
						Price:    decimal.NewFromFloat(200),
						UnitName: "份",
					},
					SetMealDetails: []*domain.OrderItemSetMealDetail{
						{
							Name:         "套餐内商品1",
							Type:         domain.ProductTypeSingle,
							SetMealPrice: decimal.NewFromFloat(100),
							SetMealID:    2,
							ProductID:    3,
							Quantity:     decimal.NewFromFloat(1),
							ProductSnapshot: domain.OrderProductInfoSnapshot{
								Price:    decimal.NewFromFloat(100),
								UnitName: "份",
							},
						},
					},
				},
			},
		}

		var newOrder *domain.Order
		var err error
		err = withTx(ctx, suite.client, func(tx *ent.Tx) error {
			repo := NewOrderRepository(tx.Client())
			newOrder, err = repo.Create(ctx, dorder)
			return err
		})
		suite.Require().NoError(err)
		suite.Positive(newOrder.ID)
		suite.Equal(dorder.No, newOrder.No)
		suite.NotZero(newOrder.CreatedAt)
		suite.NotZero(newOrder.UpdatedAt)

		// 验证订单是否创建成功
		savedOrder, err := suite.client.Order.Get(ctx, newOrder.ID)
		suite.NoError(err)
		suite.Equal(dorder.No, savedOrder.No)
		suite.Equal(dorder.TotalPrice.String(), savedOrder.TotalPrice.String())

		// 验证订单商品是否创建成功
		items, err := savedOrder.QueryItems().All(ctx)
		suite.NoError(err)
		suite.Len(items, 1)
		suite.Equal(dorder.Items[0].ProductID, items[0].ProductID)

		// 验证套餐详情是否创建成功
		details, err := items[0].QuerySetMealDetails().All(ctx)
		suite.NoError(err)
		suite.Len(details, 1)
		suite.Equal(dorder.Items[0].SetMealDetails[0].ProductID, details[0].ProductID)
		suite.Equal(dorder.Items[0].SetMealDetails[0].SetMealPrice.String(), details[0].SetMealPrice.String())
	})

	suite.Run("事务回滚测试", func() {
		// 创建一个无效的订单（商品缺少必填字段）
		dorder := &domain.Order{
			No:              util.RandomString(10),
			Type:            domain.OrderTypeDineIn,
			Source:          domain.OrderSourceOffline,
			Status:          domain.OrderStatusUnpaid,
			TotalPrice:      decimal.NewFromFloat(200),
			Discount:        decimal.Zero,
			RealPrice:       decimal.NewFromFloat(200),
			PointsAvailable: decimal.Zero,
			StoreID:         store.ID,
			StoreName:       store.Name,
			TableID:         table.ID,
			TableName:       table.Name,
			PeopleNumber:    2,
			CreatorID:       1,
			CreatorName:     "测试用户",
			Items: []*domain.OrderItem{
				{
					ProductID:     2,
					Type:          domain.ProductTypeSetMeal,
					AllowPointPay: true,
					Quantity:      decimal.NewFromFloat(1),
					Price:         decimal.NewFromFloat(200),
					Amount:        decimal.NewFromFloat(200),
					ProductSnapshot: domain.OrderProductInfoSnapshot{
						Price:    decimal.NewFromFloat(200),
						UnitName: "份",
					},
					SetMealDetails: []*domain.OrderItemSetMealDetail{
						{
							ProductID:    3,
							SetMealPrice: decimal.NewFromFloat(100),
							Quantity:     decimal.NewFromFloat(1),
							ProductSnapshot: domain.OrderProductInfoSnapshot{
								Price:    decimal.NewFromFloat(100),
								UnitName: "份",
							},
						},
					},
				},
			},
		}

		err := withTx(ctx, suite.client, func(tx *ent.Tx) error {
			repo := NewOrderRepository(tx.Client())
			_, err := repo.Create(ctx, dorder)
			return err
		})
		suite.Require().Error(err)

		// 验证订单未被创建
		count, err := suite.client.Order.Query().Where(order.NoEQ(dorder.No)).Count(ctx)
		suite.NoError(err)
		suite.Equal(0, count)
	})
}

func (suite *RepositoryTestSuite) TestOrder_GetOrders() {
	ctx := context.Background()
	store := suite.createTestStore(ctx)
	area := suite.createTestTableArea(ctx, store)
	table := suite.createTestTable(ctx, area)

	dorder1 := &domain.Order{
		No:              util.RandomString(10),
		Type:            domain.OrderTypeDineIn,
		Source:          domain.OrderSourceOffline,
		Status:          domain.OrderStatusUnpaid,
		TotalPrice:      decimal.NewFromFloat(100),
		Discount:        decimal.Zero,
		RealPrice:       decimal.NewFromFloat(100),
		PointsAvailable: decimal.Zero,
		StoreID:         store.ID,
		StoreName:       store.Name,
		TableID:         table.ID,
		TableName:       table.Name,
		PeopleNumber:    4,
		CreatorID:       1,
		CreatorName:     "测试用户",
		Items: []*domain.OrderItem{
			{
				ProductID:     1,
				Name:          "测试商品[foo]",
				Type:          domain.ProductTypeSingle,
				AllowPointPay: true,
				Quantity:      decimal.NewFromFloat(2),
				Price:         decimal.NewFromFloat(50),
				Amount:        decimal.NewFromFloat(100),
				ProductSnapshot: domain.OrderProductInfoSnapshot{
					Price:    decimal.NewFromFloat(50),
					SpecName: "默认规格",
					UnitName: "份",
				},
			},
		},
	}

	err := withTx(ctx, suite.client, func(tx *ent.Tx) error {
		repo := NewOrderRepository(tx.Client())
		_, err := repo.Create(ctx, dorder1)
		return err
	})
	suite.Require().NoError(err)

	dorder2 := &domain.Order{
		No:              util.RandomString(10),
		Type:            domain.OrderTypeDineIn,
		Source:          domain.OrderSourceOffline,
		Status:          domain.OrderStatusUnpaid,
		TotalPrice:      decimal.NewFromFloat(200),
		Discount:        decimal.Zero,
		RealPrice:       decimal.NewFromFloat(200),
		PointsAvailable: decimal.Zero,
		StoreID:         store.ID,
		StoreName:       store.Name,
		TableID:         table.ID,
		TableName:       table.Name,
		PeopleNumber:    2,
		CreatorID:       1,
		CreatorName:     "测试用户",
		MemberName:      "GopherMember",
		MemberPhone:     "13800138000",
		Items: []*domain.OrderItem{
			{
				ProductID:     2,
				Name:          "测试套餐[bar]",
				Type:          domain.ProductTypeSetMeal,
				AllowPointPay: true,
				Quantity:      decimal.NewFromFloat(1),
				Price:         decimal.NewFromFloat(200),
				Amount:        decimal.NewFromFloat(200),
				ProductSnapshot: domain.OrderProductInfoSnapshot{
					Price:    decimal.NewFromFloat(200),
					UnitName: "份",
				},
				SetMealDetails: []*domain.OrderItemSetMealDetail{
					{
						Name:         "套餐内商品1",
						Type:         domain.ProductTypeSingle,
						SetMealPrice: decimal.NewFromFloat(100),
						SetMealID:    2,
						ProductID:    3,
						Quantity:     decimal.NewFromFloat(1),
						ProductSnapshot: domain.OrderProductInfoSnapshot{
							Price:    decimal.NewFromFloat(100),
							UnitName: "份",
						},
					},
				},
			},
		},
	}

	err = withTx(ctx, suite.client, func(tx *ent.Tx) error {
		repo := NewOrderRepository(tx.Client())
		_, err := repo.Create(ctx, dorder2)
		return err
	})
	suite.Require().NoError(err)

	suite.Run("无筛选", func() {
		repo := NewOrderRepository(suite.client)
		dorders, total, err := repo.GetOrders(ctx, upagination.New(1, 10), &domain.OrderListFilter{})
		suite.Require().NoError(err)
		suite.Require().Equal(2, total)
		suite.Require().Equal(2, len(dorders))
	})

	suite.Run("商品名称筛选", func() {
		repo := NewOrderRepository(suite.client)
		dorders, total, err := repo.GetOrders(ctx, upagination.New(1, 10), &domain.OrderListFilter{
			HasItemName: "foo",
		})
		suite.Require().NoError(err)
		suite.Require().Equal(1, total)
		suite.Require().Equal(1, len(dorders))
		suite.Equal(dorder1.No, dorders[0].No)
	})

	suite.Run("会员名称或者手机号筛选", func() {
		repo := NewOrderRepository(suite.client)
		dorders, total, err := repo.GetOrders(ctx, upagination.New(1, 10), &domain.OrderListFilter{
			MemberNameOrPhone: "pher",
		})
		suite.Require().NoError(err)
		suite.Require().Equal(1, total)
		suite.Require().Equal(1, len(dorders))
		suite.Equal(dorder2.No, dorders[0].No)

		dorders, total, err = repo.GetOrders(ctx, upagination.New(1, 10), &domain.OrderListFilter{
			MemberNameOrPhone: "138",
		})
		suite.Require().NoError(err)
		suite.Require().Equal(1, total)
		suite.Require().Equal(1, len(dorders))
		suite.Equal(dorder2.No, dorders[0].No)
	})
}

func (suite *RepositoryTestSuite) TestOrder_GetOrderRange() {
	ctx := context.Background()
	store := suite.createTestStore(ctx)
	area := suite.createTestTableArea(ctx, store)
	table := suite.createTestTable(ctx, area)

	dorder1 := &domain.Order{
		No:              util.RandomString(10),
		Type:            domain.OrderTypeDineIn,
		Source:          domain.OrderSourceOffline,
		Status:          domain.OrderStatusUnpaid,
		TotalPrice:      decimal.NewFromFloat(100),
		Discount:        decimal.Zero,
		RealPrice:       decimal.NewFromFloat(100),
		PointsAvailable: decimal.Zero,
		StoreID:         store.ID,
		StoreName:       store.Name,
		TableID:         table.ID,
		TableName:       table.Name,
		PeopleNumber:    4,
		CreatorID:       1,
		CreatorName:     "测试用户",
		Items: []*domain.OrderItem{
			{
				ProductID:     1,
				Name:          "测试商品[foo]",
				Type:          domain.ProductTypeSingle,
				AllowPointPay: true,
				Quantity:      decimal.NewFromFloat(2),
				Price:         decimal.NewFromFloat(50),
				Amount:        decimal.NewFromFloat(100),
				ProductSnapshot: domain.OrderProductInfoSnapshot{
					Price:    decimal.NewFromFloat(50),
					SpecName: "默认规格",
					UnitName: "份",
				},
			},
		},
	}

	var err error
	err = withTx(ctx, suite.client, func(tx *ent.Tx) error {
		repo := NewOrderRepository(tx.Client())
		dorder1, err = repo.Create(ctx, dorder1)
		return err
	})
	suite.Require().NoError(err)

	dorder2 := &domain.Order{
		No:              util.RandomString(10),
		Type:            domain.OrderTypeDineIn,
		Source:          domain.OrderSourceOffline,
		Status:          domain.OrderStatusUnpaid,
		TotalPrice:      decimal.NewFromFloat(200),
		Discount:        decimal.Zero,
		RealPrice:       decimal.NewFromFloat(200),
		PointsAvailable: decimal.Zero,
		StoreID:         store.ID,
		StoreName:       store.Name,
		TableID:         table.ID,
		TableName:       table.Name,
		PeopleNumber:    2,
		CreatorID:       1,
		CreatorName:     "测试用户",
		MemberName:      "GopherMember",
		MemberPhone:     "13800138000",
		Items: []*domain.OrderItem{
			{
				ProductID:     2,
				Name:          "测试套餐[bar]",
				Type:          domain.ProductTypeSetMeal,
				AllowPointPay: true,
				Quantity:      decimal.NewFromFloat(1),
				Price:         decimal.NewFromFloat(200),
				Amount:        decimal.NewFromFloat(200),
				ProductSnapshot: domain.OrderProductInfoSnapshot{
					Price:    decimal.NewFromFloat(200),
					UnitName: "份",
				},
				SetMealDetails: []*domain.OrderItemSetMealDetail{
					{
						Name:         "套餐内商品1",
						Type:         domain.ProductTypeSingle,
						SetMealPrice: decimal.NewFromFloat(100),
						SetMealID:    2,
						ProductID:    3,
						Quantity:     decimal.NewFromFloat(1),
						ProductSnapshot: domain.OrderProductInfoSnapshot{
							Price:    decimal.NewFromFloat(100),
							UnitName: "份",
						},
					},
				},
			},
		},
	}

	err = withTx(ctx, suite.client, func(tx *ent.Tx) error {
		repo := NewOrderRepository(tx.Client())
		dorder2, err = repo.Create(ctx, dorder2)
		return err
	})
	suite.Require().NoError(err)

	repo := NewOrderRepository(suite.client)
	rg, err := repo.GetOrderRange(ctx, &domain.OrderListFilter{
		StoreID: store.ID,
	})
	suite.Require().NoError(err)
	suite.Require().Equal(2, rg.Count)
	suite.Require().Equal(dorder1.ID, rg.MinID)
	suite.Require().Equal(dorder2.ID, rg.MaxID)
}
