package authorizer

import (
	"context"

	"github.com/lixinio/weixin/utils"
)

const (
	apiGetAllCategories = "/cgi-bin/wxopen/getallcategories" // 获取可设置的所有类目
	apiGetCategory      = "/cgi-bin/wxopen/getcategory"      // 获取已设置的所有类目(小程序)
	apiAddCategory      = "/cgi-bin/wxopen/addcategory"      // 添加类目
	apiDeleteCategory   = "/cgi-bin/wxopen/deletecategory"   // 删除类目
	apiGetOaCategory    = "/wxaapi/newtmpl/getcategory"      // 获取已设置的所有类目(公众号)
	// apiGetJSApiTicket3  = "/cgi-bin/wxopen/getcategoriesbytype"
	// apiGetJSApiTicket6  = "/cgi-bin/wxopen/modifycategory"
	// apiGetJSApiTicket7  = "/wxa/get_category" // 获取类目名称信息
)

type Category struct {
	ID        int       `json:"id"`             // 类目 ID
	Name      string    `json:"name"`           // 类目名称
	Level     int       `json:"level"`          // 类目层级
	Father    int       `json:"father"`         // 类目父级 ID
	Children  []int     `json:"children"`       // 子级类目 ID
	Sensitive int       `json:"sensitive_type"` // 是否为敏感类目（1 为敏感类目，需要提供相应资质审核；0 为非敏感类目，无需审核）
	Qualify   *struct { // sensitive_type 为 1 的类目需要提供的资质证明,通过qualify.exter_list.inner_list.name可查看资质名称。
		Exter []struct { // 资质证明列表
			Inner []struct {
				Name string `json:"name"` // 资质文件名称
				Url  string `json:"url"`  // 资质文件示例
			} `json:"inner_list"`
		} `json:"exter_list"`
		Remark string `json:"remark"`
	} `json:"qualify"`
}

/*
获取可设置的所有类目
https://developers.weixin.qq.com/doc/oplatform/openApi/OpenApiDoc/miniprogram-management/category-management/getAllCategories.html
*/
func (api *Authorizer) GetAllCategories(ctx context.Context) ([]*Category, error) {
	result := struct {
		utils.WeixinError
		Categories struct {
			Categories []*Category `json:"categories"`
		} `json:"categories_list"`
	}{}
	err := api.Client.HTTPGet(ctx, apiGetAllCategories, &result)
	if err != nil {
		return nil, err
	}
	return result.Categories.Categories, nil
}

type MpCategoryInfo struct {
	utils.WeixinError
	Categories []*struct { // 已设置的类目信息列表
		First       int    `json:"first"`        // 一级类目 ID
		FirstName   string `json:"first_name"`   // 一级类目名称
		Second      int    `json:"second"`       // 二级类目 ID
		SecondName  string `json:"second_name"`  // 二级类目名称
		AuditStatus int    `json:"audit_status"` // 审核状态（1 审核中 2 审核不通过 3 审核通过）
		AuditReason string `json:"audit_reason"` // 审核不通过的原因
	} `json:"categories"`
	Limit         int `json:"limit"`          // 一个更改周期内可以添加类目的次数
	Quota         int `json:"quota"`          // 本更改周期内还可以添加类目的次数
	CategoryLimit int `json:"category_limit"` // 最多可以设置的类目数量
}

// 获取已设置的所有类目(小程序)
// https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/2.0/api/category/getcategory.html
func (api *Authorizer) GetCategory(ctx context.Context) (*MpCategoryInfo, error) {
	var result MpCategoryInfo
	err := api.Client.HTTPGet(ctx, apiGetCategory, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

type OaCategoryInfo struct {
	utils.WeixinError
	Categories []*struct { // 已设置的类目信息列表
		ID   int    `json:"id"`   // 类目 ID
		Name string `json:"name"` // 类目名称
	} `json:"data"`
}

// 获取已设置的所有类目(公众号)
// https://developers.weixin.qq.com/doc/service/api/notify/notify/api_getcategory.html
func (api *Authorizer) GetOaCategory(ctx context.Context) (*OaCategoryInfo, error) {
	var result OaCategoryInfo
	err := api.Client.HTTPGet(ctx, apiGetOaCategory, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// 添加类目
// https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/2.0/api/category/addcategory.html
type MpCategoryItem struct {
	Key   string `json:"key"`   // 资质名称
	Value string `json:"value"` // 资质图片
}

type MpCategoryParams struct {
	First      int              `json:"first"`                // 一级类目 ID
	Second     int              `json:"second"`               // 二级类目 ID
	Certicates []MpCategoryItem `json:"certicates,omitempty"` // 资质信息列表。如果需要资质的类目，则该字段必填
}

// https://developers.weixin.qq.com/doc/oplatform/openApi/OpenApiDoc/miniprogram-management/category-management/addCategory.html
func (api *Authorizer) AddCategory(
	ctx context.Context, param []*MpCategoryParams,
) error {
	params := map[string][]*MpCategoryParams{
		"categories": param,
	}
	return api.Client.HTTPPostJson(ctx, apiAddCategory, params, nil)
}

// https://developers.weixin.qq.com/doc/oplatform/openApi/OpenApiDoc/miniprogram-management/category-management/deleteCategory.html
func (api *Authorizer) DeleteCategory(
	ctx context.Context, first, second int,
) error {
	params := map[string]int{
		"first":  first,
		"second": second,
	}
	return api.Client.HTTPPostJson(ctx, apiDeleteCategory, params, nil)
}
