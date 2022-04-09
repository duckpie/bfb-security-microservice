package cmd

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/BurntSushi/toml"
	"github.com/duckpie/bfb-security-microservice/internal/config"
	"github.com/duckpie/bfb-security-microservice/internal/server"
	"github.com/oklog/oklog/pkg/group"
	"github.com/spf13/cobra"
	pb "github.com/wrs-news/golang-proto/pkg/proto/security"
)

func runCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Run microservice",
		Long:  `...`,
		Run: func(cmd *cobra.Command, args []string) {
			cfg := config.NewConfig()

			// os.Getenv("ENV")
			if _, err := toml.DecodeFile(
				fmt.Sprintf("config/config.%s.toml", "local"), cfg); err != nil {
				log.Printf(err.Error())
				os.Exit(1)
			}

			if err := runner(cfg); err != nil {
				log.Printf(err.Error())
				os.Exit(1)
			}
		},
	}

	return cmd
}

func runner(cfg *config.Config) (err error) {
	defer func() {
		if r := recover(); r != nil {
			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				err = errors.New("unknown panic")
			}
		}
	}()

	srv := server.InitServer(&cfg.Services.Server)
	if err := srv.ConnectToUserService(
		cfg.Microservices.UserMs.Host,
		int(cfg.Microservices.UserMs.Port),
	); err != nil {
		return err
	}

	var g group.Group
	{
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Services.Server.Port))
		if err != nil {
			return err
		}
		log.Printf("server listening a t %v", lis.Addr())

		g.Add(func() error {
			pb.RegisterSecurityServiceServer(srv.GetServer(), srv)
			return srv.GetServer().Serve(lis)
		}, func(error) {
			lis.Close()
		})
	}

	{
		cancelInterrupt := make(chan struct{})
		g.Add(func() error {
			c := make(chan os.Signal, 1)
			signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
			select {
			case sig := <-c:
				return fmt.Errorf("received signal %s", sig)
			case <-cancelInterrupt:
				return nil
			}
		}, func(error) {
			close(cancelInterrupt)
		})
	}

	return g.Run()
}
