package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strings"

	"github.com/euforic/backend-base/database"
	"github.com/euforic/backend-base/database/adapter/buntdb"
	"github.com/euforic/backend-base/database/adapter/postgres"
	"github.com/euforic/backend-base/pkg/flags"
	"github.com/euforic/backend-base/proto"
	"github.com/euforic/backend-base/server"
	"github.com/euforic/backend-base/server/todos"
	otgrpc "github.com/opentracing-contrib/go-grpc"
	"github.com/opentracing/opentracing-go"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/health/grpc_health_v1"
)

var Version = "undefined [re-build with Makefile]"

var serverCmdConfig = struct {
	Debug bool
	DbURL string
	Port  string
}{}

var serverCmd = &cobra.Command{
	Use:   "backendgrpc",
	Short: "Run the gRPC server",
	Run: func(cmd *cobra.Command, args []string) {
		healthService := server.NewHealthChecker()
		// Enable pprof profiler
		if serverCmdConfig.Debug {
			fmt.Println("Profiler Running")
			go func() {
				if err := http.ListenAndServe(":"+serverCmdConfig.Port, nil); err != nil {
					panic(err)
				}
			}()
		}
		dbURL := strings.Split(serverCmdConfig.DbURL, "://")
		if len(dbURL) != 2 {
			log.Fatal("Invalid or malformed database connection string")
		}

		// create database adapter
		var db database.Adapter
		switch dt := dbURL[0]; dt {
		case "bunt":
			db = buntdb.New(dbURL[1])
		case "postgresql":
			db = postgres.New(dbURL[1])
		default:
			log.Fatal("Invalid or missing database connection string")
		}

		// Init services
		todosService := todos.New(db)
		if err := db.Connect(); err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		// Adds gRPC internal logs. This is quite verbose, so adjust as desired!
		log := grpclog.NewLoggerV2(os.Stdout, ioutil.Discard, ioutil.Discard)
		grpclog.SetLoggerV2(log)

		addr := "0.0.0.0:" + serverCmdConfig.Port
		lis, err := net.Listen("tcp", addr)
		if err != nil {
			log.Fatalln("Failed to listen:", err)
		}

		tracer := opentracing.GlobalTracer()
		s := grpc.NewServer(
			grpc.UnaryInterceptor(
				otgrpc.OpenTracingServerInterceptor(tracer)),
			grpc.StreamInterceptor(
				otgrpc.OpenTracingStreamServerInterceptor(tracer)),
		)

		proto.RegisterTodosServiceServer(s, todosService)
		grpc_health_v1.RegisterHealthServer(s, healthService)

		// Serve gRPC Server
		log.Info("Serving gRPC on http://", addr)
		log.Fatal(s.Serve(lis))
	},
}

func init() {
	serverCmd.Flags().StringVarP(&serverCmdConfig.DbURL, "db", "d", "bunt://data.db", "database connection url")
	serverCmd.Flags().StringVarP(&serverCmdConfig.Port, "port", "p", "9000", "gRPC server port")
	serverCmd.Flags().BoolVar(&serverCmdConfig.Debug, "debug", false, "enable debug / profiler mode")
	// Set Flags from env vars
	if err := flags.SetPflagsFromEnv("BACKEND", serverCmd.Flags()); err != nil {
		os.Stderr.WriteString(err.Error())
		os.Exit(500)
		return
	}
}

func main() {
	serverCmd.Version = Version
	if err := serverCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
