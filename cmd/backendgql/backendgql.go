package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/euforic/backend-base/gql/graph"
	"github.com/euforic/backend-base/gql/graph/generated"
	"github.com/euforic/backend-base/pkg/flags"
	"github.com/euforic/backend-base/pkg/jwt"
	"github.com/euforic/backend-base/proto"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

var Version = "undefined [re-build with Makefile]"
var envPrefix = "BACKENDGQL"

var gatewayCmdConfig = struct {
	GRPCAddr     string
	Port         string
	KeyServerURL string
}{}

var gatewayCmd = &cobra.Command{
	Use:   "backendgql",
	Short: "Run the graphQL gateway server",
	Run: func(cmd *cobra.Command, args []string) {
		jwtDecoder := jwt.New()
		port := gatewayCmdConfig.Port

		conn, err := grpc.Dial(
			gatewayCmdConfig.GRPCAddr,
			grpc.WithInsecure(),
		)

		if err != nil {
			log.Fatal(err)
		}

		todosClient := proto.NewTodosServiceClient(conn)
		srv := handler.NewDefaultServer(
			generated.NewExecutableSchema(
				generated.Config{
					Resolvers: graph.NewResolver(todosClient),
				},
			),
		)

		http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
			if _, err := w.Write([]byte("ok")); err != nil {
				os.Stderr.WriteString(err.Error())
			}
		})

		http.Handle("/", playground.Handler("GraphQL playground", "/graphql"))
		http.Handle("/graphql", jwt.AuthMiddleware(jwtDecoder, srv))

		log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
		log.Fatal(http.ListenAndServe(":"+port, nil))
	},
}

func init() {
	gatewayCmd.Flags().StringVarP(&gatewayCmdConfig.Port, "port", "p", "8080", "app server port")
	gatewayCmd.Flags().StringVarP(&gatewayCmdConfig.GRPCAddr, "grpc_addr", "g", "localhost:9000", "app server port")
}

func main() {
	gatewayCmd.Version = Version
	if err := flags.SetPflagsFromEnv(envPrefix, gatewayCmd.Flags()); err != nil {
		os.Stderr.WriteString(err.Error())
		os.Exit(-1)
	}

	if err := gatewayCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
