# `esl` : ElasticSearch logs utility

Simple program to get and tail logs stored in elasticsearch by fluentd.
You can either tail the logs, or get all logs matching a filter.

For now, it is assumed that the logs comes from a Kubernetes cluster. I will make the tool more modular in the future.


## Usage

```
Usage : ./esl [flags...] <context>
  -filter string
        Overide filter in your context
  -from string
        Start timestamp. (default "now-10m")
  -to string
        End timestamp. By default there's no end timestamp, it will infinitely loop.
```
## Configuration

`esl` use configuration file `config.yaml`. These are looked up at : 

```
/etc/esl/config.yaml
$HOME/.esl/config.yaml
./config.yaml
```

Configuration defines contexts. A context allows to query elasticsearch.

Example of configuration : 

```
local:
  URL: http://localhost/es
  filter: '*'
  index: logstash*
  refresh: 200
info:
  URL: https://es.example.org/
  Username: username
  Password: password
  filter: 'kubernetes.container_name: my_app'
  index: logstash*
  refresh: 200
```

This configuration define 2 contexts : `local` and `info`. Each have their own cluster and filter. 