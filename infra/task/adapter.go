package task

import "time"

type Task interface {
	Execute(fn func(), interval time.Duration) (cancel func(), err error)
}
