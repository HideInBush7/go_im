package tool

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/HideInBush7/go_im/pkg/util"
)

// 根据uid和时间戳生成token
func CreateToken(uid int64) string {
	str := fmt.Sprintf("%d:%d", uid, time.Now().UnixNano())
	h := md5.New()
	h.Write(util.StringToBytes(str))

	return base64.RawStdEncoding.EncodeToString(h.Sum(nil))
}

// 包装token为redis key
func GetTokenKey(uid int64, token string) string {
	return fmt.Sprintf("token:%d:%s", uid, token)
}
