// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	context "context"

	model "github.com/jayvib/app/model"
	mock "github.com/stretchr/testify/mock"

	user "github.com/jayvib/app/user"
)

// SearchEngine is an autogenerated mock type for the SearchEngine type
type SearchEngine struct {
	mock.Mock
}

// Delete provides a mock function with given fields: ctx, id
func (_m *SearchEngine) Delete(ctx context.Context, id string) error {
	ret := _m.Called(ctx, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Search provides a mock function with given fields: ctx, input
func (_m *SearchEngine) Search(ctx context.Context, input user.SearchInput) (*user.SearchResult, error) {
	ret := _m.Called(ctx, input)

	var r0 *user.SearchResult
	if rf, ok := ret.Get(0).(func(context.Context, user.SearchInput) *user.SearchResult); ok {
		r0 = rf(ctx, input)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*user.SearchResult)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, user.SearchInput) error); ok {
		r1 = rf(ctx, input)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SearchByName provides a mock function with given fields: ctx, name, num, size
func (_m *SearchEngine) SearchByName(ctx context.Context, name string, num int, size int) ([]*model.User, int, error) {
	ret := _m.Called(ctx, name, num, size)

	var r0 []*model.User
	if rf, ok := ret.Get(0).(func(context.Context, string, int, int) []*model.User); ok {
		r0 = rf(ctx, name, num, size)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.User)
		}
	}

	var r1 int
	if rf, ok := ret.Get(1).(func(context.Context, string, int, int) int); ok {
		r1 = rf(ctx, name, num, size)
	} else {
		r1 = ret.Get(1).(int)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, string, int, int) error); ok {
		r2 = rf(ctx, name, num, size)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// Store provides a mock function with given fields: ctx, _a1
func (_m *SearchEngine) Store(ctx context.Context, _a1 *model.User) error {
	ret := _m.Called(ctx, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.User) error); ok {
		r0 = rf(ctx, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Update provides a mock function with given fields: ctx, _a1
func (_m *SearchEngine) Update(ctx context.Context, _a1 *model.User) error {
	ret := _m.Called(ctx, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.User) error); ok {
		r0 = rf(ctx, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
