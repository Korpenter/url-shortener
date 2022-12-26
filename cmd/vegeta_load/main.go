package main

import (
	"encoding/json"
	"fmt"
	"github.com/Mldlr/url-shortener/internal/app/utils/encoders"
	vegeta "github.com/tsenart/vegeta/v12/lib"
	"log"
	"net/http"
	"sync/atomic"
	"time"
)

func main() {
	rateP := vegeta.Rate{Freq: 200, Per: time.Second}
	rateGetOne := vegeta.Rate{Freq: 300, Per: time.Second}
	duration := 8 * time.Second
	cred := map[string][]string{"Cookie": {"user_id=user1", "signature=60e8d0babc58e796ac223a64b5e68b998de7d3b203bc8a859bc0ec15ee66f5f9"}}
	var id uint64
	targeterShortenAPI := func() vegeta.Targeter {
		type payload struct {
			URL string `json:"url"`
		}
		return func(t *vegeta.Target) (err error) {
			t.Method = http.MethodPost
			t.URL = "http://localhost:8080/api/shorten"
			t.Body, err = json.Marshal(&payload{
				URL: fmt.Sprintf("%v.ru", atomic.AddUint64(&id, 1)),
			})
			t.Header = cred

			return err
		}
	}()

	targeterShorten := func() vegeta.Targeter {
		return func(t *vegeta.Target) (err error) {
			t.Method = http.MethodPost
			t.URL = "http://localhost:8080/"
			t.Body = []byte(fmt.Sprintf("%v.ru", atomic.AddUint64(&id, 1)))
			t.Header = cred
			return err
		}
	}()

	targeterShortenBatch := func() vegeta.Targeter {
		type payloadBatch struct {
			URL string `json:"original_url"`
			ID  string `json:"correlation_id"`
		}
		return func(t *vegeta.Target) (err error) {
			t.Method = http.MethodPost
			t.URL = "http://localhost:8080/api/shorten/batch"
			t.Body, err = json.Marshal(func() []payloadBatch {
				batchURLs := make([]payloadBatch, 0, 250)
				for j := 0; j < 250; j++ {
					batchURLs = append(batchURLs, payloadBatch{URL: fmt.Sprintf("%v.ru", atomic.AddUint64(&id, 1)), ID: fmt.Sprint(id)})
				}
				return batchURLs
			}())
			t.Header = cred
			return err
		}
	}()

	targeterExpand := func(id uint64) vegeta.Targeter {
		return func(t *vegeta.Target) (err error) {
			t.Method = http.MethodGet
			t.URL = fmt.Sprintf("http://localhost:8080/%v", encoders.ToRBase62(fmt.Sprintf("%v.ru", atomic.AddUint64(&id, 1))))
			t.Header = cred
			return err
		}
	}(0)

	targeterExpandUser := func() vegeta.Targeter {
		return func(t *vegeta.Target) (err error) {
			t.Method = http.MethodGet
			t.URL = "http://localhost:8080/api/user/urls"
			t.Header = cred

			return err
		}
	}()

	targeterDeleteBatch := func(id uint64) vegeta.Targeter {
		return func(t *vegeta.Target) (err error) {
			t.Method = http.MethodDelete
			t.URL = "http://localhost:8080/api/user/urls"
			t.Body, err = json.Marshal(func() []string {
				batchURLs := make([]string, 0, 500)
				for j := 0; j < 500; j++ {
					batchURLs = append(batchURLs, encoders.ToRBase62(fmt.Sprintf("%v.ru", atomic.AddUint64(&id, 1))))
				}
				return batchURLs
			}())
			t.Header = cred
			return err
		}
	}(0)

	var metrics vegeta.Metrics
	attacker := vegeta.NewAttacker(vegeta.Redirects(vegeta.NoFollow))

	for res := range attacker.Attack(targeterShortenAPI, rateP, duration, "Shorten api") {
		metrics.Add(res)
	}
	for res := range attacker.Attack(targeterShorten, rateP, duration, "Shorten") {
		metrics.Add(res)
	}
	for res := range attacker.Attack(targeterShortenBatch, rateP, duration, "Batch Shorten api") {
		metrics.Add(res)
	}
	for res := range attacker.Attack(targeterExpand, rateGetOne, duration, "Expand") {
		metrics.Add(res)
	}
	for res := range attacker.Attack(targeterExpandUser, rateP, duration, "Expand user api") {
		metrics.Add(res)
	}
	for res := range attacker.Attack(targeterDeleteBatch, rateP, duration, "Batch delete api") {
		metrics.Add(res)
	}

	metrics.Close()
	log.Printf("99th percentile: %s\n", metrics.Latencies.P99)
}
