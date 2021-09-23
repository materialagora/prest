package main

import (
	"github.com/materialagora/prest/cmd"
	"github.com/materialagora/prest/config"
)

func main() {
	config.Load()
	cmd.Execute()
}
