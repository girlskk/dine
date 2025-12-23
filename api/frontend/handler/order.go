package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gitlab.jiguang.dev/pos-dine/dine/api/frontend/types"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/errorx"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/errorx/errcode"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/logging"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/ugin/response"
)

type OrderHandler struct {
	OrderInteractor domain.OrderInteractor
}

func NewOrderHandler(orderInteractor domain.OrderInteractor) *OrderHandler {
	return &OrderHandler{OrderInteractor: orderInteractor}
}

func (h *OrderHandler) Routes(r gin.IRouter) {
	r = r.Group("/orders")
	r.POST("", h.Create())
	r.GET("/:id", h.Get())
	r.PUT("/:id", h.Update())
	r.DELETE("/:id", h.Delete())
	r.GET("", h.List())
}

func (h *OrderHandler) NoAuths() []string {
	return []string{}
}

// Create
//
//	@Tags		订单
//	@Security	BearerAuth
//	@Summary	创建订单
//	@Accept		json
//	@Produce	json
//	@Param		data	body		types.CreateOrderReq	true	"请求信息"
//	@Success	200		{object}	types.Order				"成功"
//	@Router		/orders [post]
func (h *OrderHandler) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("OrderHandler.Create")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.CreateOrderReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		o := &domain.Order{
			ID:                uuid.New(),
			MerchantID:        req.MerchantID,
			StoreID:           req.StoreID,
			BusinessDate:      req.BusinessDate,
			ShiftNo:           req.ShiftNo,
			OrderNo:           req.OrderNo,
			OrderType:         req.OrderType,
			OriginOrderID:     req.OriginOrderID,
			DiningMode:        req.DiningMode,
			OrderStatus:       req.OrderStatus,
			PaymentStatus:     req.PaymentStatus,
			FulfillmentStatus: req.FulfillmentStatus,
			TableStatus:       req.TableStatus,
			TableID:           req.TableID,
			TableName:         req.TableName,
			TableCapacity:     req.TableCapacity,
			GuestCount:        req.GuestCount,
			OpenedBy:          req.OpenedBy,
			PlacedBy:          req.PlacedBy,
			PaidBy:            req.PaidBy,
		}

		if req.Refund != nil {
			b, err := json.Marshal(req.Refund)
			if err != nil {
				c.Error(fmt.Errorf("failed to marshal refund: %w", err))
				return
			}
			o.Refund = b
		}

		if req.OpenedAt != nil {
			o.OpenedAt = req.OpenedAt
		}
		if req.PlacedAt != nil {
			o.PlacedAt = req.PlacedAt
		}
		if req.PaidAt != nil {
			o.PaidAt = req.PaidAt
		}
		if req.CompletedAt != nil {
			o.CompletedAt = req.CompletedAt
		}

		if req.Store != nil {
			b, err := json.Marshal(req.Store)
			if err != nil {
				c.Error(fmt.Errorf("failed to marshal store: %w", err))
				return
			}
			o.Store = b
		} else {
			b, err := json.Marshal(&types.Store{StoreID: req.StoreID})
			if err != nil {
				c.Error(fmt.Errorf("failed to marshal store: %w", err))
				return
			}
			o.Store = b
		}

		if req.Channel != nil {
			b, err := json.Marshal(req.Channel)
			if err != nil {
				c.Error(fmt.Errorf("failed to marshal channel: %w", err))
				return
			}
			o.Channel = b
		}
		if req.POS != nil {
			b, err := json.Marshal(req.POS)
			if err != nil {
				c.Error(fmt.Errorf("failed to marshal pos: %w", err))
				return
			}
			o.Pos = b
		}
		if req.Cashier != nil {
			b, err := json.Marshal(req.Cashier)
			if err != nil {
				c.Error(fmt.Errorf("failed to marshal cashier: %w", err))
				return
			}
			o.Cashier = b
		}

		if req.Member != nil {
			b, err := json.Marshal(req.Member)
			if err != nil {
				c.Error(fmt.Errorf("failed to marshal member: %w", err))
				return
			}
			o.Member = b
		}
		if req.Takeaway != nil {
			b, err := json.Marshal(req.Takeaway)
			if err != nil {
				c.Error(fmt.Errorf("failed to marshal takeaway: %w", err))
				return
			}
			o.Takeaway = b
		}

		if req.Cart != nil {
			b, err := json.Marshal(req.Cart)
			if err != nil {
				c.Error(fmt.Errorf("failed to marshal cart: %w", err))
				return
			}
			o.Cart = b
		}
		if req.Products != nil {
			b, err := json.Marshal(req.Products)
			if err != nil {
				c.Error(fmt.Errorf("failed to marshal products: %w", err))
				return
			}
			o.Products = b
		}
		if req.Promotions != nil {
			b, err := json.Marshal(req.Promotions)
			if err != nil {
				c.Error(fmt.Errorf("failed to marshal promotions: %w", err))
				return
			}
			o.Promotions = b
		}
		if req.Coupons != nil {
			b, err := json.Marshal(req.Coupons)
			if err != nil {
				c.Error(fmt.Errorf("failed to marshal coupons: %w", err))
				return
			}
			o.Coupons = b
		}
		if req.TaxRates != nil {
			b, err := json.Marshal(req.TaxRates)
			if err != nil {
				c.Error(fmt.Errorf("failed to marshal tax_rates: %w", err))
				return
			}
			o.TaxRates = b
		}
		if req.Fees != nil {
			b, err := json.Marshal(req.Fees)
			if err != nil {
				c.Error(fmt.Errorf("failed to marshal fees: %w", err))
				return
			}
			o.Fees = b
		}
		if req.Payments != nil {
			b, err := json.Marshal(req.Payments)
			if err != nil {
				c.Error(fmt.Errorf("failed to marshal payments: %w", err))
				return
			}
			o.Payments = b
		}
		if req.RefundsProducts != nil {
			b, err := json.Marshal(req.RefundsProducts)
			if err != nil {
				c.Error(fmt.Errorf("failed to marshal refunds_products: %w", err))
				return
			}
			o.RefundsProducts = b
		}
		if req.Amount != nil {
			b, err := json.Marshal(req.Amount)
			if err != nil {
				c.Error(fmt.Errorf("failed to marshal amount: %w", err))
				return
			}
			o.Amount = b
		}

		created, err := h.OrderInteractor.Create(ctx, o)
		if err != nil {
			if domain.IsParamsError(err) {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
				return
			}
			if domain.IsConflict(err) {
				c.Error(errorx.New(http.StatusConflict, errcode.Conflict, err))
				return
			}
			err = fmt.Errorf("failed to create order: %w", err)
			c.Error(err)
			return
		}

		res, err := convertDomainOrderToResp(created)
		if err != nil {
			c.Error(fmt.Errorf("failed to convert order: %w", err))
			return
		}
		response.Ok(c, res)
	}
}

// Get
//
//	@Tags		订单
//	@Security	BearerAuth
//	@Summary	获取订单详情
//	@Accept		json
//	@Produce	json
//	@Param		id	path		string		true	"订单ID"
//	@Success	200	{object}	types.Order	"成功"
//	@Router		/orders/{id} [get]
func (h *OrderHandler) Get() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("OrderHandler.Get")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		o, err := h.OrderInteractor.Get(ctx, id)
		if err != nil {
			if domain.IsNotFound(err) {
				c.Error(errorx.New(http.StatusNotFound, errcode.NotFound, err))
				return
			}
			err = fmt.Errorf("failed to get order: %w", err)
			c.Error(err)
			return
		}

		res, err := convertDomainOrderToResp(o)
		if err != nil {
			c.Error(fmt.Errorf("failed to convert order: %w", err))
			return
		}
		response.Ok(c, res)
	}
}

// Update
//
//	@Tags		订单
//	@Security	BearerAuth
//	@Summary	更新订单
//	@Accept		json
//	@Produce	json
//	@Param		id		path		string					true	"订单ID"
//	@Param		data	body		types.UpdateOrderReq	true	"请求信息"
//	@Success	200		{object}	types.Order				"成功"
//	@Router		/orders/{id} [put]
func (h *OrderHandler) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("OrderHandler.Update")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		var req types.UpdateOrderReq
		if err := c.ShouldBind(&req); err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		o := &domain.Order{ID: id}

		if req.BusinessDate != nil {
			o.BusinessDate = *req.BusinessDate
		}
		if req.ShiftNo != nil {
			o.ShiftNo = *req.ShiftNo
		}
		if req.OrderNo != nil {
			o.OrderNo = *req.OrderNo
		}
		if req.OrderType != nil {
			o.OrderType = *req.OrderType
		}
		if req.OriginOrderID != nil {
			o.OriginOrderID = *req.OriginOrderID
		}
		if req.Refund != nil {
			b, err := json.Marshal(req.Refund)
			if err != nil {
				c.Error(fmt.Errorf("failed to marshal refund: %w", err))
				return
			}
			o.Refund = b
		}

		if req.DiningMode != nil {
			o.DiningMode = *req.DiningMode
		}
		if req.OrderStatus != nil {
			o.OrderStatus = *req.OrderStatus
		}
		if req.PaymentStatus != nil {
			o.PaymentStatus = *req.PaymentStatus
		}
		if req.FulfillmentStatus != nil {
			o.FulfillmentStatus = *req.FulfillmentStatus
		}
		if req.TableStatus != nil {
			o.TableStatus = *req.TableStatus
		}

		if req.TableID != nil {
			o.TableID = *req.TableID
		}
		if req.TableName != nil {
			o.TableName = *req.TableName
		}
		if req.TableCapacity != nil {
			o.TableCapacity = *req.TableCapacity
		}
		if req.GuestCount != nil {
			o.GuestCount = *req.GuestCount
		}

		if req.OpenedAt != nil {
			o.OpenedAt = req.OpenedAt
		}
		if req.PlacedAt != nil {
			o.PlacedAt = req.PlacedAt
		}
		if req.PaidAt != nil {
			o.PaidAt = req.PaidAt
		}
		if req.CompletedAt != nil {
			o.CompletedAt = req.CompletedAt
		}
		if req.OpenedBy != nil {
			o.OpenedBy = *req.OpenedBy
		}
		if req.PlacedBy != nil {
			o.PlacedBy = *req.PlacedBy
		}
		if req.PaidBy != nil {
			o.PaidBy = *req.PaidBy
		}

		if req.Store != nil {
			b, err := json.Marshal(req.Store)
			if err != nil {
				c.Error(fmt.Errorf("failed to marshal store: %w", err))
				return
			}
			o.Store = b
		}
		if req.Channel != nil {
			b, err := json.Marshal(req.Channel)
			if err != nil {
				c.Error(fmt.Errorf("failed to marshal channel: %w", err))
				return
			}
			o.Channel = b
		}
		if req.POS != nil {
			b, err := json.Marshal(req.POS)
			if err != nil {
				c.Error(fmt.Errorf("failed to marshal pos: %w", err))
				return
			}
			o.Pos = b
		}
		if req.Cashier != nil {
			b, err := json.Marshal(req.Cashier)
			if err != nil {
				c.Error(fmt.Errorf("failed to marshal cashier: %w", err))
				return
			}
			o.Cashier = b
		}
		if req.Member != nil {
			b, err := json.Marshal(req.Member)
			if err != nil {
				c.Error(fmt.Errorf("failed to marshal member: %w", err))
				return
			}
			o.Member = b
		}
		if req.Takeaway != nil {
			b, err := json.Marshal(req.Takeaway)
			if err != nil {
				c.Error(fmt.Errorf("failed to marshal takeaway: %w", err))
				return
			}
			o.Takeaway = b
		}

		if req.Cart != nil {
			b, err := json.Marshal(req.Cart)
			if err != nil {
				c.Error(fmt.Errorf("failed to marshal cart: %w", err))
				return
			}
			o.Cart = b
		}
		if req.Products != nil {
			b, err := json.Marshal(req.Products)
			if err != nil {
				c.Error(fmt.Errorf("failed to marshal products: %w", err))
				return
			}
			o.Products = b
		}
		if req.Promotions != nil {
			b, err := json.Marshal(req.Promotions)
			if err != nil {
				c.Error(fmt.Errorf("failed to marshal promotions: %w", err))
				return
			}
			o.Promotions = b
		}
		if req.Coupons != nil {
			b, err := json.Marshal(req.Coupons)
			if err != nil {
				c.Error(fmt.Errorf("failed to marshal coupons: %w", err))
				return
			}
			o.Coupons = b
		}
		if req.TaxRates != nil {
			b, err := json.Marshal(req.TaxRates)
			if err != nil {
				c.Error(fmt.Errorf("failed to marshal tax_rates: %w", err))
				return
			}
			o.TaxRates = b
		}
		if req.Fees != nil {
			b, err := json.Marshal(req.Fees)
			if err != nil {
				c.Error(fmt.Errorf("failed to marshal fees: %w", err))
				return
			}
			o.Fees = b
		}
		if req.Payments != nil {
			b, err := json.Marshal(req.Payments)
			if err != nil {
				c.Error(fmt.Errorf("failed to marshal payments: %w", err))
				return
			}
			o.Payments = b
		}
		if req.RefundsProducts != nil {
			b, err := json.Marshal(req.RefundsProducts)
			if err != nil {
				c.Error(fmt.Errorf("failed to marshal refunds_products: %w", err))
				return
			}
			o.RefundsProducts = b
		}
		if req.Amount != nil {
			b, err := json.Marshal(req.Amount)
			if err != nil {
				c.Error(fmt.Errorf("failed to marshal amount: %w", err))
				return
			}
			o.Amount = b
		}

		updated, err := h.OrderInteractor.Update(ctx, o)
		if err != nil {
			if domain.IsParamsError(err) {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
				return
			}
			if domain.IsNotFound(err) {
				c.Error(errorx.New(http.StatusNotFound, errcode.NotFound, err))
				return
			}
			if domain.IsConflict(err) {
				c.Error(errorx.New(http.StatusConflict, errcode.Conflict, err))
				return
			}
			err = fmt.Errorf("failed to update order: %w", err)
			c.Error(err)
			return
		}

		res, err := convertDomainOrderToResp(updated)
		if err != nil {
			c.Error(fmt.Errorf("failed to convert order: %w", err))
			return
		}
		response.Ok(c, res)
	}
}

// Delete
//
//	@Tags		订单
//	@Security	BearerAuth
//	@Summary	删除订单
//	@Accept		json
//	@Produce	json
//	@Param		id	path	string	true	"订单ID"
//	@Success	200	"No Content"
//	@Router		/orders/{id} [delete]
func (h *OrderHandler) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("OrderHandler.Delete")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		if err := h.OrderInteractor.Delete(ctx, id); err != nil {
			if domain.IsNotFound(err) {
				c.Error(errorx.New(http.StatusNotFound, errcode.NotFound, err))
				return
			}
			c.Error(fmt.Errorf("failed to delete order: %w", err))
			return
		}

		response.Ok(c, nil)
	}
}

// List
//
//	@Tags		订单
//	@Security	BearerAuth
//	@Summary	获取订单列表
//	@Accept		json
//	@Produce	json
//	@Param		merchant_id		query		string				true	"品牌商ID"
//	@Param		store_id		query		string				true	"门店ID"
//	@Param		business_date	query		string				false	"营业日"
//	@Param		order_no		query		string				false	"订单号"
//	@Param		order_type		query		string				false	"订单类型" Enums(SALE,REFUND,PARTIAL_REFUND)
//	@Param		order_status	query		string				false	"订单状态" Enums(DRAFT,PLACED,IN_PROGRESS,READY,COMPLETED,CANCELLED,VOIDED,MERGED)
//	@Param		payment_status	query		string				false	"支付状态" Enums(UNPAID,PAYING,PARTIALLY_PAID,PAID,PARTIALLY_REFUNDED,REFUNDED)
//	@Param		page			query		int					false	"页码"
//	@Param		size			query		int					false	"每页数量"
//	@Success	200				{object}	types.ListOrderResp	"成功"
//	@Router		/orders [get]
func (h *OrderHandler) List() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logging.FromContext(ctx).Named("OrderHandler.List")
		ctx = logging.NewContext(ctx, logger)
		c.Request = c.Request.Clone(ctx)

		var req types.ListOrderReq
		if err := c.ShouldBindQuery(&req); err != nil {
			c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
			return
		}

		params := domain.OrderListParams{
			MerchantID:    req.MerchantID,
			StoreID:       req.StoreID,
			BusinessDate:  req.BusinessDate,
			OrderNo:       req.OrderNo,
			OrderType:     req.OrderType,
			OrderStatus:   req.OrderStatus,
			PaymentStatus: req.PaymentStatus,
			Page:          req.Page,
			Size:          req.Size,
		}

		items, total, err := h.OrderInteractor.List(ctx, params)
		if err != nil {
			if domain.IsParamsError(err) {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
				return
			}
			c.Error(fmt.Errorf("failed to list orders: %w", err))
			return
		}

		resItems := make([]*types.Order, 0, len(items))
		for _, o := range items {
			it, err := convertDomainOrderToResp(o)
			if err != nil {
				c.Error(fmt.Errorf("failed to convert order: %w", err))
				return
			}
			resItems = append(resItems, it)
		}

		p := req.ToPagination()
		p.SetTotal(total)

		response.Ok(c, &types.ListOrderResp{
			Items:      resItems,
			Pagination: p,
		})
	}
}

func convertDomainOrderToResp(o *domain.Order) (*types.Order, error) {
	if o == nil {
		return nil, nil
	}

	var store types.Store
	if len(o.Store) > 0 {
		if err := json.Unmarshal(o.Store, &store); err != nil {
			return nil, err
		}
	}
	var channel types.Channel
	if len(o.Channel) > 0 {
		if err := json.Unmarshal(o.Channel, &channel); err != nil {
			return nil, err
		}
	}
	var pos types.POS
	if len(o.Pos) > 0 {
		if err := json.Unmarshal(o.Pos, &pos); err != nil {
			return nil, err
		}
	}
	var cashier types.Cashier
	if len(o.Cashier) > 0 {
		if err := json.Unmarshal(o.Cashier, &cashier); err != nil {
			return nil, err
		}
	}

	var refund types.Refund
	if len(o.Refund) > 0 {
		if err := json.Unmarshal(o.Refund, &refund); err != nil {
			return nil, err
		}
	}
	var member types.Member
	if len(o.Member) > 0 {
		if err := json.Unmarshal(o.Member, &member); err != nil {
			return nil, err
		}
	}
	var takeaway types.Takeaway
	if len(o.Takeaway) > 0 {
		if err := json.Unmarshal(o.Takeaway, &takeaway); err != nil {
			return nil, err
		}
	}

	var cart []types.Product
	if len(o.Cart) > 0 {
		if err := json.Unmarshal(o.Cart, &cart); err != nil {
			return nil, err
		}
	}
	var products []types.Product
	if len(o.Products) > 0 {
		if err := json.Unmarshal(o.Products, &products); err != nil {
			return nil, err
		}
	}
	if products == nil {
		products = make([]types.Product, 0)
	}

	var promotions []types.Promotion
	if len(o.Promotions) > 0 {
		if err := json.Unmarshal(o.Promotions, &promotions); err != nil {
			return nil, err
		}
	}
	var coupons []types.Coupon
	if len(o.Coupons) > 0 {
		if err := json.Unmarshal(o.Coupons, &coupons); err != nil {
			return nil, err
		}
	}
	var taxRates []types.TaxRate
	if len(o.TaxRates) > 0 {
		if err := json.Unmarshal(o.TaxRates, &taxRates); err != nil {
			return nil, err
		}
	}
	var fees []types.Fee
	if len(o.Fees) > 0 {
		if err := json.Unmarshal(o.Fees, &fees); err != nil {
			return nil, err
		}
	}
	var payments []types.Payment
	if len(o.Payments) > 0 {
		if err := json.Unmarshal(o.Payments, &payments); err != nil {
			return nil, err
		}
	}
	var refundsProducts []types.Product
	if len(o.RefundsProducts) > 0 {
		if err := json.Unmarshal(o.RefundsProducts, &refundsProducts); err != nil {
			return nil, err
		}
	}
	var amount types.Amount
	if len(o.Amount) > 0 {
		if err := json.Unmarshal(o.Amount, &amount); err != nil {
			return nil, err
		}
	}

	res := &types.Order{
		OrderID:           o.ID.String(),
		MerchantID:        o.MerchantID,
		Store:             store,
		BusinessDate:      o.BusinessDate,
		ShiftNo:           o.ShiftNo,
		OrderNo:           o.OrderNo,
		OrderType:         o.OrderType,
		Refund:            refund,
		DiningMode:        o.DiningMode,
		Channel:           channel,
		POS:               pos,
		Cashier:           cashier,
		OrderStatus:       o.OrderStatus,
		PaymentStatus:     o.PaymentStatus,
		FulfillmentStatus: o.FulfillmentStatus,
		TableStatus:       o.TableStatus,
		TableID:           o.TableID,
		TableName:         o.TableName,
		TableCapacity:     o.TableCapacity,
		GuestCount:        o.GuestCount,
		Member:            member,
		Takeaway:          takeaway,
		Cart:              cart,
		Products:          products,
		Promotions:        promotions,
		Coupons:           coupons,
		TaxRates:          taxRates,
		Fees:              fees,
		Payments:          payments,
		RefundsProducts:   refundsProducts,
		Amount:            amount,
	}

	if o.OpenedAt != nil {
		res.OpenedAt = *o.OpenedAt
	}
	if o.PlacedAt != nil {
		res.PlacedAt = *o.PlacedAt
	}
	if o.PaidAt != nil {
		res.PaidAt = *o.PaidAt
	}
	if o.CompletedAt != nil {
		res.CompletedAt = *o.CompletedAt
	}

	res.OpenedBy = o.OpenedBy
	res.PlacedBy = o.PlacedBy
	res.PaidBy = o.PaidBy

	return res, nil
}
