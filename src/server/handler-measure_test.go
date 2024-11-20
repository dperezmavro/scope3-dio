package server

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dperezmavro/scope3-dio/src/common"
	"github.com/dperezmavro/scope3-dio/src/storage"
)

func TestMeassure(t *testing.T) {
	tests := []struct {
		name        string
		errChan     chan error
		queries     chan []common.PropertyQuery
		results     chan []common.PropertyResponse
		pqs         []common.PropertyQuery
		prs         []common.PropertyResponse
		requestBody string
	}{
		{
			name:    "default",
			errChan: make(chan error),
			queries: make(chan []common.PropertyQuery),
			results: make(chan []common.PropertyResponse),
			pqs: []common.PropertyQuery{
				{
					Impressions: 1000,
					Weight:      10,
					InventoryID: "nytimes.com",
					UtcDateTime: "2024-10-30",
				},
				{
					Impressions: 1000,
					InventoryID: "nba.com",
					UtcDateTime: "2024-10-30",
				},
			},
			prs: []common.PropertyResponse{
				{
					Weight:         90,
					UtcDateTime:    "2024-10-30",
					PropertyName:   "nytimes.com",
					TotalEmissions: "66.09248372608445",
				},
				{
					UtcDateTime:    "2024-10-30",
					PropertyName:   "nba.com",
					TotalEmissions: "66.09248372608445",
				},
			},
			requestBody: `{"rows": [{"inventoryID":"nytimes.com","utcDateTime":"2024-10-30","impressions":1000, "weight":10}, {"inventoryID":"nba.com","utcDateTime":"2024-10-30","impressions":1000}]}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rist, err := storage.NewStorageImplementation(1e7, 200000, 64)
			if err != nil {
				t.Errorf("unable to make ristretto client: %+v", err)
			}

			cl, err := storage.New(
				rist,
				tt.errChan,
				tt.queries,
				tt.results,
				false,
			)
			if err != nil {
				t.Errorf("unable to create storage client: %+v", err)
			}

			t.Log("start client listener")
			ctx := context.WithValue(context.Background(), common.CtxKeyTraceID, "unused")
			cl.StartListening(ctx)

			request := httptest.NewRequest(
				http.MethodPost,
				"/v2/measure",
				bytes.NewBuffer([]byte(tt.requestBody)),
			).
				WithContext(ctx)
			responseRecorder := httptest.NewRecorder()

			// execute
			measure(cl)(responseRecorder, request)

			b, err := io.ReadAll(responseRecorder.Body)
			if err != nil {
				t.Errorf("unable to read body: %+v", err)
			}
			t.Logf("response %s", string(b))

			queries := <-tt.queries
			if len(queries) != len(tt.pqs) {
				t.Errorf("incorrect results received: wanted: %d, got %d: %+v", len(tt.pqs), len(queries), queries)
			}
			t.Logf("queries %+v", queries)
			tt.results <- tt.prs

			// give enough time to the goroutine in the storage client to read from the channel and store results
			time.Sleep(5 * time.Second)

			results := cl.Get(ctx, tt.pqs)
			if len(results) != len(tt.prs) {
				t.Errorf("incorrect results received: wanted: %d, got %d: %+v", len(tt.pqs), len(results), results)
			}
		})
	}
}
