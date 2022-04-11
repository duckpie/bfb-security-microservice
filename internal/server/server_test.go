package server_test

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/duckpie/bfb-security-microservice/internal/config"
	"github.com/duckpie/bfb-security-microservice/internal/core"
	"github.com/duckpie/bfb-security-microservice/internal/db/redisstore"
	"github.com/duckpie/bfb-security-microservice/internal/server"
	pb "github.com/wrs-news/golang-proto/pkg/proto/security"
	pbu "github.com/wrs-news/golang-proto/pkg/proto/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

var (
	testConfig = config.NewConfig()
	bufSize    = 1024 * 1024
	lis        *bufconn.Listener

	srv      *server.Server
	teardown func(ctx context.Context) error
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	if _, err := toml.DecodeFile(
		fmt.Sprintf("../../config/config.%s.toml", os.Getenv("ENV")),
		testConfig); err != nil {
		log.Fatalf(err.Error())
	}

	r, err := redisstore.NewClient(ctx, &testConfig.Services.Redis)
	if err != nil {
		log.Fatalf(err.Error())
	}

	teardown = func(ctx context.Context) error {
		r.FlushAll(ctx)
		return r.Close()
	}

	lis = bufconn.Listen(bufSize)
	srv := server.InitServer(&testConfig.Services.Server, redisstore.NewRedisStore(r))
	defer srv.GetServer().Stop()

	if err := srv.AddConnection(core.UMS, func() (*grpc.ClientConn, error) {
		return grpc.Dial(
			fmt.Sprintf("%s:%d", testConfig.Microservices.UserMs.Host, testConfig.Microservices.UserMs.Port),
			grpc.WithInsecure(),
		)
	}); err != nil {
		log.Fatalf(err.Error())
	}

	// Создание пользователя
	client, err := srv.GetConn(core.UMS)
	if err != nil {
		log.Fatalf(err.Error())
	}
	defer client.Close()

	conn := pbu.NewUserServiceClient(client)
	testUser, err := conn.CreateUser(ctx, &pbu.NewUserReq{
		Login:    "tester",
		Email:    "tester@gmail.com",
		Password: "12344321",
	})
	if err != nil {
		log.Fatalf(err.Error())
	}

	pb.RegisterSecurityServiceServer(srv.GetServer(), srv)
	go func() {
		if err := srv.GetServer().Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()

	code := m.Run()

	// Удаляю созданного пользователя
	if _, err := conn.DeleteUser(ctx, &pbu.UserReqUuid{Uuid: testUser.Uuid}); err != nil {
		log.Fatalf(err.Error())
	}

	os.Exit(code)

}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}
