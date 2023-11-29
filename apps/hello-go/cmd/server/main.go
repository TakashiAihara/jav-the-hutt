// Usage: go run ./cmd/server/main.go
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/bufbuild/connect-go"
	"github.com/rs/cors"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	greetv1 "hello-go/gen/greet/v1"        // generated by protoc-gen-go
	"hello-go/gen/greet/v1/greetv1connect" // generated by protoc-gen-connect-go
)

// GreetServer implements greetv1.GreetServiceServer.
type GreetServer struct{}

// Greet implements greetv1.GreetServiceServer.
func (s *GreetServer) Greet(
	ctx context.Context,
	req *connect.Request[greetv1.GreetRequest],
) (*connect.Response[greetv1.GreetResponse], error) {
	log.Println("Request headers: ", req.Header())
	res := connect.NewResponse(&greetv1.GreetResponse{
		Greeting: fmt.Sprintf("Hello, %s!", req.Msg.Name),
	})
	res.Header().Set("Greet-Version", "v1")
	return res, nil
}

// newCORS returns a CORS handler that allows all origins and methods. This is borrowed from https://github.com/bufbuild/connect-demo/blob/3a30d4de07d6ac42110acd4ebf64bb4bf8a62579/main.go#L88
func newCORS() *cors.Cors {
	// To let web developers play with the demo service from browsers, we need a
	// very permissive CORS setup.
	return cors.New(cors.Options{
		AllowedMethods: []string{
			http.MethodHead,
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
		},
		AllowOriginFunc: func(origin string) bool {
			// Allow all origins, which effectively disables CORS.
			return true
		},
		AllowedHeaders: []string{"*"},
		ExposedHeaders: []string{
			// Content-Type is in the default safelist.
			"Accept",
			"Accept-Encoding",
			"Accept-Post",
			"Connect-Accept-Encoding",
			"Connect-Content-Encoding",
			"Content-Encoding",
			"Grpc-Accept-Encoding",
			"Grpc-Encoding",
			"Grpc-Message",
			"Grpc-Status",
			"Grpc-Status-Details-Bin",
		},
		// Let browsers cache CORS information for longer, which reduces the number
		// of preflight requests. Any changes to ExposedHeaders won't take effect
		// until the cached data expires. FF caps this value at 24h, and modern
		// Chrome caps it at 2h.
		MaxAge: int(2 * time.Hour / time.Second),
	})
}

func main() {
	greeter := &GreetServer{}
	mux := http.NewServeMux()
	path, handler := greetv1connect.NewGreetServiceHandler(greeter)
	mux.Handle(path, handler)

	log.Println("Ready - Server is running at localhost:8080")
	http.ListenAndServe(
		"0.0.0.0:8080",
		newCORS().Handler(h2c.NewHandler(mux, &http2.Server{})),
	)
}