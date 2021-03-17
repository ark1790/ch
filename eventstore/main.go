package main

import (
	"github.com/ark1790/ch/eventstore/cmd"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load(".env")

	cmd.RootCmd.Execute()
}
