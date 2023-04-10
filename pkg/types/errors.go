package types

import "errors"

var ErrorRetryLimitExceeded = errors.New("retry limit exceeded")
var ErrorTyping = errors.New("invalid type")
