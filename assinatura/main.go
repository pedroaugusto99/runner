/*
Copyright © 2026 Claudio Ferreira & Pedro Augusto
*/
package main

import "github.com/pedroaugusto99/runner/assinatura/cmd"

var (
	version = "dev"
	commit  = "unknown"
)

func main() {
	cmd.Execute(version, commit)
}
