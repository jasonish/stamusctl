package run

import (
	"context"
	"os"
	"os/signal"

	"stamus-ctl/cmd/daemon/run/compose"
	docs "stamus-ctl/cmd/docs"
	"stamus-ctl/internal/logging"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	ginprometheus "github.com/zsais/go-gin-prometheus"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewPrometheusServer(ctx context.Context) {
	engineProm := gin.New()

	p := ginprometheus.NewPrometheus("gin")
	p.Use(engineProm)

	logging.LoggerWithContextToSpanContext(ctx).Info("Starting prometheus endpoint")
	engineProm.Run(":9001")
}

// Ping godoc
// @Summary ping example
// @Schemes
// @Description do ping
// @Tags example
// @Accept json
// @Produce json
// @Success 200 {string} Helloworld
// @Router /ping [get]
func ping(c *gin.Context) {
	logging.LoggerWithRequest(c.Request).Info("Ping")

	c.JSON(200, gin.H{
		"message": "pong",
	})
}

// @title           Swagger Stamusd API
// @version         1.0
// @description     Stamus daemon server.

// @BasePath /api/v1
func RunCmd() *cobra.Command {
	_, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	logging.NewTraceProvider()

	cmd := &cobra.Command{
		Use:   "run",
		Short: "Run daemon",
		RunE: func(cmd *cobra.Command, args []string) error {

			c := context.Background()
			ctx, span := logging.Tracer.Start(c, "main")
			defer span.End()
			gin.SetMode(gin.ReleaseMode)
			r := gin.New()

			go NewPrometheusServer(ctx)

			logging.LoggerWithSpanContext(span.SpanContext()).Info("Setup middleware")
			r.Use(gin.Recovery())
			r.Use(otelgin.Middleware("stamusd", otelgin.WithTracerProvider(logging.TracerProvider)))

			logging.LoggerWithSpanContext(span.SpanContext()).Info("Setup routes")
			v1 := r.Group("/api/v1")
			{
				v1.GET("/ping", ping)
				compose.NewCompose(v1)
			}

			logging.LoggerWithSpanContext(span.SpanContext()).Info("Setup swagger")
			docs.SwaggerInfo.BasePath = "/api/v1"
			r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

			logging.LoggerWithSpanContext(span.SpanContext()).Info("Starting daemon")
			// r.RunUnix("./daemon.sock")
			r.Run(":9000")
			return nil
		},
	}

	return cmd
}
