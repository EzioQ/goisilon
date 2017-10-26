package main

import (
	"context"
	"fmt"

	"github.com/thecodeteam/gournal"
	"github.com/thecodeteam/gournal/logrus"
)

// myString is a custom type that has a custom fmt.Format function.
// This function should *not* be invoked unless the log level is such that the
// log message would actually get emitted. This saves resources as fields
// and formatters are not invoked at all unless the log level allows an
// entry to be logged.
type myString string

func (s myString) Format(f fmt.State, c rune) {
	fmt.Println("* INVOKED MYSTRING FORMATTER")
	fmt.Fprint(f, string(s))
}

func main() {
	ctx := context.Background()
	ctx = context.WithValue(ctx, gournal.AppenderKey(), logrus.New())

	counter := 0

	// Set up a context fields callback that will print a loud message to the
	// console so it is very apparent when the function is invoked. This
	// function should *not* be invoked unless the log level is such that the
	// log message would actually get emitted. This saves resources as fields
	// and formatters are not invoked at all unless the log level allows an
	// entry to be logged.
	getCtxFieldsFunc := func() map[string]interface{} {
		counter++
		fmt.Println("* INVOKED CONTEXT FIELDS")
		return map[string]interface{}{"counter": counter}
	}
	ctx = context.WithValue(ctx, gournal.FieldsKey(), getCtxFieldsFunc)

	var name myString = "Bob"

	// Log "Hello Bob" at the INFO level. This log entry will not get emitted
	// because the default Gournal log level (configurable by
	// gournal.DefaultLevel) is ERROR.
	//
	// Additionally, we should *not* see the messages produced by the
	// myString.Format and getCtxFieldsFunc functions.
	gournal.Info(ctx, "Hello %s", name)

	// Keep a reference to the context that has the original log level.
	oldCtx := ctx

	// Set the context's log level to be INFO.
	ctx = context.WithValue(ctx, gournal.LevelKey(), gournal.InfoLevel)

	// Note the log level has been changed to INFO. This is also a marker to
	// show that the previous log and messages generated by the functions should
	// not have occurred prior to this statement in the terminal.
	fmt.Println("* CTX LOG LEVEL INFO")

	name = "Mary"

	fields := map[string]interface{}{
		"length":   8,
		"location": "Austin",
	}

	// Log "Hello Mary" with some field information. We should not only see
	// the messages from the myString.Format and getCtxFieldsFunc functions,
	// but the field "size" from the getCtxFieldsFunc function should add the
	// field "counter" to the fields provided directly to this call.
	gournal.WithFields(fields).Info(ctx, "Hello %s", name)

	// Log "Hello Mary" again with the exact same info, except use the original
	// context that did not have an explicit log level. Since the default log
	// level is still ERROR, nothing will be emitted, not even the messages that
	// indicate the myString.Format or getCtxFieldsFunc functions are being
	// invoked.
	gournal.WithFields(fields).Info(oldCtx, "Hello %s", name)

	// Update the default log level to INFO
	gournal.DefaultLevel = gournal.InfoLevel
	fmt.Println("* DEFAULT LOG LEVEL INFO")

	// Log "Hello Mary" again with the exact same info, even use the original
	// context that did not have an explicit log level. However, since the
	// default log level is now INFO, the entry will be emitted, along with the
	// messages from the myString.Format or getCtxFieldsFunc functions are being
	// invoked.
	//
	// Note the counter value has only be incremented once since the function
	// was not invoked when the log level did not permit the entry to be logged.
	gournal.WithFields(fields).Info(oldCtx, "Hello %s", name)
}
