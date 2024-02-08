package models

import "errors"

var ErrEmptyList = errors.New("List is empty")
var ErrKeyAlreadyExists = errors.New("Key already exists")
var ErrSubscriberAlreadyExists = errors.New("Subscriber already exists")
var ErrSubscriberDoesNotExist = errors.New("Subscriber does not exist")
