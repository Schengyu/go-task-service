package auxiliary_method

import (
	"fmt"
	"go-task-service/core/curd_methods"
	"strconv"
	"strings"
)

// 生成单个属性的拼接格式 "ParentID:AttrID"
func GenerateSkuPart(parentID, attrID int) string {
	return fmt.Sprintf("%d:%d", parentID, attrID)
}

// 生成完整的 SKUCode
func GenerateSkuCode(attrIDs []int64) (string, error) {
	var skuParts []string

	for _, attrID := range attrIDs {
		// 查询属性信息
		AttributeInfo, err := curd_methods.GetAttributeById(int(attrID))
		if err != nil {
			return "", err
		}

		// 拼接 "ParentID:AttrID"
		skuParts = append(skuParts, GenerateSkuPart(int(AttributeInfo.ParentID), int(attrID)))
	}

	// 用_连接所有部分
	return strings.Join(skuParts, "_"), nil
}

// 解析 SKUCode，提取所有的 AttrID 放入切片中
func ParseAttrIDsFromSkuCode(skuCode string) ([]int, error) {
	var attrIDs []int

	// 先用 _ 分割出每一段
	parts := strings.Split(skuCode, "_")
	for _, part := range parts {
		// 每一段是 ParentID:AttrID 的格式
		ids := strings.Split(part, ":")
		if len(ids) != 2 {
			return nil, fmt.Errorf("格式不正确: %s", part)
		}

		// 取第二个，也就是 AttrID
		attrID, err := strconv.Atoi(ids[1])
		if err != nil {
			return nil, fmt.Errorf("AttrID 转换失败: %v", err)
		}

		attrIDs = append(attrIDs, attrID)
	}

	return attrIDs, nil
}

// 调用支付链接返回数据结构体
type PaySysCallbackData struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		OrderNo string `json:"order_no"`
	} `json:"data"`
}

// 定义请求返回值结构体
type DeductRespData struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
