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

package hook

import (
	"context"
	"testing"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/goharbor/harbor/src/jobservice/common/list"
	"github.com/goharbor/harbor/src/jobservice/common/utils"
	"github.com/goharbor/harbor/src/jobservice/job"
	"github.com/goharbor/harbor/src/jobservice/tests"
	"github.com/goharbor/harbor/src/lib/errors"
)

// HookAgentTestSuite tests functions of hook agent
type HookAgentTestSuite struct {
	suite.Suite

	namespace string
	pool      *redis.Pool
	agent     *basicAgent

	event *Event
	jid   string
}

// TestHookAgentTestSuite is entry of go test
func TestHookAgentTestSuite(t *testing.T) {
	suite.Run(t, new(HookAgentTestSuite))
}

// SetupSuite prepares test suites
func (suite *HookAgentTestSuite) SetupSuite() {
	suite.pool = tests.GiveMeRedisPool()
	suite.namespace = tests.GiveMeTestNamespace()

	suite.agent = &basicAgent{
		context:   context.TODO(),
		namespace: suite.namespace,
		redisPool: suite.pool,
	}

	suite.prepareData()
}

// TearDownSuite prepares test suites
func (suite *HookAgentTestSuite) TearDownSuite() {
	conn := suite.pool.Get()
	defer func() {
		_ = conn.Close()
	}()

	_ = tests.ClearAll(suite.namespace, conn)
}

func (suite *HookAgentTestSuite) prepareData() {
	suite.jid = utils.MakeIdentifier()
	rev := time.Now().Unix()
	stats := &job.Stats{
		Info: &job.StatsInfo{
			JobID:    suite.jid,
			Status:   job.RunningStatus.String(),
			Revision: rev,
			JobKind:  job.KindGeneric,
			JobName:  job.SampleJob,
		},
	}
	t := job.NewBasicTrackerWithStats(context.TODO(), stats, suite.namespace, suite.pool, nil, list.New())
	err := t.Save()
	suite.NoError(err, "mock job stats")

	suite.event = &Event{
		URL:       "http://domian.com",
		Message:   "HookAgentTestSuite",
		Timestamp: time.Now().Unix(),
		Data: &job.StatusChange{
			JobID:  suite.jid,
			Status: job.SuccessStatus.String(),
			Metadata: &job.StatsInfo{
				JobID:    suite.jid,
				Status:   job.SuccessStatus.String(),
				Revision: rev,
				JobKind:  job.KindGeneric,
				JobName:  job.SampleJob,
			},
		},
	}
}

// TestEventSending ...
func (suite *HookAgentTestSuite) TestEventSending() {
	mc := &mockClient{}
	mc.On("SendEvent", suite.event).Return(nil)
	suite.agent.client = mc

	err := suite.agent.Trigger(suite.event)
	require.Nil(suite.T(), err, "agent trigger: nil error expected but got %s", err)

	// check
	suite.checkStatus()
}

// TestEventSending ...
func (suite *HookAgentTestSuite) TestEventSendingError() {
	mc := &mockClient{}
	mc.On("SendEvent", suite.event).Return(errors.New("internal server error: for testing"))
	suite.agent.client = mc

	err := suite.agent.Trigger(suite.event)

	suite.Error(err)
}

func (suite *HookAgentTestSuite) checkStatus() {
	t := job.NewBasicTrackerWithID(context.TODO(), suite.jid, suite.namespace, suite.pool, nil, list.New())
	err := t.Load()
	require.NoError(suite.T(), err, "load updated job stats")
	require.NotNil(suite.T(), t.Job(), "latest job stats")
	suite.Equal(job.SuccessStatus.String(), t.Job().Info.HookAck.Status, "ack status")
}

type mockClient struct {
	mock.Mock
}

func (mc *mockClient) SendEvent(evt *Event) error {
	args := mc.Called(evt)
	return args.Error(0)
}
