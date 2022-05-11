package core

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"

	error2 "github.com/quanxiang-cloud/cabin/error"
	"github.com/quanxiang-cloud/cabin/tailormade/client"
	oct "github.com/quanxiang-cloud/organizations/internal/models/octopus"
	mysql3 "github.com/quanxiang-cloud/organizations/internal/models/octopus/mysql"
	"github.com/quanxiang-cloud/organizations/internal/models/org"
	mysql2 "github.com/quanxiang-cloud/organizations/internal/models/org/mysql"
	"github.com/quanxiang-cloud/organizations/pkg/configs"
)

// Core publish
type Core interface {
}

type core struct {
	DB               *gorm.DB
	manageColumnRepo oct.ManageColumn
	useColumnsRepo   oct.UseColumnsRepo
	tableColumnsRepo oct.UserTableColumnsRepo
	redisClient      redis.UniversalClient
	userRepo         org.UserRepo
	conf             configs.Config
	client           http.Client
	columnRepo       oct.UserTableColumnsRepo
	extend           oct.ExtendRepo
}

// NewCore new
func NewCore(conf configs.Config, db *gorm.DB) Core {
	return &core{
		DB:               db,
		manageColumnRepo: mysql3.NewManageColumnRepo(),
		useColumnsRepo:   mysql3.NewUseColumnsRepo(),
		tableColumnsRepo: mysql3.NewUserTableColumnsRepo(),
		userRepo:         mysql2.NewUserRepo(),
		conf:             conf,
		client:           client.New(conf.InternalNet),
		columnRepo:       mysql3.NewUserTableColumnsRepo(),
		extend:           mysql3.NewExtendRepo(),
	}
}

// INResponse will insert or update data
type INResponse struct {
	ID string `json:"id"`
}

// DealRequest deal request
func DealRequest(c http.Client, host string, r *http.Request, data interface{}) (*http.Response, error) {
	request := r.Clone(r.Context())
	parse, _ := url.ParseRequestURI(host)
	request.URL = parse
	request.Host = parse.Host
	request.URL.Path = r.URL.Path
	request.RequestURI = ""

	request.URL.RawQuery = r.URL.RawQuery
	if r.Method != "GET" {
		marshal, _ := json.Marshal(data)
		l := len(marshal)
		itoa := strconv.Itoa(l)
		request.Header.Set("Content-Length", itoa)
		request.ContentLength = int64(l)
		request.Body = io.NopCloser(bytes.NewReader(marshal))
	}

	return c.Do(request)
}

// DealResponse deal response
func DealResponse(w http.ResponseWriter, response *http.Response) {

	defer response.Body.Close()
	all, err := io.ReadAll(response.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	header := response.Header.Clone()
	for k := range header {
		for k1 := range header[k] {
			w.Header().Add(k, header[k][k1])
		}
	}
	w.WriteHeader(response.StatusCode)
	w.Write(all)
	return
}

type respData interface{}

// R response data
type R struct {
	err  error
	Code int64    `json:"code"`
	Msg  string   `json:"msg,omitempty"`
	Data respData `json:"data"`
}

// DeserializationResp marshal response body to struct
func DeserializationResp(ctx context.Context, response *http.Response, entity interface{}) (*R, error) {
	if response.StatusCode != http.StatusOK {
		return nil, error2.New(error2.Internal)
	}
	r := new(R)
	r.Data = entity
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, r)
	if err != nil {
		return nil, err
	}
	response.Body = io.NopCloser(bytes.NewReader(body))
	return r, nil
}
