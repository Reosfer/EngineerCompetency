package models

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

type Request struct {
	ClientID     string `json:"client_id,omitempty" bson:"client_id,omitempty"`
	ClientSecret string `json:"client_secret,omitempty" bson:"client_secret,omitempty"`
	Username     string `json:"username,omitempty" bson:"username,omitempty"`
	Password     string `json:"password,omitempty" bson:"password,omitempty"`
	RedirectUri  string `json:"redirect_uri,omitempty" bson:"redirect_uri,omitempty"`
	Code         string `json:"code,omitempty" bson:"code,omitempty"`
	UserToken    string `json:"user_token,omitempty" bson:"user_token,omitempty"`
}
