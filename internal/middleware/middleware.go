package middleware

import "lommeulken/gen/dbstore"

type Middleware struct {
	queries *dbstore.Queries
}

func NewMiddleware(queries *dbstore.Queries) *Middleware {
	return &Middleware{
		queries: queries,
	}
}
