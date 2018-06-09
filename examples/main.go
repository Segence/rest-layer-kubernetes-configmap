package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/justinas/alice"
	"github.com/rs/cors"
	"github.com/rs/rest-layer/resource"
	"github.com/rs/rest-layer/rest"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
	"github.com/rs/zerolog/log"

	configmap "github.com/segence/rest-layer-kubernetes-configmap"
)

func main() {

	kubeconfig := flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	flag.Parse()

	kubeConfig := configmap.GetKubeConfig(*kubeconfig)
	log.Info().Msgf("%s", kubeConfig)

	kubernetesClient, err := configmap.LoadClientOutOfCluster(kubeConfig)
	if err != nil {
		log.Fatal().Err(err).Msg("")
	}

	kHandler := configmap.NewHandler(*kubernetesClient, "test")

	index := resource.NewIndex()

	index.Bind("config-map", configmap.ConfigMapSchema, kHandler, resource.Conf{
		AllowedModes: resource.ReadWrite,
	})

	api, err := rest.NewHandler(index)
	if err != nil {
		log.Fatal().Msgf("Invalid API configuration: %s", err)
	}

	c := alice.New()

	// Install a logger
	c = c.Append(hlog.NewHandler(log.With().Logger()))
	c = c.Append(hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
		hlog.FromRequest(r).Info().
			Str("method", r.Method).
			Str("url", r.URL.String()).
			Int("status", status).
			Int("size", size).
			Dur("duration", duration).
			Msg("")
	}))
	c = c.Append(hlog.RequestHandler("req"))
	c = c.Append(hlog.RemoteAddrHandler("ip"))
	c = c.Append(hlog.UserAgentHandler("ua"))
	c = c.Append(hlog.RefererHandler("ref"))
	c = c.Append(hlog.RequestIDHandler("req_id", "Request-Id"))
	resource.LoggerLevel = resource.LogLevelDebug
	resource.Logger = func(ctx context.Context, level resource.LogLevel, msg string, fields map[string]interface{}) {
		zerolog.Ctx(ctx).WithLevel(zerolog.Level(level)).Fields(fields).Msg(msg)
	}

	c = c.Append(cors.New(cors.Options{OptionsPassthrough: true}).Handler)

	http.Handle("/api/", http.StripPrefix("/api/", c.Then(api)))

	fmt.Println("Serving API on http://localhost:8080")

	if err := http.ListenAndServe("localhost:8080", nil); err != nil {
		log.Fatal().Err(err).Msg("")
	}
}
