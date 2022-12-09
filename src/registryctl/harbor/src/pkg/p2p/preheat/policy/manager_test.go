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

package policy

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/goharbor/harbor/src/lib/q"
	"github.com/goharbor/harbor/src/pkg/p2p/preheat/models/policy"
)

type fakeDao struct {
	mock.Mock
}

func (f *fakeDao) Count(ctx context.Context, q *q.Query) (int64, error) {
	args := f.Called()
	return int64(args.Int(0)), args.Error(1)
}
func (f *fakeDao) Create(ctx context.Context, schema *policy.Schema) (int64, error) {
	args := f.Called()
	return int64(args.Int(0)), args.Error(1)
}
func (f *fakeDao) Update(ctx context.Context, schema *policy.Schema, props ...string) error {
	args := f.Called()
	return args.Error(0)
}
func (f *fakeDao) Get(ctx context.Context, id int64) (*policy.Schema, error) {
	args := f.Called()
	var schema *policy.Schema
	if args.Get(0) != nil {
		schema = args.Get(0).(*policy.Schema)
	}
	return schema, args.Error(1)
}
func (f *fakeDao) GetByName(ctx context.Context, projectID int64, name string) (*policy.Schema, error) {
	args := f.Called()
	var schema *policy.Schema
	if args.Get(0) != nil {
		schema = args.Get(0).(*policy.Schema)
	}
	return schema, args.Error(1)
}
func (f *fakeDao) Delete(ctx context.Context, id int64) error {
	args := f.Called()
	return args.Error(0)
}
func (f *fakeDao) List(ctx context.Context, query *q.Query) ([]*policy.Schema, error) {
	args := f.Called()
	var schemas []*policy.Schema
	if args.Get(0) != nil {
		schemas = args.Get(0).([]*policy.Schema)
	}
	return schemas, args.Error(1)
}

type managerTestSuite struct {
	suite.Suite
	mgr Manager
	dao *fakeDao
}

// TestManagerTestSuite tests managerTestSuite
func TestManagerTestSuite(t *testing.T) {
	suite.Run(t, &managerTestSuite{})
}

// SetupSuite setups testing env.
func (m *managerTestSuite) SetupSuite() {
	m.dao = &fakeDao{}
	m.mgr = &manager{dao: m.dao}
}

// TearDownSuite cleans testing env.
func (m *managerTestSuite) TearDownSuite() {
	m.dao = nil
	m.mgr = nil
}

// TestCount tests Count method.
func (m *managerTestSuite) TestCount() {
	m.dao.On("Count").Return(1, nil)
	_, err := m.mgr.Count(nil, nil)
	m.Require().Nil(err)
}

// TestCreate tests Create method.
func (m *managerTestSuite) TestCreate() {
	m.dao.On("Create").Return(1, nil)
	_, err := m.mgr.Create(nil, nil)
	m.Require().Nil(err)
}

// TestUpdate tests Update method.
func (m *managerTestSuite) TestUpdate() {
	m.dao.On("Update").Return(nil)
	err := m.mgr.Update(nil, nil)
	m.Require().Nil(err)
}

// TestGet tests Get method.
func (m *managerTestSuite) TestGet() {
	m.dao.On("Get").Return(&policy.Schema{
		ID:         1,
		Name:       "mgr-policy",
		FiltersStr: `[{"type":"repository","value":"harbor*"},{"type":"tag","value":"2*"}]`,
		TriggerStr: fmt.Sprintf(`{"type":"%s", "trigger_setting":{"cron":"* * * * */1"}}`, policy.TriggerTypeScheduled),
	}, nil)
	_, err := m.mgr.Get(nil, 1)
	m.Require().Nil(err)
}

// TestGetByName tests Get method.
func (m *managerTestSuite) TestGetByName() {
	m.dao.On("GetByName").Return(&policy.Schema{
		ID:         1,
		ProjectID:  1,
		Name:       "mgr-policy",
		FiltersStr: `[{"type":"repository","value":"harbor*"},{"type":"tag","value":"2*"}]`,
		TriggerStr: fmt.Sprintf(`{"type":"%s", "trigger_setting":{"cron":"* * * * */1"}}`, policy.TriggerTypeScheduled),
	}, nil)
	_, err := m.mgr.GetByName(nil, 1, "mgr-policy")
	m.Require().Nil(err)
}

// TestDelete tests Delete method.
func (m *managerTestSuite) TestDelete() {
	m.dao.On("Delete").Return(nil)
	err := m.mgr.Delete(nil, 1)
	m.Require().Nil(err)
}

// TestListPolicies tests ListPolicies method.
func (m *managerTestSuite) TestListPolicies() {
	m.dao.On("List").Return([]*policy.Schema{}, nil)
	_, err := m.mgr.ListPolicies(nil, nil)
	m.Require().Nil(err)
}

// TestListPoliciesByProject tests ListPoliciesByProject method.
func (m *managerTestSuite) TestListPoliciesByProject() {
	m.dao.On("List").Return([]*policy.Schema{}, nil)
	_, err := m.mgr.ListPoliciesByProject(nil, 1, nil)
	m.Require().Nil(err)
}
