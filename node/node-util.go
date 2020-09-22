package node

import (
	"context"
	"fmt"
	"log"
)

var Canceled = fmt.Errorf("ctx canceled")

// LogErrors is a utility method that starts a new goroutine that logs
// Node errors.
func (n *Node) LogErrors(ctx context.Context) {
	n.Errs = make(chan error)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case err := <-n.Errs:
				log.Println(err.Error())
			}
		}
	}()
}

// LogMetrics is a utility method that starts a new goroutine that logs
// Metrics.
func (n *Node) LogMetrics(ctx context.Context) {
	n.Metrics = make(chan Metrics)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case m := <-n.Metrics:
				log.Printf("%s", JSON(m))
			}
		}
	}()
}

func (n *Node) err(ctx context.Context, e error) {
	if n.Errs == nil {
		log.Println(e.Error())
		return
	}
	select {
	case <-ctx.Done():
	case n.Errs <- e:
	}
}

type Warning struct {
	Err error
}

func (w *Warning) Error() string {
	return fmt.Sprintf("WARNING: %s", w.Err)
}

func Warningf(format string, args ...interface{}) *Warning {
	return &Warning{
		Err: fmt.Errorf(format, args...),
	}
}

func (n *Node) warnf(ctx context.Context, format string, args ...interface{}) {
	n.err(ctx, Warningf(format, args...))
}

func (n *Node) errf(ctx context.Context, format string, args ...interface{}) {
	n.err(ctx, fmt.Errorf(format, args...))
}

func (n *Node) logf(ctx context.Context, format string, args ...interface{}) {
	if n.Logging {
		log.Printf(format, args...)
	}
}
