package middlewares

import (
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/saeed903/microservice_eventsourcing_package/config"
	"github.com/saeed903/microservice_eventsourcing_package/pkg/logger"
)

type MiddlewareMetricCb func(err error)

type MiddlewareManager interface {
	RequestLoggerMiddleware(next echo.HandlerFunc) echo.HandlerFunc
}

type middlewareManager struct {
	log      logger.Logger
	cfg      config.Config
	metricCb MiddlewareMetricCb
}

func NewMiddlewareManager(log logger.Logger, cfg config.Config, metricCb MiddlewareMetricCb) *middlewareManager {
	return &middlewareManager{log: log, cfg: cfg, metricCb: metricCb}
}

func (mw *middlewareManager) RequestLoggerMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		start := time.Now()
		err := next(ctx)

		req := ctx.Request()
		res := ctx.Response()
		status := res.Status
		size := res.Size
		s := time.Since(start)

		if !mw.checkIgnoredURI(ctx.Request().RequestURI, mw.cfg.Http.IgnorLogUrls) {
			mw.log.HttpMiddlewareAccessLogger(req.Method, req.URL.String(), status, size, s)
		}

		mw.metricCb(err)
		return err
	}
}

func (mw *middlewareManager) checkIgnoredURI(requestURI string, uriList []string) bool {
	for _, v := range uriList {
		if strings.Contains(requestURI, v) {
			return true
		}
	}
	return false
}
