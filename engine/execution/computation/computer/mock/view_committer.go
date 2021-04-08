// Code generated by mockery v1.0.0. DO NOT EDIT.

package mock

import (
	state "github.com/onflow/flow-go/fvm/state"
	mock "github.com/stretchr/testify/mock"
)

// ViewCommitter is an autogenerated mock type for the ViewCommitter type
type ViewCommitter struct {
	mock.Mock
}

// CommitView provides a mock function with given fields: _a0, _a1
func (_m *ViewCommitter) CommitView(_a0 state.View, _a1 []byte) ([]byte, []byte, error) {
	ret := _m.Called(_a0, _a1)

	var r0 []byte
	if rf, ok := ret.Get(0).(func(state.View, []byte) []byte); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 []byte
	if rf, ok := ret.Get(1).(func(state.View, []byte) []byte); ok {
		r1 = rf(_a0, _a1)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).([]byte)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(state.View, []byte) error); ok {
		r2 = rf(_a0, _a1)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}
