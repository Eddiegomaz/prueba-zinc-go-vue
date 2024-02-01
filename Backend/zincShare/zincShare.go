package zincShare

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

const (
	userZinc     = "admin"
	passwordZinc = "Complexpass#123"
	urlRequest   = "http://localhost:4080/"
)

type (
	Email struct {
		ID      string `json:"_id"`
		Content string
		From    string
		To      string
		Subject string
	}
	Source struct {
		Source Email `json:"_source"`
	}
	Quantity struct {
		Value int
	}
	Emails struct {
		Total Quantity
		Hits  []Source
	}
	CountHits struct {
		Hits Emails
	}
	CreateResponse struct {
		Count int `json:"record_count"`
	}
)

func sendRequest(method, url, data string) (*http.Response, error) {
	req, err := http.NewRequest(method, url, strings.NewReader(data))
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(userZinc, passwordZinc)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.138 Safari/537.36")
	req.Close = true
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func CreateIndex() error {
	structEmail :=
		`{
		"name": "email",
		"storage_type": "disk",
		"shard_num": 1,
		"mappings": {
			"properties": {
				"from": {
					"type": "text",
					"index": true,
					"store": false
				},
				"to": {
					"type": "text",
					"index": true,
					"store": false
				},
				"subject": {
					"type": "text",
					"index": true,
					"store": false,
					"highlightable": true
				},
				"content": {
					"type": "text",
					"index": true,
					"store": false,
					"highlightable": true
				}
			}
		}
	}`
	resp, _ := sendRequest("POST", urlRequest+"api/index", structEmail)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return errors.New(string(body))
	}
	return nil
}

func DeleteIndex() (err error) {
	resp, _ := sendRequest("DELETE", urlRequest+"api/index/email", "")
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return errors.New(string(body))
	}
	return nil
}

func CreateData(data string) (int, error) {
	resp, _ := sendRequest("POST", urlRequest+"api/email/_multi", data)
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	res := CreateResponse{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return 0, err
	}
	return res.Count, nil
}

func Query(search string, from, size int) (res CountHits, err error) {

	var query string

	if search == "" {
		query = `{
			"match_all": {}
		}`
	} else {
		query = `{
			"query_string": {
				"query": "` + search + `"
			}
		}`
	}

	structQuery :=
		`{
		"query": {
		  "bool": {
			"must": [
				
				` + query + `	
			  
			]
		  }
		},
		"sort": [
		  "-@timestamp"
		],
		"from": ` + fmt.Sprint(from) + `,
		"size": ` + fmt.Sprint(size) + `,
		"aggs": {
		  "histogram": {
			"auto_date_histogram": {
			  "field": "@timestamp",
			  "buckets": 100
			}
		  }
		}
	  }`

	resp, _ := sendRequest("POST", urlRequest+"es/email/_search", structQuery)
	defer resp.Body.Close()
	log.Println(resp.StatusCode)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &res)
	return

}
