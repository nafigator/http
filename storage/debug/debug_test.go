package debug

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

const (
	msg                = "storage_test msg"
	unexpectedMsgCount = "Unexpected messages count"
	unexpectedMsg      = "Unexpected messages"
)

func TestDebugFlush(t *testing.T) {
	a := assert.New(t)
	ob, logs := observer.New(zap.DebugLevel)
	logger := zap.New(ob).Sugar()

	d := New(logger)
	d.Flush(context.Background(), msg)

	a.Len(logs.All(), 1, unexpectedMsgCount)

	expected := []observer.LoggedEntry{{
		Entry:   zapcore.Entry{Level: zap.DebugLevel, Message: msg},
		Context: []zapcore.Field{},
	}}
	actual := logs.AllUntimed()

	a.Equal(expected, actual, unexpectedMsg)
}
