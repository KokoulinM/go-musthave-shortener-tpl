package handlers

import (
	"context"
	"reflect"

	"github.com/golang/mock/gomock"

	"github.com/mkokoulin/go-musthave-shortener-tpl/internal/models"
)

// MockRepository is a mock of Repository interface.
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository.
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance.
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// AddMultipleURLs mocks base method.
func (m *MockRepository) AddMultipleURLs(ctx context.Context, user models.UserID, urls ...RequestGetURLs) ([]ResponseGetURLs, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, user}
	for _, a := range urls {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "AddMultipleURLs", varargs...)
	ret0, _ := ret[0].([]ResponseGetURLs)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddMultipleURLs indicates an expected call of AddMultipleURLs.
func (mr *MockRepositoryMockRecorder) AddMultipleURLs(ctx, user interface{}, urls ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, user}, urls...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddMultipleURLs", reflect.TypeOf((*MockRepository)(nil).AddMultipleURLs), varargs...)
}

// AddURL mocks base method.
func (m *MockRepository) AddURL(ctx context.Context, longURL models.LongURL, shortURL models.ShortURL, user models.UserID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddURL", ctx, longURL, shortURL, user)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddURL indicates an expected call of AddURL.
func (mr *MockRepositoryMockRecorder) AddURL(ctx, longURL, shortURL, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddURL", reflect.TypeOf((*MockRepository)(nil).AddURL), ctx, longURL, shortURL, user)
}

// DeleteMultipleURLs mocks base method.
func (m *MockRepository) DeleteMultipleURLs(ctx context.Context, user models.UserID, urls ...string) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, user}
	for _, a := range urls {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DeleteMultipleURLs", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteMultipleURLs indicates an expected call of DeleteMultipleURLs.
func (mr *MockRepositoryMockRecorder) DeleteMultipleURLs(ctx, user interface{}, urls ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, user}, urls...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteMultipleURLs", reflect.TypeOf((*MockRepository)(nil).DeleteMultipleURLs), varargs...)
}

// GetURL mocks base method.
func (m *MockRepository) GetURL(ctx context.Context, shortURL models.ShortURL) (models.ShortURL, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetURL", ctx, shortURL)
	ret0, _ := ret[0].(models.ShortURL)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetURL indicates an expected call of GetURL.
func (mr *MockRepositoryMockRecorder) GetURL(ctx, shortURL interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetURL", reflect.TypeOf((*MockRepository)(nil).GetURL), ctx, shortURL)
}

// GetStates mocks base method.
func (m *MockRepository) GetStates(ctx context.Context) (ResponseStates, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStates", ctx)
	ret0, _ := ret[0].(ResponseStates)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetStates indicates an expected call of GetStates.
func (mr *MockRepositoryMockRecorder) GetStates(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStates", reflect.TypeOf((*MockRepository)(nil).GetStates), ctx)
}

// GetUserURLs mocks base method.
func (m *MockRepository) GetUserURLs(ctx context.Context, user models.UserID) ([]ResponseGetURL, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserURLs", ctx, user)
	ret0, _ := ret[0].([]ResponseGetURL)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserURLs indicates an expected call of GetUserURLs.
func (mr *MockRepositoryMockRecorder) GetUserURLs(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserURLs", reflect.TypeOf((*MockRepository)(nil).GetUserURLs), ctx, user)
}

// Ping mocks base method.
func (m *MockRepository) Ping(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Ping", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// Ping indicates an expected call of Ping.
func (mr *MockRepositoryMockRecorder) Ping(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ping", reflect.TypeOf((*MockRepository)(nil).Ping), ctx)
}
