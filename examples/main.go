package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
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

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

func main() {

	inKubernetesCluster := flag.Bool("in-cluster", true, "Whether the application is running inside a Kubernetes cluster")
	kubernetesNamespace := flag.String("namespace", "default", "The Kubernetes namespace to use")
	kubeConfigPath := flag.String("kubeconfig", fmt.Sprintf("%s/.kube/config", homeDir()), "Absolute path to the kubeconfig file")
	flag.Parse()

	fmt.Println("Configuration:")
	fmt.Printf("  In Kubernetes cluster:         %t\n", *inKubernetesCluster)
	fmt.Printf("  Kubernetes namespace:          %s\n", *kubernetesNamespace)
	fmt.Printf("  Kubernetes configuration path: %s\n", *kubeConfigPath)
	fmt.Println("")

	kubernetesClient, err := configmap.NewKubernetesClient(*inKubernetesCluster, *kubeConfigPath)
	if err != nil {
		log.Fatal().Err(err).Msg("")
	}

	kHandler := configmap.NewHandler(*kubernetesClient, *kubernetesNamespace)

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
