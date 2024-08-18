package domain

import "vk_tarantool_project/internal/pkg/jsonMap"

type UserInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Data struct {
	Data jsonMap.Obj `json:"data"`
}

type DataKeys struct {
	Keys []string `json:"keys"`
}
