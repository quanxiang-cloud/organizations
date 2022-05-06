package component

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/quanxiang-cloud/cabin/logger"
	redis2 "github.com/quanxiang-cloud/cabin/tailormade/db/redis"
	ginlogger "github.com/quanxiang-cloud/cabin/tailormade/gin"
	"github.com/quanxiang-cloud/organizations/pkg/component/event"
	"github.com/quanxiang-cloud/organizations/pkg/component/publish"
	"github.com/quanxiang-cloud/organizations/pkg/component/subscribe"
	org2 "github.com/quanxiang-cloud/organizations/pkg/component/subscribe/org"
	"github.com/quanxiang-cloud/organizations/pkg/configs"
	"os"
	"os/signal"
	"syscall"
	"testing"
)

// before exec "dapr run --app-port 8001 --app-id org --app-protocol http --dapr-http-port 8002 --components-path  subscribe/samples/"

func TestSub(t *testing.T) {
	conf, err := configs.NewConfig("../../configs/config.yml")
	if err != nil {
		panic(err)
	}
	log := logger.Logger
	ctx := context.Background()
	clusterClient, err := redis2.NewClient(conf.Redis)
	if err != nil {
		panic(err)
	}
	org, err := org2.New(ctx, log, clusterClient)
	if err != nil {
		panic(err)
	}
	gin.SetMode(conf.Model)
	engine := gin.New()
	engine.Use(ginlogger.LoggerFunc(), ginlogger.LoggerFunc())
	group := engine.Group("")
	subscribe.New(ctx, org, subscribe.WithRouter(group), subscribe.WithFunc(fn))
	go engine.Run(":8001")
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			fmt.Println("stop")
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}

func fn(c context.Context, data *event.OrgSpec) {
	fmt.Println("data:", data)
}

//before exec  "dapr run --app-id org-p --app-protocol http --dapr-http-port 8003 --components-path  subscribe/samples/"
// then vim export DAPR_GRPC_PORT= dapr sidecar port
// run

func TestPub(t *testing.T) {
	ctx := context.Background()
	log := logger.Logger
	bus, err := publish.New(ctx, log,
		publish.WithPubsubName("org-redis-pubsub"),
	)
	if err != nil {
		panic(err)
	}
	defer bus.Close()

	data := event.Data{
		OrgSpec: &event.OrgSpec{
			UserID:   "1",
			SourceID: "2",
			Action:   event.ActionAdd,
		},
	}
	message := &publish.Message{}
	message.Data = data
	_, err = bus.Send(ctx, message)
	if err != nil {
		panic(err)
	}
	fmt.Println("send ok")
}
