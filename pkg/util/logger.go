package util

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
	"fmt"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
)

// ContextKey context Key
type ContextKey struct{}

// SetCtx set ctx
func SetCtx(ctx context.Context, key, value interface{}) context.Context {
	return context.WithValue(ctx, key, value)
}

// LoggerFromContext format
func LoggerFromContext(ctx context.Context) logr.Logger {
	log, ok := ctx.Value(ContextKey{}).(logr.Logger)
	if !ok {
		zapLog, err := zap.NewDevelopment()
		if err != nil {
			panic(fmt.Sprintf("who watches the watchmen (%v)?", err))
		}
		log = zapr.NewLogger(zapLog)
		log.Error(fmt.Errorf("the log processor has not been initialized"), "context")
	}

	return log
}
