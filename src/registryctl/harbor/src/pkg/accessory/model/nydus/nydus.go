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

package nydus

import (
	"github.com/goharbor/harbor/src/pkg/accessory/model"
	"github.com/goharbor/harbor/src/pkg/accessory/model/base"
)

// Nydus accelerator model
type Nydus struct {
	base.Default
}

// Kind gives the reference type of nydus accelerator.
func (ny *Nydus) Kind() string {
	return model.RefHard
}

// IsHard ...
func (ny *Nydus) IsHard() bool {
	return true
}

// New returns nydus accelerator
func New(data model.AccessoryData) model.Accessory {
	return &Nydus{base.Default{
		Data: data,
	}}
}

func init() {
	model.Register(model.TypeNydusAccelerator, New)
}
