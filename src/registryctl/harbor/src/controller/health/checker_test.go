// Copyright 2019 Project Harbor Authors
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

package health

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/docker/distribution/health"
	"github.com/stretchr/testify/assert"

	"github.com/goharbor/harbor/src/common/utils/test"
)

func TestStringOfHealthy(t *testing.T) {
	var isHealthy healthy = true
	assert.Equal(t, "healthy", isHealthy.String())
	isHealthy = false
	assert.Equal(t, "unhealthy", isHealthy.String())
}

func TestUpdater(t *testing.T) {
	updater := &updater{}
	assert.Equal(t, nil, updater.Check())
	updater.status = errors.New("unhealthy")
	assert.Equal(t, "unhealthy", updater.Check().Error())
}

func TestHTTPStatusCodeHealthChecker(t *testing.T) {
	handler := &test.RequestHandlerMapping{
		Method:  http.MethodGet,
		Pattern: "/health",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		},
	}
	server := test.NewServer(handler)
	defer server.Close()

	url := server.URL + "/health"
	checker := HTTPStatusCodeHealthChecker(
		http.MethodGet, url, map[string][]string{
			"key": {"value"},
		}, 5*time.Second, http.StatusOK)
	assert.Equal(t, nil, checker.Check())

	checker = HTTPStatusCodeHealthChecker(
		http.MethodGet, url, nil, 5*time.Second, http.StatusUnauthorized)
	assert.NotEqual(t, nil, checker.Check())
}

func TestPeriodicHealthChecker(t *testing.T) {
	firstCheck := true
	checkFunc := func() error {
		time.Sleep(2 * time.Second)
		if firstCheck {
			firstCheck = false
			return nil
		}
		return errors.New("unhealthy")
	}

	checker := PeriodicHealthChecker(health.CheckFunc(checkFunc), 1*time.Second)
	assert.Equal(t, "unknown status", checker.Check().Error())
	time.Sleep(3 * time.Second)
	assert.Equal(t, nil, checker.Check())
	time.Sleep(3 * time.Second)
	assert.Equal(t, "unhealthy", checker.Check().Error())
}

func TestCoreHealthChecker(t *testing.T) {
	checker := coreHealthChecker()
	assert.Equal(t, nil, checker.Check())
}
