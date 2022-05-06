package org

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/quanxiang-cloud/cabin/logger"
	"github.com/quanxiang-cloud/organizations/pkg/component/event"
	"github.com/quanxiang-cloud/organizations/pkg/component/subscribe"
)

// New New
func New(ctx context.Context, log logger.AdaptedLogger, redisClient redis.UniversalClient) (*Org, error) {
	return &Org{
		log:         log.WithName("org"),
		redisClient: redisClient,
	}, nil
}

// Org Org
type Org struct {
	log         logger.AdaptedLogger
	redisClient redis.UniversalClient
}

// Scaffold Scaffold
func (e *Org) Scaffold(ctx context.Context, data event.Data, fn subscribe.DealRelation) error {
	if data.OrgSpec == nil {
		return event.ErrDataIsNil
	}
	if fn == nil {
		return subscribe.ErrNoFunc
	}
	fn(ctx, data.OrgSpec)
	return nil
}
