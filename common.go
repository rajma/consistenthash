package consistenthash

import (
	"fmt"
	"time"
)

//ApplicationError ...
type ApplicationError struct {
	When time.Time
	What string
}

//Implementation of Error interface
func (e *ApplicationError) Error() string {
	return fmt.Sprintf("at %v, %s",
		e.When, e.What)
}

//InitializeMin...
func InitializeMin(actual, min int32) int32 {
	if actual < min {
		return min
	}
	return actual
}
