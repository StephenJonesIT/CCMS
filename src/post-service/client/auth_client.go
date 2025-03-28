package client

import (
	"google.golang.org/grpc"
)

type AuthClient struct {
	conn *grpc.ClientConn
	client 
}