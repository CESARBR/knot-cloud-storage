package mocks

import "github.com/stretchr/testify/mock"

// FakeLogger represents a mocking type for the logging service
type FakeLogger struct {
	mock.Mock
}

// Info provides a mock function for the debugging Info capability
func (fl *FakeLogger) Info(...interface{}) {}

// Infof provides a mock function for the debugging Info capability
func (fl *FakeLogger) Infof(string, ...interface{}) {}

// Debug provides a mock function for the debugging Info capability
func (fl *FakeLogger) Debug(...interface{}) {}

// Warn provides a mock function for the debugging Info capability
func (fl *FakeLogger) Warn(...interface{}) {}

// Error provides a mock function for the debugging Info capability
func (fl *FakeLogger) Error(...interface{}) {}

// Errorf provides a mock function for the debugging Info capability
func (fl *FakeLogger) Errorf(string, ...interface{}) {}
