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

func (s *QueryService) MakeTransfer(ctx context.Context, req *pb.TransferRequest) (*pb.TransferResponse, error) {
	err := s.Db.MakeTransfer(shared.Transfer{From: req.From, To: req.To, Amount: req.Amount})
	return &pb.TransferResponse{}, err
}
