package main

import (
	"log"

	"github.com/stdyum/api-sso/internal"
)

func main() {
	log.Fatal(internal.App().Run())
}
