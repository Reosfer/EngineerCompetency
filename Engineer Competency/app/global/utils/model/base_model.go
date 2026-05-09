package model

type Response struct {
	StatusCode     int         `json:"status_code"`
	Data           interface{} `json:"data,omitempty"`
	Error          interface{} `json:"errors,omitempty"`
	ErrorStatement error       `json:"error,omitempty"`
	Page           int         `json:"page,omitempty"`
	PerPage        int         `json:"per_page,omitempty"`
	Total          int64       `json:"total,omitempty"`
	Message        string      `json:"message,omitempty"`
	Messages       []string    `json:"messages,omitempty"`
}

type CheckField struct {
	StatusCode int         `json:"status_code"`
	Data       interface{} `json:"data,omitempty"`
	Error      interface{} `json:"errors,omitempty"`
	IsExist    bool        `json:"is_exist"`
	Message    string      `json:"message,omitempty"`
}

type ProfileResponse struct {
	StatusCode int         `json:"status_code"`
	CodeType   string      `json:"code_type,omitempty"`
	Data       interface{} `json:"data,omitempty"`
	Error      interface{} `json:"errors,omitempty"`
	ErrorLog   *ErrorLog   `json:"error_logs,omitempty"`
	Page       int         `json:"page,omitempty"`
	PerPage    int         `json:"per_page,omitempty"`
	Total      int64       `json:"total,omitempty"`
	StatCode   string      `json:"stat_code,omitempty"`
	StatMsg    string      `json:"stat_msg,omitempty"`
}

type ErrorLog struct {
	Line              string      `json:"line,omitempty"`
	Filename          string      `json:"filename,omitempty"`
	Function          string      `json:"function,omitempty"`
	Message           interface{} `json:"message,omitempty"`
	SystemMessage     string      `json:"system_message,omitempty"`
	Url               string      `json:"url,omitempty"`
	Method            string      `json:"method,omitempty"`
	Fields            interface{} `json:"fields,omitempty"`
	ConsumerTopic     string      `json:"consumer_topic,omitempty"`
	ConsumerPartition int         `json:"consumer_partition,omitempty"`
	ConsumerName      string      `json:"consumer_name,omitempty"`
	ConsumerOffset    int64       `json:"consumer_offset,omitempty"`
	ConsumerKey       string      `json:"consumer_key,omitempty"`
	Err               error       `json:"-"`
	StatusCode        int         `json:"-"`
}
type ResponseChannel struct {
	Data  interface{} `json:"data,omitempty"`
	Error error       `json:"errors,omitempty"`
}
