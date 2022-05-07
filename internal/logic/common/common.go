package common

import (
	"context"
	"github.com/quanxiang-cloud/cabin/logger"
	"github.com/quanxiang-cloud/organizations/pkg/component/event"
	"github.com/quanxiang-cloud/organizations/pkg/component/publish"
)

// SendToDapr send data by dapr
func SendToDapr(ctx context.Context, bus *publish.Bus, data ...*event.OrgSpec) {
	for k := range data {
		message := new(publish.Message)
		eventData := event.Data{}
		eventData.OrgSpec = data[k]
		message.Data = eventData
		_, err := bus.Send(ctx, message)
		if err != nil {
			logger.Logger.Error(err)
		}
	}

}
