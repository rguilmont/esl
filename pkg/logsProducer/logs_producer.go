package logsProducer

import (
	"bytes"
	"io/ioutil"
	"time"

	"github.com/Jeffail/gabs/v2"
	"github.com/elastic/go-elasticsearch/v6"
	query "github.com/rguilmont/esLogsTails/pkg/logsProducer/query"
)

const (
	responseSize = 100
)

// LogProducer represents a log producer from a channel of gabs Container. Because logs must be in some kind of JSON format.
type LogProducer interface {
	Produce() chan (*gabs.Container)
}

// TailsESLogProducer Tail logs from elasticsearch
type TailsESLogProducer struct {
	esclient    *elasticsearch.Client
	queryString string
	index       string
	refreshTime time.Duration
}

// Produce create a chan where logs will be represented as *gabs.Container
func (p TailsESLogProducer) Produce() chan (*gabs.Container) {

	producer := make(chan (*gabs.Container))

	go func() {
		timestamp := "now-10m"
		for {
			res, err := p.esclient.Search(
				p.esclient.Search.WithIndex(p.index),
				p.esclient.Search.WithSort("@timestamp:desc"),
				p.esclient.Search.WithSize(responseSize),
				p.esclient.Search.WithBody(
					bytes.NewBuffer(query.GenerateEsQuery(p.queryString, timestamp, ""))),
			)
			//fmt.Println(string(generateEsQuery(queryString, timestamp)))
			if err != nil {
				panic(err)
			}

			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				panic(err)
			}

			container, err := gabs.ParseJSON(body)
			if err != nil {
				panic(err)
			}
			//fmt.Println(container.String())
			// If there's a date to set, then check if response is empty. If yes, it means we have got
			//  everything.
			//if to != "" && len(container.Path("hits.hits").Children()) == 0 {
			//	break
			//}
			for i := range container.Path("hits.hits").Children() {
				child := container.Path("hits.hits").Children()[len(container.Path("hits.hits").Children())-1-i]
				producer <- child

				timestamp = child.Path("_source.@timestamp").Data().(string) // change that interface to string. That's not very beautiful.

			}
			time.Sleep(p.refreshTime)
		}

	}()
	return producer
}

type ESLogProducer struct {
	esclient    *elasticsearch.Client
	queryString string
	index       string
	refreshTime time.Duration

	from string
	to   string
}

// Produce create a chan where logs will be represented as *gabs.Container
func (p ESLogProducer) Produce() chan (*gabs.Container) {
	producer := make(chan (*gabs.Container))

	go func() {
		timestamp := p.from
		for {
			res, err := p.esclient.Search(
				p.esclient.Search.WithIndex(p.index),
				p.esclient.Search.WithSort("@timestamp:asc"),
				p.esclient.Search.WithSize(responseSize),
				p.esclient.Search.WithBody(
					bytes.NewBuffer(query.GenerateEsQuery(p.queryString, timestamp, p.to))),
			)
			//fmt.Println(string(generateEsQuery(queryString, timestamp)))
			if err != nil {
				panic(err)
			}

			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				panic(err)
			}

			container, err := gabs.ParseJSON(body)
			if err != nil {
				panic(err)
			}
			//fmt.Println(container.String())
			// If there's a date to set, then check if response is empty. If yes, it means we have got
			//  everything.
			if len(container.Path("hits.hits").Children()) == 0 {
				break
			}
			for _, child := range container.Path("hits.hits").Children() {
				producer <- child

				timestamp = child.Path("_source.@timestamp").Data().(string) // change that interface to string. That's not very beautiful.

			}
			time.Sleep(p.refreshTime)
		}
		close(producer)

	}()
	return producer
}

func NewTailESLogsProducer(esclient *elasticsearch.Client, queryString string, index string, refreshTime time.Duration) *TailsESLogProducer {
	return &TailsESLogProducer{
		esclient:    esclient,
		queryString: queryString,
		index:       index,
		refreshTime: refreshTime,
	}
}

func NewESLogsProducer(esclient *elasticsearch.Client, queryString string, index string, refreshTime time.Duration, from string, to string) *ESLogProducer {
	return &ESLogProducer{
		esclient:    esclient,
		queryString: queryString,
		index:       index,
		refreshTime: refreshTime,
		from:        from,
		to:          to,
	}
}
