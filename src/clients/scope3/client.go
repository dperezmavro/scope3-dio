package scope3

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
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
	results  chan common.PropertyResponse
	errors   chan error
	wg       *sync.WaitGroup
}

func New(
	token string,
	errors chan error,
	queries chan common.PropertyQuery,
	results chan common.PropertyResponse,
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
	logging.Info(
		ctx,
		logging.Data{
			"function": "client.StartListening",
			"listener": "listenForProperties",
			"package":  "scope3",
		},
		"listener starting",
	)
	go listenForProperties(s)
}

func (s *Client) fetchProperty(ctx context.Context, pq common.PropertyQuery) (common.PropertyResponse, error) {
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
		logging.Error(ctx, err, logging.Data{"properties": pq, "function": "fetchProperty"}, "error marhsaling body")
		return common.PropertyResponse{}, fmt.Errorf("unable to marshal request for properties %+v: %+v", pq, err)
	}

	req, err := http.NewRequest(
		http.MethodPost,
		"https://api.scope3.com/v2/measure?includeRows=true&latest=true&fields=emissionsBreakdown",
		bytes.NewBuffer(requestBody),
	)
	req.Header.Add(common.HeaderAuthorization, fmt.Sprintf("Bearer %s", s.apiToken))

	if err != nil {
		logging.Error(ctx, err, logging.Data{"properties": pq, "function": "fetchProperty"}, "error creating request")
		return common.PropertyResponse{}, fmt.Errorf("unable to create request for properties %+v: %+v", pq, err)
	}

	resp, err := s.hc.Do(req)
	if err != nil {
		logging.Error(ctx, err, logging.Data{"request": req, "function": "fetchProperty"}, "request-sending")
		return common.PropertyResponse{}, fmt.Errorf("unable to perform request for properties %+v: %+v", pq, err)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		logging.Error(ctx, err, logging.Data{"properties": pq, "function": "fetchProperty"}, "error reading response")
		return common.PropertyResponse{}, fmt.Errorf("unable to read response for request properties %+v: %+v", pq, err)
	}
	defer resp.Body.Close()

	m := MeasureAPIResponse{}
	err = json.Unmarshal(b, &m)
	if err != nil {
		logging.Error(
			ctx,
			err,
			logging.Data{
				"properties": pq,
				"function":   "fetchProperty",
				"response":   string(b),
			},
			"unable to unmarshal api response",
		)
		return common.PropertyResponse{}, fmt.Errorf("unable to unmarshal api response: %+v", err)
	}

	logging.Info(ctx, logging.Data{"rows": m.Rows}, "response")

	if len(m.Rows) < 1 {
		return common.PropertyResponse{
			InventoryID: pq.InventoryID,
			UtcDateTime: pq.UtcDateTime,
			Body:        "",
		}, nil
	}

	return common.PropertyResponse{
		InventoryID: pq.InventoryID,
		UtcDateTime: pq.UtcDateTime,
		Weight:      pq.Weight,
		Body:        fmt.Sprintf("%v", m.Rows[0].TotalEmissions),
	}, nil
}

func listenForProperties(c *Client) {
	for {
		property := <-c.queries
		ctx := context.WithValue(context.Background(), common.CtxKeyTraceID, "listenforproperties")
		logging.Info(ctx, logging.Data{"properties": property}, "fetching property")
		propertyResults, err := c.fetchProperty(ctx, property)
		if err != nil {
			logging.Error(ctx, err, logging.Data{"properties": property}, "error in fetching")
			c.errors <- fmt.Errorf("error fetching %+v: %+v", property, err)
		}

		propertyResults.InventoryID = property.InventoryID

		logging.Info(ctx, logging.Data{"property": property, "function": "listenForProperties"}, "store property request")
		c.results <- propertyResults
	}
}
