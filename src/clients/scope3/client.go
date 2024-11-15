package scope3

import "net/http"

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

func (s Scope3Client) FetchProperty(p string) error {
	return nil
}

/*

curl --request POST \
     --url 'https://api.scope3.com/v2/measure?includeRows=true&latest=true&fields=emissionsBreakdown' \
     --header 'accept: application/json' \
     --header 'authorization: Bearer scope3_xXF23id9h2V4AY0AoTkuM37oXt6PLhni_1YBylx6XtFYZGEGcalvAiISX4Rn4UYGUw5Kbxw3PibVZCdglH2dOP2vaDlmYSpF2' \
     --header 'content-type: application/json' \
     --data '
{
  "rows": [
    {
      "impressions": 1000,
      "utcDatetime": "2024-10-31",
      "inventoryId": "nytimes.com"
    }
  ]
}
'
*/
