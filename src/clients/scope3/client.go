package scope3

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"

	"github.com/scope3-dio/src/common"
	"github.com/scope3-dio/src/logging"
)

// Client is the client to interact with the scope3 api
type Client struct {
	hc       *http.Client
	apiToken string
	queries  chan common.PropertyQuery
	results  chan []common.PropertyResponse
	errors   chan error
	wg       *sync.WaitGroup
}

func New(
	token string,
	errors chan error,
	queries chan common.PropertyQuery,
	results chan []common.PropertyResponse,
	wg *sync.WaitGroup,
) Client {
	c := Client{
		hc:       http.DefaultClient,
		apiToken: token,
		queries:  queries,
		results:  results,
		errors:   errors,
		wg:       wg,
	}

	return c
}

func (s *Client) StartListening(ctx context.Context) {
	// wait for the listenForProperties goroutine
	s.wg.Add(1)
	logging.Info(ctx, logging.Data{"function": "client.StartListening"}, "starting scope3 goroutine listener")
	go listenForProperties(
		s,
		s.queries,
		s.results,
		s.errors,
	)
}

func (s *Client) fetchProperty(pq common.PropertyQuery) ([]common.PropertyResponse, error) {

	r := MeasureAPIRequest{
		Rows: []RowItem{
			{
				Channel:     pq.Channel,
				Country:     pq.Country,
				Impressions: pq.Impressions,
				InventoryID: pq.InventoryID,
				UtcDateTime: pq.UtcDateTime,
			},
		},
	}
	requestBody, err := json.Marshal(r)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal request for properties %+v: %+v", pq, err)
	}

	req, err := http.NewRequest(
		http.MethodPost,
		"https://api.scope3.com/v2/measure?includeRows=true&latest=true&fields=emissionsBreakdown",
		bytes.NewBuffer(requestBody),
	)
	req.Header.Add(common.HeaderAuthorization, fmt.Sprintf("Bearer %s", s.apiToken))

	if err != nil {
		return nil, fmt.Errorf("unable to create request for properties %+v: %+v", pq, err)
	}

	resp, err := s.hc.Do(req)
	if err != nil {
		return nil, fmt.Errorf("unable to perform request for properties %+v: %+v", pq, err)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read response for request properties %+v: %+v", pq, err)
	}
	defer resp.Body.Close()
	log.Println(string(b))

	return nil, nil
}

func listenForProperties(
	c *Client,
	queries chan common.PropertyQuery,
	results chan []common.PropertyResponse,
	errors chan error,
) {
	for {
		properties := <-queries
		ctx := context.WithValue(context.Background(), common.CtxKeyTraceID, "listenforproperties")
		logging.Info(ctx, logging.Data{"properties": properties}, "fetching property")
		propertyResults, err := c.fetchProperty(properties)
		if err != nil {
			logging.Error(ctx, err, logging.Data{"properties": properties}, "error in fetching")
			errors <- fmt.Errorf("error fetching %+v: %+v", properties, err)
		}

		logging.Info(ctx, logging.Data{"properties": properties}, "fetching property")
		results <- propertyResults
	}
}
