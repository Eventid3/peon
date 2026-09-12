// Package job implements the core job interface
package job

import "context"

type Job interface {
	Exec(ctx context.Context) error
}
