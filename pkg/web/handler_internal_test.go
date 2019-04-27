package web

import (
	"fmt"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
)

func TestLogError(t *testing.T) {
	err := errors.New("someting went wrong")
	ml := &mockLogger{}
	ml.Test(t)
	ml.On("Log", "level", "error", "message", fmt.Sprintf("%+v", err)).
		Return(nil)
	logError(ml, err)
	ml.AssertExpectations(t)
}

func TestLogError_DoNotLogOnNilError(t *testing.T) {
	ml := &mockLogger{}
	ml.Test(t)
	logError(ml, nil)
	ml.AssertNotCalled(t, "Log", mock.Anything)
}

type mockLogger struct {
	mock.Mock
}

func (ml *mockLogger) Log(kvs ...interface{}) error {
	args := ml.Called(kvs...)
	return args.Error(0)
}
