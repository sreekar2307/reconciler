package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"github.com/sreekar2307/reconciler/cmd"
	"github.com/sreekar2307/reconciler/cmd/migrate"
	"github.com/sreekar2307/reconciler/cmd/recon"
	"github.com/sreekar2307/reconciler/cmd/seed"
	"syscall"
)

var (
	dbName   = flag.String("dbName", "recon", "Database name to use for reconciliation")
	mongoURI = flag.String("mongoURI", "mongodb://localhost:27017/?replicaSet=rs0", "MongoDB URI to connect to")
)

func init() {
	flag.Parse()
	if *dbName == "" || *mongoURI == "" {
		panic("dbName and mongoURI must be provided")
	}
}

func main() {
	ctx := context.Background()
	ctx, cancelFunc := signal.NotifyContext(ctx, syscall.SIGKILL, syscall.SIGTERM)
	defer cancelFunc()
	deps, err := cmd.NewDeps(ctx, *dbName, *mongoURI)
	if err != nil {
		panic("failed to initialize dependencies: " + err.Error())
	}
	if len(os.Args) < 2 {
		panic("no command provided")
	}
	switch os.Args[1] {
	case "recon":
		if err := recon.Run(ctx, deps); err != nil {
			panic("reconciliation failed: " + err.Error())
		}
	case "migrate":
		if err := migrate.Run(ctx, deps); err != nil {
			panic("migration failed: " + err.Error())
		}
	case "seed":
		if err := seed.Run(ctx, deps); err != nil {
			panic("seeding failed: " + err.Error())
		}
	default:
		panic("unknown command: " + os.Args[1])

	}
}
