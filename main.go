package main

import (
	"golang-crud/route"
)

func main() {
	r := route.SetupRoute()
	r.Run("localhost:8181")
}
