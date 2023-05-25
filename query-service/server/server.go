package server

import (
	"context"

	"github.com/JakubDaleki/transfer-app/query-service/db"
	"github.com/JakubDaleki/transfer-app/shared-dependencies"
	pb "github.com/JakubDaleki/transfer-app/shared-dependencies/grpc"
)

type QueryService struct {
	pb.UnimplementedQueryServiceServer
	Db *db.Database
}

func (s *QueryService) GetBalance(ctx context.Context, req *pb.BalanceRequest) (*pb.BalanceReponse, error) {
	b := s.Db.GetBalance(req.Username)
	return &pb.BalanceReponse{Username: b.Username, Balance: b.Balance}, nil
}

func (s *QueryService) UpdateBalance(ctx context.Context, req *pb.UpdateBalanceRequest) (*pb.UpdateBalanceResponse, error) {
	err := s.Db.UpdateBalance(shared.Balance{Username: req.User, Balance: req.Amount})
	return &pb.UpdateBalanceResponse{}, err
}
