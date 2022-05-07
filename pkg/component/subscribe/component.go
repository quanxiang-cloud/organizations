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

// DealRelation deal rlation func
type DealRelation func(context.Context, *event.OrgSpec)

// Subscribe Subscribe
type Subscribe interface {
	Scaffold(context.Context, event.Data, DealRelation) error
}

// Component Component
type Component struct {
	e *gin.RouterGroup

	sub Subscribe
	fn  DealRelation
}

// New new
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

// Option Option
type Option func(*Component)

// WithRouter WithRouter
func WithRouter(group *gin.RouterGroup) Option {
	return func(c *Component) {
		c.e = group
	}
}

// WithFunc WithFunc
func WithFunc(fn DealRelation) Option {
	return func(c *Component) {
		c.fn = fn
	}
}

// Error
var (
	ErrNoFunc = errors.New("no func")
)
