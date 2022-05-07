package publish

import (
	"context"
	"fmt"
	"github.com/quanxiang-cloud/cabin/logger"

	daprd "github.com/dapr/go-sdk/client"
	"github.com/quanxiang-cloud/organizations/pkg/component/event"
)

//go:generate stringer -type Channel
// Channel channel
type Channel int

// channel
const (
	None Channel = iota
	Org
)

// Message message
type Message struct {
	event.Data `json:",omitempty"`
}

// SendResp send resp
type SendResp struct{}

// Bus bus
type Bus struct {
	daprClient daprd.Client
	log        logger.AdaptedLogger

	pubsubName string
	tenant     string
}

// New new
func New(ctx context.Context, log logger.AdaptedLogger, opts ...Option) (*Bus, error) {
	client, err := daprd.NewClient()
	if err != nil {
		return nil, err
	}
	bus := &Bus{
		daprClient: client,
		log:        log.WithName("bus"),
	}

	for _, fn := range opts {
		fn(bus)
	}
	return bus, nil
}

// Option option
type Option func(*Bus) error

// WithPubsubName WithPubsubName
func WithPubsubName(pubsubName string) Option {
	return func(b *Bus) error {
		b.pubsubName = pubsubName
		return nil
	}
}

// WithTenant WithTenant
func WithTenant(tenant string) Option {
	return func(b *Bus) error {
		b.tenant = tenant
		return nil
	}
}

// Send send
func (b *Bus) Send(ctx context.Context, req *Message) (*SendResp, error) {
	var topic string

	if req.Data.OrgSpec != nil {
		if b.tenant == "" {
			b.tenant = "default"
		}
		topic = fmt.Sprintf("%s.%s", b.tenant, Org.String())
		if err := b.publish(ctx, topic, req.Data); err != nil {
			b.log.Error(err, "push org", "userID", req.Data)
			return &SendResp{}, err
		}
	}

	b.log.Info("publish success")
	return &SendResp{}, nil
}

func (b *Bus) publish(ctx context.Context, topic string, data interface{}) error {
	b.log.Info("send org", " topic ", topic)
	if err := b.daprClient.PublishEvent(context.Background(), b.pubsubName, topic, data); err != nil {
		b.log.Error(err, "publishEvent", "topic", topic, "pubsubName", b.pubsubName)
		return err
	}
	return nil
}

// Close Close
func (b *Bus) Close() error {
	b.daprClient.Close()
	return nil
}
