package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

type Session interface {
	WithTransaction(context.Context, func(context.Context) (any, error), *TransactionOptions) (any, error)
	End(context.Context)
}

type TransactionOptions struct {
	WriteConcern string
	ReadConcern  string
}

type session struct {
	mongoSession mongo.Session
}

func (s *session) End(ctx context.Context) {
	s.mongoSession.EndSession(ctx)
}

func (s *session) WithTransaction(ctx context.Context, fn func(context.Context) (any, error), opts *TransactionOptions) (any, error) {
	txOptions := &options.TransactionOptions{
		ReadConcern:  &readconcern.ReadConcern{Level: opts.ReadConcern},
		WriteConcern: &writeconcern.WriteConcern{W: opts.WriteConcern},
	}
	return s.mongoSession.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		return fn(sessCtx)
	}, txOptions)
}
