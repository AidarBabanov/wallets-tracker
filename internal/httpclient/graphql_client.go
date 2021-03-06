package httpclient

import (
	"net/http"
	"time"
)

type GraphqlClient struct {
	*Client
	url string
}

func NewGraphQLClient(url string, requestTimeout time.Duration, delay time.Duration) *GraphqlClient {
	client := new(GraphqlClient)
	client.url = url
	client.Client = New(requestTimeout, delay)
	return client
}

func (g *GraphqlClient) DoGraphqlRequest(query string, response interface{}) error {
	headers := map[string]string{
		"Content-Type": "application/json",
	}
	return g.Client.DoRequest(http.MethodPost, g.url, nil, headers, []byte(query), response)
}
