package client

import (
	"context"
	"fmt"
	"time"

	pb "github.com/StephenJonesIT/CCMS/src/user-service/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthClient struct {
	conn    *grpc.ClientConn
	client  pb.AuthServiceClient
	timeout time.Duration
}



func NewAuthClient(addr string, timeout time.Duration) (*AuthClient, error) {
	// Set up a connection to the server
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := pb.NewAuthServiceClient(conn)

	return &AuthClient{
		conn:    conn,
		client:  client,
		timeout: timeout,
	}, nil
}

func (c *AuthClient) Close() error {
	return c.conn.Close()
}

// VerifyToken implementation
func (c *AuthClient) VerifyToken(token string) (*pb.TokenResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	return c.client.VerifyToken(ctx, &pb.TokenRequest{
		Token: token,
	})
}

func (c *AuthClient) CheckPermission(token, resource string) (*pb.PermissionResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	response, err := c.client.CheckPermission(ctx, &pb.PermissionRequest{
		Token:    token,
		Resource: resource,
	})

	if err != nil {
		return nil, fmt.Errorf("permission check failed: %w", err)
	}

	if !response.Allowed {
		return response, fmt.Errorf("permission denied for :%s", resource)
	}

	return response, nil
}