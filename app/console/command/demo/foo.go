package demo

import (
	"fmt"
	"github.com/sunranlike/sunr/framework/cobra"
	_ "log"
)

var FooCommand = &cobra.Command{
	Use:   "foo",
	Short: "foo",
	RunE: func(c *cobra.Command, args []string) error {
		container := c.GetContainer()
		fmt.Println(container)
		return nil
	},
}
