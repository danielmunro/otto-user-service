/*
 * Otto user service
 */

package main

import (
	"github.com/akrylysov/algnhsa"
	"github.com/danielmunro/otto-user-service/internal"
)

func main() {
	router := internal.NewRouter()
	algnhsa.ListenAndServe(router, nil)
}
