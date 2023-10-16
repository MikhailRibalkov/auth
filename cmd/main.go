package main

import (
	context "context"
	deps "github.com/MikhailRibalkov/auth/pkg/auth_v1/pkg/auth_v1"
	sq "github.com/MikhailRibalkov/auth/pkg/auth_v1/postgres/query_with_squirrel"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"net"
)

const grpcPort = ":8081"

type server struct {
	deps.UnimplementedAuthV1Server
}

func (s *server) Get(ctx context.Context, req *deps.GetRequest) (*deps.GetResponse, error) {
	log.Printf("Client id: %d", req.GetId())

	client := sq.PgClient{}
	userInfo, err := client.GetUserInfo(req)
	if err != nil {
		log.Printf("error of getting info of user: %s", err)
		return nil, err
	}
	return &deps.GetResponse{
		Info: &deps.UserInfo{
			Id:        req.GetId(),
			Name:      userInfo.GetName(),
			Email:     userInfo.GetEmail(),
			CreatedAt: userInfo.GetCreatedAt(),
			Role:      userInfo.GetRole(),
			UpdatedAt: userInfo.GetUpdatedAt(),
		},
	}, nil
}

func (s *server) Create(ctx context.Context, req *deps.CreateRequest) (*deps.CreateResponse, error) {
	client := sq.PgClient{}
	_, err := client.CreateUser(req)
	if err != nil {
		log.Printf("can not create user: %s", err)
		return nil, err
	}

	log.Printf("Created user: %v", req.GetUser())

	return &deps.CreateResponse{}, nil
}

func (s *server) Update(ctx context.Context, req *deps.UpdateRequest) (*emptypb.Empty, error) {
	client := sq.PgClient{}
	id, err := client.UpdateUser(req)
	if err != nil {
		log.Printf("Update error: %s", err)
		return nil, err
	}
	log.Printf("Updated user, id: %d", id)

	return &emptypb.Empty{}, nil

}

func (s *server) Delete(ctx context.Context, req *deps.DeleteRequest) (*emptypb.Empty, error) {
	client := sq.PgClient{}
	id, err := client.DeleteUser(req)
	if err != nil {
		log.Printf("Delete error: %s", err)
		return nil, err
	}

	log.Printf("Deleted user, id: %d", id)

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
