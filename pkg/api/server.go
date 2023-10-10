package api

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	clientinterfaces "videoStreaming/pkg/client/clientInterfaces"
	"videoStreaming/pkg/config"
	"videoStreaming/pkg/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Server struct {
	gs     *grpc.Server
	Lis    net.Listener
	Port   string
	Client clientinterfaces.MonitClient
}

func NewgrpcServer(c *config.Config, service pb.VideoServiceServer, monitClient clientinterfaces.MonitClient) (*Server, error) {
	// Listen on the specified port
	lis, err := net.Listen("tcp", c.Port)
	if err != nil {
		return nil, fmt.Errorf("failed to listen on port %s: %w", c.Port, err)
	}

	// Load CA certificate
	caPem, err := os.ReadFile("cert/ca-cert.pem")
	if err != nil {
		return nil, fmt.Errorf("failed to read CA certificate: %w", err)
	}

	// Create certificate pool and append CA certificate
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(caPem) {
		return nil, errors.New("failed to append CA certificate to certificate pool")
	}

	// Load server certificate and key
	serverCert, err := tls.LoadX509KeyPair("cert/server-cert.pem", "cert/server-key.pem")
	if err != nil {
		return nil, fmt.Errorf("failed to load server certificate and key: %w", err)
	}

	// Configure the TLS certificate
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certPool,
	}

	// Create TLS credentials
	tlsCredentials := credentials.NewTLS(tlsConfig)

	// Create a gRPC server with TLS credentials
	grpcServer := grpc.NewServer(grpc.Creds(tlsCredentials))

	// Register the service with the gRPC server
	pb.RegisterVideoServiceServer(grpcServer, service)

	// Log the server's address
	log.Printf("listening at %v", lis.Addr())

	return &Server{
		gs:     grpcServer,
		Lis:    lis,
		Port:   c.Port,
		Client: monitClient,
	}, nil
}

func (s *Server) Start() error {
	if err := s.gs.Serve(s.Lis); err != nil {
		return fmt.Errorf("gRPC server failed to serve: %w", err)
	}
	return nil
}

// func NewgrpcServe(c *config.Config, service pb.VideoServiceServer, monitClient clientinterfaces.MonitClient) (*Server, error) {
// 	altsTC := alts.NewServerCreds(alts.DefaultServerOptions())
// 	grpcserver := grpc.NewServer(grpc.Creds(altsTC))

// 	// Register your gRPC service
// 	pb.RegisterVideoServiceServer(grpcserver, service)

// 	lis, err := net.Listen("tcp", c.Port)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &Server{
// 		gs:     grpcserver,
// 		Lis:    lis,
// 		Port:   c.Port,
// 		Client: monitClient,
// 	}, nil
// }
