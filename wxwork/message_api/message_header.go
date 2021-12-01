package message_api

import "strings"

type MessageHeader struct {
	ToUser                 string `json:"touser,omitempty"`
	ToParty                string `json:"toparty,omitempty"`
	ToTag                  string `json:"totag,omitempty"`
	Safe                   int    `json:"safe,omitempty"`
	EnableIDTrans          int    `json:"enable_id_trans,omitempty"`
	EnableDuplicateCheck   int    `json:"enable_duplicate_check,omitempty"`
	DuplicateCheckInterval int    `json:"duplicate_check_interval,omitempty"`
}

func (h *MessageHeader) SetSafe(safe int) *MessageHeader {
	h.Safe = safe
	return h
}

func (h *MessageHeader) SetEnableIDTrans(enableIDTrans int) *MessageHeader {
	h.EnableIDTrans = enableIDTrans
	return h
}

func (h *MessageHeader) SetEnableDuplicateCheck(enableDuplicateCheck int) *MessageHeader {
	h.EnableDuplicateCheck = enableDuplicateCheck
	return h
}

func (h *MessageHeader) SetDuplicateCheckInterval(duplicateCheckInterval int) *MessageHeader {
	h.DuplicateCheckInterval = duplicateCheckInterval
	return h
}

/**
touser	否	成员ID列表（消息接收者，多个接收者用‘|’分隔，最多支持1000个）。特殊情况：指定为@all，则向关注该企业应用的全部成员发送
toparty	否	部门ID列表，多个接收者用‘|’分隔，最多支持100个。当touser为@all时忽略本参数
totag	否	标签ID列表，多个接收者用‘|’分隔，最多支持100个。当touser为@all时忽略本参数
**/
func NewMessageHeaderByUser(user string) *MessageHeader {
	return &MessageHeader{
		ToUser: user,
	}
}

func NewMessageHeaderByUsers(users []string) *MessageHeader {
	return &MessageHeader{
		ToUser: strings.Join(users, "|"),
	}
}

func NewMessageHeaderByParty(party string) *MessageHeader {
	return &MessageHeader{
		ToParty: party,
	}
}

func NewMessageHeaderByParties(parties []string) *MessageHeader {
	return &MessageHeader{
		ToParty: strings.Join(parties, "|"),
	}
}

func NewMessageHeaderByTag(tag string) *MessageHeader {
	return &MessageHeader{
		ToTag: tag,
	}
}

func NewMessageHeaderByTags(tags []string) *MessageHeader {
	return &MessageHeader{
		ToTag: strings.Join(tags, "|"),
	}
}

func NewMessageHeaderByAll() *MessageHeader {
	return &MessageHeader{
		ToUser: "@all",
	}
}
