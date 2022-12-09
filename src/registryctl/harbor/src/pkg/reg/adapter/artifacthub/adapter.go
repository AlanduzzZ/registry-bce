// Copyright Project Harbor Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package artifacthub

import (
	"errors"

	"github.com/goharbor/harbor/src/lib/log"
	adp "github.com/goharbor/harbor/src/pkg/reg/adapter"
	"github.com/goharbor/harbor/src/pkg/reg/model"
)

func init() {
	if err := adp.RegisterFactory(model.RegistryTypeArtifactHub, new(factory)); err != nil {
		log.Errorf("failed to register factory for %s: %v", model.RegistryTypeArtifactHub, err)
		return
	}
	log.Infof("the factory for adapter %s registered", model.RegistryTypeArtifactHub)
}

type factory struct {
}

// Create ...
func (f *factory) Create(r *model.Registry) (adp.Adapter, error) {
	return newAdapter(r)
}

// AdapterPattern ...
func (f *factory) AdapterPattern() *model.AdapterPattern {
	return &model.AdapterPattern{
		EndpointPattern: &model.EndpointPattern{
			EndpointType: model.EndpointPatternTypeFix,
			Endpoints: []*model.Endpoint{
				{
					Key:   "artifacthub.io",
					Value: "https://artifacthub.io",
				},
			},
		},
	}
}

var (
	_ adp.Adapter       = (*adapter)(nil)
	_ adp.ChartRegistry = (*adapter)(nil)
)

type adapter struct {
	registry *model.Registry
	client   *Client
}

func newAdapter(registry *model.Registry) (*adapter, error) {
	return &adapter{
		registry: registry,
		client:   newClient(registry),
	}, nil
}

func (a *adapter) Info() (*model.RegistryInfo, error) {
	return &model.RegistryInfo{
		Type: model.RegistryTypeArtifactHub,
		SupportedResourceTypes: []string{
			model.ResourceTypeChart,
		},
		SupportedResourceFilters: []*model.FilterStyle{
			{
				Type:  model.FilterTypeName,
				Style: model.FilterStyleTypeText,
			},
			{
				Type:  model.FilterTypeTag,
				Style: model.FilterStyleTypeText,
			},
		},
		SupportedTriggers: []string{
			model.TriggerTypeManual,
			model.TriggerTypeScheduled,
		},
	}, nil
}

func (a *adapter) PrepareForPush(resources []*model.Resource) error {
	return errors.New("not supported")
}

// HealthCheck checks health status of a registry
func (a *adapter) HealthCheck() (string, error) {
	err := a.client.checkHealthy()
	if err == nil {
		return model.Healthy, nil
	}
	return model.Unhealthy, err
}
