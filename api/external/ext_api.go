package external

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"demo_project/api/external/demo"
	"demo_project/api/internal"
	"demo_project/config"
	"demo_project/helpers"
	"demo_project/proto"
	"demo_project/static"
	"fmt"
	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/pkg/errors"
	"github.com/tmc/grpc-websocket-proxy/wsproxy"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	"net/http"
	"strings"
	log "github.com/sirupsen/logrus"
	"time"
	"github.com/gorilla/mux"
)

var (
	bind            string
	tlsCert         string
	tlsKey          string
	jwtSecret       string
	corsAllowOrigin string
)

func Setup(conf config.Config) error {
	bind = conf.DemoServer.ExternalAPI.Bind
	tlsCert = conf.DemoServer.ExternalAPI.TLSCert
	tlsKey = conf.DemoServer.ExternalAPI.TLSKey
	jwtSecret = conf.DemoServer.ExternalAPI.JWTSecret
	corsAllowOrigin = conf.DemoServer.ExternalAPI.CORSAllowOrigin

	return setupAPI(conf)
}

func setupAPI(conf config.Config) error {
	grpcOpts := helpers.GetgRPCServerOptions()
	grpcServer := grpc.NewServer(grpcOpts...)
	proto.RegisterDemoServiceServer(grpcServer, demo.NewDemoAPI())
	proto.RegisterGrpcServiceServer(grpcServer, internal.NewGrpcAPI())

	// setup the client http interface variable
	// we need to start the gRPC service first, as it is used by the
	// grpc-gateway
	var clientHTTPHandler http.Handler

	// switch between gRPC and "plain" http handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			if clientHTTPHandler == nil {
				w.WriteHeader(http.StatusNotImplemented)
				return
			}

			if corsAllowOrigin != "" {
				w.Header().Set("Access-Control-Allow-Origin", corsAllowOrigin)
				w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
				w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Grpc-Metadata-Authorization")

				if r.Method == "OPTIONS" {
					return
				}
			}

			clientHTTPHandler.ServeHTTP(w, r)
		}
	})

	// start the API server
	go func() {
		log.WithFields(log.Fields{
			"bind":     bind,
			"tls-cert": tlsCert,
			"tls-key":  tlsKey,
		}).Info("api/external: starting api server")

		if tlsCert == "" || tlsKey == "" {
			log.Fatal(http.ListenAndServe(bind, h2c.NewHandler(handler, &http2.Server{})))
		} else {
			log.Fatal(http.ListenAndServeTLS(
				bind,
				tlsCert,
				tlsKey,
				h2c.NewHandler(handler, &http2.Server{}),
			))
		}
	}()

	// give the http server some time to start
	time.Sleep(time.Millisecond * 100)

	// setup the HTTP handler
	clientHTTPHandler, err := setupHTTPAPI(conf)
	if err != nil {
		return err
	}

	return nil
}

func setupHTTPAPI(conf config.Config) (http.Handler, error) {
	r := mux.NewRouter()

	// setup json api handler
	jsonHandler, err := getJSONGateway(context.Background())
	if err != nil {
		return nil, err
	}

	log.WithField("path", "/api").Info("api/external: registering rest api handler and documentation endpoint")
	r.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		data, err := static.Asset("swagger/index.html")
		if err != nil {
			log.WithError(err).Error("get swagger template error")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(data)
	}).Methods("get")
	r.PathPrefix("/api").Handler(jsonHandler)

	// setup static file server
	r.PathPrefix("/").Handler(http.FileServer(&assetfs.AssetFS{
		Asset:     static.Asset,
		AssetDir:  static.AssetDir,
		AssetInfo: static.AssetInfo,
		Prefix:    "",
	}))

	return wsproxy.WebsocketProxy(r), nil
}

func getJSONGateway(ctx context.Context) (http.Handler, error) {
	// dial options for the grpc-gateway
	var grpcDialOpts []grpc.DialOption

	if tlsCert == "" || tlsKey == "" {
		grpcDialOpts = append(grpcDialOpts, grpc.WithInsecure())
	} else {
		b, err := ioutil.ReadFile(tlsCert)
		if err != nil {
			return nil, errors.Wrap(err, "read external api tls cert error")
		}
		cp := x509.NewCertPool()
		if !cp.AppendCertsFromPEM(b) {
			return nil, errors.Wrap(err, "failed to append certificate")
		}
		grpcDialOpts = append(grpcDialOpts, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
			// given the grpc-gateway is always connecting to localhost, does
			// InsecureSkipVerify=true cause any security issues?
			InsecureSkipVerify: true,
			RootCAs:            cp,
		})))
	}

	bindParts := strings.SplitN(bind, ":", 2)
	if len(bindParts) != 2 {
		log.Fatal("get port from bind failed")
	}
	apiEndpoint := fmt.Sprintf("localhost:%s", bindParts[1])

	mux := runtime.NewServeMux(runtime.WithMarshalerOption(
		runtime.MIMEWildcard,
		&runtime.JSONPb{
			EnumsAsInts:  false,
			EmitDefaults: true,
		},
	))

	if err := proto.RegisterDemoServiceHandlerFromEndpoint(ctx, mux, apiEndpoint, grpcDialOpts); err != nil {
		return nil, errors.Wrap(err, "register application handler error")
	}

	return mux, nil
}