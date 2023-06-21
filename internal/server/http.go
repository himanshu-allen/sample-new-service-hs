package server

import (
	"github.com/Allen-Career-Institute/go-kratos-commons/otel/middleware/metric"
	"github.com/Allen-Career-Institute/go-kratos-commons/otel/middleware/trace"
	v1 "github.com/Allen-Career-Institute/go-kratos-sample/api/user/v1"
	"github.com/Allen-Career-Institute/go-kratos-sample/internal/conf"
	"github.com/Allen-Career-Institute/go-kratos-sample/internal/healthcheck"
	"github.com/Allen-Career-Institute/go-kratos-sample/internal/request"
	"github.com/Allen-Career-Institute/go-kratos-sample/internal/service"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/validate"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/go-kratos/swagger-api/openapiv2"
	"go.opentelemetry.io/otel"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, user *service.UserService, logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			validate.Validator(),
			metric.Meter(otel.Meter("meter-go-kratos-sample")),
			trace.Trace(otel.Tracer("tracer-go-kratos-sample")),
			request.Populate(),
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)
	v1.RegisterUserHTTPServer(srv, user)
	openAPIhandler := openapiv2.NewHandler()
	healthHandler := healthcheck.NewHandler()

	srv.HandlePrefix("/q/", openAPIhandler)
	srv.HandlePrefix("/health/", healthHandler)

	return srv
}
