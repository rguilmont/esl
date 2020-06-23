package query

import (
	"encoding/json"
	"fmt"
)

// This package is used to query elasticsearch to get logs.

const (
	initialTimestamp = `"now-10m"`

	// This is an hardcoded elasticsearch query, because create a query with the library,
	//  based on map[string]interface{} is really ugly, and we want something simple here.
	//  gabs .String() method anyway returns valid json.
	esQuery = `{
		"query": {
			"bool": {
				"must": [{"query_string": { "query": %v}}],
				"filter": [
					{
						"range": {
							"@timestamp": %v
						}
					}
				]	
			}
		}
	}`
)

func GenerateEsQuery(queryString string, tsFrom, tsAfter string) []byte {

	qsJSON, err := json.Marshal(queryString)
	if err != nil {
		panic(err)
	}

	timestampQuery := map[string]string{
		"gt": tsFrom,
	}

	if tsAfter != "" {
		timestampQuery["lt"] = tsAfter
	}

	tsQueryJSON, err := json.Marshal(timestampQuery)
	if err != nil {
		panic(err)
	}

	res := []byte(fmt.Sprintf(esQuery, string(qsJSON), string(tsQueryJSON)))
	return res

}
