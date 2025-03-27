package request

type Sensitive struct {
	Keywords string `json:"keywords" binding:"required"`
}
