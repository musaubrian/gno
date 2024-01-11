package main

import "github.com/musaubrian/gno"

func main() {
	g := gno.New()
	g.BootstrapBuild("out", "example", ".")
	g.AddCommand("./out/example")
	g.Build()
	g.RunCommandsSync()
}
