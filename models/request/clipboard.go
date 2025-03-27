package request

type Clipboard struct {
	Title              string `json:"title" binding:"required"`
	Content            string `json:"content" binding:"required"`
	SyntaxHighlighting string `json:"syntax_highlighting" binding:"required"`
	Expire             string `json:"expire" binding:"required"`
	Snapchat           bool   `json:"snapchat"`
}
