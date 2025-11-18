package app

import (
	"context"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
	fb "github.com/voita-pro/http-grpc-api/internal/adapters/firebase"
	grpcadapter "github.com/voita-pro/http-grpc-api/internal/adapters/grpc"
	httpadapter "github.com/voita-pro/http-grpc-api/internal/adapters/http"
	"github.com/voita-pro/http-grpc-api/internal/config"
	"github.com/voita-pro/http-grpc-api/internal/domain"
	helperjwt "github.com/voita-pro/http-grpc-api/internal/helpers/jwt"
	"github.com/voita-pro/http-grpc-api/internal/repository"
	"github.com/voita-pro/http-grpc-api/internal/repository/pgdb"
	"github.com/voita-pro/http-grpc-api/internal/usecase"
)

type Service struct {
	ctx context.Context
	cfg *config.Cfg
}

func NewApp(ctx context.Context) (*Service, error) {
	conf, err := config.Init()
	if err != nil {
		return nil, err
	}

	return &Service{
		ctx: ctx,
		cfg: conf,
	}, nil
}

func (a *Service) Run() error {
	// DB connection (optional for stub)
	db, err := repository.Connect(a.ctx, &a.cfg.DB)
	if err != nil {
		return err
	}
	defer db.Close()

	// Main business logic and firebase initiate
	var ucSrv domain.Usecase
	var authSrv domain.Authenticator
	if a.cfg.Firebase.WebAPIKey != "" ||
		a.cfg.Firebase.CredsFile != "" ||
		a.cfg.Firebase.ProjectID != "" {
		authSrv, err = fb.New(a.ctx, &a.cfg.Firebase)
		if err != nil {
			log.Error().Err(err).Msg("failed to init firebase; falling back to stub auth")
		} else {
			log.Info().Msg("firebase auth initialized")
		}
	} else {
		log.Info().Msg("filed to read firebase config; using JWT auth")
		authSrv, err = helperjwt.New(a.ctx)
	}

	ucSrv = usecase.New(a.ctx, authSrv, pgdb.New(db.Pool))

	//gRPC server
	grpcHostPort := net.JoinHostPort(a.cfg.GRPC.Host, strconv.Itoa(a.cfg.GRPC.Port))
	grpcLis, err := net.Listen("tcp", grpcHostPort)
	if err != nil {
		return err
	}
	grpcServer := grpcadapter.NewGRPCServer(authSrv, ucSrv)
	go func() {
		// Start GRPC server
		log.Info().Str("addr", grpcHostPort).Msg("gRPC server listening")
		if err := grpcServer.Serve(grpcLis); err != nil {
			log.Fatal().Err(err).Msg("gRPC serve failed")
		}
	}()

	// REST gateway and HTTP server
	mux := httpadapter.NewMux()
	gwCtx, cancel := context.WithCancel(a.ctx)
	defer cancel()
	httpHostPort := net.JoinHostPort(a.cfg.HTTP.Host, strconv.Itoa(a.cfg.HTTP.Port))
	gw, err := httpadapter.NewGatewayMux(gwCtx, grpcHostPort)
	if err != nil {
		return err
	}
	mux.Handle("/", gw)
	var sMux http.Handler
	if a.cfg.Debug {
		// swagger
		httpadapter.MountSwagger(mux, httpHostPort)
		sMux = allowCORS(mux)
	} else {
		sMux = mux
	}
	httpSrv := &http.Server{
		Addr:    httpHostPort,
		Handler: sMux,
	}
	go func() {
		// Start HTTP server
		log.Info().Str("addr", httpHostPort).Msg("HTTP server listening")
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("http serve failed")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	<-a.ctx.Done()
	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), a.cfg.ShutdownTimeout)
	defer cancelShutdown()
	err = httpSrv.Shutdown(shutdownCtx)
	if err != nil {
		log.Err(err).Msg("failed to gracefully shutdown the HTTP server")
	}
	grpcServer.GracefulStop()
	log.Info().Msg("servers stopped")
	return nil
}

// allowCORS allows Cross Origin Resource Sharing from any origin.
// Don't do this without consideration in production systems.
func allowCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			if r.Method == "OPTIONS" && r.Header.Get("Access-Control-Request-Method") != "" {
				preflightHandler(w, r)
				return
			}
		}
		h.ServeHTTP(w, r)
	})
}

// preflightHandler adds the necessary headers in order to serve
// CORS from any origin using the methods "GET", "HEAD", "POST", "PUT", "DELETE"
// We insist, don't do this without consideration in production systems.
func preflightHandler(w http.ResponseWriter, r *http.Request) {
	headers := []string{"Content-Type", "Accept", "Authorization"}
	w.Header().Set("Access-Control-Allow-Headers", strings.Join(headers, ","))
	methods := []string{"GET", "HEAD", "POST", "PUT", "DELETE"}
	w.Header().Set("Access-Control-Allow-Methods", strings.Join(methods, ","))
}
