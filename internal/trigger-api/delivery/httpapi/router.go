package httpapi

import (
	"bytes"
	"encoding/base64"
	"net/http"

	"github.com/go-chi/chi"
	chimid "github.com/go-chi/chi/middleware"
	"github.com/krostar/httpinfo"
	"github.com/krostar/httpw"
	"github.com/krostar/httpx/midx"
	"github.com/krostar/logger/logmid"
	"go.opencensus.io/trace"

	"github.com/krostar/r10k-trigger/internal/trigger-api/delivery/httpapi/handler"
	"github.com/krostar/r10k-trigger/internal/trigger-api/delivery/httpapi/middleware"
)

func (h *HTTP) initRouter(usecases Usecases) http.Handler {
	router, wrapper, tracer := h.defaultRouter()

	router.Route("/v1", func(v1 chi.Router) {
		v1.
			With(tracer.Trace(midx.TracerWithAlwaysSampler())).
			Route("/deploy", func(deploy chi.Router) {
				var gh = deploy
				if h.cfg.EnsureDeployerOrigin {
					gh = deploy.With(middleware.EnsureGithubOrigin(h.cfg.DeployerMACSecret))
				}
				gh.Post("/from-github", wrapper.Wrap(handler.DeployFromGithub(h.log, usecases)))
			})
	})

	// print a debug for each routes
	if err := chi.Walk(router,
		func(method string, route string, _ http.Handler, _ ...func(http.Handler) http.Handler) error {
			h.log.WithFields(map[string]interface{}{"method": method, "route": route}).Debug("enabled http route")
			return nil
		},
	); err != nil {
		h.log.WithError(err).Warn("unable to walk throught routes")
	}

	return router
}

func (h *HTTP) defaultRouter() (*chi.Mux, *httpw.Wrapper, *midx.Tracer) {
	var (
		wrapper = httpw.New(httpw.WithOnErrorCallback(onHandlerError))
		tracer  = midx.NewTracer(midx.TracerWithCallback(logAddTracerFields))
		router  = chi.NewRouter()
	)

	// set default handler
	var notFoundHandler = wrapper.WrapF(handler.NotFound)
	router.NotFound(notFoundHandler)
	router.MethodNotAllowed(notFoundHandler)

	// set middlewares
	router.Use(
		httpinfo.Record(httpinfo.WithRouteGetterFunc(routeFromRequest)),
		chimid.RealIP,
		logmid.New(h.log, logmid.WithDefaultFields()),
		midx.Recover(onHandlerPanic),
		chimid.Timeout(h.cfg.RequestTimeout),
		midx.Monitor(),
	)

	if h.statsEndpoint != "" && h.statsHandler != nil {
		router.
			With(tracer.Trace(midx.TracerWithProbabilitySampler((1 / 24.0)))).
			Get(h.statsEndpoint, h.statsHandler)
	}

	return router, wrapper, tracer
}

func routeFromRequest(r *http.Request) string {
	var routeCtx = chi.RouteContext(r.Context())

	if len(routeCtx.RoutePatterns) == 0 {
		return "404"
	}

	return r.Method + " " + routeCtx.RoutePattern()
}

// in addition to the classic panic handler
func onHandlerPanic(_ http.ResponseWriter, r *http.Request, reason interface{}, stack []byte) {
	var (
		ctx    = r.Context()
		bStack bytes.Buffer
	)

	base64.NewEncoder(base64.StdEncoding, &bStack).Write(stack) // nolint: errcheck, gosec

	logmid.AddFieldInContext(ctx, "panic", reason)
	logmid.AddFieldInContext(ctx, "stack", bStack.String())
}

// in addition to the classic error handler
func onHandlerError(r *http.Request, err error) {
	logmid.AddErrorInContext(r.Context(), err)
}

func logAddTracerFields(r *http.Request, span *trace.Span) {
	var (
		ctx      = r.Context()
		traceCtx = span.SpanContext()
	)

	logmid.AddFieldInContext(ctx, "trace-id", traceCtx.TraceID)
	logmid.AddFieldInContext(ctx, "trace-sampled", traceCtx.IsSampled())
}
