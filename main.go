package main

import (
	"flag"

	"github.com/RyanMccaffrey/ExampleServer/server"
	"github.com/RyanMccaffrey/ExampleServer/web"
)

var (
	addr = flag.String("addr", "localhost:8080", "address to serve")
)

func main() {
	helloWorld := web.NewHelloWorldService("/hello")
	if err := helloWorld.ServeRoutes(); err != nil {
		panic(err)
	}
	server := server.NewExampleServerService(*addr)
	if err := server.Serve(); err != nil {
		panic(err)
	}
}
