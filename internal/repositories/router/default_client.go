package router

import "context"

// IClient interface
type IClient interface {
	ValidateIMEI(ctx context.Context, request Request) (*Response, error)
}
