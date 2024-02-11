package client

import "errors"

var (
	ErrClientNotFound     = errors.New("client not found")
	ErrClientCannotAfford = errors.New("client cannot afford this transaction")
)
