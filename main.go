package main

import (
	"fmt"

	"github.com/gistapp/api/gists"
	"github.com/gistapp/api/server"
	"github.com/gistapp/api/utils"
)

func main() {
	port := utils.Get("PORT")
	s := server.NewServer(fmt.Sprintf(":%s", port))

	gistRouter := gists.GistRouter{
		Controller: gists.GistController,
	}

	// Start the server
	s.Ignite(&gistRouter)
}
