package health

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/realHoangHai/awesome/config"
	"github.com/realHoangHai/awesome/pkg/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
	"net/http"
	"time"
)

type (
	// HServer is the simple implementation of the Server
	HServer struct {
		checkers map[string]Checker
		ticker   *time.Ticker
		log      log.Logger

		server *health.Server
		cfg    *config.Config
	}

	// Option is a function to provide additional options for server
	Option func(*HServer)

	// Server provides health check services via both gRPC and HTTP
	// The implementation must follow protocol defined in https://github.com/grpc/grpc/blob/master/doc/health-checking.md
	Server interface {
		// HealthServer implements grpc_health_v1.HealthServer for general health check
		// and load balancing according to gRPC protocol.
		grpc_health_v1.HealthServer
		// Handler implements http.Handler for checking API via HTTP.
		http.Handler
		// Register registers the Server with the grpc.Server.
		Register(srv *grpc.Server)
		// Init initializes status, perform necessary setup and start a first health check
		// immediately overall status and all dependent services's status.
		Init(status Status) error
		// Close closes the underlying resources.
		// It sets all serving status to NOT_SERVING and configures the server to ignore all future status changes.
		Close() error
	}

	// StatusSetter is an interface to set status for a service according to gRPC Health Checkl protocol.
	StatusSetter interface {
		// SetStatus resets the serving status of a service or insert a new service entry into statusMap.
		// Use empty service name will set overall status.
		SetStatus(service string, status Status)
	}

	// Status is an alias of grpc_health_v1.HealthCheckResponse_ServingStatus
	Status = grpc_health_v1.HealthCheckResponse_ServingStatus
	// CheckRequest is an alias of grpc_health_v1.HealthCheckRequest
	CheckRequest = grpc_health_v1.HealthCheckRequest
	// CheckResponse is an alias of grpc_health_v1.HealthCheckResponse
	CheckResponse = grpc_health_v1.HealthCheckResponse
	// WatchServer is an alias of grpc_health_v1.Health_WatchServer
	WatchServer = grpc_health_v1.Health_WatchServer
)

// Status const defines short name for grpc_health_v1.HealthCheckResponse_ServingStatus.
const (
	StatusUnknown        = grpc_health_v1.HealthCheckResponse_UNKNOWN
	StatusServing        = grpc_health_v1.HealthCheckResponse_SERVING
	StatusNotServing     = grpc_health_v1.HealthCheckResponse_NOT_SERVING
	StatusServiceUnknown = grpc_health_v1.HealthCheckResponse_SERVICE_UNKNOWN
)

const (
	// OverallServiceName is service name of server's overall status.
	OverallServiceName = ""
)

var (
	// force HServer implements required interfaces.
	_ Server       = &HServer{}
	_ StatusSetter = &HServer{}
)

func NewServer(m map[string]Checker, opts ...Option) *HServer {
	hs := &HServer{
		checkers: m,
		server:   health.NewServer(),
		cfg:      &config.Config{},
	}
	for _, opt := range opts {
		opt(hs)
	}
	// default if not set
	if hs.cfg.Health.Interval == 0 {
		hs.cfg.Health.Interval = time.Second * 60
	}
	if hs.log == nil {
		hs.log = log.Root()
	}
	if hs.cfg.Health.Timeout == 0 {
		hs.cfg.Health.Timeout = time.Second * 1
	}
	hs.ticker = time.NewTicker(hs.cfg.Health.Interval)
	return hs
}

// Register implements health.Server.
func (hs *HServer) Register(srv *grpc.Server) {
	grpc_health_v1.RegisterHealthServer(srv, hs)
}

// Init implements health.Server.
func (hs *HServer) Init(status Status) error {
	hs.server.SetServingStatus(OverallServiceName, status)
	// if there is no dependent services, don't need to do anything else.
	if len(hs.checkers) == 0 {
		return nil
	}
	// if there are dependent services, set overall status and all dependent services
	// to NotServing as we don't know their status yet.
	hs.server.SetServingStatus(OverallServiceName, StatusNotServing)
	for name := range hs.checkers {
		hs.server.SetServingStatus(name, StatusNotServing)
	}
	// start a first check immediately.
	hs.checkAll()
	// schedule the check
	go func() {
		for range hs.ticker.C {
			hs.checkAll()
		}
	}()
	return nil
}

func (hs *HServer) checkAll() {
	logger := hs.log.Fields(log.CorrelationID, uuid.New().String())
	bg := time.Now()
	overall := StatusServing
	for service, check := range hs.checkers {
		state := StatusServing
		if err := hs.check(service, check); err != nil {
			overall = StatusNotServing
			state = StatusNotServing
			logger.Infof("health check failed, service: %s, err: %v", service, err)
		}
		hs.server.SetServingStatus(service, state)
	}
	hs.server.SetServingStatus(OverallServiceName, overall)
	logger.Fields("status", overall, "duration", time.Since(bg)).Info("health check completed")
}

func (hs *HServer) check(service string, check Checker) error {
	ctx, cancel := context.WithTimeout(context.Background(), hs.cfg.Health.Timeout)
	defer cancel()
	ch := make(chan error)
	go func() {
		ch <- check.CheckHealth(ctx)
	}()
	select {
	case err := <-ch:
		return err
	case <-ctx.Done():
		return errors.New("health: check exceed timeout")
	}
}

// Check implements health.Server.
func (hs *HServer) Check(ctx context.Context, req *CheckRequest) (*CheckResponse, error) {
	return hs.server.Check(ctx, req)
}

// Watch implements health.Server.
func (hs *HServer) Watch(req *CheckRequest, srv WatchServer) error {
	return hs.server.Watch(req, srv)
}

// ServeHTTP implements health.Server.
func (hs *HServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	service := r.URL.Query().Get("service")
	data := make(map[string]interface{})
	check := func(service string) Status {
		resp, err := hs.Check(r.Context(), &CheckRequest{Service: service})
		if status.Code(err) == codes.NotFound {
			return StatusServiceUnknown
		}
		if err != nil {
			return StatusNotServing
		}
		return resp.Status
	}
	// overall - check all dependent services.
	if service == OverallServiceName {
		overall := check(OverallServiceName)
		services := make(map[string]Status)
		for svc := range hs.checkers {
			state := check(svc)
			services[svc] = state
			if state != StatusServing {
				overall = StatusNotServing
			}
		}
		data["status"] = overall
		data["services"] = services
	} else {
		data["status"] = check(service)
	}
	marshal, err := json.Marshal(data)
	if err != nil {
		marshal = []byte(fmt.Sprintf(`{"status":%d}`, StatusNotServing))
		hs.log.Errorf("failed to marshal json: %v", err)
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(marshal)
}

func (hs *HServer) SetStatus(service string, status Status) {
	hs.server.SetServingStatus(service, status)
}

func (hs *HServer) Close() error {
	if hs.ticker != nil {
		hs.ticker.Stop()
	}
	if hs.server != nil {
		hs.server.Shutdown()
	}
	return nil
}

//Interval is an option to set interval for health check.
func Interval(d time.Duration) Option {
	return func(srv *HServer) {
		srv.cfg.Health.Interval = d
	}
}

// Timeout is an option to set timeout for each service health check.
func Timeout(d time.Duration) Option {
	return func(srv *HServer) {
		srv.cfg.Health.Timeout = d
	}
}

// Logger is an option to set logger for the health check server.
func Logger(l log.Logger) Option {
	return func(srv *HServer) {
		srv.log = l
	}
}
