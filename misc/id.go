package misc

import (
	"math/rand"
	"strings"
	"time"
	"unicode"

	"github.com/google/uuid"
)

// IDGenerator id 生成器
type IDGenerator struct {
}

// GenerateRandomString 生成随机字符串
func (gen IDGenerator) GenerateRandomString(length int, alphanum []byte) string {
	if alphanum == nil {
		alphanum = []byte("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")
	}

	rand.Seed(time.Now().UTC().UnixNano())
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = alphanum[rand.Intn(len(alphanum))]
	}
	return string(result)
}

// GenerateUUID 生成UUID
func (gen IDGenerator) GenerateUUID(includeSeparator bool) string {
	u4 := uuid.New()

	vid := u4.String()
	if !includeSeparator {
		vid = strings.Replace(vid, "-", "", -1)
	}

	return vid
}

// BighumpToUnderscore 大驼峰参数转换成下划线
func (gen IDGenerator) BighumpToUnderscore(param string) string {
	// ParamValue -> param_value
	transParam := ""
	for i, ch := range param {
		if unicode.IsUpper(ch) {
			if i > 0 && !unicode.IsUpper(rune(param[i-1])) {
				transParam += "_"
			}
		}
		transParam += string(unicode.ToLower(ch))
	}
	return transParam
}

// UnderscoreToBighump 下划线参数转换成大驼峰
func (gen IDGenerator) UnderscoreToBighump(param string) string {
	// param_value -> ParamValue
	transParam := ""
	capital := false
	for _, ch := range param {
		if ch == '_' {
			capital = true
			continue
		}
		if capital {
			transParam += string(unicode.ToUpper(ch))
		} else {
			transParam += string(unicode.ToLower(ch))
		}
	}
	return transParam
}
