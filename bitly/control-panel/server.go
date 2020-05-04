package main

import (
	"net/http"

	"github.com/cmpe281-sshekhar93/bitly/control-panel/controller"

	router "github.com/cmpe281-sshekhar93/bitly/control-panel/http"
)

var (
	httpRouter     router.Router             = router.NewMuxRouter()
	linkController controller.LinkController = controller.NewController()
)

const port string = ":8000"

func main() {
	httpRouter.GET("/", func(resp http.ResponseWriter, req *http.Request) {
		resp.Write([]byte("up and running"))
	})
	httpRouter.POST("/create", linkController.CreateLink)
	httpRouter.GET("/link/{id}", linkController.GetLink)
	httpRouter.GET("/links", linkController.GetLinks)
	httpRouter.PUT("/link", linkController.UpdateLink)
	httpRouter.DELETE("/link/{id}", linkController.DeleteLink)

	httpRouter.SERVE(port)
}
