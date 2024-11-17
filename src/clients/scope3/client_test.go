package scope3

import (
	"context"
	"io"
	"net/http"
	"strings"
	"sync"
	"testing"

	"github.com/dperezmavro/scope3-dio/src/common"
)

func TestGet(t *testing.T) {
	tests := []struct {
		name         string
		pq           common.PropertyQuery
		responseBody string
	}{
		{
			name: "default",
			pq: common.PropertyQuery{
				Impressions: 1000,
				Weight:      10,
				InventoryID: "nytimes.com",
				UtcDateTime: "2024-10-28",
			},
			responseBody: `{"rows": [{"InventoryID":"nytimes.com","UtcDateTime":"2024-10-28","Body":"66.09248372608445","Weight":10}]}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cl := New("dummy", nil, nil, nil, nil)
			m := &ClientMock{
				t:    t,
				body: tt.responseBody,
			}
			cl.hc = m

			ctx := context.WithValue(context.Background(), common.CtxKeyTraceID, "unused")
			resp, err := cl.fetchProperty(ctx, tt.pq)
			if err != nil {
				t.Errorf("unexpected error: %+v", err)
			}

			if resp.InventoryID != tt.pq.InventoryID {
				t.Errorf("invalid InventoryID: wanted %s, got %s", tt.pq.InventoryID, resp.InventoryID)
			}

			if resp.UtcDateTime != tt.pq.UtcDateTime {
				t.Errorf("invalid UtcDateTime: wanted %s, got %s", tt.pq.UtcDateTime, resp.UtcDateTime)
			}
		})
	}
}

func TestChannels(t *testing.T) {
	tests := []struct {
		name         string
		errChan      chan error
		queries      chan common.PropertyQuery
		results      chan common.PropertyResponse
		wg           *sync.WaitGroup
		pq           common.PropertyQuery
		responseBody string
	}{
		{
			name:    "default",
			errChan: make(chan error),
			queries: make(chan common.PropertyQuery),
			results: make(chan common.PropertyResponse),
			wg:      &sync.WaitGroup{},
			pq: common.PropertyQuery{
				Impressions: 1000,
				Weight:      10,
				InventoryID: "nytimes.com",
				UtcDateTime: "2024-10-28",
			},
			responseBody: `{"rows": [{"InventoryID":"nytimes.com","UtcDateTime":"2024-10-28","Body":"66.09248372608445","Weight":10}]}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// setup
			cl := New(
				"dummy",
				tt.errChan,
				tt.queries,
				tt.results,
				tt.wg,
			)

			m := &ClientMock{
				t:    t,
				body: tt.responseBody,
			}
			cl.hc = m

			// execute
			t.Log("start client listener")
			ctx := context.WithValue(context.Background(), common.CtxKeyTraceID, "unused")
			cl.StartListening(ctx)

			t.Log("sending query")
			tt.queries <- tt.pq

			t.Log("getting result")
			result := <-tt.results

			// assert
			if result.InventoryID != tt.pq.InventoryID {
				t.Errorf("invalid InventoryID: wanted %s, got %s", tt.pq.InventoryID, result.InventoryID)
			}
		})
	}
}

type ClientMock struct {
	t    *testing.T
	body string
}

func (c *ClientMock) Do(req *http.Request) (*http.Response, error) {
	c.t.Helper()
	return &http.Response{
		Body: io.NopCloser(strings.NewReader(c.body)),
	}, nil
}
