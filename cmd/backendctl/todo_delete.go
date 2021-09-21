package main

import (
	"context"
	"errors"

	"github.com/euforic/backend-base/proto"
	"github.com/spf13/cobra"
)

var deleteTodoCmd = &cobra.Command{
	Use:   "delete [id]",
	Short: "delete the todo for the given id",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			exitError(errors.New("no todo id provided"), 1)
		}
		c := client()
		_, err := c.DeleteTodo(context.Background(), &proto.DeleteTodoReq{Id: args[0]})
		if err != nil {
			exitError(err, 1)
		}
	},
}

func init() {
	todoCmd.AddCommand(deleteTodoCmd)
}
