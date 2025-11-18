package httpadapter

import (
	"context"
	"fmt"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"github.com/voita-pro/http-grpc-api/api"
	"github.com/voita-pro/http-grpc-api/pkg/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// NewMux returns a basic http mux with health endpoint.
func NewMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	return mux
}

// NewGatewayMux creates a grpc-gateway mux that proxies REST to a gRPC endpoint.
func NewGatewayMux(ctx context.Context, grpcEndpoint string, dialOpts ...grpc.DialOption) (*runtime.ServeMux, error) {
	mux := runtime.NewServeMux()
	if len(dialOpts) == 0 {
		dialOpts = []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	}
	if err := pb.RegisterServicePBHandlerFromEndpoint(ctx, mux, grpcEndpoint, dialOpts); err != nil {
		return nil, err
	}
	return mux, nil
}

// MountSwagger mounts handlers for serving the swagger JSON and a simple UI.
func MountSwagger(mux *http.ServeMux, httpHostPort string) {
	// SWAGGER SERVER
	swgFs := http.StripPrefix("/openapi/",
		http.FileServer(http.FS(api.EmbedSwagger)),
	)
	// Proxy to static files in the embed directory
	mux.Handle("/openapi/", swgFs)
	// Full host path to swagger JSON
	swgPath := fmt.Sprintf(`http://%s/openapi/service.swagger.json`, httpHostPort)
	hswg := httpSwagger.Handler(
		httpSwagger.URL(swgPath),
	)
	mux.Handle("/swagger/", hswg)
}
