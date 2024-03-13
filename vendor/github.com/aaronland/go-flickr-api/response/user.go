package response

type User struct {
	Id       string    `json:"id"`
	Username *Username `json:"username"`
}

type Username struct {
	Value string `json:"_content"`
}
