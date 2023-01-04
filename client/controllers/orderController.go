package controllers

import (
	"context"
	"go-grpc-restaurant-client/grpc_client"
	"go-grpc-restaurant-client/logger"
	"net/http"
	"sync"

	"github.com/cloudwego/hertz/pkg/app"
	hutils "github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/go-playground/validator"
)

var (
	orderController     *orderControllerStruct
	orderControllerOnce sync.Once

	orderServiceGRPC grpc_client.IOrderServiceGRPC = grpc_client.GetOrderServiceGRPC()

	validate = validator.New()
)

type IOrderController interface {
	GreetCustomer(ctx context.Context, c *app.RequestContext)
	CreateOrderByIngredients(ctx context.Context, c *app.RequestContext)
	GenerateRecommendedIngredientsCombo(ctx context.Context, c *app.RequestContext)
	ProcessOrder(ctx context.Context, c *app.RequestContext)
}

type orderControllerStruct struct{}

type GreetRequest struct {
	Name string `json:"name" validate:"required"`
}

type CreateOrderRequests struct {
	IngredientList []string `json:"ingredientList" validate:"required"`
}

func GetOrderController() *orderControllerStruct {
	if orderController == nil {
		orderControllerOnce.Do(func() {
			orderController = &orderControllerStruct{}
		})
	}
	return orderController
}

func (o *orderControllerStruct) GreetCustomer(ctx context.Context, c *app.RequestContext) {
	var req GreetRequest

	err := c.Bind(&req)

	if err != nil {
		logger.Logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, hutils.H{"error": err.Error()})
		return
	}

	err = validate.Struct(req)

	if err != nil {
		logger.Logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, hutils.H{"error": err.Error()})
		return
	}

	msg, err := orderServiceGRPC.CallGreetCustomer(ctx, req.Name)

	if err != nil {
		logger.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, hutils.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, hutils.H{"message": msg})
}

func (o *orderControllerStruct) CreateOrderByIngredients(ctx context.Context, c *app.RequestContext) {
	var coReq CreateOrderRequests

	err := c.Bind(&coReq)

	if err != nil {
		logger.Logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, hutils.H{"error": err.Error()})
		return
	}

	err = validate.Struct(coReq)

	if err != nil {
		logger.Logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, hutils.H{"error": err.Error()})
		return
	}

	if len(coReq.IngredientList) < 1 {
		logger.Logger.Warn("ingredient list should not be empty")
		c.JSON(http.StatusBadRequest, hutils.H{"error": "ingredient list should not be empty"})
		return
	}

	orderMessage, err := orderServiceGRPC.CallCreateOrderByIngredients(ctx, coReq.IngredientList)

	if err != nil {
		logger.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, hutils.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, hutils.H{"message": orderMessage})
}

func (o *orderControllerStruct) GenerateRecommendedIngredientsCombo(ctx context.Context, c *app.RequestContext) {
	msg, err := orderServiceGRPC.CallGenerateRecommendedIngredientsCombo(ctx)

	if err != nil {
		logger.Logger.Error(err.Error())
		c.JSON(http.StatusInternalServerError, hutils.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, hutils.H{"message": msg})
}

func (o *orderControllerStruct) ProcessOrder(ctx context.Context, c *app.RequestContext) {
	err := orderServiceGRPC.CallProcessOrder(ctx, c)
	if err != nil {
		logger.Logger.Error(err.Error())
	}
}
