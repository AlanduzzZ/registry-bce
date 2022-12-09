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

package rds

import (
	"testing"

	"github.com/gomodule/redigo/redis"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/goharbor/harbor/src/jobservice/tests"
)

// RdsUtilsTestSuite tests functions located in rds package
type RdsUtilsTestSuite struct {
	suite.Suite
	pool      *redis.Pool
	namespace string
	conn      redis.Conn
}

// SetupSuite prepares test suite
func (suite *RdsUtilsTestSuite) SetupSuite() {
	suite.pool = tests.GiveMeRedisPool()
	suite.namespace = tests.GiveMeTestNamespace()
}

// SetupTest prepares test cases
func (suite *RdsUtilsTestSuite) SetupTest() {
	suite.conn = suite.pool.Get()
}

// TearDownTest clears test cases
func (suite *RdsUtilsTestSuite) TearDownTest() {
	err := suite.conn.Close()
	assert.NoError(suite.T(), err, "close conn: nil error expected but got %s", err)
}

// TearDownSuite clears test suite
func (suite *RdsUtilsTestSuite) TearDownSuite() {
	conn := suite.pool.Get()
	defer func() {
		_ = conn.Close()
	}()

	err := tests.ClearAll(suite.namespace, conn)
	assert.NoError(suite.T(), err, "clear all: nil error expected but got %s", err)
}

// TestHmGetAndSet ...
func (suite *RdsUtilsTestSuite) TestHmGetAndSet() {
	key := KeyJobStats(suite.namespace, "fake_job_id")
	err := HmSet(suite.conn, key, "a", "hello", "b", 100)
	require.Nil(suite.T(), err, "nil error should be returned for HmSet")

	values, err := HmGet(suite.conn, key, "a", "b")
	require.Nil(suite.T(), err, "nil error should be returned for HmGet")
	assert.Equal(suite.T(), 2, len(values), "two values should be returned")
	assert.Equal(suite.T(), string(values[0].([]byte)), "hello")
	assert.Equal(suite.T(), string(values[1].([]byte)), "100")
}

// TestAcquireAndReleaseLock ...
func (suite *RdsUtilsTestSuite) TestAcquireAndReleaseLock() {
	key := KeyPeriodicLock(suite.namespace)
	err := AcquireLock(suite.conn, key, "RdsUtilsTestSuite", 60)
	assert.Nil(suite.T(), err, "nil error should be returned for 1st acquiring lock")

	err = AcquireLock(suite.conn, key, "RdsUtilsTestSuite", 60)
	assert.NotNil(suite.T(), err, "non nil error should be returned for 2nd acquiring lock")

	err = ReleaseLock(suite.conn, key, "RdsUtilsTestSuite")
	assert.Nil(suite.T(), err, "nil error should be returned for releasing lock")
}

// TestGetZsetByScore ...
func (suite *RdsUtilsTestSuite) TestGetZsetByScore() {
	key := KeyPeriod(suite.namespace)

	count, err := suite.conn.Do("ZADD", key, 1, "hello", 2, "world")
	require.Nil(suite.T(), err, "nil error should be returned when adding prepared data by ZADD")
	require.Equal(suite.T(), int64(2), count.(int64), "two items should be added")

	datas, err := GetZsetByScore(suite.conn, key, []int64{0, 2})
	require.Nil(suite.T(), err, "nil error should be returned when getting data with scores")
	assert.Equal(suite.T(), 2, len(datas), "expected 2 items but got %d", len(datas))
}

// TestRdsUtilsTestSuite is suite entry for 'go test'
func TestRdsUtilsTestSuite(t *testing.T) {
	suite.Run(t, new(RdsUtilsTestSuite))
}
