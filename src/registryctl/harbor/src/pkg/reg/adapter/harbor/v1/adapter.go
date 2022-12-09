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

package v1

import (
	"fmt"

	"github.com/goharbor/harbor/src/common/utils"
	adp "github.com/goharbor/harbor/src/pkg/reg/adapter"
	"github.com/goharbor/harbor/src/pkg/reg/adapter/harbor/base"
	"github.com/goharbor/harbor/src/pkg/reg/filter"
	"github.com/goharbor/harbor/src/pkg/reg/model"
)

var _ adp.Adapter = &adapter{}
var _ adp.ArtifactRegistry = &adapter{}
var _ adp.ChartRegistry = &adapter{}

// New creates a Adapter for Harbor 1.x
func New(base *base.Adapter) adp.Adapter {
	return &adapter{
		Adapter: base,
		client:  &client{Client: base.Client},
	}
}

type adapter struct {
	*base.Adapter
	client *client
}

func (a *adapter) FetchArtifacts(filters []*model.Filter) ([]*model.Resource, error) {
	projects, err := a.ListProjects(filters)
	if err != nil {
		return nil, err
	}

	var resources []*model.Resource
	for _, project := range projects {
		repositories, err := a.listRepositories(project, filters)
		if err != nil {
			return nil, err
		}
		if len(repositories) == 0 {
			continue
		}

		var rawResources = make([]*model.Resource, len(repositories))
		runner := utils.NewLimitedConcurrentRunner(adp.MaxConcurrency)

		for i, r := range repositories {
			index := i
			repo := r
			runner.AddTask(func() error {
				artifacts, err := a.listArtifacts(repo.Name, filters)
				if err != nil {
					return fmt.Errorf("failed to list artifacts of repository '%s': %v", repo.Name, err)
				}
				if len(artifacts) == 0 {
					rawResources[index] = nil
					return nil
				}

				rawResources[index] = &model.Resource{
					Type:     model.ResourceTypeImage,
					Registry: a.Registry,
					Metadata: &model.ResourceMetadata{
						Repository: &model.Repository{
							Name:     repo.Name,
							Metadata: project.Metadata,
						},
						Artifacts: artifacts,
					},
				}
				return nil
			})
		}
		if err = runner.Wait(); err != nil {
			return nil, fmt.Errorf("failed to fetch artifacts: %v", err)
		}
		for _, r := range rawResources {
			if r != nil {
				resources = append(resources, r)
			}
		}
	}

	return resources, nil
}

// override the default implementation by calling Harbor API directly
func (a *adapter) DeleteManifest(repository, reference string) error {
	return a.client.deleteManifest(repository, reference)
}

func (a *adapter) listRepositories(project *base.Project, filters []*model.Filter) ([]*model.Repository, error) {
	repositories, err := a.client.listRepositories(project)
	if err != nil {
		return nil, err
	}
	return filter.DoFilterRepositories(repositories, filters)
}

func (a *adapter) listArtifacts(repository string, filters []*model.Filter) ([]*model.Artifact, error) {
	artifacts, err := a.client.listArtifacts(repository)
	if err != nil {
		return nil, err
	}
	return filter.DoFilterArtifacts(artifacts, filters)
}
