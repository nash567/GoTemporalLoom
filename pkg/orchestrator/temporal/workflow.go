package temporal

import (
	"fmt"
	logModel "github.com/nash-567/goTemporalLoom/pkg/logger/model"
	"github.com/nash-567/goTemporalLoom/pkg/orchestrator/temporal/model"
	"go.temporal.io/sdk/workflow"
	"time"
)

// workflowEngine encapsulates workflow functionalities that are provided by Temporal through dangling functions.
type workflowEngine struct{}

// GetLogger returns a logger to be used in the workflow's context.
func (we *workflowEngine) GetLogger(ctx model.Context) logModel.KeyValLogger {
	return workflow.GetLogger(model.ToTemporalContext(ctx))
}

// ExecuteActivity executes a workflow activity.
func (we *workflowEngine) ExecuteActivity(
	ctx model.Context,
	activity interface{},
	args ...interface{},
) model.Future {
	return newFuture(workflow.ExecuteActivity(model.ToTemporalContext(ctx), activity, args...))
}

// SetQueryHandler sets a query handler for the workflow.
func (we *workflowEngine) SetQueryHandler(ctx model.Context, queryType string, handler interface{}) error {
	if err := workflow.SetQueryHandler(model.ToTemporalContext(ctx), queryType, handler); err != nil {
		return fmt.Errorf("set query handler: %w", err)
	}
	return nil
}

// Sleep pauses the workflow for the specified duration.
func (we *workflowEngine) Sleep(ctx model.Context, d time.Duration) error {
	if err := workflow.Sleep(model.ToTemporalContext(ctx), d); err != nil {
		return fmt.Errorf("error sleeping workflow: %w", err)
	}
	return nil
}

// GetSignalChannel returns a channel to receive signals for the workflow.
func (we *workflowEngine) GetSignalChannel(ctx model.Context, signalName string) model.ReceiveChannel {
	return newReceiveChannel(workflow.GetSignalChannel(model.ToTemporalContext(ctx), signalName))
}

// WithActivityOptions sets the options for the workflow's activities.
func (we *workflowEngine) WithActivityOptions(
	ctx model.Context, options model.ActivityOptions,
) model.Context {
	return newContext(
		workflow.WithActivityOptions(
			model.ToTemporalContext(ctx),
			model.ToTemporalActivityOptions(&options),
		),
	)
}
