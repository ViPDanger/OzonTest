package app

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/ViPDanger/OzonTest/internal/infastructure/mongodb"
	"github.com/ViPDanger/OzonTest/internal/interfaces/gRPC"
	"github.com/ViPDanger/OzonTest/internal/interfaces/handlers"
	"github.com/ViPDanger/OzonTest/internal/interfaces/handlers/middleware"
	"github.com/ViPDanger/OzonTest/internal/proto"
	"github.com/ViPDanger/OzonTest/internal/usecase"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

const CircuitBreakerMax = 10000           // число при котором Circuit Breaker прекратит принимать запросы
const CircuitBreakerMin = 500             // число при котором Circuit Breaker возобновит приём запросов
const TimeouterMaxTime = 10 * time.Second // время до автоматического TimeOut <-ctx.Done

type Application struct {
	ginServer  *http.Server
	grpcServer *grpc.Server
}

func (a *Application) GetGinServer() *http.Server {
	return a.ginServer
}

func (a *Application) GetGRPCServer() *grpc.Server {
	return a.grpcServer
}

func Run(ctx context.Context, host string, mongoURI string, database string, user string, password string, grpcAdress string) (*Application, error) {
	//	MONGO SETUP
	cred := options.Credential{
		Username: user,
		Password: password,
	}
	clientOpts := options.Client().
		ApplyURI(mongoURI).
		SetAuth(cred)

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, fmt.Errorf("Run()/%w", err)
	}
	if err := client.Ping(pingCtx, nil); err != nil {
		return nil, fmt.Errorf("Run()/%w", err)
	}
	db := client.Database(database)

	//===================GIN setup================
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(
		middleware.NewCurcuitBreaker(CircuitBreakerMax, CircuitBreakerMin).CircuitBreakerHandler,
		middleware.NewRetrier().RetryHandler,
		middleware.NewTimeouter(TimeouterMaxTime).TimeoutHandler)
	usecase := usecase.NewValCursUseCase(mongodb.NewValCursRepository(db))
	handler := handlers.NewValCursHandler(usecase)
	r.GET("/", handler.GetByDateAndName)
	//===============gRPC server setup=============
	grpcServer := grpc.NewServer()
	proto.RegisterMockXMLDailyServer(
		grpcServer,
		gRPC.NewMockXMLDailyServer(usecase, handler))
	// GRACEFULL SHUTDOWN CTX---------
	ginServer := &http.Server{Addr: host, Handler: r.Handler()}
	ctx, cancel = context.WithCancel(ctx)
	go func() {
		defer cancel()
		fmt.Printf("app is started on %v\n", host)
		if err := ginServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("server error: %v\n", err)
		}

	}()
	go func() {
		defer cancel()
		lis, err := net.Listen("tcp", grpcAdress)
		if err != nil {
			fmt.Printf("failed to listen: %v\n", err)
			return
		}
		if err := grpcServer.Serve(lis); err != nil {
			fmt.Printf("failed to run gRPC: %v\n", err)
		}
	}()
	go func() {
		<-ctx.Done()
		ginServer.Close()
		grpcServer.GracefulStop()
		time.Sleep(1 * time.Second)
		fmt.Println("app is closed")
	}()
	return &Application{ginServer: ginServer, grpcServer: grpcServer}, nil
}
