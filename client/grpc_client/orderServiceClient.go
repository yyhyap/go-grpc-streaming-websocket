package grpc_client

import (
	"context"
	"encoding/json"
	"fmt"
	"go-grpc-restaurant-client/logger"
	"go-grpc-restaurant-client/proto"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/hertz-contrib/websocket"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

const (
	// use "host.docker.internal" if running on docker container
	grpcServerHost = "localhost"
	grpcServerPort = ":8001"
)

var (
	orderServiceGRPC *orderServiceGRPCStruct

	orderServiceGRPCOnce sync.Once

	upgrader = websocket.HertzUpgrader{
		ReadBufferSize:  512,
		WriteBufferSize: 512,
		CheckOrigin: func(c *app.RequestContext) bool {
			return string(c.Method()[:]) == http.MethodGet
		},
	}

	clients     = make(map[*websocket.Conn]struct{})
	orderIdPool = make(map[string]struct{})

	clientsMux = &sync.RWMutex{}
	poolMux    = &sync.RWMutex{}
)

type IOrderServiceGRPC interface {
	CallGreetCustomer(ctx context.Context, name string) (string, error)
	CallCreateOrderByIngredients(ctx context.Context, ingredients []string) (string, error)
	CallGenerateRecommendedIngredientsCombo(ctx context.Context) (string, error)
	CallProcessOrder(ctx context.Context, c *app.RequestContext) error
}

type orderServiceGRPCStruct struct{}

type ProcessOrderRequest struct {
	OrderId string `json:"orderId"`
}

type ProcessOrderResponse struct {
	Total_order int                 `json:"totalOrder"`
	Orders      map[string]struct{} `json:"orders"`
}

type PBProcessOrderClient struct {
	Conn   *grpc.ClientConn
	Stream *proto.OrderService_ProcessOrderClient
}

func GetOrderServiceGRPC() *orderServiceGRPCStruct {
	if orderServiceGRPC == nil {
		orderServiceGRPCOnce.Do(func() {
			orderServiceGRPC = &orderServiceGRPCStruct{}
		})
	}
	return orderServiceGRPC
}

func dialGrpcServer() (*grpc.ClientConn, error) {
	return grpc.Dial(fmt.Sprintf("%s%s", grpcServerHost, grpcServerPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
}

func (o *orderServiceGRPCStruct) CallGreetCustomer(ctx context.Context, name string) (string, error) {
	conn, err := dialGrpcServer()

	if err != nil {
		return "", fmt.Errorf("did not connect to server: %v", err)
	}

	defer conn.Close()

	client := proto.NewOrderServiceClient(conn)

	res, err := client.GreetCustomer(ctx, &proto.GreetRequest{Name: name})

	if err != nil {
		return "", fmt.Errorf("error sending greet request to grpc request: %v", err)
	}

	return res.Name, nil
}

func (o *orderServiceGRPCStruct) CallCreateOrderByIngredients(ctx context.Context, ingredients []string) (string, error) {
	conn, err := dialGrpcServer()

	if err != nil {
		return "", fmt.Errorf("did not connect to server: %v", err)
	}

	defer conn.Close()

	client := proto.NewOrderServiceClient(conn)

	stream, err := client.CreateOrderByIngredients(ctx)

	if err != nil {
		return "", fmt.Errorf("error when starting client streaming for create order by ingredients: %v", err)
	}

	for _, ingredient := range ingredients {
		req := &proto.Ingredient{
			Name: ingredient,
		}

		err = stream.Send(req)

		if err != nil {
			return "", fmt.Errorf("error while client streaming for create order by ingredients: %v", err)
		}
	}

	res, err := stream.CloseAndRecv()

	if err != nil {
		return "", fmt.Errorf("error while closing client streaming and receive response: %v", err)
	}

	logger.Logger.Info("Client streaming finished")
	return res.Messsage, nil
}

func (o *orderServiceGRPCStruct) CallGenerateRecommendedIngredientsCombo(ctx context.Context) (string, error) {
	conn, err := dialGrpcServer()

	if err != nil {
		return "", fmt.Errorf("did not connect to server: %v", err)
	}

	defer conn.Close()

	client := proto.NewOrderServiceClient(conn)

	stream, err := client.GenerateRecommendedIngredientsCombo(ctx, &proto.EmptyParam{})

	if err != nil {
		return "", fmt.Errorf("error when starting server streaming for generating recommended ingredients combo: %v", err)
	}

	ingredientCombo := []string{}
	for {
		res, err := stream.Recv()

		if err == io.EOF {
			break
		}
		if err != nil {
			return "", fmt.Errorf("error while server streaming for generating recommended ingredients combo: %v", err)
		}

		ingredientCombo = append(ingredientCombo, res.Name)
	}

	var sb strings.Builder
	for i, ingredient := range ingredientCombo {
		if i == 0 {
			sb.WriteString(ingredient)
		} else {
			sb.WriteString(fmt.Sprintf(", %s", ingredient))
		}
	}
	return fmt.Sprintf("The ingredients combo of the days is: %s", sb.String()), nil
}

func (o *orderServiceGRPCStruct) CallProcessOrder(ctx context.Context, c *app.RequestContext) error {
	err := upgrader.Upgrade(c, func(conn *websocket.Conn) {
		defer conn.Close()
		client := &PBProcessOrderClient{}
		grpcConn, err := dialGrpcServer()

		if err != nil {
			logger.Logger.Error(fmt.Errorf("did not connect to grpc server: %v", err).Error())
			return
		}

		defer func() {
			err := grpcConn.Close()
			if err != nil {
				logger.Logger.Error(fmt.Errorf("unable to close grpc server connection: %v", err).Error())
			}
		}()

		client.Conn = grpcConn

		reconnectCh := make(chan struct{})

		defer close(reconnectCh)

		var isConnectedWebSocket *bool
		t := true
		isConnectedWebSocket = &t

		defer func() {
			*isConnectedWebSocket = false
		}()

		err = o.generateNewProcessOrderStream(client, ctx)

		defer func() {
			// close send
			if client.Stream != nil {
				err := (*client.Stream).CloseSend()
				if err != nil {
					logger.Logger.Error(fmt.Errorf("unable to close send bidirectional stream: %v", err).Error())
				}
			}
		}()

		// reestablish stream connection
		websocketCtx, websocketCancel := context.WithCancel(context.Background())
		defer websocketCancel()
		go o.reestablishStreamConnection(reconnectCh, client, websocketCtx, isConnectedWebSocket)

		if err != nil {
			if *isConnectedWebSocket {
				reconnectCh <- struct{}{}
			}
			logger.Logger.Error(fmt.Errorf("error when initializea a bi-directional streaming for processing order: %v", err).Error())
		}

		o.addClient(conn)
		logger.Logger.Info(fmt.Sprintf("Added a new connection. Total connection: %v", len(clients)))

		defer func() {
			o.disconnectClient(conn)
			logger.Logger.Info(fmt.Sprintf("A client has been disconnected. Total connection: %v", len(clients)))
		}()

		poolMux.RLock()
		firstProcessOrderResponse := ProcessOrderResponse{
			Total_order: len(orderIdPool),
			Orders:      orderIdPool,
		}
		poolMux.RUnlock()

		data, err := json.Marshal(firstProcessOrderResponse)

		if err != nil {
			logger.Logger.Error(err.Error())
			return
		}

		err = conn.WriteMessage(websocket.TextMessage, data)

		if err != nil {
			logger.Logger.Error(err.Error())
			return
		}

		// listening server side streaming
		go o.listenProcessOrderServerSide(client, reconnectCh, websocketCtx, isConnectedWebSocket)

		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				logger.Logger.Warn(err.Error())
				return
			}

			var processOrderReq ProcessOrderRequest
			err = json.Unmarshal(msg, &processOrderReq)

			if err != nil {
				logger.Logger.Error(err.Error())
				continue
			}

			poolMux.Lock()
			orderIdPool[processOrderReq.OrderId] = struct{}{}
			poolMux.Unlock()

			poolMux.RLock()
			processOrderResponse := ProcessOrderResponse{
				Total_order: len(orderIdPool),
				Orders:      orderIdPool,
			}
			poolMux.RUnlock()

			data, err := json.Marshal(processOrderResponse)

			if err != nil {
				logger.Logger.Error(err.Error())
				continue
			}

			go o.broadcast(data)

			// grpc client side bidirection streaming
			grpcReq := &proto.OrderId{
				OrderId: processOrderReq.OrderId,
			}

			if client.Conn.GetState() == connectivity.Ready {
				if client.Stream != nil {
					err = (*client.Stream).Send(grpcReq)

					if err == io.EOF {
						logger.Logger.Info("EOF: Stream ended, attempting to terminate websocket connection...")
					}
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
						logger.Logger.Info("Stream ended, attempting to terminate websocket connection...")
					}
				} else {
					logger.Logger.Warn("grpc server connection is ready, but the stream is nil...")
				}
			} else {
				logger.Logger.Warn("grpc server connection is not ready yet, unable to send request to server...")
			}
		}
	})
	if err != nil {
		return fmt.Errorf("error on websocket connection: %v", err)
	}

	return nil
}

func (o *orderServiceGRPCStruct) broadcast(msg []byte) {
	clientsMux.Lock()
	wg := new(sync.WaitGroup)
	for client := range clients {
		wg.Add(1)
		go func(conn *websocket.Conn, waitGroup *sync.WaitGroup) {
			defer waitGroup.Done()
			err := conn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				logger.Logger.Error(err.Error())
			}
		}(client, wg)
	}
	wg.Wait()
	clientsMux.Unlock()
}

func (o *orderServiceGRPCStruct) addClient(conn *websocket.Conn) {
	clientsMux.Lock()
	clients[conn] = struct{}{}
	clientsMux.Unlock()
}

func (o *orderServiceGRPCStruct) disconnectClient(conn *websocket.Conn) {
	clientsMux.Lock()
	delete(clients, conn)
	clientsMux.Unlock()
}

func (o *orderServiceGRPCStruct) generateNewProcessOrderStream(pbClient *PBProcessOrderClient, ctx context.Context) error {
	client := proto.NewOrderServiceClient(pbClient.Conn)
	stream, err := client.ProcessOrder(ctx)

	if err != nil {
		pbClient.Stream = nil
		return err
	}
	pbClient.Stream = &stream
	return nil
}

func (o *orderServiceGRPCStruct) reestablishStreamConnection(reconnectCh chan struct{}, client *PBProcessOrderClient,
	ctx context.Context, isConnectedWebSocket *bool) {
	defer func() {
		logger.Logger.Info("reestablishStreamConnection goroutine killed")
	}()
	for {
		select {
		case <-ctx.Done():
			return
		case <-reconnectCh:
			if client.Conn.GetState() != connectivity.Ready && *isConnectedWebSocket {
				if o.waitUntilReady(client, isConnectedWebSocket, ctx) {
					err := o.generateNewProcessOrderStream(client, ctx)
					if err != nil {
						logger.Logger.Error("failed to establish stream connection to grpc server ...")
					}

					// re-listening server side streaming
					go o.listenProcessOrderServerSide(client, reconnectCh, ctx, isConnectedWebSocket)
				}
			}
		}
	}
}

func (o *orderServiceGRPCStruct) listenProcessOrderServerSide(client *PBProcessOrderClient,
	reconnectCh chan struct{}, ctx context.Context, isConnectedWebSocket *bool) {
	defer func() {
		logger.Logger.Info("listenProcessOrderServerSide goroutine killed")
	}()
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if client.Stream != nil {
				msg, err := (*client.Stream).Recv()

				if err == io.EOF {
					logger.Logger.Info(fmt.Sprintf("Stream EOF: %v", err))
					if *isConnectedWebSocket {
						reconnectCh <- struct{}{}
					}
					return
				}
				if err != nil {
					e, ok := status.FromError(err)
					if ok {
						switch e.Code() {
						case codes.Canceled:
							logger.Logger.Warn(fmt.Sprintf("Context canceled: %v", err))
						case codes.Unavailable:
							logger.Logger.Warn(fmt.Sprintf("Stream is unavailable: %v", err))
						default:
							logger.Logger.Error(err.Error())
						}

					} else {
						logger.Logger.Error(fmt.Errorf("error occured while listening bi-directional server side stream: %v", err).Error())
					}
					if *isConnectedWebSocket {
						reconnectCh <- struct{}{}
					}
					return
				}

				processedOrderId := msg.OrderId

				poolMux.Lock()
				delete(orderIdPool, processedOrderId)
				poolMux.Unlock()

				poolMux.RLock()
				processOrderResponse := ProcessOrderResponse{
					Total_order: len(orderIdPool),
					Orders:      orderIdPool,
				}
				poolMux.RUnlock()

				data, err := json.Marshal(processOrderResponse)

				if err != nil {
					logger.Logger.Error(err.Error())
					continue
				}

				go o.broadcast(data)
			} else {
				logger.Logger.Warn("client stream is nil")
				if *isConnectedWebSocket {
					reconnectCh <- struct{}{}
				}
				return
			}
		}
	}
}

func (o *orderServiceGRPCStruct) waitUntilReady(client *PBProcessOrderClient, isConnectedWebSocket *bool, ctx context.Context) bool {
	for {
		select {
		case <-ctx.Done():
			return false
		default:
			if *isConnectedWebSocket {
				logger.Logger.Info(fmt.Sprintf("current state: %v", client.Conn.GetState()))

				if client.Conn.GetState() != connectivity.Ready {
					client.Conn.Connect()
				}

				// reserve a short duration (customizable) for conn to change state from idle to ready if grpc server is up
				time.Sleep(500 * time.Millisecond)

				if client.Conn.GetState() == connectivity.Ready {
					logger.Logger.Info(fmt.Sprintf("current state is ready!!!!: %v", client.Conn.GetState()))
					return true
				}

				// define reconnect time interval (backoff) or/and reconnect attempts here
				time.Sleep(2 * time.Second)
			} else {
				return false
			}
		}
	}
}
