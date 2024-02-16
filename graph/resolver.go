package graph

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct{}

type Int64 int64
type Time time.Time

func (i Int64) MarshalGQL() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%d"`, i)), nil
}
func (i *Int64) UnmarshalGQL(input interface{}) error {
	switch value := input.(type) {
	case int:
		*i = Int64(value)
		return nil
	case int64:
		*i = Int64(value)
		return nil
	case float64:
		*i = Int64(value)
		return nil
	case string:
		val, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		*i = Int64(val)
		return nil
	default:
		return errors.New("invalid Int64 value")
	}
}

func (t Time) MarshalTime() time.Time {
	return time.Time(t)
}
func (t *Time) UnmarshalTime(value time.Time) {
	*t = Time(value)
}
