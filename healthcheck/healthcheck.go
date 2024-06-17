package healthcheck

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/docker/docker/api/types/container"

	"github.com/docker/docker/client"
)

func GetHealthCheck(w http.ResponseWriter, r *http.Request) {
	var ctx context.Context = r.Context()

	log.Println("Performing healthcheck...")
	apiClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		log.Fatal(err)
	}
	defer gracefulShutdown(apiClient)

	var options container.ListOptions = container.ListOptions{}

	containers, err := apiClient.ContainerList(ctx, options)
	if err != nil {
		log.Fatal(err)
	}

	// holds the container id
	var id string = os.Getenv("HOSTNAME")
	var checks map[string]bool = make(map[string]bool)

	var deps []string = make([]string, 0)

	for _, ctr := range containers {
		ctr, err := apiClient.ContainerInspect(ctx, ctr.ID)

		// continue anyway
		if err != nil {
			continue
		}

		var labels map[string]string = ctr.Config.Labels
		if strings.HasPrefix(ctr.ID, id) {
			// this is the container we're running in, check for label depends-on
			var dependsOn string = labels["healthcheck.depends-on"]

			// https://stackoverflow.com/a/18595217
			if dependsOn != "" {
				deps = append(deps, dependsOn)
			}

		} else {
			// check for docker compose
			var service string = labels["com.docker.compose.service"]
			if service != "" {
				var enabled string = labels["healthcheck.enable"]

				if enabled == "true" {
					checks[service] = ctr.State.Status == "running"
					if ctr.State.Health != nil {
						checks[service] = ctr.State.Health.Status == "healthy"
					}
				}
			}
		}
	}
	var healthy bool = true
	for _, dep := range deps {
		healthy = healthy && checks[dep]
	}

	if healthy {
		log.Println("healthy")
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, "healthy")
		return
	}
	log.Println("unhealthy")
	w.WriteHeader(http.StatusBadRequest)
	_, _ = io.WriteString(w, "unhealthy")

}
