package data

import (
	"fmt"
	"strconv"
)

type runtime int32

// The rule about pointers vs. values for receivers is that value methods can be invoked on
// pointers and values, but pointer methods can only be invoked on pointers.
func (r runtime) MarshalJSON() ([]byte, error) {
	// Generate a string containing the movie runtime in the required format.
	jsonValue := fmt.Sprintf("%d mins", r)

	// Use the strconv.Quote() function on the string to wrap it in double quotes. It
	// needs to be surrounded by double quotes in order to be a valid *JSON string*.
	quotedJsonValue := strconv.Quote(jsonValue)

	return []byte(quotedJsonValue), nil
}
