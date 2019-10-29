package model

type Request struct {
	// Header config values can be of type string or []string.Hence using generic interface{} type
	RequestHeaders map[string]interface{}
	// request body can have many types as string,numeric,array,another struct etc.Hence using generic interface{} type
	RequestJsonBody map[string]interface{}
}

type Response struct {
	HttpCode int
	// Header config values can be of type string or []string.Hence using generic interface{} type
	ResponseHeaders map[string]interface{}
	// response body can have many types as string,numeric,array,another struct etc.Hence using generic interface{} type
	ResponseJsonBody map[string]interface{}
}

type ApiConfig struct {
	Id                     uint64
	Url                    string
	Method                 string
	ResponseDelayInSeconds int
	RequestConfig          Request
	ResponseConfig         Response
}
