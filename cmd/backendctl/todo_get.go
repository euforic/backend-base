package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/euforic/backend-base/proto"
	"github.com/spf13/cobra"
)

var getTodoCmd = &cobra.Command{
	Use:   "get [id]",
	Short: "get a todo by id",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			exitError(errors.New("no todo id provided"), 1)
		}
		c := client()
		res, err := c.GetTodo(context.Background(), &proto.GetTodoReq{Id: args[0]})
		if err != nil {
			exitError(err, 1)
		}
		out, _ := json.MarshalIndent(res, "", "  ")
		fmt.Println(string(out))
	},
}

func init() {
	todoCmd.AddCommand(getTodoCmd)
}
