package model

import "time"

type ActivityOptions struct {
	TaskQueue              string
	ScheduleToCloseTimeout time.Duration
	ScheduleToStartTimeout time.Duration
	StartToCloseTimeout    time.Duration
	HeartbeatTimeout       time.Duration
	WaitForCancellation    bool
	ActivityID             string
	RetryPolicy            *RetryPolicy
	DisableEagerExecution  bool
}

type RetryPolicy struct {
	InitialInterval        time.Duration
	BackoffCoefficient     float64
	MaximumInterval        time.Duration
	MaximumAttempts        int32
	NonRetryableErrorTypes []string
}

type ChildWorkflowOptions struct {
	TaskQueue                string
	WorkflowExecutionTimeout time.Duration
	WorkflowRunTimeout       time.Duration
	WorkflowTaskTimeout      time.Duration
	RetryPolicy              *RetryPolicy
	CronSchedule             string
	Memo                     map[string]interface{}
	StartDelay               time.Duration
}
