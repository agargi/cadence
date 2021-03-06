// Copyright (c) 2017-2020 Uber Technologies Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package resource

import (
	"github.com/uber/cadence/common/resource"
	"github.com/uber/cadence/common/service"
	"github.com/uber/cadence/service/history/config"
	"github.com/uber/cadence/service/history/events"
)

// Resource is the interface which expose common history resources
type Resource interface {
	resource.Resource
	GetEventCache() events.Cache
}

// Impl contains all common resources shared across history
type Impl struct {
	resource.Resource
	eventCache events.Cache
}

// GetEventCache return event cache
func (h *Impl) GetEventCache() events.Cache {
	return h.eventCache
}

// New create a new resource containing common history dependencies
func New(
	params *service.BootstrapParams,
	serviceName string,
	config *config.Config,
	visibilityManagerInitializer resource.VisibilityManagerInitializer,
) (impl *Impl, retError error) {
	serviceResource, err := resource.New(
		params,
		serviceName,
		config.PersistenceMaxQPS,
		config.PersistenceGlobalMaxQPS,
		config.ThrottledLogRPS,
		visibilityManagerInitializer,
	)
	if err != nil {
		return nil, err
	}

	impl = &Impl{
		Resource: serviceResource,
		eventCache: events.NewGlobalCache(
			config.EventsCacheGlobalInitialCount(),
			config.EventsCacheGlobalMaxCount(),
			config.EventsCacheTTL(),
			serviceResource.GetHistoryManager(),
			params.Logger,
			params.MetricsClient,
			uint64(config.EventsCacheMaxSize()),
		),
	}
	return impl, nil
}
