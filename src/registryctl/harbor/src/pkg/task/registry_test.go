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

package task

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegisterCheckInProcessor(t *testing.T) {
	err := RegisterCheckInProcessor("test", nil)
	assert.Nil(t, err)

	// already exist
	err = RegisterCheckInProcessor("test", nil)
	assert.NotNil(t, err)
}

func TestRegisterTaskStatusChangePostFunc(t *testing.T) {
	err := RegisterTaskStatusChangePostFunc("test", nil)
	assert.Nil(t, err)

	// already exist
	err = RegisterTaskStatusChangePostFunc("test", nil)
	assert.NotNil(t, err)
}

func TestRegisterExecutionStatusChangePostFunc(t *testing.T) {
	err := RegisterExecutionStatusChangePostFunc("test", nil)
	assert.Nil(t, err)

	// already exist
	err = RegisterExecutionStatusChangePostFunc("test", nil)
	assert.NotNil(t, err)
}
