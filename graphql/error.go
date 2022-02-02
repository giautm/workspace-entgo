package graphql

import (
	"context"
	"errors"

	"entgo.io/ent/privacy"
	"giautm.dev/awesome/pkg/logging"
	"github.com/99designs/gqlgen/graphql"
	"github.com/getsentry/sentry-go"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"go.uber.org/zap"
)

func recoverFunc(ctx context.Context, err interface{}) (userMessage error) {
	log := logging.FromContext(ctx)
	if e, ok := err.(error); ok {
		hub := sentry.GetHubFromContext(ctx)
		hub.ConfigureScope(createContextGQL(ctx, e, true))
		hub.CaptureException(e)

		log.Error("graphql internal failure", zap.Error(e))
	} else if s, ok := err.(string); ok {
		hub := sentry.GetHubFromContext(ctx)
		hub.ConfigureScope(createContextGQL(ctx, errors.New(s), true))
		hub.CaptureMessage(s)

		log.Error("graphql internal failure", zap.String("err", s))
	}

	return gqlerror.Errorf("Sorry, something went wrong")
}

func errorPresenter(ctx context.Context, err error) (gqlErr *gqlerror.Error) {
	hub := sentry.CurrentHub()

	// We trying to unwrap one level to check if there is an internal error.
	// Due to the bellow commit, GQLGen always wraps the error with `gqlerror.Error`
	// See: https://github.com/99designs/gqlgen/commit/e821b97bfbb922589c9eea649f0415ec3454e446
	if errInternal := errors.Unwrap(err); errInternal != nil {
		hub.ConfigureScope(createContextGQL(ctx, err, true))
		err = errInternal
	} else {
		hub.ConfigureScope(createContextGQL(ctx, err, false))
	}

	defer func() {
		if errors.Is(err, privacy.Deny) {
			gqlErr.Message = "Permission denied"
		}
	}()

	if errors.As(err, &gqlErr) {
		if gqlErr.Path == nil {
			gqlErr.Path = graphql.GetPath(ctx)
		}

		hub.CaptureException(err)
		return gqlErr
	}

	path := graphql.GetPath(ctx)
	log := logging.FromContext(ctx)
	log.Error("graphql internal failure",
		zap.Error(err),
		zap.String("path", path.String()),
	)

	hub.CaptureException(err)
	return gqlerror.ErrorPathf(path, "Sorry, something went wrong")
}

func createContextGQL(ctx context.Context, err error, markAsInternal bool) func(scope *sentry.Scope) {
	return func(scope *sentry.Scope) {
		scope.SetTag("component", "graphql")

		gql := map[string]interface{}{
			"Internal Error": markAsInternal,
		}

		var gqlErr *gqlerror.Error
		if errors.As(err, &gqlErr) {
			gql["Error Message"] = gqlErr.Message
			gql["Error Extensions"] = gqlErr.Extensions
			if p := gqlErr.Path; p != nil {
				gql["Error Path"] = p.String()
			} else if p = graphql.GetPath(ctx); p != nil {
				gql["Error Path"] = p.String()
			}
		} else {
			if p := graphql.GetPath(ctx); p != nil {
				gql["Error Path"] = p.String()
			}
		}

		if graphql.HasOperationContext(ctx) {
			o := graphql.GetOperationContext(ctx)
			gql["Raw Query"] = o.RawQuery
			gql["Variables"] = o.Variables

			if op := o.Operation; op != nil {
				if name := op.Name; name != "" {
					scope.SetTag("graphql.operation_name", name)
				}
				if kind := string(op.Operation); kind != "" {
					scope.SetTag("graphql.operation_kind", kind)
				}
			}
		}

		scope.SetContext("GraphQL", gql)
	}
}
