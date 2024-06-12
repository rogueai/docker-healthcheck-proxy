package main

import (
	"errors"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/docker/docker/client"
)

func main() {
	http.HandleFunc("/healthcheck", getHealthCheck)

	err := http.ListenAndServe(":3333", nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}

func getHealthCheck(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	log.Println("Performing healthcheck...")
	apiClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	defer func(apiClient *client.Client) {
		err := apiClient.Close()
		if err != nil {
			log.Println(err)
		}
	}(apiClient)

	options := container.ListOptions{}
	containers, err := apiClient.ContainerList(ctx, options)
	if err != nil {
		panic(err)
	}
	// holds the container id
	id := os.Getenv("HOSTNAME")
	checks := make(map[string]bool)
	deps := make([]string, 0)
	for _, ctr := range containers {
		ctr, err := apiClient.ContainerInspect(ctx, ctr.ID)
		if err != nil {
			continue
		}
		labels := ctr.Config.Labels
		if strings.HasPrefix(ctr.ID, id) {
			// this is the container we're running in, check for label depends-on
			dependsOn := labels["healthcheck.depends-on"]
			if dependsOn != "" {
				deps = append(deps, dependsOn)
			}
		} else {
			// check for docker compose
			service := labels["com.docker.compose.service"]
			if service != "" {
				enabled := labels["healthcheck.enable"]
				if enabled == "true" {
					if ctr.State.Health != nil {
						checks[service] = ctr.State.Health.Status == "healthy"
					} else {
						checks[service] = ctr.State.Status == "running"
					}
				}
			}
		}
	}
	healthy := true
	for _, dep := range deps {
		healthy = healthy && checks[dep]
	}
	if healthy {
		log.Println("healthy")
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, "healthy")
	} else {
		log.Println("unhealthy")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = io.WriteString(w, "unhealthy")
	}

}
