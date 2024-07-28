package temporal

import (
	"fmt"
	"github.com/nash-567/goTemporalLoom/pkg/orchestrator/temporal/model"
	"go.temporal.io/sdk/workflow"
	"time"
)

type (
	// channelWrapper is a wrapper around a Temporal Channel.
	// it inherits the Channel interface from the temporal package and overrides the method that have Context as a parameter.
	channelWrapper struct {
		workflow.Channel
	}

	// receiveChannelWrapper is a wrapper around a Temporal ReceiveChannel.
	// it inherits the ReceiveChannel interface from the temporal package and overrides the method that have Context as a parameter.
	receiveChannelWrapper struct {
		workflow.ReceiveChannel
	}

	// futureWrapper is a wrapper around a Temporal Future.
	futureWrapper struct {
		workflow.Future
	}

	// contextWrapper is a wrapper around a Temporal Context.
	contextWrapper struct {
		workflow.Context
	}
)

// === Channel Methods Start ===

func (c *channelWrapper) Send(ctx model.Context, v interface{}) {
	c.Channel.Send(model.ToTemporalContext(ctx), v)
}

func (c *channelWrapper) Receive(ctx model.Context, valuePtr interface{}) bool {
	return c.Channel.Receive(model.ToTemporalContext(ctx), valuePtr)
}

func (c *channelWrapper) ReceiveWithTimeout(ctx model.Context, timeout time.Duration, valuePtr interface{}) (bool, bool) {
	return c.Channel.ReceiveWithTimeout(model.ToTemporalContext(ctx), timeout, valuePtr)
}

func newChannel(ch workflow.Channel) model.Channel {
	return &channelWrapper{ch}
}

// === Channel Methods End ===

// === Future Methods Start ===
func newFuture(f workflow.Future) model.Future {
	return &futureWrapper{f}
}
func (f futureWrapper) Get(ctx model.Context, valuePtr interface{}) error {
	if err := f.Future.Get(model.ToTemporalContext(ctx), valuePtr); err != nil {
		return fmt.Errorf("decode output value: %w", err)
	}
	return nil
}

// === Future Methods End ===

// === ReceiveChannel Methods Start ===
func (r *receiveChannelWrapper) Receive(ctx model.Context, valuePtr interface{}) bool {
	return r.ReceiveChannel.Receive(model.ToTemporalContext(ctx), valuePtr)
}

func (r *receiveChannelWrapper) ReceiveWithTimeout(
	ctx model.Context,
	timeout time.Duration,
	valuePtr interface{},
) (bool, bool) {
	return r.ReceiveChannel.ReceiveWithTimeout(model.ToTemporalContext(ctx), timeout, valuePtr)
}

func newReceiveChannel(ch workflow.ReceiveChannel) model.ReceiveChannel {
	return &receiveChannelWrapper{ch}
}

// === ReceiveChannel Methods End ===

// === Context Methods Start ===

func (c *contextWrapper) GetParentContext() interface{} { return c.Context }

func (c *contextWrapper) Done() model.Channel { return newChannel(c.Context.Done()) }

func newContext(ctx workflow.Context) model.Context { return &contextWrapper{ctx} }

// === Context Methods End ===
