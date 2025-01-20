package main

import (
	"os"
	"time"

	"github.com/chyiyaqing/gmicro-user/config"
	"github.com/chyiyaqing/gmicro-user/internal/adapters/db"
	"github.com/chyiyaqing/gmicro-user/internal/adapters/grpc"
	"github.com/chyiyaqing/gmicro-user/internal/application/core/api"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/trace"
)

const (
	service     = "user"
	environment = "dev"
	id          = 4
)

type customLogger struct {
	formatter log.JSONFormatter
}


// Format(*Entry) ([]byte, error)
func (l *customLogger) Format(entry *log.Entry) ([]byte, error) {
	span := trace.SpanFromContext(entry.Context)
	entry.Data["trace_id"] = span.SpanContext().TraceID().String()
	entry.Data["span_id"] = span.SpanContext().SpanID().String()
	// Below injection is Just to understand what Context has
	entry.Data["Context"] = span.SpanContext()
	return l.formatter.Format(entry)
}

func init() {
	log.SetFormatter(&customLogger{
		formatter: log.JSONFormatter{
			FieldMap: log.FieldMap{
				"msg": "message",
			},
		},
	})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

func main() {
	dbAdapter, err := db.NewAdapter(config.GetSqliteDB())
	if err != nil {
		log.Fatalf("Failed to connect to database. Error: %v", err.Error())
	}

	application := api.NewApplication(dbAdapter)
	grpcAdapter := grpc.NewAdaptor(application, config.GetApplicationGrpcPort(), config.GetApplicationHttpPort(), config.GetJwtSecret(), time.Minute*time.Duration(config.GetJwtTokenDurationMinute()))
	grpcAdapter.Run()
}
