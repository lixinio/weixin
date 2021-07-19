package user_api

import "github.com/lixinio/weixin/utils"

type UserInfo struct {
	utils.CommonError

	UserID     string `json:"userid"`
	Name       string `json:"name"`
	Alias      string `json:"alias"`     // 别名；第三方仅通讯录应用可获取
	Mobile     string `json:"mobile"`    // 手机号码；第三方仅通讯录应用可获取
	Email      string `json:"email"`     // 邮箱；第三方仅通讯录应用可获取
	Position   string `json:"position"`  // 职务信息；第三方仅通讯录应用可获取
	AvatarURL  string `json:"avatar"`    // NOTE：如果要获取小图将url最后的”/0”改成”/100”即可。
	Telephone  string `json:"telephone"` // 座机；第三方仅通讯录应用可获取
	Gender     string `json:"gender"`    // 性别
	Status     int    `json:"status"`    // 成员激活状态
	Department []int  `json:"department"`
}

// UserGender 用户性别
const (
	// UserGenderUnspecified 性别未定义
	UserGenderUnspecified string = "0"
	UserGenderMale        string = "1"
	UserGenderFemale      string = "2"
)

// UserStatus 用户激活信息
//
// 已激活代表已激活企业微信或已关注微工作台（原企业号）。
// 未激活代表既未激活企业微信又未关注微工作台（原企业号）。
const (
	UserStatusActivated   int = 1 // 已激活
	UserStatusDeactivated int = 2 // 已禁用
	UserStatusUnactivated int = 4 // 未激活
)
