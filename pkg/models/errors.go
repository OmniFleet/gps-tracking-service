package models

import "errors"

var ErrNoRecord = errors.New("models: no matching records found")
var DecodeError = errors.New("models: unable to decode payload")
var ValidationError = errors.New("models: validation error")
var ReadyError = errors.New("datastore is not ready")
var AliveError = errors.New("datastoer is not alive")
