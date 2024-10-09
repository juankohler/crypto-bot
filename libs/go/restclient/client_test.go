package restclient

import (
	"context"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestCustomTransport(t *testing.T) {
	mockedTransport := httpmock.NewMockTransport()
	mockedTransport.RegisterResponder("GET", "/get", httpmock.NewStringResponder(201, "Test"))

	client := New(Config{
		BaseUrl:         "https://example.com",
		CustomTransport: mockedTransport,
	})

	endpoint := client.GET("/get")

	res1 := endpoint.DoRequest(context.Background())
	assert.Nil(t, res1.Err())
	assert.Equal(t, "Test", string(res1.Body()))

	res2 := endpoint.DoRequest(context.Background())
	assert.Nil(t, res2.Err())
	assert.Equal(t, "Test", string(res2.Body()))

	assert.Equal(
		t,
		map[string]int{
			"GET /get": 2,
		},
		mockedTransport.GetCallCountInfo(),
	)
}
