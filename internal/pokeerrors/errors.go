package pokeerrors

import "fmt"

type NotFoundError struct {
	Entity string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s not found\n", e.Entity)
}
