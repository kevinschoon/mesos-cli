package filter

import (
	"errors"
	"github.com/gogo/protobuf/proto"
)

var (
	ErrNotFound       = errors.New("Not found")
	ErrTooManyResults = errors.New("Too many results")
)

// Messages is a filterable array of protobuf.Message.
type Messages []proto.Message

// FindAny will return the first message
// where all filters return true. If no
// messages match we will return ErrNotFound.
func (messages Messages) FindAny(filters ...Filter) (proto.Message, error) {
	var match bool
loop:
	for _, message := range messages {
		for _, filter := range filters {
			match = filter(message)
			if !match {
				continue loop
			}
		}
		if match {
			return message, nil
		}
	}
	return nil, ErrNotFound
}

// FindOne will a single message where all
// filters return true. If more than one
// message matches return ErrTooManyResults.
// If no messages match all filters we will
// return ErrNotFound.
func (messages Messages) FindOne(filters ...Filter) (proto.Message, error) {
	var (
		position int
		match    bool
	)
loop:
	for i, message := range messages {
		for _, filter := range filters {
			if !filter(message) {
				continue loop
			}
		}
		// Already matched a message
		if match {
			return nil, ErrTooManyResults
		}
		// Mark this message as matched
		match = true
		position = i
	}
	if !match {
		return nil, ErrNotFound
	}
	return messages[position], nil
}

// FindMany will return as many messages
// in which all of the filters match true.
// If no messages match all the filters
// we will return an empty []proto.Message.
func (messages Messages) FindMany(filters ...Filter) Messages {
	matches := Messages{}
loop:
	for _, message := range messages {
		for _, filter := range filters {
			if !filter(message) {
				continue loop
			}
		}
		matches = append(matches, message)
	}
	return matches
}
