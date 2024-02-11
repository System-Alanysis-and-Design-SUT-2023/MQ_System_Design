package models

import "errors"

var ErrQueueIsEmpty = errors.New("queue is empty")
var ErrKeyAlreadyExistsInQueue = errors.New("key already exists in queue")

var ErrSubscriberAlreadyExists = errors.New("subscriber already exists")
var ErrSubscriberDoesNotExist = errors.New("subscriber does not exist")
var ErrNoSubscriberExists = errors.New("no subscriber exists")
