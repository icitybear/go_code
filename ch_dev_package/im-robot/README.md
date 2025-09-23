# im-robot

> feishu 是飞书机器人, 支持**加签**安全设置，支持**链式语法**创建消息；支持**文本（text）、富文本（post）、图片（image）、群名片（share_chat）、消息卡片（interactive）** 消息类型。
   

> 钉钉(待完善)

# feishu

- [飞书文档](https://www.feishu.cn/hc/zh-CN/articles/360024984973)

- [飞书可视化搭建工具 cardbuilder](https://open.feishu.cn/cardkit) 自动生成。

## 支持类型
- [x] 支持加签

<img src="https://p6-hera.byteimg.com/tos-cn-i-jbbdkfciu3/fb5e1dd375684dd2b9b6037d86f862b0~tplv-jbbdkfciu3-image:0:0.image" width = 50% />

- [x] 文本（text）消息

<img src="https://p1-hera.byteimg.com/tos-cn-i-jbbdkfciu3/c9c86efea1754e269dbdc5517b4d958a~tplv-jbbdkfciu3-image:0:0.image" width = 50% />

- [x] 富文本（post）消息

<img src="https://p3-hera.byteimg.com/tos-cn-i-jbbdkfciu3/661d8ee4446c47bca5ac61bfb2ef1a6f~tplv-jbbdkfciu3-image:0:0.image" width = 50% />

- [x] 图片（image）消息

<img src="https://p1-hera.byteimg.com/tos-cn-i-jbbdkfciu3/5607aa65324e4e14bd94192ba81fe0b3~tplv-jbbdkfciu3-image:0:0.image" width = 50% />

- [x] 群名片（share_chat）消息

<img src="https://p9-hera.byteimg.com/tos-cn-i-jbbdkfciu3/ba60b1c2835a4950926bb86687e183a8~tplv-jbbdkfciu3-image:0:0.image" width = 50% />

- [x] 消息卡片（interactive）消息

<img src="https://p6-hera.byteimg.com/tos-cn-i-jbbdkfciu3/4bf5072377cf4c02990ce28731634e6a~tplv-jbbdkfciu3-image:0:0.image" width = 50% />



# feishu使用
```
go get github.com/CatchZeng/feishu
```

```go
package main

import (
	"log"

	"github.com/icitybear/im-robot/feishu"
)

func main() {
	token := "6cxxxx80-xxxx-49e2-ac86-7f378xxxx960"
	secret := "k6usknqxxxxazNxxxx443d"

	client := feishu.NewClient(token, secret)

	text := feishu.NewText("文本")
	a := feishu.NewA("链接", "https://www.baidu.com/")
	at := feishu.NewAT("all")
	line := []feishu.PostItem{text, a, at}
	msg := feishu.NewPostMessage()
	msg.SetZHTitle("测试富文本 @all").
		AppendZHContent(line)

	req, resp, err := client.Send(msg)
	if err != nil {
		log.Print(err)
		return
	}
	log.Print(resp)
}
```

## 相关json数据
``` json
{
    "content": {
        "post": {
            "zh_cn": {
                "content": [
                    [
                        {
                            "tag": "text",
                            "text": "信息"
                        },
                        {
                            "tag": "a",
                            "text": "链接文本",
                            "href": "https://makeoptim.com/"
                        },
                        {
                            "tag": "at",
                            "user_id": "all"
                        }
                    ]
                ],
                "title": "标题"
            }
        }
    },
    "msg_type": "post",
    "sign": "HR7kQhgapScmp/2bfLWdYmC7C6pUV3C/pQUiS3OQDIA=",
    "timestamp": "1642561080"
}
```
``` json
{
    "config": {
        "wide_screen_mode": true,
        "enable_forward": true
    },
    "elements": [
        {
            "tag": "div",
            "text": {
                "content": "**西湖**，位于浙江省杭州市西湖区龙井路1号，杭州市区西部，景区总面积49平方千米，汇水面积为21.22平方千米，湖面面积为6.38平方千米。",
                "tag": "lark_md"
            }
        },
        {
            "actions": [
                {
                    "tag": "button",
                    "text": {
                        "content": "更多景点介绍 :玫瑰:",
                        "tag": "lark_md"
                    },
                    "url": "https://www.example.com",
                    "type": "default",
                    "value": {}
                }
            ],
            "tag": "action"
        }
    ],
    "header": {
        "title": {
            "content": "今日旅游推荐",
            "tag": "plain_text"
        }
    }
}
```

# dingding