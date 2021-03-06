// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	context "context"

	model "github.com/jayvib/app/model"
	mock "github.com/stretchr/testify/mock"

	search "github.com/jayvib/app/internal/app/search"
)

// SearchEngine is an autogenerated mock type for the SearchEngine type
type SearchEngine struct {
	mock.Mock
}

// Search provides a mock function with given fields: ctx, input
func (_m *SearchEngine) Search(ctx context.Context, input search.Input) (*search.Result, error) {
	ret := _m.Called(ctx, input)

	var r0 *search.Result
	if rf, ok := ret.Get(0).(func(context.Context, search.Input) *search.Result); ok {
		r0 = rf(ctx, input)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*search.Result)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, search.Input) error); ok {
		r1 = rf(ctx, input)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Store provides a mock function with given fields: ctx, a
func (_m *SearchEngine) Store(ctx context.Context, a *model.Article) error {
	ret := _m.Called(ctx, a)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.Article) error); ok {
		r0 = rf(ctx, a)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
