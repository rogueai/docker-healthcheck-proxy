package healthcheck

import (
	"log"

	"github.com/docker/docker/client"
)

/*
Close the client.
*/
func gracefulShutdown(apiClient *client.Client) {
	var err error = apiClient.Close()
	if err != nil {
		log.Println(err)
	}
}
