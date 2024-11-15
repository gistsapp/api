package main

import (
	"fmt"
	"os"

	"github.com/gistapp/api/gists"
	"github.com/gistapp/api/organizations"
	"github.com/gistapp/api/server"
	"github.com/gistapp/api/storage"
	"github.com/gistapp/api/user"
	"github.com/gistapp/api/utils"
	"github.com/gofiber/fiber/v2/log"
)

func main() {
	if len(os.Args) > 1 {
		args := os.Args[1]

		if args == "migrate" {
			err := storage.Migrate()
			if err != nil {
				log.Error(err)
				return
			}
			log.Info("Migration successful")
		} else {
			log.Error("unknown command")
		}
		return
	}

	port := utils.Get("PORT")
	s := server.NewServer(fmt.Sprintf(":%s", port))

	gistRouter := gists.GistRouter{
		Controller: gists.GistController,
	}

	authRouter := user.AuthRouter{
		Controller: &user.AuthControllerImpl{
			AuthService: &user.AuthService,
		},
	}

	userRouter := user.UserRouter{
		Controller: &user.UserControllerImpl{},
	}

	orgRouter := organizations.OrganizationRouter{
		Controller: organizations.OrganizationControllerImpl{},
	}

	user.AuthService.RegisterProviders() //register goth providers for authentication

	// Start the server
	s.Setup(&gistRouter, &authRouter, &orgRouter, &userRouter)
	s.Ignite()
}
