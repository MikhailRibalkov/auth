package main

import (
	"context"
	deps "github.com/MikhailRibalkov/auth/pkg/auth_v1/pkg/auth_v1"
	"github.com/brianvoe/gofakeit"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"net"
)

const grpcPort = ":8081"

type server struct {
	deps.UnimplementedAuthV1Server
}

func (s *server) Get(ctx context.Context, req *deps.GetRequest) (*deps.GetResponse, error) {
	log.Printf("Client id: %d", req.GetId())

	return &deps.GetResponse{
		Info: &deps.UserInfo{
			Id:        req.GetId(),
			Name:      gofakeit.Name(),
			Email:     gofakeit.Email(),
			CreatedAt: timestamppb.New(gofakeit.Date()),
			Role:      deps.UserRole_USER,
			UpdatedAt: timestamppb.New(gofakeit.Date()),
		},
	}, nil
}

func (s *server) Create(ctx context.Context, req *deps.CreateRequest) (*deps.CreateResponse, error) {
	log.Printf("Created user: %v", req.GetUser())

	return &deps.CreateResponse{}, nil
}

func (s *server) Update(ctx context.Context, req *deps.UpdateRequest) (*emptypb.Empty, error) {
	log.Printf("Updated user, id: %d", req.GetId())

	return &emptypb.Empty{}, nil

}

func (s *server) Delete(ctx context.Context, req *deps.DeleteRequest) (*emptypb.Empty, error) {
	log.Printf("Deleted user, id: %d", req.GetId())

	return &emptypb.Empty{}, nil
}

func main() {
	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatal("Failed to listen: ", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	deps.RegisterAuthV1Server(s, &server{})

	log.Printf("server listening at #{lis.Addr()}")
	if err = s.Serve(lis); err != nil {
		log.Fatal("failed to serve: #{err}")
	}
}
