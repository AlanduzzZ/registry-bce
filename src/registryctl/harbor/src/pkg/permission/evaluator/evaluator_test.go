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

package evaluator

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/goharbor/harbor/src/pkg/permission/types"
)

type mockEvaluator struct {
	name string
}

func (e *mockEvaluator) HasPermission(ctx context.Context, resource types.Resource, action types.Action) bool {
	return true
}

func TestEvaluatorsAdd(t *testing.T) {
	eva1 := &mockEvaluator{name: "eva1"}
	eva2 := &mockEvaluator{name: "eva2"}
	eva3 := Evaluators{eva1, eva2}

	var es1 Evaluators
	assert.Len(t, es1.Add(eva3), 2)
	assert.Len(t, es1.Add(eva1, eva2, eva3), 2)
}
