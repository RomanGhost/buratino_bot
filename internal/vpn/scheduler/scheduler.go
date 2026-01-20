package scheduler

import "context"

type Scheduler interface {
	Run(ctx context.Context)
}
