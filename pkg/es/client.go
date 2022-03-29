package es

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
	"errors"
	"fmt"

	"github.com/olivere/elastic/v7"

	"github.com/quanxiang-cloud/cabin/logger"
	es2 "github.com/quanxiang-cloud/cabin/tailormade/db/elastic"
	ginheader "github.com/quanxiang-cloud/cabin/tailormade/header"
	"github.com/quanxiang-cloud/search/pkg/apis/v1alpha1"
)

// Client client
type Client struct {
	esClient *elastic.Client
}

func new(conf *es2.Config, log logger.AdaptedLogger) *Client {
	client, err := es2.NewClient(conf, log)
	if err != nil {
		logger.Logger.Error(err)
		return nil
	}
	return &Client{
		esClient: client,
	}
}

// AddUser add user to es
func (e *Client) AddUser(ctx context.Context, entiy []v1alpha1.User) error {
	if entiy == nil {
		return errors.New("nil data")
	}
	for k := range entiy {
		_, err := e.esClient.Index().Index(v1alpha1.UserIndex).BodyJson(entiy[k]).Do(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

// DelUser del user from es
func (e *Client) DelUser(ctx context.Context, entiy []v1alpha1.User) error {
	if entiy == nil {
		return errors.New("nil data")
	}

	_, tenantID := ginheader.GetTenantID(ctx).Wreck()
	for k := range entiy {
		queries := make([]elastic.Query, 0)
		queries = append(queries, elastic.NewTermQuery("id.keyword", entiy[k].ID))
		if tenantID != "" {
			queries = append(queries, elastic.NewTermQuery("tenantID.keyword", tenantID))
		} else {
			queries = append(queries, elastic.NewExistsQuery("tenantID.keyword"))
		}
		_, err := e.esClient.DeleteByQuery().Index(v1alpha1.UserIndex).Query(elastic.NewBoolQuery().Must(queries...)).Do(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

// AddDepartment add user to es
func (e *Client) AddDepartment(ctx context.Context, entiy []v1alpha1.Department) error {
	if len(entiy) == 0 {
		return errors.New("nil data")
	}
	for k := range entiy {
		_, err := e.esClient.Index().Index(v1alpha1.DepartmentIndex).BodyJson(entiy[k]).Do(ctx)
		if err != nil {
			logger.Logger.Error(err)
			return err
		}
	}

	return nil
}

// DelDepartment del user from es
func (e *Client) DelDepartment(ctx context.Context) error {
	_, tenantID := ginheader.GetTenantID(ctx).Wreck()
	queries := make([]elastic.Query, 0)
	if tenantID != "" {
		queries = append(queries, elastic.NewTermQuery("tenantID.keyword", tenantID))
	} else {
		queries = append(queries, elastic.NewExistsQuery("tenantID.keyword"))
	}
	_, err := e.esClient.DeleteByQuery().Index(v1alpha1.DepartmentIndex).Query(elastic.NewBoolQuery().Must(queries...)).Do(ctx)
	if err != nil {
		return err
	}
	return nil
}

var se *Search

// Search search
type Search struct {
	ctx    context.Context
	client *Client
	user   chan *SearchUser
	dep    chan *SearchDepartment
}

// SearchUser es user
type SearchUser struct {
	User []v1alpha1.User
	Ctx  context.Context
	Sig  chan int
}

// SearchDepartment es department
type SearchDepartment struct {
	Deps []v1alpha1.Department
	Ctx  context.Context
	Sig  chan int
}

// GetSearch get search
func GetSearch() *Search {
	if se == nil {
		return nil
	}
	return se
}

// New new es for es
func New(conf *es2.Config, log logger.AdaptedLogger) {
	se = &Search{
		ctx:    context.Background(),
		client: new(conf, log),
		user:   make(chan *SearchUser),
		dep:    make(chan *SearchDepartment),
	}
	go se.process(se.ctx)
}

// AddUserSearch add data to es
func (s *Search) AddUserSearch(entity *SearchUser) {
	s.user <- entity
}

// AddDepartmentSearch add data to es
func (s *Search) AddDepartmentSearch(entity *SearchDepartment) {
	s.dep <- entity
}

func (s *Search) process(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("done**")
			return
		case entity := <-s.user:
			s.client.DelUser(entity.Ctx, entity.User)
			s.client.AddUser(entity.Ctx, entity.User)
			if entity.Sig != nil {
				entity.Sig <- 1
			}
		case entity := <-s.dep:
			s.client.DelDepartment(entity.Ctx)
			s.client.AddDepartment(entity.Ctx, entity.Deps)
			if entity.Sig != nil {
				entity.Sig <- 1
			}
		}
	}
}
