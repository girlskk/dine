package repository

import (
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/ent"
)

func convertOrder(od *ent.Order) *domain.Order {
	if od == nil {
		return nil
	}

	tableID := 0
	if od.TableID != nil {
		tableID = *od.TableID
	}

	dorder := &domain.Order{
		ID:                   od.ID,
		No:                   od.No,
		Type:                 od.Type,
		Source:               od.Source,
		Status:               od.Status,
		TotalPrice:           od.TotalPrice,
		Discount:             od.Discount,
		RealPrice:            od.RealPrice,
		PointsAvailable:      od.PointsAvailable,
		MemberID:             od.MemberID,
		MemberName:           od.MemberName,
		MemberPhone:          od.MemberPhone,
		StoreID:              od.StoreID,
		StoreName:            od.StoreName,
		TableID:              tableID,
		TableName:            od.TableName,
		PeopleNumber:         od.PeopleNumber,
		CreatorID:            od.CreatorID,
		CreatorName:          od.CreatorName,
		Paid:                 od.Paid,
		Refunded:             od.Refunded,
		PaidChannels:         od.PaidChannels,
		CashPaid:             od.CashPaid,
		WechatPaid:           od.WechatPaid,
		WechatRefunded:       od.WechatRefunded,
		AlipayPaid:           od.AlipayPaid,
		AlipayRefunded:       od.AlipayRefunded,
		PointsPaid:           od.PointsPaid,
		PointsRefunded:       od.PointsRefunded,
		PointsWalletPaid:     od.PointsWalletPaid,
		PointsWalletRefunded: od.PointsWalletRefunded,
		LastPaidAt:           od.LastPaidAt,
		FinishedAt:           od.FinishedAt,
		CreatedAt:            od.CreatedAt,
		UpdatedAt:            od.UpdatedAt,
		Items:                convertOrderItems(od.Edges.Items),
		Logs:                 convertOrderLogs(od.Edges.Logs),
	}

	return dorder
}

func convertOrderItems(items []*ent.OrderItem) []*domain.OrderItem {
	var ditems []*domain.OrderItem
	for _, item := range items {
		ditems = append(ditems, convertOrderItem(item))
	}
	return ditems
}

func convertOrderItem(item *ent.OrderItem) *domain.OrderItem {
	return &domain.OrderItem{
		ID:              item.ID,
		OrderID:         item.OrderID,
		ProductID:       item.ProductID,
		Name:            item.Name,
		Type:            domain.ProductType(item.Type),
		AllowPointPay:   item.AllowPointPay,
		Quantity:        item.Quantity,
		Price:           item.Price,
		Amount:          item.Amount,
		ProductSnapshot: item.ProductSnapshot,
		SetMealDetails:  convertOrderItemSetMealDetails(item.Edges.SetMealDetails),
		Remark:          item.Remark,
		CreatedAt:       item.CreatedAt,
		UpdatedAt:       item.UpdatedAt,
	}
}

func convertOrderItemSetMealDetails(details []*ent.OrderItemSetMealDetail) []*domain.OrderItemSetMealDetail {
	var setMealDetails []*domain.OrderItemSetMealDetail
	for _, detail := range details {
		setMealDetails = append(setMealDetails, &domain.OrderItemSetMealDetail{
			ID:              detail.ID,
			OrderItemID:     detail.OrderItemID,
			Name:            detail.Name,
			Type:            domain.ProductType(detail.Type),
			SetMealPrice:    detail.SetMealPrice,
			SetMealID:       detail.SetMealID,
			ProductID:       detail.ProductID,
			Quantity:        detail.Quantity,
			ProductSnapshot: detail.ProductSnapshot,
			CreatedAt:       detail.CreatedAt,
			UpdatedAt:       detail.UpdatedAt,
		})
	}
	return setMealDetails
}

func convertOrderLogs(logs []*ent.OrderLog) []*domain.OrderLog {
	var dlogs []*domain.OrderLog
	for _, log := range logs {
		dlogs = append(dlogs, convertOrderLog(log))
	}
	return dlogs
}

func convertOrderLog(log *ent.OrderLog) *domain.OrderLog {
	return &domain.OrderLog{
		ID:           log.ID,
		OrderID:      log.OrderID,
		Event:        log.Event,
		OperatorType: log.OperatorType,
		OperatorID:   log.OperatorID,
		OperatorName: log.OperatorName,
		CreatedAt:    log.CreatedAt,
	}
}

func convertOrderFinanceLog(log *ent.OrderFinanceLog) *domain.OrderFinanceLog {
	return &domain.OrderFinanceLog{
		ID:          log.ID,
		OrderID:     log.OrderID,
		Amount:      log.Amount,
		Type:        log.Type,
		Channel:     log.Channel,
		SeqNo:       log.SeqNo,
		CreatorType: log.CreatorType,
		CreatorID:   log.CreatorID,
		CreatorName: log.CreatorName,
		CreatedAt:   log.CreatedAt,
		UpdatedAt:   log.UpdatedAt,
	}
}
