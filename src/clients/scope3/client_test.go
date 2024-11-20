package scope3

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/dperezmavro/scope3-dio/src/common"
)

func TestGet(t *testing.T) {
	tests := []struct {
		name                  string
		pq                    common.PropertyQuery
		responseBody          string
		numOfResponseReecords int
	}{
		{
			name: "default",
			pq: common.PropertyQuery{
				Impressions: 1000,
				Weight:      10,
				InventoryID: "nytimes.com",
				UtcDateTime: "2024-10-28",
			},
			numOfResponseReecords: 1,
			responseBody:          `{"rows": [{"InventoryID":"nytimes.com","UtcDateTime":"2024-10-28","Body":"66.09248372608445","Weight":10}]}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cl := New("dummy", nil, nil, nil, 5)
			m := &ClientMock{
				t:    t,
				body: tt.responseBody,
			}
			cl.hc = m

			ctx := context.WithValue(context.Background(), common.CtxKeyTraceID, "unused")
			resp, err := cl.fetchProperty(ctx, []common.PropertyQuery{tt.pq})
			if err != nil {
				t.Errorf("unexpected error: %+v", err)
			}

			if len(resp) != tt.numOfResponseReecords {
				t.Errorf("incorrect response length: wanted %d, got %d", tt.numOfResponseReecords, len(resp))
			}

			if resp[0].PropertyName != tt.pq.InventoryID {
				t.Errorf("invalid InventoryID: wanted %s, got %s", tt.pq.InventoryID, resp[0].PropertyName)
			}

			if resp[0].UtcDateTime != tt.pq.UtcDateTime {
				t.Errorf("invalid UtcDateTime: wanted %s, got %s", tt.pq.UtcDateTime, resp[0].UtcDateTime)
			}
		})
	}
}

func TestChannels(t *testing.T) {
	tests := []struct {
		name         string
		errChan      chan error
		queries      chan []common.PropertyQuery
		results      chan []common.PropertyResponse
		pq           common.PropertyQuery
		responseBody string
	}{
		{
			name:    "default",
			errChan: make(chan error),
			queries: make(chan []common.PropertyQuery),
			results: make(chan []common.PropertyResponse),
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
				5,
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
			tt.queries <- []common.PropertyQuery{tt.pq}

			t.Log("getting result")
			result := <-tt.results

			// assert
			if result[0].PropertyName != tt.pq.InventoryID {
				t.Errorf("invalid InventoryID: wanted %s, got %s", tt.pq.InventoryID, result[0].PropertyName)
			}
		})
	}
}

func BenchmarkChannels(b *testing.B) {
	// setup
	errChan := make(chan error)
	queries := make(chan []common.PropertyQuery)
	results := make(chan []common.PropertyResponse)
	pq := []common.PropertyQuery{
		{
			Impressions: 1000,
			Weight:      10,
			InventoryID: "nytimes.com",
			UtcDateTime: "2024-10-28",
		},
		{
			Impressions: 1000,
			Weight:      10,
			InventoryID: "nba.com",
			UtcDateTime: "2024-10-28",
		},
	}
	responseBody := `{"rows": [{"InventoryID":"nytimes.com","UtcDateTime":"2024-10-28","Body":"66.09248372608445","Weight":10},{"InventoryID":"nba.com","UtcDateTime":"2024-10-28","Body":"66.09248372608445","Weight":10}]}`
	cl := New(
		"dummy",
		errChan,
		queries,
		results,
		5,
	)

	m := &ClientMock{
		t:    nil,
		b:    b,
		body: responseBody,
	}
	cl.hc = m

	// execute
	b.Log("start client listener")
	ctx := context.WithValue(context.Background(), common.CtxKeyTraceID, "unused")
	cl.StartListening(ctx)

	// simulate storage listener
	go func() {
		for {
			result := <-results
			b.Logf("get result %+v", result)
		}
	}()

	for n := 0; n < b.N; n++ {
		b.Log("sending query")
		queries <- pq
	}
}

type ClientMock struct {
	t    *testing.T
	b    *testing.B
	body string
}

func (c *ClientMock) Do(req *http.Request) (*http.Response, error) {
	if c.t != nil {
		c.t.Helper()
	} else {
		c.b.Helper()
	}
	return &http.Response{
		Body: io.NopCloser(strings.NewReader(c.body)),
	}, nil
}
