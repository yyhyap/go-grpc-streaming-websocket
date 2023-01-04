package grpc_server

import (
	"context"
	"fmt"
	"go-grpc-restaurant-server/logger"
	"go-grpc-restaurant-server/proto"
	"io"
	"strings"
	"sync"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	processedOrderIdPool = make(map[string]struct{})
	mux                  = &sync.RWMutex{}
)

type OrderServiceServer struct {
	proto.OrderServiceServer
}

func (s *OrderServiceServer) GreetCustomer(ctx context.Context, req *proto.GreetRequest) (*proto.GreetResponse, error) {
	return &proto.GreetResponse{
		Name: fmt.Sprintf("Hello, welcome to the restaurant, %s", req.Name),
	}, nil
}

func (s *OrderServiceServer) CreateOrderByIngredients(stream proto.OrderService_CreateOrderByIngredientsServer) error {
	ingredientSlice := []string{}
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// return order message
			var sb strings.Builder
			for i, ingredient := range ingredientSlice {
				if i == 0 {
					sb.WriteString(ingredient)
				} else {
					sb.WriteString(fmt.Sprintf(", %s", ingredient))
				}
			}
			orderMessage := &proto.OrderMessage{
				Messsage: fmt.Sprintf("Your order with indgredients: %s is now preparing. Thank you.", sb.String()),
			}
			return stream.SendAndClose(orderMessage)
		}
		if err != nil {
			return err
		}
		// simulate heavy task
		time.Sleep(1 * time.Second)
		ingredientSlice = append(ingredientSlice, req.Name)
	}
}

func (o *OrderServiceServer) GenerateRecommendedIngredientsCombo(emptyParam *proto.EmptyParam,
	stream proto.OrderService_GenerateRecommendedIngredientsComboServer) error {
	ingredientsOfTheDay := []string{"Beef Burger Patty", "Fish Fillet", "Lettuce", "Pickle", "Tomato", "Cheddar cheese"}
	for _, ingredient := range ingredientsOfTheDay {
		// simulate some heavy calculation
		time.Sleep(1 * time.Second)
		res := &proto.Ingredient{
			Name: ingredient,
		}

		err := stream.Send(res)

		if err != nil {
			return err
		}
	}
	return nil
}

func (o *OrderServiceServer) ProcessOrder(stream proto.OrderService_ProcessOrderServer) error {
	for {
		req, err := stream.Recv()

		if err == io.EOF {
			logger.Logger.Info("Bi-directional streaming end of file")
			return nil
		}
		if err != nil {
			if e, ok := status.FromError(err); ok {
				switch e.Code() {
				case codes.Canceled:
					logger.Logger.Warn(fmt.Sprintf("Context canceled: %v", err))
				default:
					logger.Logger.Error(err.Error())
				}
			} else {
				logger.Logger.Error(err.Error())
			}
			return err
		}

		orderId := req.OrderId

		// using go func here to spawn another goroutine to handle heavy task (simulated)
		go func(releaseOrderId string) {
			var res *proto.OrderId = nil
			mux.RLock()
			_, exist := processedOrderIdPool[releaseOrderId]
			mux.RUnlock()
			if exist {
				res = &proto.OrderId{
					OrderId: releaseOrderId,
				}
				mux.Lock()
				delete(processedOrderIdPool, releaseOrderId)
				mux.Unlock()
			} else {
				// simulate processing task
				time.Sleep(10 * time.Second)
				res = &proto.OrderId{
					OrderId: releaseOrderId,
				}
			}

			err := stream.Send(res)
			/*
				Can implement broadcast to every streaming client here if necessary
			*/

			if err != nil {
				if e, ok := status.FromError(err); ok {
					switch e.Code() {
					case codes.Canceled:
						logger.Logger.Warn(fmt.Sprintf("Context canceled: %v", err))
					case codes.Unavailable:
						logger.Logger.Warn(fmt.Sprintf("Stream is unavailable: %v", err))
					default:
						logger.Logger.Error(err.Error())
					}
				} else {
					logger.Logger.Error(err.Error())
				}
				mux.Lock()
				processedOrderIdPool[releaseOrderId] = struct{}{}
				mux.Unlock()
			}
		}(orderId)
	}
}
