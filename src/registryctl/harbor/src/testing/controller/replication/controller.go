// Code generated by mockery v2.12.3. DO NOT EDIT.

package replication

import (
	context "context"

	model "github.com/goharbor/harbor/src/controller/replication/model"
	mock "github.com/stretchr/testify/mock"

	q "github.com/goharbor/harbor/src/lib/q"

	regmodel "github.com/goharbor/harbor/src/pkg/reg/model"

	replication "github.com/goharbor/harbor/src/controller/replication"
)

// Controller is an autogenerated mock type for the Controller type
type Controller struct {
	mock.Mock
}

// CreatePolicy provides a mock function with given fields: ctx, policy
func (_m *Controller) CreatePolicy(ctx context.Context, policy *model.Policy) (int64, error) {
	ret := _m.Called(ctx, policy)

	var r0 int64
	if rf, ok := ret.Get(0).(func(context.Context, *model.Policy) int64); ok {
		r0 = rf(ctx, policy)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *model.Policy) error); ok {
		r1 = rf(ctx, policy)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeletePolicy provides a mock function with given fields: ctx, id
func (_m *Controller) DeletePolicy(ctx context.Context, id int64) error {
	ret := _m.Called(ctx, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ExecutionCount provides a mock function with given fields: ctx, query
func (_m *Controller) ExecutionCount(ctx context.Context, query *q.Query) (int64, error) {
	ret := _m.Called(ctx, query)

	var r0 int64
	if rf, ok := ret.Get(0).(func(context.Context, *q.Query) int64); ok {
		r0 = rf(ctx, query)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *q.Query) error); ok {
		r1 = rf(ctx, query)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetExecution provides a mock function with given fields: ctx, executionID
func (_m *Controller) GetExecution(ctx context.Context, executionID int64) (*replication.Execution, error) {
	ret := _m.Called(ctx, executionID)

	var r0 *replication.Execution
	if rf, ok := ret.Get(0).(func(context.Context, int64) *replication.Execution); ok {
		r0 = rf(ctx, executionID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*replication.Execution)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, executionID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetPolicy provides a mock function with given fields: ctx, id
func (_m *Controller) GetPolicy(ctx context.Context, id int64) (*model.Policy, error) {
	ret := _m.Called(ctx, id)

	var r0 *model.Policy
	if rf, ok := ret.Get(0).(func(context.Context, int64) *model.Policy); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Policy)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTask provides a mock function with given fields: ctx, taskID
func (_m *Controller) GetTask(ctx context.Context, taskID int64) (*replication.Task, error) {
	ret := _m.Called(ctx, taskID)

	var r0 *replication.Task
	if rf, ok := ret.Get(0).(func(context.Context, int64) *replication.Task); ok {
		r0 = rf(ctx, taskID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*replication.Task)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, taskID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTaskLog provides a mock function with given fields: ctx, taskID
func (_m *Controller) GetTaskLog(ctx context.Context, taskID int64) ([]byte, error) {
	ret := _m.Called(ctx, taskID)

	var r0 []byte
	if rf, ok := ret.Get(0).(func(context.Context, int64) []byte); ok {
		r0 = rf(ctx, taskID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, taskID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListExecutions provides a mock function with given fields: ctx, query
func (_m *Controller) ListExecutions(ctx context.Context, query *q.Query) ([]*replication.Execution, error) {
	ret := _m.Called(ctx, query)

	var r0 []*replication.Execution
	if rf, ok := ret.Get(0).(func(context.Context, *q.Query) []*replication.Execution); ok {
		r0 = rf(ctx, query)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*replication.Execution)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *q.Query) error); ok {
		r1 = rf(ctx, query)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListPolicies provides a mock function with given fields: ctx, query
func (_m *Controller) ListPolicies(ctx context.Context, query *q.Query) ([]*model.Policy, error) {
	ret := _m.Called(ctx, query)

	var r0 []*model.Policy
	if rf, ok := ret.Get(0).(func(context.Context, *q.Query) []*model.Policy); ok {
		r0 = rf(ctx, query)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Policy)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *q.Query) error); ok {
		r1 = rf(ctx, query)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListTasks provides a mock function with given fields: ctx, query
func (_m *Controller) ListTasks(ctx context.Context, query *q.Query) ([]*replication.Task, error) {
	ret := _m.Called(ctx, query)

	var r0 []*replication.Task
	if rf, ok := ret.Get(0).(func(context.Context, *q.Query) []*replication.Task); ok {
		r0 = rf(ctx, query)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*replication.Task)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *q.Query) error); ok {
		r1 = rf(ctx, query)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PolicyCount provides a mock function with given fields: ctx, query
func (_m *Controller) PolicyCount(ctx context.Context, query *q.Query) (int64, error) {
	ret := _m.Called(ctx, query)

	var r0 int64
	if rf, ok := ret.Get(0).(func(context.Context, *q.Query) int64); ok {
		r0 = rf(ctx, query)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *q.Query) error); ok {
		r1 = rf(ctx, query)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Start provides a mock function with given fields: ctx, policy, resource, trigger
func (_m *Controller) Start(ctx context.Context, policy *model.Policy, resource *regmodel.Resource, trigger string) (int64, error) {
	ret := _m.Called(ctx, policy, resource, trigger)

	var r0 int64
	if rf, ok := ret.Get(0).(func(context.Context, *model.Policy, *regmodel.Resource, string) int64); ok {
		r0 = rf(ctx, policy, resource, trigger)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *model.Policy, *regmodel.Resource, string) error); ok {
		r1 = rf(ctx, policy, resource, trigger)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Stop provides a mock function with given fields: ctx, executionID
func (_m *Controller) Stop(ctx context.Context, executionID int64) error {
	ret := _m.Called(ctx, executionID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) error); ok {
		r0 = rf(ctx, executionID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// TaskCount provides a mock function with given fields: ctx, query
func (_m *Controller) TaskCount(ctx context.Context, query *q.Query) (int64, error) {
	ret := _m.Called(ctx, query)

	var r0 int64
	if rf, ok := ret.Get(0).(func(context.Context, *q.Query) int64); ok {
		r0 = rf(ctx, query)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *q.Query) error); ok {
		r1 = rf(ctx, query)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdatePolicy provides a mock function with given fields: ctx, policy, props
func (_m *Controller) UpdatePolicy(ctx context.Context, policy *model.Policy, props ...string) error {
	_va := make([]interface{}, len(props))
	for _i := range props {
		_va[_i] = props[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, policy)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.Policy, ...string) error); ok {
		r0 = rf(ctx, policy, props...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type NewControllerT interface {
	mock.TestingT
	Cleanup(func())
}

// NewController creates a new instance of Controller. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewController(t NewControllerT) *Controller {
	mock := &Controller{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
