package server

import (
	"context"

	"github.com/JakubDaleki/transfer-app/query-service/db"
	pb "github.com/JakubDaleki/transfer-app/shared-dependencies/grpc"
)

type GreeterService struct {
	pb.UnimplementedGreeterServer
	Db *db.Database
}

// GetFeature returns the feature at the given point.
func (s *GreeterService) GetBalance(ctx context.Context, req *pb.BalanceRequest) (*pb.BalanceReponse, error) {
	b := s.Db.GetBalance(req.Username)
	return &pb.BalanceReponse{Username: b.Username, Balance: b.Balance}, nil
}
