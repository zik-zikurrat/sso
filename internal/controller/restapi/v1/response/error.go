package response

type Error struct {
	Code  int    `json:"code" example:"404"`
	Error string `json:"error" example:"message"`
}
