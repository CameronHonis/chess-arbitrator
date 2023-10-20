package server

type Prompt struct {
	Type    PromptType
	Content interface{}
}

type InitClientPromptContent struct {
	PublicKey string
}
