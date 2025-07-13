package location

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const ipSbEndpoint = "https://api.ip.sb/geoip"

type IpSbQueryResponse struct {
	Organization    string  `json:"organization"`
	Latitude        float64 `json:"latitude"`
	Isp             string  `json:"isp"`
	ContinentCode   string  `json:"continent_code"`
	AsnOrganization string  `json:"asn_organization"`
	Country         string  `json:"country"`
	Asn             int     `json:"asn"`
	Ip              string  `json:"ip"`
	Offset          int     `json:"offset"`
	Timezone        string  `json:"timezone"`
	Longitude       float64 `json:"longitude"`
	CountryCode     string  `json:"country_code"`
}

type ipSbPositionLocator struct {
	client *http.Client
}

func NewIpSbPositionLocator(
	client *http.Client,
) PositionLocator {
	return &ipSbPositionLocator{client: client}
}

func (i *ipSbPositionLocator) DetectIP(ctx context.Context, ip string) (result *Address, err error) {
	requestURL, parseErr := url.Parse(fmt.Sprintf("%s/%s", ipSbEndpoint, ip))
	if parseErr != nil {
		return nil, parseErr
	}

	request, buildRequestErr := http.NewRequestWithContext(ctx, http.MethodGet, requestURL.String(), nil)
	if buildRequestErr != nil {
		return nil, buildRequestErr
	}
	request.Header.Set("Accept", "application/json")
	request.Header.Set("User-Agent", "Mozilla/5.0")

	response, executeErr := i.client.Do(request)
	if executeErr != nil {
		return nil, executeErr
	}

	defer func() { _ = response.Body.Close() }()
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s", response.Status)
	}

	responsePayload := &IpSbQueryResponse{}
	if unmarshalErr := json.NewDecoder(response.Body).Decode(responsePayload); unmarshalErr != nil {
		return nil, unmarshalErr
	}

	result = &Address{
		Region:    responsePayload.CountryCode,
		Longitude: responsePayload.Longitude,
		Latitude:  responsePayload.Latitude,
	}

	return result, nil
}
