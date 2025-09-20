package main

import (
	"context"
	"fmt"
	"lavanilla/service/custom"
	"lavanilla/service/shopify"
	"log"
	"os"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
	"github.com/vektah/gqlparser/v2/ast"

	graph "lavanilla/graphql/backoffice"
)

var ginLambda *ginadapter.GinLambda

const defaultPort = "8192"

func ginHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return ginLambda.ProxyWithContext(ctx, request)
}

func graphqlHandler() gin.HandlerFunc {

	ctx := context.Background()
	awsS3Config, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("ap-southeast-1"),
	)
	if err != nil {
		log.Fatal(err)
	}
	s3Client := s3.NewFromConfig(awsS3Config)

	c := graph.Config{Resolvers: &graph.Resolver{
		CustomClient:    custom.NewClient(),
		ShopifyClient:   shopify.NewClient(),
		S3PresignClient: s3.NewPresignClient(s3Client),
		S3Client:        s3Client,
	}}

	h := handler.New(graph.NewExecutableSchema(c))
	h.AddTransport(transport.Options{})
	h.AddTransport(transport.POST{})

	// Add the introspection middleware.
	h.Use(extension.Introspection{})

	h.SetQueryCache(lru.New[*ast.QueryDocument](1000))
	h.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	h.AroundOperations(func(ctx context.Context, next graphql.OperationHandler) graphql.ResponseHandler {
		graphql.GetOperationContext(ctx).DisableIntrospection = false
		return next(ctx)
	})

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func main() {

	isInLambda := os.Getenv("AWS_LAMBDA_RUNTIME_API")
	if isInLambda != "" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	err := router.SetTrustedProxies(nil)
	if err != nil {
		log.Fatal("unable to set trusted proxies")
	}
	router.POST("/graphql-bo",
		graphqlHandler(),
	)

	if isInLambda != "" {
		log.Println("starting server")
		ginLambda = ginadapter.New(router)
		lambda.Start(ginHandler)
	} else {
		log.Printf("connect to http://localhost:%s", defaultPort)
		log.Fatal(router.Run(fmt.Sprintf(":%v", defaultPort)))
	}
}
