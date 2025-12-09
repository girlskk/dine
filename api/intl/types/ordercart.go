package types

import (
	"gitlab.jiguang.dev/pos-dine/dine/api/intl/pb"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func OrderCartToProto(item *domain.OrderCart) *pb.OrderCartItem {
	pbItem := &pb.OrderCartItem{
		Id:            int64(item.ID),
		TableId:       int64(item.TableID),
		ProductId:     int64(item.ProductID),
		Name:          item.Name,
		Price:         DecimalToProtoDecimal(item.Price),
		Images:        item.Images,
		ProductSpecId: int64(item.ProductSpecID),
		SpecName:      item.SpecName,
		AttrId:        int64(item.AttrID),
		AttrName:      item.AttrName,
		RecipeId:      int64(item.RecipeID),
		RecipeName:    item.RecipeName,
		CategoryId:    int64(item.CategoryID),
		CategoryName:  item.CategoryName,
		Quantity:      DecimalToProtoDecimal(item.Quantity),
		CreatedAt:     timestamppb.New(item.CreatedAt),
		UpdatedAt:     timestamppb.New(item.UpdatedAt),
	}

	if item.SetMealDetails != nil {
		pbItem.SetMealDetails = make([]*pb.SetMealDetail, len(item.SetMealDetails))
		for i, detail := range item.SetMealDetails {
			pbItem.SetMealDetails[i] = &pb.SetMealDetail{
				Id:          int64(detail.ID),
				ProductId:   int64(detail.ProductID),
				ProductName: detail.Name,
				Quantity:    DecimalToProtoDecimal(detail.Quantity),
				CreatedAt:   timestamppb.New(detail.CreatedAt),
				UpdatedAt:   timestamppb.New(detail.UpdatedAt),
			}
		}
	}
	return pbItem
}
