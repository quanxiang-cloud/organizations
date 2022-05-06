package subscribe

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/quanxiang-cloud/organizations/pkg/component/event"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DealRelation func(context.Context, *event.OrgSpec)

type Subscribe interface {
	Scaffold(context.Context, event.Data, DealRelation) error
}

type Component struct {
	e *gin.RouterGroup

	sub Subscribe
	fn  DealRelation
}

func New(ctx context.Context, sender Subscribe, opts ...Option) *Component {
	c := &Component{
		sub: sender,
	}

	for _, opt := range opts {
		opt(c)
	}
	c.init(ctx, sender)
	return c
}

func (c *Component) init(ctx context.Context, sub Subscribe) {
	c.e.POST("/user/relation", func(ctx *gin.Context) {
		body, err := ioutil.ReadAll(ctx.Request.Body)
		if err != nil {
			errHandle(ctx, err)
			return
		}

		event := new(event.DaprEvent)
		err = json.Unmarshal(body, event)
		if err != nil {
			errHandle(ctx, err)
			return
		}

		err = c.sub.Scaffold(ctx, event.Data, c.fn)
		if err != nil {
			errHandle(ctx, err)
			return
		}
	})
}

func errHandle(c *gin.Context, err error) {
	log.Println(err.Error())
	c.JSON(http.StatusOK, nil)
}

type Option func(*Component)

func WithRouter(group *gin.RouterGroup) Option {
	return func(c *Component) {
		c.e = group
	}
}

func WithFunc(fn DealRelation) Option {
	return func(c *Component) {
		c.fn = fn
	}
}

var (
	ErrNoFunc = errors.New("no func")
)
