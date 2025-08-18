package wxwork_suite

const (
	EventSubTypeCreateUser  = "create_user"
	EventSubTypeUpdateUser  = "update_user"
	EventSubTypeDeleteUser  = "delete_user"
	EventSubTypeCreateParty = "create_party"
	EventSubTypeUpdateParty = "update_party"
	EventSubTypeDeleteParty = "delete_party"
	EventSubTypeUpdateTag   = "update_tag"
)

type EventCreateUser struct {
	Event
	AuthCorpId     string
	UserID         string
	OpenUserID     string
	Department     string
	MainDepartment int
	IsLeaderInDept string
	Gender         int
	Status         int
	Alias          string
	// 职位信息。长度为0~64个字节，仅通讯录应用可获取
	// Position       string
	// 座机，仅通讯录管理应用可获取
	// Telephone      string
	// 头像url。注：如果要获取小图将url最后的”/0”改成”/100”即可，仅通讯录管理应用可获取
	// Avatar         string
	// 邮箱，仅通讯录管理应用可获取
	// Email          string
	// 手机号码，仅通讯录管理应用可获取
	// Mobile         string
	// 2020年6月30日起，对所有历史第三方应用不再返回真实name，使用userid代替name，后续第三方仅通讯录应用可获取，第三方页面需要通过通讯录展示组件来展示名字
	// Name           string
}

type EventUpdateUser struct {
	Event
	AuthCorpId     string
	UserID         string
	OpenUserID     string
	NewUserID      string // new
	Department     string
	MainDepartment int
	IsLeaderInDept string
	Gender         int
	Status         int
	Alias          string
}

type EventDeleteUser struct {
	Event
	AuthCorpId string
	UserID     string
	OpenUserID string
}

type EventCreateParty struct {
	Event
	AuthCorpId string
	// 部门名称，此字段从2019年12月30日起，对新创建第三方应用不再返回，2020年6月30日起，对所有历史第三方应用不再返回真实Name字段，
	// 使用Id字段代替Name字段，后续第三方仅通讯录应用可获取，第三方页面需要通过通讯录展示组件来展示名字。回收后普通第三方应用name变更不再回调
	Name     string
	ParentId int
	Order    int
	Id       int
}

type EventUpdateParty struct {
	Event
	AuthCorpId string
	Name       string
	ParentId   int
	Order      int
	Id         int
}

type EventDeleteParty struct {
	Event
	AuthCorpId string
	Id         int
}

type EventUpdateTag struct {
	Event
	AuthCorpId    string
	Name          string
	TagId         int
	AddUserItems  string
	DelUserItems  string
	AddPartyItems string
	DelPartyItems string
}
