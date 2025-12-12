package service

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"gitlab.jiguang.dev/pos-dine/dine/api/intl/pb"
	"google.golang.org/grpc"
)

var _ pb.IntlServer = (*IntlService)(nil)

type IntlService struct {
	pb.UnimplementedIntlServer
}

func NewIntlService() pb.IntlServer {
	return &IntlService{}
}

func (s *IntlService) RegisterService(server *grpc.Server) {
	pb.RegisterIntlServer(server, s)
}

func (s *IntlService) Ping(ctx context.Context, req *pb.PingRequest) (*pb.PingReply, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "IntlService.Ping")
	defer span.Finish()

	return &pb.PingReply{Message: "pong"}, nil
}
