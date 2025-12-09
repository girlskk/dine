package repository

import (
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/ent"
)

// 转换购物车项列表
func convertOrderCarts(items []*ent.OrderCart) domain.OrderCarts {
	var result domain.OrderCarts
	for _, item := range items {
		result = append(result, convertOrderCart(item))
	}
	return result
}

// 转换单个购物车项
func convertOrderCart(item *ent.OrderCart) *domain.OrderCart {
	if item == nil {
		return nil
	}

	cart := &domain.OrderCart{
		ID:            item.ID,
		TableID:       item.TableID,
		ProductID:     item.ProductID,
		ProductSpecID: item.ProductSpecID,
		AttrID:        item.AttrID,
		RecipeID:      item.RecipeID,
		Quantity:      item.Quantity,
		CreatedAt:     item.CreatedAt,
		UpdatedAt:     item.UpdatedAt,
	}

	// 填充商品信息
	if item.Edges.Product != nil {
		product := item.Edges.Product

		cart.Name = product.Name
		cart.Price = product.Price
		cart.Images = product.Images

		if product.Edges.Category != nil {
			cart.CategoryID = product.Edges.Category.ID
			cart.CategoryName = product.Edges.Category.Name
		}

		if item.Edges.ProductSpec != nil && item.Edges.ProductSpec.Edges.Spec != nil {
			cart.SpecName = item.Edges.ProductSpec.Edges.Spec.Name
		}
		if item.Edges.Attr != nil {
			cart.AttrName = item.Edges.Attr.Name
		}
		if item.Edges.Recipe != nil {
			cart.RecipeName = item.Edges.Recipe.Name
		}

		// 如果是多规格商品，使用规格价格
		if product.Type == int(domain.ProductTypeMulti) && item.Edges.ProductSpec != nil {
			cart.Price = item.Edges.ProductSpec.Price
		}
		// 填充套餐详情（如果是套餐商品）
		if product.Type == int(domain.ProductTypeSetMeal) && len(product.Edges.SetMealDetails) > 0 {
			cart.SetMealDetails = convertSetMealDetails(product.Edges.SetMealDetails)
		}
	}
	return cart
}

// 转换套餐详情列表
func convertSetMealDetails(details []*ent.SetMealDetail) domain.SetMealDetails {
	var result domain.SetMealDetails
	for _, detail := range details {
		result = append(result, convertSetMealDetail(detail))
	}
	return result
}

// 转换单个套餐详情
func convertSetMealDetail(detail *ent.SetMealDetail) *domain.SetMealDetail {
	if detail == nil {
		return nil
	}
	setMealDetail := &domain.SetMealDetail{
		ID:        detail.ID,
		ProductID: detail.ProductID,
		Quantity:  detail.Quantity,
		CreatedAt: detail.CreatedAt,
		UpdatedAt: detail.UpdatedAt,
	}

	// 填充商品名称
	if detail.Edges.Product != nil {
		setMealDetail.Name = detail.Edges.Product.Name
	}
	return setMealDetail
}
