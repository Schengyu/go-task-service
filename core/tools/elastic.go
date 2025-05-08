package tools

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"go-task-service/cmd/global"
	"log"
)

type InventoryResponse []Inventory

type Inventory struct {
	Assets       []Asset       `json:"assets"`
	Descriptions []Description `json:"descriptions"`
}

type Asset struct {
	Amount     string `json:"amount"`
	AppID      int    `json:"appid"`
	AssetID    string `json:"assetid"`
	ClassID    string `json:"classid"`
	ContextID  string `json:"contextid"`
	InstanceID string `json:"instanceid"`
}

type Description struct {
	AppID          int         `json:"appid"`
	ClassID        string      `json:"classid"`
	InstanceID     string      `json:"instanceid"`
	Name           string      `json:"name"`
	NameColor      string      `json:"name_color"`
	MarketName     string      `json:"market_name"`
	MarketHashName string      `json:"market_hash_name"`
	IconURL        string      `json:"icon_url"`
	Type           string      `json:"type"`
	Tradable       int         `json:"tradable"`
	TradableTime   *string     `json:"tradable_time"`
	Marketable     int         `json:"marketable"`
	FraudWarnings  interface{} `json:"fraudwarnings"`
	Actions        []Action    `json:"actions"`
	MarketActions  []Action    `json:"market_actions"`
	Tags           []Tag       `json:"tags"`
}

type Action struct {
	Link string `json:"link"`
	Name string `json:"name"`
}

type Tag struct {
	Category              string `json:"category"`
	InternalName          string `json:"internal_name"`
	LocalizedCategoryName string `json:"localized_category_name"`
	LocalizedTagName      string `json:"localized_tag_name"`
	Color                 string `json:"color,omitempty"` // 有些Tag不包含 color 字段
}

// BuildMultiTermQuery 构建多个 term 条件的查询语句
func BuildMultiTermQuery(terms map[string]interface{}) map[string]interface{} {
	mustClauses := []map[string]interface{}{}
	for field, value := range terms {
		mustClauses = append(mustClauses, map[string]interface{}{
			"term": map[string]interface{}{
				field: map[string]interface{}{
					"value": value,
				},
			},
		})
	}
	return map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": mustClauses,
			},
		},
	}
}

func ExecuteESQuery(index string, query map[string]interface{}) (InventoryResponse, error) {
	// 构造查询体
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, fmt.Errorf("encode query error: %w", err)
	}

	// 执行查询
	res, err := global.ESClient.Search(
		global.ESClient.Search.WithContext(context.Background()),
		global.ESClient.Search.WithIndex(index),
		global.ESClient.Search.WithBody(&buf),
		global.ESClient.Search.WithTrackTotalHits(true),
		global.ESClient.Search.WithPretty(),
	)
	if err != nil {
		return nil, fmt.Errorf("search error: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("response error: %s", res.String())
	}

	// 解析响应
	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response error: %w", err)
	}

	hitData := result["hits"].(map[string]interface{})["hits"].([]interface{})
	results := make([]map[string]interface{}, 0, len(hitData))
	for _, h := range hitData {
		source := h.(map[string]interface{})["_source"].(map[string]interface{})
		results = append(results, source)
	}
	marshal, err := json.Marshal(results)
	if err != nil {
		return nil, err
	}

	var inventory InventoryResponse
	err = json.Unmarshal(marshal, &inventory)
	if err != nil {
		log.Fatal(err)
	}

	return inventory, nil
}
