package dchest_test

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/dchest/captcha"
)

func TestYzm2(t *testing.T) {
	// captcha.New() 是指定NewLen(默认长度)
	id := captcha.NewLen(4)
	// id := "csx"
	digits := []byte{1, 5, 7, 9, 3}              // 指定生成的字符串
	img := captcha.NewImage(id, digits, 120, 36) // 返回图片*Image  NewAudio
	// 因为id可以自己指定 所以Reload支持指定刷新
	var buf bytes.Buffer
	// captcha.WriteImage(&buf, id, 120, 40) // 调用的就是 NewImage 再以png写入
	img.WriteTo(&buf)
	b64 := base64.StdEncoding.EncodeToString(buf.Bytes())

	fmt.Println(id)
	fmt.Println("data:image/png;base64," + b64)

	res := captcha.Verify(id, digits) // 要使用缓存配合才能验证 // 建议使用mojocn/base64Captcha
	fmt.Println(res)
}
