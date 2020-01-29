package model

type MockRequest struct {
	// Header config values can be of type string or []string.Hence using generic interface{} type
	RequestHeadersMockValues map[string]interface{}
	// request body can have many types as string,numeric,array,another struct etc.Hence using generic interface{} type
	RequestBodyMockValues map[string]interface{}
}

type MockResponse struct {
	HttpCode int
	// Header config values can be of type string or []string.Hence using generic interface{} type
	ResponseHeadersMockValues map[string]interface{}
	// response body can have many types as string,numeric,array,another struct etc.Hence using generic interface{} type
	ResponseBodyMockValues map[string]interface{}
}

type ApiConfig struct {
	Id                     uint64
	Url                    string
	Method                 string
	ResponseDelayInSeconds int
	RequestMockValues      MockRequest
	ResponseMockValues     MockResponse
}
