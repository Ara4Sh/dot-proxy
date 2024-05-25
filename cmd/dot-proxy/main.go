package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	dnsproxy "github.com/ara4sh/go-dot-proxy/pkg/dnsproxy"
	log "github.com/sirupsen/logrus"
)

func main() {
	var (
		listenTCP         bool
		hostTCP           string
		portTCP           int
		listenUDP         bool
		hostUDP           string
		portUDP           int
		cloudFlareDoTAddr string
		debug             bool
	)
	flag.BoolVar(&listenTCP, "tcp", true, "Enable TCP listening mode.")
	flag.StringVar(&hostTCP, "host-tcp", "0.0.0.0", "Host to listen on for TCP.")
	flag.IntVar(&portTCP, "port-tcp", 8053, "Port to listen on for TCP.")
	flag.BoolVar(&listenUDP, "udp", true, "Enables UDP listening mode.")
	flag.StringVar(&hostUDP, "host-udp", "0.0.0.0", "Host to listen on for UDP.")
	flag.IntVar(&portUDP, "port-up", 8053, "Port to listen on for UDP.")
	flag.StringVar(&cloudFlareDoTAddr, "cloudflare-dot-addr", "1.1.1.1:853", "Cloudflare DoT address.")
	flag.BoolVar(&debug, "debug", false, "Enables Debug mode.")

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [options]\n\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	if debug {
		log.SetLevel(log.DebugLevel)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	log.Info(fmt.Sprintf("Using %s as Cloudflare DoT server", cloudFlareDoTAddr))
	dnsProxy := dnsproxy.NewDNSProxy(cloudFlareDoTAddr)

	go func() {
		if listenTCP {
			// TODO Panic if address is already open
			log.Info(fmt.Sprintf("Starting TCP server on %s:%d", hostTCP, portTCP))
			if err := dnsProxy.ListenTCP(ctx, hostTCP, portTCP); err != nil {
				log.Errorf("TCP listener error: %v", err)
				panic(err)
			}
		}
	}()

	go func() {
		if listenUDP {
			log.Info(fmt.Sprintf("Starting UDP server on %s:%d", hostUDP, portUDP))
			if err := dnsProxy.ListenUDP(ctx, hostUDP, portUDP); err != nil {
				log.Errorf("UDP listener error: %v", err)
				panic(err)
			}
		}
	}()
	<-sigChan
	log.Info("Received shutdown signal. Shutting down gracefully...")
	cancel()
}
