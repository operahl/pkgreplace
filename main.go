package main

import (
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "pkgreplace"
	app.Usage = "替换go项目中的包名"
	app.Version = "1.0.0"
	app.Action = DoAction
	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}
