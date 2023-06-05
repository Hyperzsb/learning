package jsonio

// Response defines standard format of the outbound HTTP response to the client.
// The Code indicate the HTTP status code of the response, while the Status and
// the Message describes the result of request either briefly or in detail.
type Response struct {
	Code    int         `json:"-"`
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
