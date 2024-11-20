package scope3

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/dperezmavro/scope3-dio/src/common"
	"github.com/dperezmavro/scope3-dio/src/logging"
)

// Client is the client to interact with the scope3 api
type Client struct {
	hc              HTTPClient
	apiToken        string
	queries         chan []common.PropertyQuery
	results         chan []common.PropertyResponse
	errors          chan error
	numberOfWorkers int
}

func New(
	token string,
	errors chan error,
	queries chan []common.PropertyQuery,
	results chan []common.PropertyResponse,
	workerNo int,
) Client {
	c := Client{
		hc:              http.DefaultClient,
		apiToken:        token,
		queries:         queries,
		results:         results,
		errors:          errors,
		numberOfWorkers: workerNo,
	}

	return c
}

func (s *Client) StartListening(ctx context.Context) {
	// wait for the listenForProperties goroutine
	logging.Info(
		ctx,
		logging.Data{
			"function": "client.StartListening",
			"listener": "listenForProperties",
			"package":  "scope3",
		},
		"listener starting",
	)
	for i := 0; i < s.numberOfWorkers; i++ {
		go listenForProperties(ctx, s, i)
	}
}

func (s *Client) fetchProperty(
	ctx context.Context,
	pq []common.PropertyQuery,
) ([]common.PropertyResponse, error) {
	rowItems := make([]RowItem, len(pq))
	for idx, p := range pq {
		rowItems[idx] = RowItem{
			Channel:     p.Channel,
			Country:     p.Country,
			Impressions: p.Impressions,
			InventoryID: p.InventoryID,
			UtcDateTime: p.UtcDateTime,
		}
	}

	r := MeasureAPIRequest{
		Rows: rowItems,
	}
	requestBody, err := json.Marshal(r)
	if err != nil {
		logging.Error(ctx, err, logging.Data{"properties": pq, "function": "fetchProperty"}, "error marhsaling body")
		return []common.PropertyResponse{}, fmt.Errorf("unable to marshal request for properties %+v: %+v", pq, err)
	}

	req, err := http.NewRequest(
		http.MethodPost,
		"https://api.scope3.com/v2/measure?includeRows=true&latest=true&fields=emissionsBreakdown",
		bytes.NewBuffer(requestBody),
	)
	req.Header.Add(common.HeaderAuthorization, fmt.Sprintf("Bearer %s", s.apiToken))

	if err != nil {
		logging.Error(ctx, err, logging.Data{"properties": pq, "function": "fetchProperty"}, "error creating request")
		return []common.PropertyResponse{}, fmt.Errorf("unable to create request for properties %+v: %+v", pq, err)
	}

	resp, err := s.hc.Do(req)
	if err != nil {
		logging.Error(ctx, err, logging.Data{"request": req, "function": "fetchProperty"}, "request-sending")
		return []common.PropertyResponse{}, fmt.Errorf("unable to perform request for properties %+v: %+v", pq, err)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		logging.Error(ctx, err, logging.Data{"properties": pq, "function": "fetchProperty"}, "error reading response")
		return []common.PropertyResponse{}, fmt.Errorf("unable to read response for request properties %+v: %+v", pq, err)
	}
	defer resp.Body.Close()

	m := MeasureAPIResponse{
		// Rows: make([]Row, len(pq)),
	}
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
		return []common.PropertyResponse{}, fmt.Errorf("unable to unmarshal api response: %+v", err)
	}

	logging.Info(ctx, logging.Data{"rows": m.Rows, "pq": pq}, "response")

	if len(m.Rows) < 1 {
		return []common.PropertyResponse{}, nil
	}

	responses := make([]common.PropertyResponse, len(m.Rows))
	for idx, r := range m.Rows {
		responses[idx] = common.PropertyResponse{
			PropertyName:   r.InventoryID,
			UtcDateTime:    pq[idx].UtcDateTime,
			Weight:         pq[idx].Weight,
			TotalEmissions: fmt.Sprintf("%v", r.TotalEmissions),
		}
	}

	return responses, nil
}

func listenForProperties(ctx context.Context, c *Client, workerID int) {
	logging.Info(ctx, logging.Data{"id": workerID}, "scope3 worker started")
	for {
		properties := <-c.queries
		ctx := context.WithValue(context.Background(), common.CtxKeyTraceID, "listenforproperties")
		logData := logging.Data{"properties": properties, "workerId": workerID}
		logging.Info(ctx, logData, "fetching properties")
		propertyResults, err := c.fetchProperty(ctx, properties)
		if err != nil {
			logging.Error(ctx, err, logData, "error in fetching")
			c.errors <- fmt.Errorf("error fetching %+v: %+v", properties, err)
		}

		logging.Info(ctx,
			logging.Data{
				"properties": properties,
				"function":   "listenForProperties",
				"workerId":   workerID,
			},
			"store properties request",
		)
		c.results <- propertyResults
	}
}
