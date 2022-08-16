package main

import (
	"eat-and-go/service/api"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"eat-and-go/config"
	log "eat-and-go/pkg/logger"
	v "eat-and-go/pkg/version"
	"eat-and-go/router"
	"eat-and-go/router/middleware"

	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"
)

var (
	version = pflag.BoolP("version", "v", false, "show version info.")
)

// pingServer pings the http server to make sure the router is working.
func pingServer() error {

	for i := 0; i < config.GetConfig().Service.MaxPingCount; i++ {
		// Ping the server by sending a GET request to `/health`.
		resp, err := http.Get(config.GetConfig().Service.Url + "/sd/health")
		if err == nil && resp.StatusCode == 200 {
			return nil
		}

		// Sleep for a second to continue the next ping.
		fmt.Printf("Waiting for the router, retry in 1 second.\n")
		time.Sleep(time.Second)
	}
	return errors.New("Cannot connect to the router.\n")
}

func main() {
	pflag.Parse()

	if *version {
		v := v.Get()
		marshalled, err := json.MarshalIndent(&v, "", "  ")
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}
		fmt.Println(string(marshalled))
		return
	}

	if err := config.SetConfig(); err != nil {
		panic(err)
	}

	if err := log.SetLog(); err != nil {
		panic(err)
	}

	gin.SetMode(config.GetConfig().Service.Runmode)
	//g := gin.New() // Create the Gin engine.
	g := router.InitEngine()
	// Routes.
	router.Load(g, middleware.Logging(), middleware.RequestID())
	api.InitSlack()
	// HealthCheck: Ping the server to make sure the router is working.
	go func() {
		if err := pingServer(); err != nil {
			fmt.Println("Starting...\n", err)
		}
		fmt.Printf("Started Successfully.\n")
		log.MetricsEmit("eat-and-go.pingServer", "", "eat-and-go Started Successfully.", true)
	}()

	fmt.Printf("Listening Address: %s \n", config.GetConfig().Service.Addr)
	log.MetricsEmit("eat-and-go.main",
		"",
		fmt.Sprintf("Listening Address: %s", config.GetConfig().Service.Addr),
		true)
	fmt.Printf(http.ListenAndServe(config.GetConfig().Service.Addr, g).Error())
	//api.SendAttachments()
}
