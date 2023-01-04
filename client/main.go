package main

import (
	"go-grpc-restaurant-client/routes"
	"log"
	"math/rand"
	"time"

	"github.com/cloudwego/hertz/pkg/app/middlewares/server/recovery"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/hertz-contrib/cors"
)

func init() {
	log.SetPrefix("[LOG] ")
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Llongfile)
	rand.Seed(time.Now().UnixNano())
}

func main() {
	h := server.Default(server.WithHostPorts("0.0.0.0:8000"))
	// default cors for localhost
	h.Use(cors.Default())
	h.Use(recovery.Recovery())

	h.NoHijackConnPool = true
	routes.RegisterRoutes(h)

	h.Spin()
}
