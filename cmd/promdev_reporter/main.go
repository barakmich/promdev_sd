package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/barakmich/promdev_sd"
	"github.com/spf13/pflag"
)

var (
	hostport          = pflag.StringP("hostport", "H", "127.0.0.1:9111", "Hostport of the promdev_sd server")
	heartbeatInterval = pflag.DurationP("heartbeat-interval", "I", 15*time.Second, "Interval to refresh the advertisement")
	targets           = pflag.StringArrayP("targets", "t", nil, "Targets to Report")
	labels            = pflag.StringArrayP("labels", "l", nil, "Labels to add to the targets")
	namespace         = pflag.StringP("namespace", "n", "", "Namespace to report to")
	debug             = pflag.Bool("debug", false, "Debug output")
)

func main() {
	pflag.Parse()
	if len(*targets) == 0 {
		log.Fatalln("Must include at least one target")
	}

	if *namespace == "" {
		log.Fatalln("Namespace is required (--namespace or -n)")
	}

	ts := buildTargetSet()
	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(ts)
	if err != nil {
		log.Fatalln("Error marshalling target set:", err)
	}
	baseURL, err := url.Parse("http://" + *hostport)
	if err != nil {
		log.Fatalln("Couldn't parse hostport:", err)
	}
	initial, err := http.NewRequest(http.MethodPut, baseURL.JoinPath("register", *namespace).String(), buf)
	if err != nil {
		log.Fatalln("Initial request could not be created")
	}
	resp, err := http.DefaultClient.Do(initial)
	if err != nil && !errors.Is(err, io.EOF) {
		log.Fatalln("Error in initial request", err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln("Error reading token body")
	}
	resp.Body.Close()
	token := string(body)

	ctx, cancel := CtrlCContext()

outer:
	for {
		select {
		case <-ctx.Done():
			cancel()
			break outer
		case <-time.Tick(*heartbeatInterval):
			final, err := http.NewRequest(http.MethodPut, baseURL.JoinPath("heartbeat", *namespace, token).String(), nil)
			if err != nil {
				log.Fatalln("Update request could not be created")
			}
			resp, err = http.DefaultClient.Do(final)
			if err != nil {
				log.Fatalln("Error in update request", err)
			}
			resp.Body.Close()
		}
	}

	final, err := http.NewRequest(http.MethodDelete, baseURL.JoinPath("heartbeat", *namespace, token).String(), buf)
	if err != nil {
		log.Fatalln("Final request could not be created")
	}
	resp, err = http.DefaultClient.Do(final)
	if err != nil {
		log.Fatalln("Error in final request", err)
	}
}

func buildTargetSet() *promdev_sd.TargetSet {
	out := &promdev_sd.TargetSet{
		Targets: *targets,
		Labels:  make(promdev_sd.LabelSet),
	}
	for _, v := range *labels {
		pair := strings.Split(v, "=")
		if len(pair) != 2 {
			log.Fatalln("Invalid label:", v)
		}
		out.Labels[pair[0]] = pair[1]
	}
	return out
}

func CtrlCContext() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func(cancel context.CancelFunc) {
		sig := <-sigs
		log.Println("Caught signal:", sig)
		cancel()
	}(cancel)
	return ctx, cancel
}
