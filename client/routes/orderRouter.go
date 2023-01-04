package routes

import (
	"go-grpc-restaurant-client/controllers"

	"github.com/cloudwego/hertz/pkg/app/server"
)

var (
	orderController controllers.IOrderController = controllers.GetOrderController()
)

func RegisterOrderRoutes(h *server.Hertz) {
	apiGroup := h.Group("/api/order")
	{
		apiGroup.POST("/greet-customer", orderController.GreetCustomer)
		apiGroup.POST("/create-order", orderController.CreateOrderByIngredients)
		apiGroup.GET("/generate-combo", orderController.GenerateRecommendedIngredientsCombo)
	}

	wsGroup := h.Group("/ws")
	{
		wsGroup.GET("/process-order", orderController.ProcessOrder)
	}
}
