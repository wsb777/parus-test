package logger

import (
	"net/http"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Loggers struct {
	App     *zap.Logger
	Service *zap.Logger
	HTTP    *zap.Logger
	DB      *zap.Logger
}

func New(env string) (*Loggers, error) {
	encoderCfg := zap.NewProductionEncoderConfig()

	encoderCfg.TimeKey = "@timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderCfg.EncodeCaller = zapcore.ShortCallerEncoder

	var encoder zapcore.Encoder

	if env == "production" {
		encoder = zapcore.NewJSONEncoder(encoderCfg)
	} else {
		encoder = zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	}

	ws := zapcore.AddSync(os.Stdout)

	core := zapcore.NewCore(
		encoder,
		ws,
		zap.InfoLevel,
	)

	base := zap.New(core, zap.AddCaller())

	return &Loggers{
		Service: base.With(zap.String("layer", "service")),
		HTTP:    base.With(zap.String("layer", "http")),
		App:     base.With(zap.String("layer", "app")),
		DB:      base.With(zap.String("layer", "db")),
	}, nil
}

func HTTPLoggerMiddleware(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			start := time.Now()

			next.ServeHTTP(w, r)

			logger.Info("http_request",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.Duration("latency", time.Since(start)),
				zap.String("ip", r.RemoteAddr),
			)
		})
	}
}
