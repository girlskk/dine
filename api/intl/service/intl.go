package service

import (
	"context"
	"errors"

	"gitlab.jiguang.dev/pos-dine/dine/api/intl/types"

	"github.com/opentracing/opentracing-go"
	"github.com/samber/lo"
	"gitlab.jiguang.dev/pos-dine/dine/api/intl/pb"
	"gitlab.jiguang.dev/pos-dine/dine/domain"
	"gitlab.jiguang.dev/pos-dine/dine/pkg/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ pb.IntlServer = (*IntlService)(nil)

type IntlService struct {
	pb.UnimplementedIntlServer

	OrderCartInteractor domain.OrderCartInteractor
	UserInteractor      domain.FrontendUserInteractor
	TableInteractor     domain.TableInteractor
	OrderInteractor     domain.OrderInteractor
}

func NewIntlService(
	orderCartInteractor domain.OrderCartInteractor,
	userInteractor domain.FrontendUserInteractor,
	tableInteractor domain.TableInteractor,
	orderInteractor domain.OrderInteractor,
) pb.IntlServer {
	return &IntlService{
		OrderCartInteractor: orderCartInteractor,
		UserInteractor:      userInteractor,
		TableInteractor:     tableInteractor,
		OrderInteractor:     orderInteractor,
	}
}

func (s *IntlService) RegisterService(server *grpc.Server) {
	pb.RegisterIntlServer(server, s)
}

func (s *IntlService) Ping(ctx context.Context, req *pb.PingRequest) (*pb.PingReply, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "IntlService.Ping")
	defer span.Finish()

	return &pb.PingReply{Message: "pong"}, nil
}

func (s *IntlService) OrderCartGetByTableId(ctx context.Context, req *pb.OrderCartGetByTableIdRequest) (*pb.OrderCartGetByTableIdReply, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "IntlService.OrderCartGetByTableId")
	defer span.Finish()
	logger := logging.FromContext(ctx).Named("IntlService.OrderCartGetByTableId")
	ctx = logging.NewContext(ctx, logger)

	items, err := s.OrderCartInteractor.ListByTable(ctx, int(req.TableId))
	if err != nil {
		if domain.IsParamsError(err) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		} else {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}
	return &pb.OrderCartGetByTableIdReply{
		Items: lo.Map(items, func(item *domain.OrderCart, _ int) *pb.OrderCartItem {
			return types.OrderCartToProto(item)
		}),
	}, nil
}

func (s *IntlService) OrderCartAddItem(ctx context.Context, req *pb.OrderCartAddItemRequest) (*pb.OrderCartAddItemReply, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "IntlService.OrderCartAddItem")
	defer span.Finish()
	logger := logging.FromContext(ctx).Named("IntlService.OrderCartAddItem")
	ctx = logging.NewContext(ctx, logger)

	items, err := s.OrderCartInteractor.AddItem(ctx, domain.OrderCartAddParams{
		TableID:       int(req.TableId),
		ProductID:     int(req.ProductId),
		ProductSpecID: int(req.ProductSpecId),
		AttrID:        int(req.AttrId),
		RecipeID:      int(req.RecipeId),
	})
	if err != nil {
		if domain.IsParamsError(err) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		} else {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}
	return &pb.OrderCartAddItemReply{
		Items: lo.Map(items, func(item *domain.OrderCart, _ int) *pb.OrderCartItem {
			return types.OrderCartToProto(item)
		}),
	}, nil
}

func (s *IntlService) OrderCartRemoveItem(ctx context.Context, req *pb.OrderCartRemoveItemRequest) (*pb.OrderCartRemoveItemReply, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "IntlService.OrderCartRemoveItem")
	defer span.Finish()
	logger := logging.FromContext(ctx).Named("IntlService.OrderCartRemoveItem")
	ctx = logging.NewContext(ctx, logger)

	items, err := s.OrderCartInteractor.RemoveItem(ctx, int(req.ItemId), int(req.TableId))
	if err != nil {
		if domain.IsParamsError(err) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		} else {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}
	return &pb.OrderCartRemoveItemReply{
		Items: lo.Map(items, func(item *domain.OrderCart, _ int) *pb.OrderCartItem {
			return types.OrderCartToProto(item)
		}),
	}, nil
}

func (s *IntlService) GetFrontendUserByToken(ctx context.Context, req *pb.GetFrontendUserByTokenRequest) (*pb.FrontendUser, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "IntlService.GetFrontendUserByToken")
	defer span.Finish()
	logger := logging.FromContext(ctx).Named("IntlService.GetFrontendUserByToken")
	ctx = logging.NewContext(ctx, logger)

	user, err := s.UserInteractor.Authenticate(ctx, req.Token)
	if err != nil {
		if errors.Is(err, domain.ErrTokenInvalid) {
			return nil, status.Error(codes.Unauthenticated, err.Error())
		}

		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.FrontendUser{
		Id:       int64(user.ID),
		Username: user.Username,
		Nickname: user.Nickname,
		StoreId:  int64(user.StoreID),
	}, nil
}

func (s *IntlService) GetTableInfo(ctx context.Context, req *pb.GetTableInfoRequest) (*pb.GetTableInfoReply, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "IntlService.GetTableInfo")
	defer span.Finish()
	logger := logging.FromContext(ctx).Named("IntlService.GetTableInfo")
	ctx = logging.NewContext(ctx, logger)

	table, err := s.TableInteractor.GetWithOrder(ctx, int(req.TableId))
	if err != nil {
		if domain.IsParamsError(err) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		} else {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	orderCartItems, err := s.OrderCartInteractor.ListByTable(ctx, table.ID)
	if err != nil {
		if domain.IsParamsError(err) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		} else {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	var orderNo string
	if table.Order != nil {
		orderNo = table.Order.No
	}

	return &pb.GetTableInfoReply{
		Id:      int64(table.ID),
		StoreId: int64(table.StoreID),
		OrderNo: orderNo,
		Items: lo.Map(orderCartItems, func(item *domain.OrderCart, _ int) *pb.OrderCartItem {
			return types.OrderCartToProto(item)
		}),
	}, nil
}
