package main

import (
	"fmt"
	"log"
	"os"

	"github.com/euforic/backend-base/proto"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

var addr string

var todoCmd = &cobra.Command{
	Use:   "backendctl",
	Short: "cli client for backendgrpc server",
}

var client = func() proto.TodosServiceClient {
	conn, err := grpc.Dial(
		addr,
		grpc.WithInsecure(),
	)

	if err != nil {
		log.Fatal(err)
	}

	return proto.NewTodosServiceClient(conn)
}

func main() {
	if err := todoCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	todoCmd.PersistentFlags().StringVar(&addr, "addr", "0.0.0.0:9000", "address for service. i.e 0.0.0.0:9000")
}

func exitError(err error, code int) {
	fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
	os.Exit(code)
}
