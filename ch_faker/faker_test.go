package faker_test

import (
	"fmt"
	"testing"

	"github.com/go-faker/faker/v4"
)

type SomeStructWithTags struct {
	Dt        string  `json:"dt"  faker:"date"`
	Name      string  `faker:"name"`
	Latitude  float32 `faker:"lat"`
	Longitude float32 `faker:"long"`
	Val       string  `faker:"oneof: 1, 2"`
	// Val  string `faker:"year"`
	Val1 string `faker:"year"`
	// 新版faker才支持faker标签 oneof
}

type SomeTT struct {
	List []*SomeStructWithTags `faker:"slice_len=5"` // 不支持
}

func TestFaker(t *testing.T) {
	a := SomeStructWithTags{}
	err := faker.FakeData(&a)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v", a)

	b := SomeTT{}
	err = faker.FakeData(&b)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v", b)
}
