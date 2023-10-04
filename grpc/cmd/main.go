package main

import (
	"context"
	desc "github.com/MikhailRibalkov/auth/grpc/pkg/auth_v1/pkg/auth_v1"
	"github.com/brianvoe/gofakeit"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"net"
)

const grpcPort = ":8081"

type server struct {
	desc.UnimplementedAuthV1Server
}

func (s *server) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	log.Printf("Client id: #{req.GetId}")

	return &desc.GetResponse{
		Info: &desc.UserInfo{
			Id:        req.GetId(),
			Name:      gofakeit.Name(),
			Email:     gofakeit.Email(),
			CreatedAt: timestamppb.New(gofakeit.Date()),
			Role:      desc.UserRole_ADMIN_ROLE,
			UpdatedAt: timestamppb.New(gofakeit.Date()),
		},
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatal("Failed to listen: ", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterAuthV1Server(s, &server{})

	log.Printf("server listening at #{lis.Addr()}")
	if err = s.Serve(lis); err != nil {
		log.Fatal("failed to serve: #{err}")
	}
}
