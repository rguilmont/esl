package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	elasticsearch "github.com/elastic/go-elasticsearch/v6"
	"github.com/rguilmont/esLogsTails/pkg/configuration"
	"github.com/rguilmont/esLogsTails/pkg/formater"
	"github.com/rguilmont/esLogsTails/pkg/logsProducer"
	"golang.org/x/crypto/ssh/terminal"
)

const esJSONPrefix = "_source."

// Export variables for version
var (
	GitSummary string
	GitCommit  string
	BuildDate  string
)

func currentFilter(contextFilter string, cliFilter string) string {
	if cliFilter != "" {
		return cliFilter
	}
	return contextFilter
}

func main() {

	flag.Usage = func() {
		fmt.Println("esl, elasticsearch logs query utility.")
		fmt.Println("By default, logs will be tailed, unless from and to flag are set.")

		fmt.Printf("Usage : %v [flags...] <context>\n", os.Args[0])
		flag.PrintDefaults()

	}

	version := flag.Bool("v", false, "Display version.")

	from := flag.String("from", "now-10m", "Start timestamp.")
	to := flag.String("to", "", "End timestamp. By default there's no end timestamp, it will infinitely loop.")
	filter := flag.String("filter", "", "Overide filter in your context")
	flag.Parse()

	if *version {
		fmt.Printf("Version:     %v\n", GitSummary)
		fmt.Printf("Commit hash: %v\n", GitCommit)
		fmt.Printf("Build date:  %v\n", BuildDate)
		os.Exit(0)
	}

	if len(flag.Args()) != 1 {
		fmt.Println("Invalid usage")

		fmt.Printf("Usage : %v [flags...] <context>\n", os.Args[0])

		flag.PrintDefaults()
		os.Exit(1)
	}
	c := flag.Args()[0]

	config, err := configuration.LoadConf(nil)
	if err != nil {
		panic(err)
	}

	context, err := config.Context(c)
	if err != nil {
		panic(err)
	}
	client, err := elasticsearch.NewClient(elasticsearch.Config{
		Username:  *context.Username,
		Password:  *context.Password,
		Addresses: []string{*context.URL},
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: *context.VerifySSLCertificate,
			},
		},
	})

	if err != nil {
		panic(err)
	}

	var producer logsProducer.LogProducer
	if *to != "" {
		producer = logsProducer.NewESLogsProducer(client, currentFilter(*context.Filter, *filter), *context.Index, time.Duration(*context.Refresh), *from, *to)

	} else {
		producer = logsProducer.NewTailESLogsProducer(client, currentFilter(*context.Filter, *filter), *context.Index, time.Duration(*context.Refresh))
	}

	f, _ := formater.NewFormater(
		func() formater.Formater {
			if terminal.IsTerminal(int(os.Stdout.Fd())) {
				return formater.ColoredFormater
			}
			return formater.DefaultFormater

		}(), formater.LogFieldsMapping{
			DateField:    esJSONPrefix + "@timestamp",
			LogField:     esJSONPrefix + "log",
			ServiceField: esJSONPrefix + "kubernetes.pod_name",
			ExtraField:   nil,
		})

	logProducer := producer.Produce()
	for log := range logProducer {
		f.PrintLog(log)
	}

}
