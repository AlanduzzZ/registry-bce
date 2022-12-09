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

package apiversion

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/goharbor/harbor/src/lib"
)

func TestMiddleware(t *testing.T) {
	version := ""
	middleware := Middleware("1.0")
	req := httptest.NewRequest("GET", "http://localhost", nil)
	handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		version = lib.GetAPIVersion(req.Context())
	})
	middleware(handler).ServeHTTP(nil, req)
	assert.Equal(t, "1.0", version)
}
