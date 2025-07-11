package cmd

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"github.com/sreekar2307/reconciler/internal/recon"
	"github.com/sreekar2307/reconciler/internal/repository"
	"github.com/sreekar2307/reconciler/pkg/db/mongodb"
)

type Deps struct {
	Client        *mongo.Client
	ReconDatabase string
	Reconciler    recon.Reconcile
	Repository    repository.Repository
}

func NewDeps(
	ctx context.Context,
	dbName,
	mongoURI string,
) (*Deps, error) {
	client, err := mongodb.Client(ctx, mongoURI)
	if err != nil {
		return nil, err
	}
	deps := new(Deps)
	deps.Client = client
	deps.ReconDatabase = dbName
	deps.Repository = repository.NewRepository(client, dbName)
	deps.Reconciler = recon.NewReconcile(deps.Repository)
	return deps, nil
}
