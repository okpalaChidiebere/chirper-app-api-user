package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	grpchealth "github.com/bufbuild/connect-grpchealth-go"
	grpcreflect "github.com/bufbuild/connect-grpcreflect-go"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	api "github.com/okpalaChidiebere/chirper-app-api-user/api"
	"github.com/okpalaChidiebere/chirper-app-api-user/config"
	usersservice "github.com/okpalaChidiebere/chirper-app-api-user/v0/users/business_logic"
	usersrepo "github.com/okpalaChidiebere/chirper-app-api-user/v0/users/data_access"
	imagefilterconnect "github.com/okpalaChidiebere/chirper-app-gen-protos/image_filter/v1/image_filterv1connect"
)

var (
	mConfig    = config.NewConfig()

	imageFilterServiceName   string 
	imageFilterServicePort   string 
	ServerPort   = 8000

	cfg aws.Config
	err error
)

func init() {
	imageFilterServiceName = config.MustGetenv("IMAGE_FILTER_SERVICE_NAME")
	imageFilterServicePort = config.MustGetenv("IMAGE_FILTER_SERVICE_PORT")
}

func main() {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	if mConfig.IsLocal() {
		cfg, err = awsconfig.LoadDefaultConfig(ctx, 
		awsconfig.WithSharedConfigProfile(mConfig.Aws.Aws_profile))
		if err != nil {
			log.Fatalf("unable to load local SDK config, %v\n", err)
		}
	} else {
		cfg, err = awsconfig.LoadDefaultConfig(ctx, awsconfig.WithRegion(mConfig.Aws.Aws_region))
		if err != nil {
			log.Fatalf("unable to load SDK config, %v\n", err)
		}
	}

	imageFilterServiceClient := imagefilterconnect.NewImagefilterServiceClient(
		http.DefaultClient,
		//we are using the service name because its more reliable than using the exact IP address of the container running the service
		//IP address can change when the container goes down and re-generate
		//If you are running locally you can inspect the chirper-app-image-filter-service container running and go to the `NetworkSettings` to see all network details
		fmt.Sprintf("http://%s:%s/", imageFilterServiceName, imageFilterServicePort),
	)
	dynamodbClient := dynamodb.NewFromConfig(cfg)

	usersRepo := usersrepo.NewDynamoDbRepo(dynamodbClient, mConfig.Dev.UserTable)

	usersService := usersservice.New(usersRepo)

	port := os.Getenv("PORT")
	if port == "" {
		port = strconv.Itoa(ServerPort)
	}

	s := api.Servers{
		UserServer: api.NewUserServer(usersService, imageFilterServiceClient),
	}

	mux := http.NewServeMux()

	apiServer := s.NewAPIServer(mux)

	sh := s.GetAllServiceHandlers()

	apiServer.RegisterEndpoints(sh)

	var services []string
	for key := range sh {
		services = append(services, key)
	}
	reflector := grpcreflect.NewStaticReflector(services...)
	mux.Handle(grpcreflect.NewHandlerV1(reflector))
  	mux.Handle(grpcreflect.NewHandlerV1Alpha(reflector))

	//add support for gRPC style health checks and http style health checks as well
	//documentation: https://github.com/bufbuild/connect-grpchealth-go#readme
	checker := grpchealth.NewStaticChecker(services...)
	mux.Handle(grpchealth.NewHandler(checker))

	httpLis, err := net.Listen("tcp", "0.0.0.0:"+port)
	if err != nil {
		log.Printf("HTTP server: failed to listen: error %v", err)
		os.Exit(2)
	}

	httpServer := &http.Server{
		Handler: h2c.NewHandler(mux, &http2.Server{}),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 20 * time.Second,
	}

	go func() {
		log.Printf("server listening at %v", httpLis.Addr())
		err = httpServer.Serve(httpLis)
		if errors.Is(err, http.ErrServerClosed) {
			fmt.Println("server closed")
		} else if err != nil {
			panic(err)
		}
	}()

	// Listen for the interrupt signal.
	<-ctx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()
	fmt.Println("shutting down gracefully, press Ctrl+C again to force")

	// Perform application shutdown with a maximum timeout of 20 seconds.
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(timeoutCtx); err != nil {
		fmt.Println(err)
	}
}
