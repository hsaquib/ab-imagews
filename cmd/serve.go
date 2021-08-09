/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"github.com/hsaquib/ab-imagews/api"
	"github.com/hsaquib/ab-imagews/config"
	"github.com/hsaquib/ab-imagews/service"
	"github.com/hsaquib/ab-imagews/utils"
	rLog "github.com/hsaquib/rest-log"
	"github.com/spf13/cobra"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: serve,
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	//serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}

func serve(cmd *cobra.Command, args []string) error {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("starting server")
	app, err := startServer()
	if err != nil {
		log.Println(err.Error())
		return err
	}

	err = stopServer(app)
	if err != nil {
		return err
	}

	return nil
}

func startServer() (*http.Server, error) {
	err := config.LoadConfig()
	cfg := config.GetConfig()
	rLog.Init(cfg.Env == utils.DEV_ENV, "image-server")
	if err != nil {
		log.Println("could not load one or more config")
		return nil, err
	}
	//rLogger
	rLogger := rLog.GetLogger()

	serviceProvider := service.InitProvider(cfg, rLogger)
	server, err := api.Start(cfg, serviceProvider, rLogger)
	if err != nil {
		log.Println("err:", err)
		return nil, err
	}
	return server, nil
}

func stopServer(server *http.Server) error {
	//defer db.Close(context.Background())
	//defer cache.Client.Close()
	var err error
	graceful := func() error {
		log.Println("Shutting down server gracefully")
		return nil
	}

	forced := func() error {
		log.Println("Shutting down server forcefully")
		return nil
	}

	sigs := []os.Signal{syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGTERM}
	errCh := make(chan error)
	go func() {
		errCh <- HandleSignals(sigs, graceful, forced)
	}()
	if err = <-errCh; err != nil {
		log.Println(err)
		return err
	}

	err = api.Stop(server)
	if err != nil {
		log.Println("server stop err:", err)
		return err
	}

	return nil
}

// HandleSignals listen on the registered signals and fires the gracefulHandler for the
// first signal and the forceHandler (if any) for the next this function blocks and
// return any error that returned by any of the api first
func HandleSignals(sigs []os.Signal, gracefulHandler, forceHandler func() error) error {
	sigCh := make(chan os.Signal)
	errCh := make(chan error, 1)

	signal.Notify(sigCh, sigs...)
	defer signal.Stop(sigCh)

	grace := true
	for {
		select {
		case err := <-errCh:
			return err
		case <-sigCh:
			if grace {
				grace = false
				go func() {
					errCh <- gracefulHandler()
				}()
			} else if forceHandler != nil {
				err := forceHandler()
				errCh <- err
			}
		}
	}
}
