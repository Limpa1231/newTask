package web

import (
	"log"

	"github.com/labstack/echo-contrib/echoprometheus"
	echo "github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	legalEntitiesRequestsTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "legal_entities_requests_total",
			Help: "Total number of requests to /legal-entities endpoint",
		},
	)

	bankAccountsRequestsTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "bank_accounts_requests_total",
			Help: "Total number of requests to /bank-accounts endpoint",
		},
	)
)

func init() {
	prometheus.MustRegister(legalEntitiesRequestsTotal)
	prometheus.MustRegister(bankAccountsRequestsTotal)
}

func metricsMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := next(c)

		switch c.Path() {
		case "/legal-entities":
			legalEntitiesRequestsTotal.Inc()

		case "/bank-accounts":
			bankAccountsRequestsTotal.Inc()

		}

		return err
	}
}

func initMetricsRoutes(a *Web, e *echo.Echo) {
	e.Use(metricsMiddleware)

	e.Use(echoprometheus.NewMiddleware(a.Options.APP_NAME))

	customCounter := prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "custom_requests_total",
			Help: "How many HTTP requests processed, partitioned by status code and HTTP method.",
		},
	)

	if err := prometheus.Register(customCounter); err != nil {
		log.Fatal(err)
	}

	e.Use(echoprometheus.NewMiddlewareWithConfig(echoprometheus.MiddlewareConfig{
		AfterNext: func(c echo.Context, err error) {
			customCounter.Inc()
		},
	}))

	e.GET("/metrics", echoprometheus.NewHandler())
}
