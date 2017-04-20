package message

import "errors"

// InvalidSignatureError is returned when the signature on a received message does not
// validate.
var InvalidSignatureError = errors.New("A message had an invalid signature")
