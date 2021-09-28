package postback

import "encoding/json"

// 將 Postback data 轉為字典
func ToMap(data string) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := json.Unmarshal([]byte(data), &result)
	return result, err
}
