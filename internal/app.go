package internal

import (
	"log"

	"github.com/stdyum/api-sso/internal/use"
)

func App() *use.Use {
	useCase, err := use.Default()
	if err != nil {
		log.Fatal(err)
	}

	return useCase
}
