package misc

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/tidwall/gjson"
)

// traverseMap 遍历, 如果string太长折叠
func traverseMap(mp map[string]interface{}) (res map[string]interface{}) {
	res = map[string]interface{}{}
	for key, element := range mp {
		switch value := element.(type) {
		case string:
			if len(value) >= 1000 {
				res[key] = fmt.Sprintf("[%d bytes data too log to be folded]", len(value))
			} else {
				res[key] = element
			}
		case map[string]interface{}:
			res[key] = traverseMap(value)
		case []interface{}:
			temp := []interface{}{}
			for _, ele := range value {
				switch v := ele.(type) {
				case map[string]interface{}:
					temp = append(temp, traverseMap(v))
				case string:
					if len(v) >= 1000 {
						temp = append(temp, fmt.Sprintf("[%d bytes data too log to be folded]", len(v)))
					} else {
						temp = append(temp, v)
					}
				default:
					temp = append(temp, v)
				}
			}
			res[key] = temp
		default:
			res[key] = value
		}
	}
	return res
}

// FoldLog 折叠log中的长数据，比如图片 base64 等等
func FoldLog(content []byte) string {
	data := gjson.ParseBytes(content).Map()
	mp := map[string]interface{}{}
	for key, value := range data {
		mp[key] = value.Value()
	}
	res := traverseMap(mp)
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	encoder.Encode(res)
	bs := buffer.Bytes()
	return string(bs[0 : len(bs)-1])
}
