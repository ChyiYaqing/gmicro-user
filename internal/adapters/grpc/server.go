package grpc

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	userv1 "github.com/chyiyaqing/gmicro-proto/golang/user/v1"
	"github.com/chyiyaqing/gmicro-user/config"
	"github.com/chyiyaqing/gmicro-user/internal/adapters/auth"
	"github.com/chyiyaqing/gmicro-user/internal/adapters/interceptor"
	"github.com/chyiyaqing/gmicro-user/internal/adapters/role"
	"github.com/chyiyaqing/gmicro-user/internal/ports"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
)

var allowedHeaders = map[string]struct{}{
	"x-request-id": {},
}

type Adapter struct {
	api        ports.APIPort
	port       int
	gwPort     int
	server     *grpc.Server
	jwtManager *auth.JWTManager
	userv1.UnimplementedUserServer
}

func NewAdaptor(api ports.APIPort, port int, gwPort int, jwtSecret string, tokenDuration time.Duration) *Adapter {
	// auth
	if len(jwtSecret) == 0 {
		log.Fatal("jwtSecret must be provided")
	}
	jwtManager, err := auth.NewJWTManager(jwtSecret, tokenDuration)
	if err != nil {
		log.Fatalf("failed to create JWT manager: %v", err)
	}
	return &Adapter{
		api:        api,
		port:       port,
		gwPort:     gwPort,
		jwtManager: jwtManager,
	}
}

func (a *Adapter) Run() {
	// Start gRPC Server
	go a.runGRPCServer()
	// Start gRPC-Gateway server
	a.runGatewayServer()
}

func (a *Adapter) runGRPCServer() {
	var err error

	authInterceptor, err := interceptor.NewAuthInterceptor(a.jwtManager, role.AccessibleRoles())
	if err != nil {
		log.Fatalf("failed to create auth interceptor: %v", err)
	}
	// start listening on port: for a tcp connection
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		log.Fatalf("failed to listen on port %d, error: %v", a.port, err)
	}
	var opts []grpc.ServerOption
	opts = append(opts,
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
		grpc.ChainUnaryInterceptor(
			interceptor.UnaryServerInterceptor,  // for Unary RPC
			authInterceptor.UnaryAuthMiddleware, // auth
		),
		grpc.StreamInterceptor(interceptor.StreamServerInterceptor), // for Stream RPC
	)
	// create new gRPC server
	grpcServer := grpc.NewServer(opts...)
	a.server = grpcServer

	// register the GreeterServerImpl on the gRPC server
	userv1.RegisterUserServer(grpcServer, a)

	if config.GetEnv() == "development" {
		reflection.Register(grpcServer)
	}

	log.Printf("starting user service on port %d...", a.port)
	// the gRPC server
	if err := grpcServer.Serve(listen); err != nil {
		log.Fatalf("failed to serve grpc on port %d", a.port)
	}
}

func (a *Adapter) runGatewayServer() {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// creating mux for gRPC gateway, this will multiplex or route request different gRPC service
	mux := runtime.NewServeMux(
		// convert header in response(going from gateway) from metadata received.
		runtime.WithOutgoingHeaderMatcher(isHeaderAllowed),

		// WithMetadata handles incoming HTTP headers
		// HTTP headers extracted here will be sent to the gRPC service in metedata.
		runtime.WithMetadata(func(ctx context.Context, request *http.Request) metadata.MD {
			header := request.Header.Get("Authrization")
			// send all the headers received from the client
			md := metadata.Pairs("auth", header)
			return md
		}),

		runtime.WithErrorHandler(
			func(
				ctx context.Context, // hold metadata
				mux *runtime.ServeMux, //
				marshaler runtime.Marshaler, // converts protocol buffer response to JSON response
				writer http.ResponseWriter,
				request *http.Request,
				err error,
			) { // error sent by the gRPC service
				// creating a new HTTPStatusError with a custom status, and passing error
				newError := runtime.HTTPStatusError{
					HTTPStatus: 400,
					Err:        err,
				}
				// using default handler to do the rest of heavy lifting of marshaling error and adding headers
				runtime.DefaultHTTPErrorHandler(ctx, mux, marshaler, writer, request, &newError)
			}))
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	// setting up a dail up for gRPC service by specifying endpoint/target url
	err := userv1.RegisterUserHandlerFromEndpoint(ctx, mux, fmt.Sprintf("localhost:%d", a.port), opts)
	if err != nil {
		log.Fatalf("Failed to register gateway: %v", err)
	}

	// Start HTTP server (and proxy calls to gRPC server)
	gwServer := &http.Server{
		Addr:         fmt.Sprintf(":%d", a.gwPort),
		Handler:      interceptor.WithLoggerInterceptor(mux),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	serverError := make(chan error, 1)

	// running the HTTP server in a separate groutine, so that it doesn't block the main thread
	// the main goroutine to handle shutdown signal
	go func() {
		log.Printf("starting gRPC-Gateway server on port %d...", a.gwPort)
		if err := gwServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Failed to serve gRPC-Gateway: %v", err)
			serverError <- err
		}
	}()

	stop := make(chan os.Signal, 1)
	// notify the stop channel for SIGINT(Ctrl+C) and SIGTERM
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverError:
		log.Printf("Server error: %v", err)
	case sig := <-stop:
		log.Printf("Received shutdown signal: %v", sig)
	}

	log.Println("gRPC-Gateway server is shutting done...")

	ctx, cancel = context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// closing database connections, saving state, cleaning up temporary files
	// TODO: Add your cleanup code here

	// Shutdown method stops accepting new connections immediately while allowing existing requests to complete
	if err := gwServer.Shutdown(ctx); err != nil {
		log.Fatalf("Failed to shutdown gRPC-Gateway server: %v", err)
		return
	}
	log.Println("gRPC-Gateway Server exited properly")
}

func (a *Adapter) Stop() {
	a.server.Stop()
}

func isHeaderAllowed(s string) (string, bool) {
	// check if allowedHeaders contain the header
	if _, isAllowed := allowedHeaders[s]; isAllowed {
		// send uppercase header
		return strings.ToUpper(s), true
	}
	// if not in the allowed header, don't send the header
	return s, false
}
