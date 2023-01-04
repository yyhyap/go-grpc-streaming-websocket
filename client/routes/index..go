package routes

import "github.com/cloudwego/hertz/pkg/app/server"

func RegisterRoutes(h *server.Hertz) {
	RegisterOrderRoutes(h)
}
