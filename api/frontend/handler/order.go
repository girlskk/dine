package handler

import (
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
//	@Success	200		{object}	types.OrderResp			"成功"
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

		merchantID := req.MerchantID
		storeID := req.StoreID

		o := &domain.Order{
			ID:                uuid.New(),
			MerchantID:        merchantID,
			StoreID:           storeID,
			BusinessDate:      req.BusinessDate,
			ShiftNo:           req.ShiftNo,
			OrderNo:           req.OrderNo,
			OrderType:         domain.OrderType(req.OrderType),
			OriginOrderID:     req.OriginOrderID.String(),
			DiningMode:        domain.DiningMode(req.DiningMode),
			OrderStatus:       domain.OrderStatus(req.OrderStatus),
			PaymentStatus:     domain.PaymentStatus(req.PaymentStatus),
			FulfillmentStatus: domain.FulfillmentStatus(req.FulfillmentStatus),
			TableStatus:       domain.TableStatus(req.TableStatus),
			TableID:           req.TableID,
			TableName:         req.TableName,
			TableCapacity:     req.TableCapacity,
			GuestCount:        req.GuestCount,
			OpenedBy:          req.OpenedBy,
			PlacedBy:          req.PlacedBy,
			PaidBy:            req.PaidBy,
		}

		o.Refund = toDomainRefund(req.Refund)

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

		o.Store = toDomainStore(req.Store)
		if o.Store == nil {
			o.Store = &domain.OrderStore{StoreID: req.StoreID}
		}
		o.Channel = toDomainChannel(req.Channel)
		o.Pos = toDomainPOS(req.POS)
		o.Cashier = toDomainCashier(req.Cashier)
		o.Member = toDomainMember(req.Member)
		o.Takeaway = toDomainTakeaway(req.Takeaway)
		o.Cart = toDomainProducts(req.Cart)
		o.Products = toDomainProducts(req.Products)
		o.Promotions = toDomainPromotions(req.Promotions)
		o.Coupons = toDomainCoupons(req.Coupons)
		o.TaxRates = toDomainTaxRates(req.TaxRates)
		o.Fees = toDomainFees(req.Fees)
		o.Payments = toDomainPayments(req.Payments)
		o.RefundsProducts = toDomainProducts(req.RefundsProducts)
		o.Amount = toDomainAmount(req.Amount)

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
//	@Param		id	path		string			true	"订单ID"
//	@Success	200	{object}	types.OrderResp	"成功"
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
//	@Success	200		{object}	types.OrderResp			"成功"
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
			o.OrderType = domain.OrderType(*req.OrderType)
		}
		if req.OriginOrderID != nil {
			o.OriginOrderID = req.OriginOrderID.String()
		}
		if req.Refund != nil {
			o.Refund = toDomainRefund(req.Refund)
		}

		if req.DiningMode != nil {
			o.DiningMode = domain.DiningMode(*req.DiningMode)
		}
		if req.OrderStatus != nil {
			o.OrderStatus = domain.OrderStatus(*req.OrderStatus)
		}
		if req.PaymentStatus != nil {
			o.PaymentStatus = domain.PaymentStatus(*req.PaymentStatus)
		}
		if req.FulfillmentStatus != nil {
			o.FulfillmentStatus = domain.FulfillmentStatus(*req.FulfillmentStatus)
		}
		if req.TableStatus != nil {
			o.TableStatus = domain.TableStatus(*req.TableStatus)
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
			o.Store = toDomainStore(req.Store)
		}
		if req.Channel != nil {
			o.Channel = toDomainChannel(req.Channel)
		}
		if req.POS != nil {
			o.Pos = toDomainPOS(req.POS)
		}
		if req.Cashier != nil {
			o.Cashier = toDomainCashier(req.Cashier)
		}
		if req.Member != nil {
			o.Member = toDomainMember(req.Member)
		}
		if req.Takeaway != nil {
			o.Takeaway = toDomainTakeaway(req.Takeaway)
		}
		if req.Cart != nil {
			o.Cart = toDomainProducts(req.Cart)
		}
		if req.Products != nil {
			o.Products = toDomainProducts(req.Products)
		}
		if req.Promotions != nil {
			o.Promotions = toDomainPromotions(req.Promotions)
		}
		if req.Coupons != nil {
			o.Coupons = toDomainCoupons(req.Coupons)
		}
		if req.TaxRates != nil {
			o.TaxRates = toDomainTaxRates(req.TaxRates)
		}
		if req.Fees != nil {
			o.Fees = toDomainFees(req.Fees)
		}
		if req.Payments != nil {
			o.Payments = toDomainPayments(req.Payments)
		}
		if req.RefundsProducts != nil {
			o.RefundsProducts = toDomainProducts(req.RefundsProducts)
		}
		if req.Amount != nil {
			o.Amount = toDomainAmount(req.Amount)
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
//	@Param		order_type		query		string				false	"订单类型"	Enums(SALE,REFUND,PARTIAL_REFUND)
//	@Param		order_status	query		string				false	"订单状态"	Enums(DRAFT,PLACED,IN_PROGRESS,READY,COMPLETED,CANCELLED,VOIDED,MERGED)
//	@Param		payment_status	query		string				false	"支付状态"	Enums(UNPAID,PAYING,PARTIALLY_PAID,PAID,PARTIALLY_REFUNDED,REFUNDED)
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

		merchantUUID := req.MerchantID
		storeUUID := req.StoreID

		params := domain.OrderListParams{
			MerchantID:    merchantUUID,
			StoreID:       storeUUID,
			BusinessDate:  req.BusinessDate,
			OrderNo:       req.OrderNo,
			OrderType:     domain.OrderType(req.OrderType),
			OrderStatus:   domain.OrderStatus(req.OrderStatus),
			PaymentStatus: domain.PaymentStatus(req.PaymentStatus),
			Page:          req.Page,
			Size:          req.Size,
		}

		params.MerchantID = merchantUUID
		params.StoreID = storeUUID

		items, total, err := h.OrderInteractor.List(ctx, params)
		if err != nil {
			if domain.IsParamsError(err) {
				c.Error(errorx.New(http.StatusBadRequest, errcode.InvalidParams, err))
				return
			}
			c.Error(fmt.Errorf("failed to list orders: %w", err))
			return
		}

		resItems := make([]*types.OrderResp, 0, len(items))
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

func convertDomainOrderToResp(o *domain.Order) (*types.OrderResp, error) {
	if o == nil {
		return nil, nil
	}

	store := toTypesStore(o.Store)
	channel := toTypesChannel(o.Channel)
	pos := toTypesPOS(o.Pos)
	cashier := toTypesCashier(o.Cashier)
	refund := toTypesRefund(o.Refund)
	member := toTypesMember(o.Member)
	takeaway := toTypesTakeaway(o.Takeaway)
	cart := toTypesProducts(o.Cart)
	products := toTypesProducts(o.Products)
	if products == nil {
		products = make([]types.Product, 0)
	}
	promotions := toTypesPromotions(o.Promotions)
	coupons := toTypesCoupons(o.Coupons)
	taxRates := toTypesTaxRates(o.TaxRates)
	fees := toTypesFees(o.Fees)
	payments := toTypesPayments(o.Payments)
	refundsProducts := toTypesProducts(o.RefundsProducts)
	amount := toTypesAmount(o.Amount)

	res := &types.OrderResp{
		OrderID:           o.ID,
		MerchantID:        o.MerchantID,
		Store:             store,
		BusinessDate:      o.BusinessDate,
		ShiftNo:           o.ShiftNo,
		OrderNo:           o.OrderNo,
		OrderType:         string(o.OrderType),
		Refund:            refund,
		DiningMode:        string(o.DiningMode),
		Channel:           channel,
		POS:               pos,
		Cashier:           cashier,
		OrderStatus:       string(o.OrderStatus),
		PaymentStatus:     string(o.PaymentStatus),
		FulfillmentStatus: string(o.FulfillmentStatus),
		TableStatus:       string(o.TableStatus),
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

func toDomainRefund(r *types.Refund) *domain.OrderRefund {
	if r == nil {
		return nil
	}
	return &domain.OrderRefund{
		OriginOrderID: r.OriginOrderID,
		OriginOrderNo: r.OriginOrderNo,
		Reason:        r.Reason,
	}
}

func toDomainMember(m *types.Member) *domain.OrderMember {
	if m == nil {
		return nil
	}
	return &domain.OrderMember{
		MemberID:        m.MemberID,
		MemberNo:        m.MemberNo,
		MemberName:      m.MemberName,
		MemberPhone:     m.MemberPhone,
		MemberLevelName: m.MemberLevelName,
	}
}

func toDomainStore(s *types.Store) *domain.OrderStore {
	if s == nil {
		return nil
	}
	return &domain.OrderStore{
		StoreID:      s.StoreID,
		StoreNo:      s.StoreNo,
		StoreName:    s.StoreName,
		StorePhone:   s.StorePhone,
		StoreAddress: s.StoreAddress,
	}
}

func toDomainChannel(ch *types.Channel) *domain.OrderChannel {
	if ch == nil {
		return nil
	}
	return &domain.OrderChannel{Code: ch.Code, Name: ch.Name}
}

func toDomainCashier(ca *types.Cashier) *domain.OrderCashier {
	if ca == nil {
		return nil
	}
	return &domain.OrderCashier{CashierID: ca.CashierID, CashierName: ca.CashierName}
}

func toDomainPOS(p *types.POS) *domain.OrderPOS {
	if p == nil {
		return nil
	}
	return &domain.OrderPOS{PosID: p.PosID, PosCode: p.PosCode, DeviceID: p.DeviceID}
}

func toDomainTakeaway(tk *types.Takeaway) *domain.OrderTakeaway {
	if tk == nil {
		return nil
	}
	return &domain.OrderTakeaway{
		TakeawayType:       tk.TakeawayType,
		ContactName:        tk.ContactName,
		ContactPhone:       tk.ContactPhone,
		PickupNo:           tk.PickupNo,
		PickupEtaAt:        tk.PickupEtaAt,
		DeliveryAddress:    tk.DeliveryAddress,
		DeliveryFee:        tk.DeliveryFee,
		DeliveryPlatform:   tk.DeliveryPlatform,
		DeliveryOrderNo:    tk.DeliveryOrderNo,
		DeliveryTrackingNo: tk.DeliveryTrackingNo,
		DeliveryStatus:     tk.DeliveryStatus,
		DeliveryRiderName:  tk.DeliveryRiderName,
		DeliveryRiderPhone: tk.DeliveryRiderPhone,
		DeliveryRemark:     tk.DeliveryRemark,
	}
}

func toDomainPromotion(p types.Promotion) domain.OrderPromotion {
	return domain.OrderPromotion{
		PromotionID:    p.PromotionID,
		PromotionName:  p.PromotionName,
		PromotionType:  p.PromotionType,
		DiscountAmount: p.DiscountAmount,
		Meta:           p.Meta,
	}
}

func toDomainProduct(p types.Product) domain.OrderProduct {
	promotions := make([]domain.OrderPromotion, 0, len(p.Promotions))
	for _, pr := range p.Promotions {
		promotions = append(promotions, toDomainPromotion(pr))
	}
	return domain.OrderProduct{
		OrderItemID:       p.OrderItemID,
		Index:             p.Index,
		RefundReason:      p.RefundReason,
		RefundedBy:        p.RefundedBy,
		RefundedAt:        p.RefundedAt,
		Promotions:        promotions,
		PromotionDiscount: p.PromotionDiscount,
		ProductID:         p.ProductID,
		ProductName:       p.ProductName,
		SkuID:             p.SkuID,
		SkuName:           p.SkuName,
		Qty:               p.Qty,
		Price:             p.Price,
		Subtotal:          p.Subtotal,
		DiscountAmount:    p.DiscountAmount,
		AmountBeforeTax:   p.AmountBeforeTax,
		TaxRate:           p.TaxRate,
		Tax:               p.Tax,
		AmountAfterTax:    p.AmountAfterTax,
		Total:             p.Total,
		VoidQty:           p.VoidQty,
		VoidAmount:        p.VoidAmount,
		Note:              p.Note,
		Options:           p.Options,
	}
}

func toDomainProducts(ps *[]types.Product) *[]domain.OrderProduct {
	if ps == nil {
		return nil
	}
	res := make([]domain.OrderProduct, 0, len(*ps))
	for _, p := range *ps {
		res = append(res, toDomainProduct(p))
	}
	return &res
}

func toDomainPromotions(ps *[]types.Promotion) *[]domain.OrderPromotion {
	if ps == nil {
		return nil
	}
	res := make([]domain.OrderPromotion, 0, len(*ps))
	for _, p := range *ps {
		res = append(res, toDomainPromotion(p))
	}
	return &res
}

func toDomainCoupons(cs *[]types.Coupon) *[]domain.OrderCoupon {
	if cs == nil {
		return nil
	}
	res := make([]domain.OrderCoupon, 0, len(*cs))
	for _, c := range *cs {
		res = append(res, domain.OrderCoupon{
			CouponID:       c.CouponID,
			CouponName:     c.CouponName,
			CouponType:     c.CouponType,
			CouponCode:     c.CouponCode,
			DiscountAmount: c.DiscountAmount,
			Meta:           c.Meta,
		})
	}
	return &res
}

func toDomainTaxRates(ts *[]types.TaxRate) *[]domain.OrderTaxRate {
	if ts == nil {
		return nil
	}
	res := make([]domain.OrderTaxRate, 0, len(*ts))
	for _, t := range *ts {
		res = append(res, domain.OrderTaxRate{
			TaxRateID:     t.TaxRateID,
			TaxRateName:   t.TaxRateName,
			Rate:          t.Rate,
			TaxableAmount: t.TaxableAmount,
			TaxAmount:     t.TaxAmount,
			Meta:          t.Meta,
		})
	}
	return &res
}

func toDomainFees(fs *[]types.Fee) *[]domain.OrderFee {
	if fs == nil {
		return nil
	}
	res := make([]domain.OrderFee, 0, len(*fs))
	for _, f := range *fs {
		res = append(res, domain.OrderFee{
			FeeID:   f.FeeID,
			FeeName: f.FeeName,
			FeeType: f.FeeType,
			Amount:  f.Amount,
			Meta:    f.Meta,
		})
	}
	return &res
}

func toDomainPayments(ps *[]types.Payment) *[]domain.OrderPayment {
	if ps == nil {
		return nil
	}
	res := make([]domain.OrderPayment, 0, len(*ps))
	for _, p := range *ps {
		res = append(res, domain.OrderPayment{
			PaymentNo:     p.PaymentNo,
			PaymentMethod: p.PaymentMethod,
			PaymentAmount: p.PaymentAmount,
			POS:           domain.OrderPOS{PosID: p.POS.PosID, PosCode: p.POS.PosCode, DeviceID: p.POS.DeviceID},
			Cashier:       domain.OrderCashier{CashierID: p.Cashier.CashierID, CashierName: p.Cashier.CashierName},
			PaidAt:        p.PaidAt,
		})
	}
	return &res
}

func toDomainAmount(a *types.Amount) *domain.OrderAmount {
	if a == nil {
		return nil
	}
	return &domain.OrderAmount{
		ItemsSubtotal:          a.ItemsSubtotal,
		DiscountTotal:          a.DiscountTotal,
		PromotionDiscountTotal: a.PromotionDiscountTotal,
		VoucherDiscountTotal:   a.VoucherDiscountTotal,
		TaxTotal:               a.TaxTotal,
		ServiceFeeTotal:        a.ServiceFeeTotal,
		DeliveryFee:            a.DeliveryFee,
		FeeTotal:               a.FeeTotal,
		RoundingAmount:         a.RoundingAmount,
		AmountDue:              a.AmountDue,
		AmountPaid:             a.AmountPaid,
		ChangeAmount:           a.ChangeAmount,
		AmountRefunded:         a.AmountRefunded,
	}
}

func toTypesStore(s *domain.OrderStore) types.Store {
	if s == nil {
		return types.Store{}
	}
	return types.Store{StoreID: s.StoreID, StoreNo: s.StoreNo, StoreName: s.StoreName, StorePhone: s.StorePhone, StoreAddress: s.StoreAddress}
}

func toTypesChannel(ch *domain.OrderChannel) types.Channel {
	if ch == nil {
		return types.Channel{}
	}
	return types.Channel{Code: ch.Code, Name: ch.Name}
}

func toTypesCashier(ca *domain.OrderCashier) types.Cashier {
	if ca == nil {
		return types.Cashier{}
	}
	return types.Cashier{CashierID: ca.CashierID, CashierName: ca.CashierName}
}

func toTypesPOS(p *domain.OrderPOS) types.POS {
	if p == nil {
		return types.POS{}
	}
	return types.POS{PosID: p.PosID, PosCode: p.PosCode, DeviceID: p.DeviceID}
}

func toTypesMember(m *domain.OrderMember) types.Member {
	if m == nil {
		return types.Member{}
	}
	return types.Member{MemberID: m.MemberID, MemberNo: m.MemberNo, MemberName: m.MemberName, MemberPhone: m.MemberPhone, MemberLevelName: m.MemberLevelName}
}

func toTypesTakeaway(tk *domain.OrderTakeaway) types.Takeaway {
	if tk == nil {
		return types.Takeaway{}
	}
	return types.Takeaway{
		TakeawayType:       tk.TakeawayType,
		ContactName:        tk.ContactName,
		ContactPhone:       tk.ContactPhone,
		PickupNo:           tk.PickupNo,
		PickupEtaAt:        tk.PickupEtaAt,
		DeliveryAddress:    tk.DeliveryAddress,
		DeliveryFee:        tk.DeliveryFee,
		DeliveryPlatform:   tk.DeliveryPlatform,
		DeliveryOrderNo:    tk.DeliveryOrderNo,
		DeliveryTrackingNo: tk.DeliveryTrackingNo,
		DeliveryStatus:     tk.DeliveryStatus,
		DeliveryRiderName:  tk.DeliveryRiderName,
		DeliveryRiderPhone: tk.DeliveryRiderPhone,
		DeliveryRemark:     tk.DeliveryRemark,
	}
}

func toTypesRefund(r *domain.OrderRefund) types.Refund {
	if r == nil {
		return types.Refund{}
	}
	return types.Refund{OriginOrderID: r.OriginOrderID, OriginOrderNo: r.OriginOrderNo, Reason: r.Reason}
}

func toTypesPromotion(p domain.OrderPromotion) types.Promotion {
	return types.Promotion{PromotionID: p.PromotionID, PromotionName: p.PromotionName, PromotionType: p.PromotionType, DiscountAmount: p.DiscountAmount, Meta: p.Meta}
}

func toTypesProduct(p domain.OrderProduct) types.Product {
	promotions := make([]types.Promotion, 0, len(p.Promotions))
	for _, pr := range p.Promotions {
		promotions = append(promotions, toTypesPromotion(pr))
	}
	return types.Product{
		OrderItemID:       p.OrderItemID,
		Index:             p.Index,
		RefundReason:      p.RefundReason,
		RefundedBy:        p.RefundedBy,
		RefundedAt:        p.RefundedAt,
		Promotions:        promotions,
		PromotionDiscount: p.PromotionDiscount,
		ProductID:         p.ProductID,
		ProductName:       p.ProductName,
		SkuID:             p.SkuID,
		SkuName:           p.SkuName,
		Qty:               p.Qty,
		Price:             p.Price,
		Subtotal:          p.Subtotal,
		DiscountAmount:    p.DiscountAmount,
		AmountBeforeTax:   p.AmountBeforeTax,
		TaxRate:           p.TaxRate,
		Tax:               p.Tax,
		AmountAfterTax:    p.AmountAfterTax,
		Total:             p.Total,
		VoidQty:           p.VoidQty,
		VoidAmount:        p.VoidAmount,
		Note:              p.Note,
		Options:           p.Options,
	}
}

func toTypesProducts(ps *[]domain.OrderProduct) []types.Product {
	if ps == nil {
		return nil
	}
	res := make([]types.Product, 0, len(*ps))
	for _, p := range *ps {
		res = append(res, toTypesProduct(p))
	}
	return res
}

func toTypesPromotions(ps *[]domain.OrderPromotion) []types.Promotion {
	if ps == nil {
		return nil
	}
	res := make([]types.Promotion, 0, len(*ps))
	for _, p := range *ps {
		res = append(res, toTypesPromotion(p))
	}
	return res
}

func toTypesCoupons(cs *[]domain.OrderCoupon) []types.Coupon {
	if cs == nil {
		return nil
	}
	res := make([]types.Coupon, 0, len(*cs))
	for _, c := range *cs {
		res = append(res, types.Coupon{CouponID: c.CouponID, CouponName: c.CouponName, CouponType: c.CouponType, CouponCode: c.CouponCode, DiscountAmount: c.DiscountAmount, Meta: c.Meta})
	}
	return res
}

func toTypesTaxRates(ts *[]domain.OrderTaxRate) []types.TaxRate {
	if ts == nil {
		return nil
	}
	res := make([]types.TaxRate, 0, len(*ts))
	for _, t := range *ts {
		res = append(res, types.TaxRate{TaxRateID: t.TaxRateID, TaxRateName: t.TaxRateName, Rate: t.Rate, TaxableAmount: t.TaxableAmount, TaxAmount: t.TaxAmount, Meta: t.Meta})
	}
	return res
}

func toTypesFees(fs *[]domain.OrderFee) []types.Fee {
	if fs == nil {
		return nil
	}
	res := make([]types.Fee, 0, len(*fs))
	for _, f := range *fs {
		res = append(res, types.Fee{FeeID: f.FeeID, FeeName: f.FeeName, FeeType: f.FeeType, Amount: f.Amount, Meta: f.Meta})
	}
	return res
}

func toTypesPayments(ps *[]domain.OrderPayment) []types.Payment {
	if ps == nil {
		return nil
	}
	res := make([]types.Payment, 0, len(*ps))
	for _, p := range *ps {
		res = append(res, types.Payment{
			PaymentNo:     p.PaymentNo,
			PaymentMethod: p.PaymentMethod,
			PaymentAmount: p.PaymentAmount,
			POS:           types.POS{PosID: p.POS.PosID, PosCode: p.POS.PosCode, DeviceID: p.POS.DeviceID},
			Cashier:       types.Cashier{CashierID: p.Cashier.CashierID, CashierName: p.Cashier.CashierName},
			PaidAt:        p.PaidAt,
		})
	}
	return res
}

func toTypesAmount(a *domain.OrderAmount) types.Amount {
	if a == nil {
		return types.Amount{}
	}
	return types.Amount{
		ItemsSubtotal:          a.ItemsSubtotal,
		DiscountTotal:          a.DiscountTotal,
		PromotionDiscountTotal: a.PromotionDiscountTotal,
		VoucherDiscountTotal:   a.VoucherDiscountTotal,
		TaxTotal:               a.TaxTotal,
		ServiceFeeTotal:        a.ServiceFeeTotal,
		DeliveryFee:            a.DeliveryFee,
		FeeTotal:               a.FeeTotal,
		RoundingAmount:         a.RoundingAmount,
		AmountDue:              a.AmountDue,
		AmountPaid:             a.AmountPaid,
		ChangeAmount:           a.ChangeAmount,
		AmountRefunded:         a.AmountRefunded,
	}
}
