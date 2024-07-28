package model

import (
	"go.temporal.io/sdk/temporal"
	sdkWorkflow "go.temporal.io/sdk/workflow"
)

func ToTemporalContext(ctx Context) sdkWorkflow.Context {
	return ctx.GetParentContext().(sdkWorkflow.Context) //nolint:forcetypeassert
}

func toTemporalRetryPolicy(r *RetryPolicy) *temporal.RetryPolicy {
	return &temporal.RetryPolicy{
		InitialInterval:        r.InitialInterval,
		BackoffCoefficient:     r.BackoffCoefficient,
		MaximumInterval:        r.MaximumInterval,
		MaximumAttempts:        r.MaximumAttempts,
		NonRetryableErrorTypes: r.NonRetryableErrorTypes,
	}
}
func ToTemporalActivityOptions(o *ActivityOptions) sdkWorkflow.ActivityOptions {
	retryPolicy := &temporal.RetryPolicy{}
	if o.RetryPolicy != nil {
		retryPolicy = toTemporalRetryPolicy(o.RetryPolicy)
	}
	return sdkWorkflow.ActivityOptions{
		TaskQueue:              o.TaskQueue,
		ScheduleToCloseTimeout: o.ScheduleToCloseTimeout,
		ScheduleToStartTimeout: o.ScheduleToStartTimeout,
		StartToCloseTimeout:    o.StartToCloseTimeout,
		HeartbeatTimeout:       o.HeartbeatTimeout,
		WaitForCancellation:    o.WaitForCancellation,
		ActivityID:             o.ActivityID,
		RetryPolicy:            retryPolicy,
		DisableEagerExecution:  o.DisableEagerExecution,
	}
}
