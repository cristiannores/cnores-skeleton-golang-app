package http_client

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"cnores-skeleton-golang-app/app/infrastructure/constant"
	"cnores-skeleton-golang-app/app/shared/utils/config"
	utils_context "cnores-skeleton-golang-app/app/shared/utils/context"
)

type HttpClientSettings struct {
	Host                  string
	Headers               map[string]string
	Timeout               int
	SecurityPolicyEnabled bool
	Config                config.Config
	ServiceName           string
	Secure                bool
}

type HttpClient[Request any, Response any] struct {
	settings HttpClientSettings
}

type HttpResponse[Response any] struct {
	StatusCode int
	Body       Response
}

type HttpClientError struct {
	StatusCodeMessage int
	ErrorMessage      string
}

func (h *HttpClientError) Error() string {
	return h.ErrorMessage
}

func (h *HttpClientError) StatusCode() int {
	return h.StatusCodeMessage
}
func NewHttpClient[Request any, Response any](settings HttpClientSettings) HttpClientInterface[Request, Response] {
	return &HttpClient[Request, Response]{settings}
}

func (httpClient *HttpClient[Request, Response]) Post(ctx context.Context, url string, body Request) (*HttpResponse[Response], error) {
	log := utils_context.GetLogFromContext(ctx, constant.InfrastructureLayer, "http_client.Post")
	startTime := time.Now()
	var response Response
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(body)
	if err != nil {
		msg := fmt.Sprintf("Error encoding a message, error: %v", err.Error())
		log.Error(msg)
		return nil, returnHttpError(nil, msg)
	}
	endpoint := fmt.Sprintf("%s/%s", httpClient.settings.Host, url)
	log.Info("calling endpoint %s", endpoint)

	jsonData, _ := json.Marshal(body)
	log.Info(fmt.Sprintf("calling with body %s", string(jsonData)))

	req, errRequest := http.NewRequest("POST", endpoint, &buf)
	if errRequest != nil {
		msg := fmt.Sprintf("Error Creating the http post request, error: %v", errRequest.Error())
		log.Error(msg)
		return nil, returnHttpError(nil, msg)
	}

	for key, value := range httpClient.settings.Headers {
		req.Header.Set(key, value)
	}
	var client *http.Client
	if httpClient.settings.Secure == true {
		client = &http.Client{Timeout: time.Duration(httpClient.settings.Timeout) * time.Millisecond}
	} else {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
		client = &http.Client{Transport: tr, Timeout: time.Duration(httpClient.settings.Timeout) * time.Millisecond}
	}

	res, errClient := client.Do(req)
	if errClient != nil {
		msg := fmt.Sprintf("Error while executing client Do for http post, error: %v", errClient.Error())
		return nil, returnHttpError(res, msg)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Error("Error closing the Http Client")
		}
	}(res.Body)

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated {
		bodyBytes, err := io.ReadAll(res.Body)
		errResponse := string(bodyBytes)
		if err != nil {
			msg := fmt.Sprintf("Status Code: %d with error: %s", res.StatusCode, "Eror reading response body")
			log.Error(msg)
			return nil, returnHttpError(res, msg)
		}
		msg := fmt.Sprintf("Status Code: %d with error: %#v", res.StatusCode, errResponse)
		log.Error(msg)
		return nil, returnHttpError(res, msg)
	}

	if res.Header.Get("content-type") == "application/json" || res.Header.Get("content-type") == "application/json; charset=utf-8" {
		bodyBytes, errReading := io.ReadAll(res.Body)
		if errReading != nil {
			msg := fmt.Sprintf("Status Code: %d with error: %s", res.StatusCode, errReading.Error())
			log.Error(msg)
			return nil, returnHttpError(res, msg)
		}
		bodyString := string(bodyBytes)

		errDecoding := json.Unmarshal(bodyBytes, &response)
		if errDecoding != nil {
			msg := fmt.Sprintf("Error parsing json response %s with msg: %v", bodyString, errDecoding)
			return nil, returnHttpError(res, msg)
		}
		log.Info(fmt.Sprintf("success parsing response bodyString %s", bodyString))
		httpClient.logElapsedTime(ctx, startTime)
		return &HttpResponse[Response]{
			Body:       response,
			StatusCode: res.StatusCode,
		}, nil
	}

	readerResponse, errReading := io.ReadAll(res.Body)

	if errReading != nil {
		msg := fmt.Sprintf("Error reading response with msg: %v", errReading)
		return nil, returnHttpError(res, msg)
	}

	var asInterface interface{} = string(readerResponse)
	httpClient.logElapsedTime(ctx, startTime)
	return &HttpResponse[Response]{
		Body:       asInterface.(Response),
		StatusCode: res.StatusCode,
	}, nil

}

func (httpClient *HttpClient[Request, Response]) Get(ctx context.Context, url string) (*HttpResponse[Response], error) {
	log := utils_context.GetLogFromContext(ctx, constant.InfrastructureLayer, "http_client.Post")
	startTime := time.Now()
	var response Response
	endpoint := fmt.Sprintf("%s/%s", httpClient.settings.Host, url)
	log.Info("calling endpoint %s", endpoint)
	req, errRequest := http.NewRequest("GET", endpoint, nil)
	if errRequest != nil {
		msg := fmt.Sprintf("Error Creating the http post request, error: %v", errRequest.Error())
		log.Error(msg)
		return nil, returnHttpError(nil, msg)
	}

	var client *http.Client
	if httpClient.settings.Secure == true {
		client = &http.Client{Timeout: time.Duration(httpClient.settings.Timeout) * time.Millisecond}
	} else {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
		client = &http.Client{Transport: tr, Timeout: time.Duration(httpClient.settings.Timeout) * time.Millisecond}
	}
	res, errClient := client.Do(req)
	if errClient != nil {
		msg := fmt.Sprintf("Error while executing client Do for http post, error: %v", errClient.Error())
		log.Error(msg)
		return nil, returnHttpError(res, msg)
	}
	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated {

		bodyBytes, err := io.ReadAll(res.Body)
		errResponse := string(bodyBytes)
		if err != nil {
			msg := fmt.Sprintf("Status Code: %d with error: %s", res.StatusCode, "Eror reading response body")
			log.Error(msg)
			return nil, returnHttpError(res, msg)
		}
		msg := fmt.Sprintf("Status Code: %d with error: %#v", res.StatusCode, errResponse)
		log.Error(msg)
		return nil, returnHttpError(res, msg)
	}

	if res.Header.Get("content-type") == "application/json" || res.Header.Get("content-type") == "application/json; charset=utf-8" {
		errDecoding := json.NewDecoder(res.Body).Decode(&response)
		if errDecoding != nil {
			msg := fmt.Sprintf("Error parsing json response with msg: %v", errDecoding)
			return nil, returnHttpError(res, msg)
		}
		httpClient.logElapsedTime(ctx, startTime)
		return &HttpResponse[Response]{
			Body:       response,
			StatusCode: res.StatusCode,
		}, nil
	}

	readerResponse, errReading := io.ReadAll(res.Body)

	if errReading != nil {

		msg := fmt.Sprintf("Error reading response with msg: %v", errReading)
		return nil, returnHttpError(res, msg)
	}

	var asInterface interface{} = string(readerResponse)
	httpClient.logElapsedTime(ctx, startTime)
	return &HttpResponse[Response]{
		Body:       asInterface.(Response),
		StatusCode: res.StatusCode,
	}, nil
}

func returnHttpError(res *http.Response, msg string) *HttpClientError {

	statusCode := http.StatusInternalServerError
	if res != nil {
		statusCode = res.StatusCode
	}
	return &HttpClientError{
		StatusCodeMessage: statusCode,
		ErrorMessage:      msg,
	}

}

func (httpClient *HttpClient[Request, Response]) logElapsedTime(ctx context.Context, startTime time.Time) {
	log := utils_context.GetLogFromContext(ctx, constant.InfrastructureLayer, "http_client.logElapsedTime")
	endTime := time.Since(startTime).Seconds()
	log.Info(fmt.Sprintf("service %s time elapsed  %f", httpClient.settings.ServiceName, endTime))
}
