package methods

type InitializeParams struct {
	ProcessId  int        `json:"processId"`
	ClientInfo ClientInfo `json:"clientInfo"`
	Locale     string     `json:"locale"`
	RootPath   string     `json:"rootPath"`
	RootUri    string     `json:"rootUri"`
}

type IntializeResult struct {
	Capabilities Capabilities `json:"capabilities"`
	ServerInfo   ServerInfo   `json:"serverInfo"`
}

type ClientInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type CompletionProvider struct {
}

type Capabilities struct {
	CompletionProvider CompletionProvider `json:"completionProvider"`
	TextDocumentSync   int                `json:"textDocumentSync"`
}

func Initialize(params RequestMessage) (IntializeResult, error) {
	return IntializeResult{
		Capabilities: Capabilities{
			CompletionProvider: CompletionProvider{},
			TextDocumentSync:   1,
		},
		ServerInfo: ServerInfo{
			Name:    "templ-html-htmx-lsp",
			Version: "0.0.1",
		},
	}, nil
}
