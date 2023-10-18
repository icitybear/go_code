package ch_wenjian

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"
)

const filename = "./preferences.json"

type Preferences struct {
	Name    string /* 配置别名 */
	Home    string /* 工作目录 */
	Workers int32  /* 工作单元数 */
}

func TestWenjian(t *testing.T) {

	// 结构体转换为 JSON 对象，并且只有结构体内变量必须首字母大写，才可被导出的字段转化输出
	p := Preferences{Name: "ccc", Home: "/home", Workers: 16}
	// 可以结构体指针&p 也可结构体实例p  接收任意类型的接口数据 v any 不想改变结构体的值 就不用指针
	marshal, err := json.Marshal(p) //序列化结构体 成json []type
	if err != nil {
		fmt.Printf("json marshal err: %v", err)
		return
	}
	// 直接用ioutil.WriteFile  整个文件 会创建并保存  err因为前面已经声明了 所以不用err:
	err = ioutil.WriteFile(filename, marshal, 0644)
	if err != nil {
		fmt.Printf("write file err: %v", err)
		return
	}
	// 直接用ioutil.ReadFile 返回的是[]byte
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("read file err: %v", err)
		return
	}

	rp := Preferences{}
	// 这里必须指针去接收，因为想改变这个结构体的值
	err = json.Unmarshal(data, &rp)
	if err != nil {
		fmt.Printf("json unmarshal err: %v", err)
		return
	}
	fmt.Printf("preferences is: %v", rp)
}
