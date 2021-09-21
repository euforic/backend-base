//go:build matr
// +build matr

package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/matr-builder/matr/matr"
	"github.com/pkg/errors"
)

// Build will build the requested binary
func Build(ctx context.Context, args []string) error {
	fs := flag.NewFlagSet("build", flag.ExitOnError)
	var platform = fs.String("p", "linux", "platform")
	fs.Parse(args)

	if len(args) < 1 {
		// loop through cmd dirctory to get all build paths
		files, err := ioutil.ReadDir("./cmd")
		if err != nil {
			log.Fatal(err)
		}

		for _, f := range files {
			if f.IsDir() {
				args = append(args, f.Name())
			}
		}
	}

	for _, b := range args {
		err := matr.Sh("GOOS=%s CGO_ENABLED=0 go build -a -ldflags '-s' -installsuffix cgo -o build/%s ./cmd/%s", *platform, b, b).Run()
		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
		}
	}

	return nil
}

// Docker will build the docker images for all services
func Docker(ctx context.Context, args []string) error {
	dockerfile := "."
	imgName := "backendgrpc"
	if len(args) > 0 {
		dockerfile = args[0] + ".Dockerfile"
		imgName = args[0]
	}
	fmt.Println("Building Docker Image")
	err := matr.Sh(`docker build %s -t euforic/%s:latest`, dockerfile, imgName).Run()
	return err
}

// Proto generates all the protobuf based artifacts
func Proto(ctx context.Context, args []string) error {
	if err := matr.Sh(`buf generate --path proto/base.proto`).Run(); err != nil {
		fmt.Println("Proto Gen")
		return errors.Wrap(err, "generate proto")
	}

	return nil
}

// Gql generates the graphql model after updating the graphqls file
func Gql(ctx context.Context, args []string) error {
	c := matr.Sh(`gqlgen generate --config gqlgen.yml`)
	c.Dir = "gql"
	if err := c.Run(); err != nil {
		fmt.Println("Gql Gen")
		return errors.Wrap(err, "generate gql")
	}

	return nil
}

// Generate regenerates the Protobuf schema and GraphQL model
func Generate(ctx context.Context, args []string) error {
	err := matr.Deps(ctx, Proto, Gql)
	if err != nil {
		fmt.Println("Provider error")
		return errors.Wrap(err, "regenerate provider")
	}

	return nil
}

// Run will run the requested program (backendgrpc, backendgql or mactl) if left empty backendrpc and backendgql will run
func Run(ctx context.Context, args []string) error {
	if len(args) == 0 {
		go func() {
			fmt.Println(matr.Sh(`go run ./cmd/backendgrpc --db bunt://data.db`).Run())
		}()
		go func() {
			fmt.Println(matr.Sh(`go run ./cmd/backendgql --port 8080`).Run())
		}()

		// keep things going
		select {}
	}

	if err := matr.Sh(`go run ./cmd/` + strings.Join(args, " ")).Run(); err != nil {
		return err
	}

	return nil
}
