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

package tencentcr

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	commonhttp "github.com/goharbor/harbor/src/common/http"
	"github.com/goharbor/harbor/src/common/utils/test"
	"github.com/goharbor/harbor/src/pkg/reg/model"
)

func mockChartClient(registry *model.Registry) *adapter {
	return &adapter{
		registry: registry,
		client: commonhttp.NewClient(
			&http.Client{
				Transport: commonhttp.GetHTTPTransport(commonhttp.WithInsecure(registry.Insecure)),
			},
		),
	}
}

func TestFetchCharts(t *testing.T) {
	server := test.NewServer([]*test.RequestHandlerMapping{
		{
			Method:  http.MethodGet,
			Pattern: "/api/chartrepo/library/charts/harbor",
			Handler: func(w http.ResponseWriter, r *http.Request) {
				data := `[{
				"name": "harbor",
				"version":"1.0"
			},{
				"name": "harbor",
				"version":"2.0"
			}]`
				w.Write([]byte(data))
			},
		},
		{
			Method:  http.MethodGet,
			Pattern: "/api/chartrepo/library/charts",
			Handler: func(w http.ResponseWriter, r *http.Request) {
				data := `[{
				"name": "harbor"
			}]`
				w.Write([]byte(data))
			},
		},
	}...)
	defer server.Close()
	var adapter = mockChartClient(&model.Registry{URL: server.URL})

	// nil filter
	resources, err := adapter.fetchCharts([]string{"library"}, nil)
	require.Nil(t, err)
	assert.Equal(t, 2, len(resources))
	assert.Equal(t, model.ResourceTypeChart, resources[0].Type)
	assert.Equal(t, "library/harbor", resources[0].Metadata.Repository.Name)
	assert.Equal(t, 1, len(resources[0].Metadata.Artifacts))
	assert.Equal(t, "1.0", resources[0].Metadata.Artifacts[0].Tags[0])
	// not nil filter
	filters := []*model.Filter{
		{
			Type:  model.FilterTypeName,
			Value: "library/*",
		},
		{
			Type:  model.FilterTypeTag,
			Value: "1.0",
		},
	}
	resources, err = adapter.fetchCharts([]string{"library"}, filters)
	require.Nil(t, err)
	require.Equal(t, 1, len(resources))
	assert.Equal(t, model.ResourceTypeChart, resources[0].Type)
	assert.Equal(t, "library/harbor", resources[0].Metadata.Repository.Name)
	assert.Equal(t, 1, len(resources[0].Metadata.Artifacts))
	assert.Equal(t, "1.0", resources[0].Metadata.Artifacts[0].Tags[0])
}

func TestChartExist(t *testing.T) {
	server := test.NewServer(&test.RequestHandlerMapping{
		Method:  http.MethodGet,
		Pattern: "/api/chartrepo/library/charts/harbor/1.0",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			data := `{
				"metadata": {
					"urls":["http://127.0.0.1/charts"]
				}
			}`
			w.Write([]byte(data))
		},
	})
	defer server.Close()
	var adapter = mockChartClient(&model.Registry{URL: server.URL})
	var exist, err = adapter.ChartExist("library/harbor", "1.0")
	require.Nil(t, err)
	require.True(t, exist)
}

func TestDownloadChart(t *testing.T) {
	server := test.NewServer([]*test.RequestHandlerMapping{
		{
			Method:  http.MethodGet,
			Pattern: "/api/chartrepo/library/charts/harbor/1.0",
			Handler: func(w http.ResponseWriter, r *http.Request) {
				data := `{
				"metadata": {
					"urls":["charts/harbor-1.0.tgz"]
				}
			}`
				w.Write([]byte(data))
			},
		},
		{
			Method:  http.MethodGet,
			Pattern: "/chartrepo/library/charts/harbor-1.0.tgz",
			Handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			},
		},
	}...)
	defer server.Close()
	var adapter = mockChartClient(&model.Registry{URL: server.URL})
	var _, err = adapter.DownloadChart("library/harbor", "1.0", "")
	require.Nil(t, err)
}

func TestUploadChart(t *testing.T) {
	server := test.NewServer(&test.RequestHandlerMapping{
		Method:  http.MethodPost,
		Pattern: "/api/chartrepo/library/charts",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		},
	})
	defer server.Close()
	var adapter = mockChartClient(&model.Registry{URL: server.URL})
	var err = adapter.UploadChart("library/harbor", "1.0", bytes.NewBuffer(nil))
	require.Nil(t, err)
}

func TestDeleteChart(t *testing.T) {
	server := test.NewServer(&test.RequestHandlerMapping{
		Method:  http.MethodDelete,
		Pattern: "/api/chartrepo/library/charts/harbor/1.0",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		},
	})
	defer server.Close()
	var adapter = mockChartClient(&model.Registry{URL: server.URL})
	var err = adapter.DeleteChart("library/harbor", "1.0")
	require.Nil(t, err)
}
