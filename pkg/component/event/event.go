package event

import "errors"

// Action action
type Action int

// Action do some for logic
const (
	ActionAdd Action = iota + 1
	ActionDel
)

// DaprEvent dapr event
type DaprEvent struct {
	Topic           string `json:"topic"`
	Pubsubname      string `json:"pubsubname"`
	Traceid         string `json:"traceid"`
	ID              string `json:"id"`
	Datacontenttype string `json:"datacontenttype"`
	Data            Data   `json:"data"`
	Type            string `json:"type"`
	Specversion     string `json:"specversion"`
	Source          string `json:"source"`
}

// Data data
type Data struct {
	*OrgSpec `json:"org,omitempty"`
}

// OrgSpec org data
type OrgSpec struct {
	UserID   string `json:"userID"`
	SourceID string `json:"sourceID"`
	Action   Action `json:"action"`
}

// Error err
var (
	ErrDataIsNil = errors.New("data is nil")
)
