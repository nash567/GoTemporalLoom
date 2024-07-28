package model

import (
	logModel "github.com/nash-567/goTemporalLoom/pkg/logger/model"
	"time"
)

type (
	Params interface {
		Marshal() ([]byte, error)
		String() string
	}
	// WorkflowEngine defines the operations available in a workflow,
	// it is a wrapper over dangling workflow functions.
	WorkflowEngine interface {
		// GetLogger returns a logger to be used in the workflow's context.
		GetLogger(ctx Context) logModel.KeyValLogger

		// ExecuteActivity executes a workflow activity.
		// The activity is executed asynchronously and the result is returned as a Future.
		// The Future.Get method should be used to block until the result is available.
		// The activity can be a function or the name of a registered activity.
		ExecuteActivity(ctx Context, activity interface{}, args ...interface{}) Future

		// SetQueryHandler sets a query handler for the workflow.
		// The query handler is a function that is called when a query is made to the workflow.
		// The query handler should be a function with the signature:
		//  func (ctx Context, queryArgs ...interface{}) (interface{}, error)
		// The queryArgs are the arguments passed to the query.
		// The query handler should return the result of the query and an error if any.
		// The query handler should be registered before the query is made.
		// The query handler can be registered multiple times for different query types.
		SetQueryHandler(ctx Context, queryType string, handler interface{}) error

		// Sleep pauses the workflow for the specified duration.
		// The workflow will be paused and will not consume any resources during the sleep.
		// The workflow will be resumed after the specified duration.
		Sleep(ctx Context, d time.Duration) error

		// GetSignalChannel returns a channel to receive signals for the workflow.
		// The signalName is the name of the signal to receive.
		// The signal channel can be used to receive signals from the workflow.
		GetSignalChannel(ctx Context, signalName string) ReceiveChannel

		// WithActivityOptions returns a new context with the provided activity options.
		// The activity options are used to configure the behavior of activities executed in the workflow.
		// The activity options can be used to configure the task queue, retry policy, and other activity options.
		WithActivityOptions(ctx Context, options ActivityOptions) Context

		// ExecuteChildWorkflow starts a new child workflow execution.
		ExecuteChildWorkflow(
			ctx Context, options ChildWorkflowOptions, childWorkflow WorkflowDescriptor, args Params,
		) (ChildWorkflowFuture, error)
	}

	Future interface {
		Get(ctx Context, valuePtr interface{}) error
		IsReady() bool
	}

	Context interface {
		Deadline() (deadline time.Time, ok bool)
		Done() Channel
		Err() error
		Value(key interface{}) interface{}
		GetParentContext() interface{}
	}

	// Channel must be used instead of native go channel by workflow code.
	// Use workflow.NewChannel(ctx) method to create Channel instance.
	Channel interface {
		SendChannel
		ReceiveChannel
	}

	SendChannel interface {
		Name() string
		// Send blocks until the data is sent.
		Send(ctx Context, v interface{})

		// SendAsync try to send without blocking. It returns true if the data was sent, otherwise it returns false.
		SendAsync(v interface{}) (ok bool)

		// Close the Channel, and prohibit subsequent sends.
		Close()
	}

	// ReceiveChannel is a read-only view of the Channel.
	ReceiveChannel interface {
		Name() string
		// Receive blocks until it receives a value, and then assigns the received value to the provided pointer.
		// Returns false when Channel is closed.
		// Parameter valuePtr is a pointer to the expected data structure to be received. For example,
		//  var v string
		//  c.Receive(ctx, &v)
		//
		// Note, values should not be reused for extraction here because merging on
		// top of existing values may result in unexpected behavior similar to
		// json.Unmarshal.
		Receive(ctx Context, valuePtr interface{}) (more bool)

		// ReceiveWithTimeout blocks up to timeout until it receives a value, and then assigns the received value to the
		// provided pointer.
		// Returns more value of false when Channel is closed.
		// Returns ok value of false when no value was found in the channel for the duration of timeout or
		// the ctx was canceled.
		// The valuePtr is not modified if ok is false.
		// Parameter valuePtr is a pointer to the expected data structure to be received. For example,
		//  var v string
		//  c.ReceiveWithTimeout(ctx, time.Minute, &v)
		//
		// Note, values should not be reused for extraction here because merging on
		// top of existing values may result in unexpected behavior similar to
		// json.Unmarshal.
		ReceiveWithTimeout(ctx Context, timeout time.Duration, valuePtr interface{}) (ok, more bool)

		// ReceiveAsync try to receive from Channel without blocking. If there is data available from the Channel, it
		// assigns the data to valuePtr and returns true. Otherwise, it returns false immediately.
		//
		// Note, values should not be reused for extraction here because merging on
		// top of existing values may result in unexpected behavior similar to
		// json.Unmarshal.
		ReceiveAsync(valuePtr interface{}) (ok bool)

		// ReceiveAsyncWithMoreFlag is the same as ReceiveAsync with extra return value more to indicate if there could be
		// more value from the Channel. The more is false when the Channel is closed.
		//
		// Note, values should not be reused for extraction here because merging on
		// top of existing values may result in unexpected behavior similar to
		// json.Unmarshal.
		ReceiveAsyncWithMoreFlag(valuePtr interface{}) (ok bool, more bool)

		// Len returns the number of buffered messages plus the number of blocked Send calls.
		Len() int
	}

	// WorkflowDescriptor defines the metadata and identification methods for a workflow.
	WorkflowDescriptor interface {
		descriptor
		// GenerateWorkflowID is used to generate a unique workflow ID.
		GenerateWorkflowID(in Params) (string, error)
	}

	descriptor interface {
		// Name returns the name of the workflow/activity.
		Name() string
		// Description returns the description of the workflow/activity.
		Description() string
	}

	// ChildWorkflowFuture represents the result of a child workflow execution.
	ChildWorkflowFuture interface {
		Future
		// GetChildWorkflowExecution returns a future that will be ready when child workflow execution started. You can
		// get the WorkflowExecution of the child workflow from the future. Then you can use Workflow ID and RunID of
		// child workflow to cancel or send signal to child workflow.
		//  childWorkflowFuture := workflow.ExecuteChildWorkflow(ctx, child, ...)
		//  var childWE workflow.Execution
		//  if err := childWorkflowFuture.GetChildWorkflowExecution().Get(ctx, &childWE); err == nil {
		//      // child workflow started, you can use childWE to get the WorkflowID and RunID of child workflow
		//  }
		GetChildWorkflowExecution() Future

		// SignalChildWorkflow sends a signal to the child workflow. This call will block until child workflow is started.
		SignalChildWorkflow(ctx Context, signalName string, data interface{}) Future
	}
)
