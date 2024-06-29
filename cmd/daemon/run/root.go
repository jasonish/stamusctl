package run

import (
	"context"
	"os"
	"os/signal"

	"stamus-ctl/cmd/daemon/run/compose"
	"stamus-ctl/internal/logging"

	"github.com/docker/distribution/context"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

// Commands
func RunCmd() *cobra.Command {
	_, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	logging.NewTraceProvider()

	// Create command
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Run daemon",
		RunE: func(cmd *cobra.Command, args []string) error {

			c := context.Background()
			_, span := logging.Tracer.Start(c, "main")
			defer span.End()

			gin.SetMode(gin.ReleaseMode)
			r := gin.New()

			logging.LoggerWithSpanContext(span.SpanContext()).Info("Setup middleware")
			r.Use(gin.Recovery())
			r.Use(otelgin.Middleware("service-name", otelgin.WithTracerProvider(logging.TracerProvider)))

			logging.LoggerWithSpanContext(span.SpanContext()).Info("Setup routes")
			r.GET("/ping", func(c *gin.Context) {
				logging.LoggerWithRequest(c.Request).Info("Ping")

				c.JSON(200, gin.H{
					"message": "pong",
				})
			})
			compose.NewCompose(r.Group("/compose"))

			// r.RunUnix("./daemon.sock")
			logging.LoggerWithSpanContext(span.SpanContext()).Info("Starting daemon")
			r.Run(":9000")
			return nil
		},
	}

	return cmd
}
