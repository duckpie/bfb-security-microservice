package server_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	pb "github.com/wrs-news/golang-proto/pkg/proto/security"
	"google.golang.org/grpc"
)

func Test_Server_Handlers(t *testing.T) {
	ctx := context.TODO()
	defer teardown(ctx)

	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	client := pb.NewSecurityServiceClient(conn)

	var tokens *pb.TokensPair

	t.Run("login", func(t *testing.T) {
		tokens, err = client.Login(ctx, &pb.LoginReq{
			Login:    "tester",
			Password: "12344321",
		})

		assert.NoError(t, err)
		assert.NotNil(t, tokens)
	})

	t.Run("auth_check", func(t *testing.T) {
		resp, err := client.AuthCheck(ctx,
			&pb.AuthCheckReq{AccessToken: tokens.AccessToken})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
	})

	t.Run("refresh_token", func(t *testing.T) {
		resp, err := client.RefreshToken(ctx,
			&pb.RefreshTokenReq{Token: tokens.RefreshToken})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
	})

	t.Run("logout", func(t *testing.T) {
		resp, err := client.Logout(ctx,
			&pb.LogoutReq{Token: tokens.AccessToken})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
	})

}
