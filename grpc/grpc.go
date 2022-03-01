package grpc

import (
	"fmt"
        "log"
	"net"
	"context"

        "google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	pb "vatz-plugin-solana/plugin"
	execute "vatz-plugin-solana/execute"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

const (
        grpcPort = 9091
)

type service struct {
	pb.UnimplementedManagerPluginServer
}

func (s *service) Init(ctx context.Context, in *emptypb.Empty) (*pb.PluginInfo, error) {
	log.Printf("Init currently not implemented")

	return nil, nil
}

func (s *service) Verify(ctx context.Context, in *emptypb.Empty) (*pb.VerifyInfo, error) {
	log.Printf("Verify currently not implemented")

	return nil, nil
}

func (s *service) Execute(ctx context.Context, in *pb.ExecuteRequest) (*pb.ExecuteResponse, error) {
	log.Printf("Execute service")

	//req := in.GetExecuteInfo().
	//switch (req)
	//case getHealth:
	res, err := execute.GetHealth(execute.Testnet)

	return &pb.ExecuteResponse {
		State:		pb.ExecuteResponse_SUCCESS,
		Message:	res,
		ResourceType:	"Go",
	}, err
}

func StartServer() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
        if err != nil {
                log.Fatalf("failed to listen: %v", err)
        }

        s := grpc.NewServer()
	pb.RegisterManagerPluginServer(s, &service{})
	reflection.Register(s)
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
