package yzm_test

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image/color"
	"image/png"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/fogleman/gg"
	"github.com/mojocn/base64Captcha"
)

// 设置自带的store
var store = base64Captcha.DefaultMemStore

// 生成验证码
func CaptMake() (id, b64s, ans string, err error) {
	var driver base64Captcha.Driver
	var driverString base64Captcha.DriverString

	// 配置验证码信息
	captchaConfig := base64Captcha.DriverString{
		Height:          60,
		Width:           200,
		NoiseCount:      0,
		ShowLineOptions: 2 | 4,
		Length:          4,
		Source:          "1234567890qwertyuioplkjhgfdsazxcvbnm",
		BgColor: &color.RGBA{
			R: 3,
			G: 102,
			B: 214,
			A: 125,
		},
		// Fonts: []string{"wqy-microhei.ttc"},
	}

	driverString = captchaConfig
	driver = driverString.ConvertFonts()
	captcha := base64Captcha.NewCaptcha(driver, store)
	lid, lb64s, lans, lerr := captcha.Generate()
	return lid, lb64s, lans, lerr
}

func TestYzm(t *testing.T) {

	id, b64s, ans, _ := CaptMake()
	spew.Println(id)
	spew.Println(b64s)
	fmt.Println("ans")
	spew.Println(ans)

}

// 图片64位问题
func TestBase64png(t *testing.T) {

	// data:image/png;base64
	base64Str := "data:image/png;base64,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" // 示例 Base64 字符串

	// 清理前缀
	parts := strings.Split(base64Str, ",")
	if len(parts) == 2 {
		base64Str = parts[1]
	}

	// 解码
	imgBytes, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		fmt.Println("解码失败:", err)
		return
	}

	// 验证 PNG 文件头
	pngHeader := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}
	if len(imgBytes) < 8 || !bytes.Equal(imgBytes[:8], pngHeader) {
		fmt.Println("无效的 PNG 文件头")
		return
	}

	// 生成文件名并保存
	fileName := fmt.Sprintf("%d.png", time.Now().UnixNano())
	if err := os.WriteFile(fileName, imgBytes, 0666); err != nil {
		fmt.Println("保存失败:", err)
		return
	}

	fmt.Println("保存成功:", fileName)

	// 解码 Base64 转图片  base64要去掉前缀
	img, err := png.Decode(bytes.NewReader(imgBytes))
	if err != nil {
		fmt.Println("无效的img指针", err)
		return
	}
	fmt.Println("新图片")
	// 创建新画布并添加边距 默认透明画布
	// 创建新画布（宽度增加边距）
	newWidth := img.Bounds().Dx() + 20 // 左侧留20像素边距
	dc := gg.NewContext(newWidth, img.Bounds().Dy())

	// 填充与系统一致的背景色（示例为浅蓝色RGBA: 3,102,214,255）
	dc.SetRGBA255(3, 102, 214, 255)
	dc.Clear() // 必须调用Clear填充背景

	// 绘制原图到右侧（偏移20像素）
	dc.DrawImage(img, 20, 0)
	// dc.SavePNG("output.png")

	// 重新编码为 Base64
	buffer := new(bytes.Buffer)
	dc.EncodePNG(buffer)
	b64 := buffer.Bytes()
	b64s := base64.StdEncoding.EncodeToString(b64)
	fmt.Println(b64s)
	fileName2 := fmt.Sprintf("new%d.png", time.Now().UnixNano())
	if err := os.WriteFile(fileName2, b64, 0666); err != nil {
		fmt.Println("保存失败:", err)
		return
	}
}

// 自定义驱动（解决贴边问题） 进阶使用
type CustomDriver struct {
	base64Captcha.DriverString
}

func (d *CustomDriver) DrawCaptcha(content string) (item base64Captcha.Item, err error) {
	paddedContent := " " + content // 左侧添加字符空白

	// 调用父类方法生成原始 Item
	item, err = d.DriverString.DrawCaptcha(paddedContent)
	return item, err
}

// 生成验证码
func TestYzm2(t *testing.T) {
	// 配置验证码信息
	driverConfig := base64Captcha.DriverString{
		Height:          36,
		Width:           200,
		NoiseCount:      0,
		ShowLineOptions: 0, // 2 | 4
		Length:          4,
		Source:          "1234567890qwertyuioplkjhgfdsazxcvbnm",
		BgColor: &color.RGBA{
			R: 3,
			G: 102,
			B: 214,
			A: 125,
		},
		Fonts: []string{"wqy-microhei.ttc"},
	}
	// 使用自定义驱动
	customDriver := &CustomDriver{driverConfig}
	customDriver.ConvertFonts()

	captcha := base64Captcha.NewCaptcha(customDriver, store)
	lid, lb64s, lans, _ := captcha.Generate()

	spew.Println(lid)
	spew.Println(lb64s)
	fmt.Println("答案")
	spew.Println(lans)
	res := captcha.Verify(lid, lans, true) // 如果没使用redis 使用内存就要同一个实例验证
	fmt.Println(res)
}
