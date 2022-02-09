package clog

import (
	"context"
	"runtime"

	"github.com/sirupsen/logrus"
)

// HandlePanic structurally logs a panic before optionally propagating.
//
// Propagating a panic can be important for cases like calls made via net/http where the whole process isn't required
// to fail because one request panics. Propagating makes sure we dont disturb upstream panic handling.
func HandlePanic(ctx context.Context, propagate bool) {
	r := recover()
	if r == nil {
		return
	}

	st := make([]byte, 1<<16) // create a 2 byte stack trace buffer
	st = st[:runtime.Stack(st, false)]

	var logger *logrus.Entry
	ctxLogger := getContextLogger(ctx)
	if ctxLogger != nil {
		logger = ctxLogger.entry
	} else {
		logger = Config{Format: "json", Debug: false}.Configure(ctx)
	}

	logger.WithFields(logrus.Fields{
		"error":       "panic",
		"panic":       r,
		"stack_trace": string(st),
	}).Error("request")

	if propagate {
		panic(r)
	}
}
