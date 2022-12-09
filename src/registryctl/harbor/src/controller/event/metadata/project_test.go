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

package metadata

import (
	"testing"

	"github.com/stretchr/testify/suite"

	event2 "github.com/goharbor/harbor/src/controller/event"
	"github.com/goharbor/harbor/src/pkg/notifier/event"
)

type projectEventTestSuite struct {
	suite.Suite
}

func (p *projectEventTestSuite) TestResolveOfCreateProjectEventMetadata() {
	e := &event.Event{}
	metadata := &CreateProjectEventMetadata{
		Project:  "library",
		Operator: "admin",
	}
	err := metadata.Resolve(e)
	p.Require().Nil(err)
	p.Equal(event2.TopicCreateProject, e.Topic)
	p.Require().NotNil(e.Data)
	data, ok := e.Data.(*event2.CreateProjectEvent)
	p.Require().True(ok)
	p.Equal("library", data.Project)
	p.Equal("admin", data.Operator)
}

func (p *projectEventTestSuite) TestResolveOfDeleteProjectEventMetadata() {
	e := &event.Event{}
	metadata := &DeleteProjectEventMetadata{
		Project:  "library",
		Operator: "admin",
	}
	err := metadata.Resolve(e)
	p.Require().Nil(err)
	p.Equal(event2.TopicDeleteProject, e.Topic)
	p.Require().NotNil(e.Data)
	data, ok := e.Data.(*event2.DeleteProjectEvent)
	p.Require().True(ok)
	p.Equal("library", data.Project)
	p.Equal("admin", data.Operator)
}

func TestProjectEventTestSuite(t *testing.T) {
	suite.Run(t, &projectEventTestSuite{})
}
