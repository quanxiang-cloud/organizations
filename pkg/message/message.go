package message

/*
Copyright 2022 QuanxiangCloud Authors
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
     http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
import (
	"context"
	"net/http"

	"github.com/quanxiang-cloud/cabin/tailormade/client"
)

const (
	host = "http://message/api/v1/message"

	sendMessageURI = "/manager/create/batch"
)

// Message 消息服务提供
type Message interface {
	SendMessage(ctx context.Context, req []*CreateReq) error
}

type message struct {
	client http.Client
}

// NewMessage new
func NewMessage(conf client.Config) Message {
	return &message{
		client: client.New(conf),
	}
}

// CreateReq CreateReq
type CreateReq struct {
	data `json:",omitempty"`
}

type data struct {
	Letter *Letter `json:"letter"`
	Email  *Email  `json:"email"`
	Phone  *Phone  `json:"phone"`
}

// Phone Phone
type Phone struct {
}

// Letter Letter
type Letter struct {
	ID      string   `json:"id,omitempty"`
	UUID    []string `json:"uuid,omitempty"`
	Content *Content `json:"contents"`
}

// Email Email
type Email struct {
	To          []string `json:"To"`
	Title       string   `json:"title"`
	Content     *Content `json:"contents"`
	ContentType string   `json:"content_type,omitempty"`
}

// Content Content
type Content struct {
	Content     string            `json:"content"`
	TemplateID  string            `json:"templateID"`
	KeyAndValue map[string]string `json:"keyAndValue"`
}

// Resp 未知返回体
type Resp interface {
}

func (u *message) SendMessage(ctx context.Context, req []*CreateReq) error {
	var resp = new(Resp)
	err := client.POST(ctx, &u.client, host+sendMessageURI, req, resp)
	return err
}
