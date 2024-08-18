package domain

type UserInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Data struct {
	Data map[string]interface{} `json:"data"`
}

type DataKeys struct {
	Keys []string `json:"keys"`
}
