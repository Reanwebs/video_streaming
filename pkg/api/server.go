package api

import (
	"fmt"
	"net"
	clientinterfaces "videoStreaming/pkg/client/clientInterfaces"
	"videoStreaming/pkg/config"
	"videoStreaming/pkg/pb"

	"google.golang.org/grpc"
)

type Server struct {
	gs     *grpc.Server
	Lis    net.Listener
	Port   string
	Client clientinterfaces.MonitClient
}

func NewgrpcServe(c *config.Config, service pb.VideoServiceServer, monitClient clientinterfaces.MonitClient) (*Server, error) {
	grpcserver := grpc.NewServer()
	pb.RegisterVideoServiceServer(grpcserver, service)
	lis, err := net.Listen("tcp", c.Port)
	if err != nil {
		return nil, err
	}
	return &Server{
		gs:     grpcserver,
		Lis:    lis,
		Port:   c.Port,
		Client: monitClient,
	}, nil
}

func (s *Server) Start() error {
	fmt.Println("Video service on:", s.Port)
	return s.gs.Serve(s.Lis)
}
