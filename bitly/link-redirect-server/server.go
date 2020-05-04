package main

import (
	"net/http"

	"github.com/cmpe281-sshekhar93/bitly/link-redirect-server/controller"

	router "github.com/cmpe281-sshekhar93/bitly/link-redirect-server/http"
)

var (
	httpRouter     router.Router             = router.NewMuxRouter()
	linkController controller.LinkController = controller.NewController()
)

const port string = ":2000"

func main() {
	httpRouter.GET("/", func(resp http.ResponseWriter, req *http.Request) {
		resp.Header().Set("Content-Type", "application/json")
		resp.WriteHeader(http.StatusOK)
		resp.Write([]byte("LRS up and running"))
	})
	httpRouter.GET("/{shortLink}", linkController.RedirectLink)

	httpRouter.SERVE(port)
}
