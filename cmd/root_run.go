package cmd

import (
	"context"
	"demo_project/api/external"
	"demo_project/config"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	log "github.com/sirupsen/logrus"
	"syscall"
)

func run(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	tasks := []func() error {
		setLogLevel,
		printStartMessage,
		setupExternalAPI,
	}

	for _, t := range tasks {
		if err := t(); err != nil {
			log.Fatal(err)
		}
	}

	sigChan := make(chan os.Signal)
	exitChan := make(chan struct{})
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	log.WithField("signal", <-sigChan).Info("signal received")
	go func() {
		log.Warning("stopping server")
		// todo: handle graceful shutdown?
		exitChan <- struct{}{}
	}()
	select {
	case <-exitChan:
	case s := <-sigChan:
		log.WithField("signal", s).Info("signal received, stopping immediately")
	}

	return nil
}

func setLogLevel() error {
	log.SetLevel(log.Level(uint8(config.Conf.General.LogLevel)))
	return nil
}

func printStartMessage() error {
	log.WithFields(log.Fields{
		"version": version,
		"docs":    "",
	}).Info("starting server")
	return nil
}

func setupExternalAPI() error {
	if err := external.Setup(config.Conf); err != nil {
		return errors.Wrap(err, "setup proto error")
	}
	return nil
}