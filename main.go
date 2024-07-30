package main

import (
	"fmt"

	"github.com/gistapp/api/server"
	"github.com/gistapp/api/utils"
)

func main() {
	port := utils.Get("PORT")
	s := server.NewServer(fmt.Sprintf(":%s", port))
	// Start the server
	s.Start()
}
