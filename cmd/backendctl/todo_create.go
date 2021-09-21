package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/euforic/backend-base/proto"
	"github.com/spf13/cobra"
)

var createTodoCmd = &cobra.Command{
	Use:   "create",
	Short: "create an todo",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			exitError(errors.New("No todo data provided"), 0)
		}

		c := client()

		req := proto.CreateTodoReq{}
		if err := json.Unmarshal([]byte(args[0]), &req); err != nil {
			os.Stderr.WriteString(err.Error())
			os.Exit(500)
			return
		}

		res, err := c.CreateTodo(context.Background(), &req)
		if err != nil {
			exitError(err, 1)
		}
		fmt.Println(json.MarshalIndent(res, "", "  "))
	},
}

func init() {
	todoCmd.AddCommand(createTodoCmd)
}
