// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Code generated by MockGen. DO NOT EDIT.
// Source: logger/common.go

// Package mock_logger is a generated GoMock package.
package mock_logger

import (
	logger "github.com/docker/docker/daemon/logger"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// Mockclient is a mock of client interface
type Mockclient struct {
	ctrl     *gomock.Controller
	recorder *MockclientMockRecorder
}

// MockclientMockRecorder is the mock recorder for Mockclient
type MockclientMockRecorder struct {
	mock *Mockclient
}

// NewMockclient creates a new mock instance
func NewMockclient(ctrl *gomock.Controller) *Mockclient {
	mock := &Mockclient{ctrl: ctrl}
	mock.recorder = &MockclientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *Mockclient) EXPECT() *MockclientMockRecorder {
	return m.recorder
}

// Log mocks base method
func (m *Mockclient) Log(arg0 *logger.Message) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Log", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Log indicates an expected call of Log
func (mr *MockclientMockRecorder) Log(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Log", reflect.TypeOf((*Mockclient)(nil).Log), arg0)
}

// MockLogDriver is a mock of LogDriver interface
type MockLogDriver struct {
	ctrl     *gomock.Controller
	recorder *MockLogDriverMockRecorder
}

// MockLogDriverMockRecorder is the mock recorder for MockLogDriver
type MockLogDriverMockRecorder struct {
	mock *MockLogDriver
}

// NewMockLogDriver creates a new mock instance
func NewMockLogDriver(ctrl *gomock.Controller) *MockLogDriver {
	mock := &MockLogDriver{ctrl: ctrl}
	mock.recorder = &MockLogDriverMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockLogDriver) EXPECT() *MockLogDriverMockRecorder {
	return m.recorder
}

// Start mocks base method
func (m *MockLogDriver) Start(arg0 func() error) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Start", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Start indicates an expected call of Start
func (mr *MockLogDriverMockRecorder) Start(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Start", reflect.TypeOf((*MockLogDriver)(nil).Start), arg0)
}
