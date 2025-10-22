package router

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/jmontesinos91/ologs/logger"
	"github.com/jmontesinos91/osecurity/services/omnibackend"
	"github.com/sirupsen/logrus"
	"go.elastic.co/apm/module/apmhttp/v2"
	"golang.org/x/net/context/ctxhttp"
)

type DefaultWebClient struct {
	log        *logger.ContextLogger
	httpClient *http.Client
	baseURL    string
	reqPath    string
}

// NewRouterService to call api client
func NewRouterService(l *logger.ContextLogger, omniview omnibackend.OmniViewConfigurations) *DefaultWebClient {
	client := retryablehttp.NewClient()
	// Set timeout for client
	client.HTTPClient = &http.Client{
		Timeout: time.Duration(omniview.TimeoutInSeconds) * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	// Configure retry policy
	client.RetryWaitMin = time.Duration(omniview.RetryWaitMinInSeconds) * time.Second
	client.RetryWaitMax = time.Duration(omniview.RetryWaitMaxInSeconds) * time.Second
	// Set MaxNumber of retries
	client.RetryMax = omniview.MaxRetries

	return &DefaultWebClient{
		log:        l,
		baseURL:    omniview.Server,
		httpClient: apmhttp.WrapClient(client.StandardClient()),
		reqPath:    "/v1/colector/alarm",
	}
}

func (c *DefaultWebClient) ValidateIMEI(ctx context.Context, request Request) (*Response, error) {
	var res Response

	url := c.baseURL + c.reqPath

	bodyReq, err := json.Marshal(request)
	if err != nil {
		return nil, errors.New("alarm_client: error parsing request body: " + err.Error())
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(bodyReq))
	if err != nil {
		return nil, errors.New("alarm_client: error creating http request: " + err.Error())
	}

	// Set required headers
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	// Do request
	response, err := ctxhttp.Do(ctx, c.httpClient, req)

	if err != nil {
		return nil, errors.New("alarm_client: error doing request: " + err.Error())
	}
	// Defer closing of the response body
	defer func() {
		err := response.Body.Close()
		if err != nil {
			c.log.Error(logrus.ErrorLevel, "ValidateToken", "alarm_client: %v", err)
		}
	}()

	// Read response
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("alarm_client: error reading response: %v", err)
	}
	// If status is not a 2xx response return an error
	if response.StatusCode < 200 || response.StatusCode >= 500 {
		return nil, errors.New("alarm_client: service response with " + string(rune(response.StatusCode)) + " status code, response body: " + string(body))
	}

	if err := json.Unmarshal(body, &res); err != nil {
		return nil, errors.New("alarm_client: error parsing response body: " + err.Error())
	}

	return &res, nil
}
