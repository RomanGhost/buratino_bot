package data

type AccessKey struct {
	ID        int    `json:"id,string"`
	Name      string `json:"name,omitempty"`
	Password  string `json:"password"`
	Port      int    `json:"port"`
	Method    string `json:"method"`
	AccessURL string `json:"accessUrl"`
	DataLimit *struct {
		Bytes int64 `json:"bytes"`
	} `json:"dataLimit,omitempty"`
	UsedBytes int64 `json:"usedBytes,omitempty"`
}

type KeyConnectData struct {
	ID          int
	Name        string
	ConnectData string
}
