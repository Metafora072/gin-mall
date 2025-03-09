package serializer

// Response 是最基础的返回结构体
type Response struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
	Msg    string      `json:"msg"`
	Error  string      `json:"error"`
}

// TokenData 是带用户 token 信息的 Data (作为 Response 结构体的 Data 字段)
type TokenData struct {
	User  interface{} `json:"user"`
	Token string      `json:"token"`
}
