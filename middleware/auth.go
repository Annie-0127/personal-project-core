package middleware

import (
	"context"
	"github.com/getkin/kin-openapi/openapi3filter"
	serverContext "personal-project-core/context"
)

func AuthMiddleware(csc *serverContext.ServiceContext) openapi3filter.AuthenticationFunc {
	return func(ctx context.Context, input *openapi3filter.AuthenticationInput) error {
		return nil
	}
}
