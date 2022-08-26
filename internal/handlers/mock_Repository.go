package handlers

import (
	"context"
	"net"
	"reflect"

	"github.com/golang/mock/gomock"

	"github.com/mkokoulin/go-musthave-shortener-tpl/internal/models"
)

// MockURLServiceInterface is a mock of Repository interface.
type MockURLServiceInterface struct {
	ctrl     *gomock.Controller
	recorder *MockURLServiceInterfaceMockRecorder
}

// MockURLServiceInterfaceMockRecorder is the mock recorder for MockURLServiceInterface.
type MockURLServiceInterfaceMockRecorder struct {
	mock *MockURLServiceInterface
}

// NewMockURLServiceInterface creates a new mock instance.
func NewMockURLServiceInterface(ctrl *gomock.Controller) *MockURLServiceInterface {
	mock := &MockURLServiceInterface{ctrl: ctrl}
	mock.recorder = &MockURLServiceInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockURLServiceInterface) EXPECT() *MockURLServiceInterfaceMockRecorder {
	return m.recorder
}

// CreateBatch mocks base method.
func (m *MockURLServiceInterface) CreateBatch(ctx context.Context, urls []RequestGetURLs, user models.UserID) ([]ResponseGetURLs, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, user}
	for _, a := range urls {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "CreateBatch", varargs...)
	ret0, _ := ret[0].([]ResponseGetURLs)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateBatch indicates an expected call of AddMultipleURLs.
func (mr *MockURLServiceInterfaceMockRecorder) CreateBatch(ctx interface{}, urls []RequestGetURLs, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, urls, user})
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateBatch", reflect.TypeOf((*MockURLServiceInterface)(nil).CreateBatch), varargs...)
}

// CreateURL mocks base method.
func (m *MockURLServiceInterface) CreateURL(ctx context.Context, longURL models.LongURL, user models.UserID) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateURL", ctx, longURL, user)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateURL indicates an expected call of AddURL.
func (mr *MockURLServiceInterfaceMockRecorder) CreateURL(ctx, longURL, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateURL", reflect.TypeOf((*MockURLServiceInterface)(nil).CreateURL), ctx, longURL, user)
}

// DeleteBatch mocks base method.
func (m *MockURLServiceInterface) DeleteBatch(urls []string, user models.UserID) {
	m.ctrl.T.Helper()
	varargs := []interface{}{user}
	for _, a := range urls {
		varargs = append(varargs, a)
	}
}

// DeleteBatch indicates an expected call of DeleteMultipleURLs.
func (mr *MockURLServiceInterfaceMockRecorder) DeleteBatch(ctx, user interface{}, urls ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, user}, urls...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteBatch", reflect.TypeOf((*MockURLServiceInterface)(nil).DeleteBatch), varargs...)
}

// GetURL mocks base method.
func (m *MockURLServiceInterface) GetURL(ctx context.Context, shortURL models.ShortURL) (models.ShortURL, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetURL", ctx, shortURL)
	ret0, _ := ret[0].(models.ShortURL)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetURL indicates an expected call of GetURL.
func (mr *MockURLServiceInterfaceMockRecorder) GetURL(ctx, shortURL interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetURL", reflect.TypeOf((*MockURLServiceInterface)(nil).GetURL), ctx, shortURL)
}

//GetStates(ctx context.Context, ip net.IP) (bool, ResponseStates, error)

// GetStates mocks base method.
func (m *MockURLServiceInterface) GetStates(ctx context.Context, ip net.IP) (bool, ResponseStates, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStates", ctx, ip)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(ResponseStates)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetStates indicates an expected call of GetStates.
func (mr *MockURLServiceInterfaceMockRecorder) GetStates(ctx interface{}, ip net.IP) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStates", reflect.TypeOf((*MockURLServiceInterface)(nil).GetStates), ctx, ip)
}

// GetUserURLs mocks base method.
func (m *MockURLServiceInterface) GetUserURLs(ctx context.Context, user models.UserID) ([]ResponseGetURL, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserURLs", ctx, user)
	ret0, _ := ret[0].([]ResponseGetURL)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserURLs indicates an expected call of GetUserURLs.
func (mr *MockURLServiceInterfaceMockRecorder) GetUserURLs(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserURLs", reflect.TypeOf((*MockURLServiceInterface)(nil).GetUserURLs), ctx, user)
}

// Ping mocks base method.
func (m *MockURLServiceInterface) Ping(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Ping", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// Ping indicates an expected call of Ping.
func (mr *MockURLServiceInterfaceMockRecorder) Ping(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ping", reflect.TypeOf((*MockURLServiceInterface)(nil).Ping), ctx)
}
