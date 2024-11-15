package scope3

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/scope3-dio/common"
)

// Scope3Client is the client to interact with the scope3 api
type Scope3Client struct {
	hc       *http.Client
	apiToken string
}

func New(t string) Scope3Client {
	return Scope3Client{
		hc:       http.DefaultClient,
		apiToken: t,
	}
}

func (s Scope3Client) FetchProperty(pq []PropertyQuery) ([]PropertyResponse, error) {

	r := MeasureAPIRequest{}

	for _, p := range pq {
		r.Rows = append(r.Rows, RowItem{
			Channel:     p.Channel,
			Country:     p.Country,
			Impressions: 1000, // default
			InventoryID: p.InventoryID,
			UtcDateTime: p.UtcDateTime,
		})
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
