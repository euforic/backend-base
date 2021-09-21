package main

import (
	"context"
	"encoding/json"
	"errors"
	"os"

	"github.com/euforic/backend-base/proto"
	"github.com/spf13/cobra"
)

var updateTodoCmd = &cobra.Command{
	Use:   "update",
	Short: "update an todo",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			exitError(errors.New("No todo data provided"), 0)
		}

		c := client()

		req := proto.UpdateTodoReq{}
		if err := json.Unmarshal([]byte(args[0]), &req); err != nil {
			os.Stderr.WriteString(err.Error())
			os.Exit(500)
			return
		}

		_, err := c.UpdateTodo(context.Background(), &req)
		if err != nil {
			exitError(err, 1)
		}
	},
}

func init() {
	todoCmd.AddCommand(updateTodoCmd)
}
