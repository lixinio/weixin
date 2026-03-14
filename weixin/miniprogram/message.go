package miniprogram

type DataValue struct {
	Value string `json:"value"`
}

type JsonMessage struct {
	// 订阅模板id
	TemplateID string `json:"template_id"`
	// 跳转页面
	Page string `json:"page"`
	// 接收者（用户）的 openid
	ToUser string `json:"touser"`
	// 模板内容
	Data map[string]DataValue `json:"data"`
	// 跳转小程序类型
	MiniProgramState MiniProgramState `json:"miniprogram_state"`
	// 进入小程序查看”的语言类型
	Lang string `json:"lang"`
}

func defaultJsonMessage(tmplID string, toUser string, data map[string]DataValue) *JsonMessage {
	return &JsonMessage{
		TemplateID: tmplID,
		ToUser:     toUser,
		Data:       data,
	}
}
