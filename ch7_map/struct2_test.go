package my_map

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/mitchellh/mapstructure"
)

type CallResult struct {
	Code string      `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func TestStructToMapX(t *testing.T) {

	body := `{"success":true,"code":"0","msg":"","data":{"page":{"pageNum":1,"pageSize":10,"size":1,"startRow":1,"endRow":1,"total":1,"pages":1,"list":[{"id":12591,"name":"advertisement2025-03-25","fileSize":656,"fileNum":1,"remark":"","status":1,"createUser":"citybear09","isDelete":0,"createTime":1742869775,"fileType":0,"ossUrl":"advertisement2025-03-2520250325102936.CSV","fileSizeStr":"0.64KB","dataType":0,"loginname":"citybear09","uuid":null,"approval":false,"approvalStatus":null,"systemFlag":null}],"prePage":0,"nextPage":0,"isFirstPage":true,"isLastPage":true,"hasPreviousPage":false,"hasNextPage":false,"navigatePages":8,"navigatepageNums":[1],"navigateFirstPage":1,"navigateLastPage":1,"firstPage":1,"lastPage":1},"permissionLevel":0},"extra":""}`

	var result CallResult
	err := json.Unmarshal([]byte(body), &result)
	fmt.Println(0)
	if err != nil {
		spew.Println(err)
		return
	}
	spew.Println(result)
	fmt.Println(1)
	if result.Code != "" && result.Code != "0" {
		spew.Println("code != 0")
		return
	}
	fmt.Println(2)
	out := make(map[string]any, 0)
	err = mapToStruct(result.Data, out)
	if err != nil {
		spew.Println(err)
		return
	}
	spew.Println(out)
}

func mapToStruct(input, output interface{}) error {
	config := &mapstructure.DecoderConfig{
		Metadata:         nil,
		Result:           &output,
		TagName:          "json",
		WeaklyTypedInput: true,
	}
	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}
	return decoder.Decode(input)
}
