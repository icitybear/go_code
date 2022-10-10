package testing

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSpec(t *testing.T) {

	// Only pass t into top-level Convey calls
	Convey("Given 2 even numbers", t, func() {
		a := 3 //2
		b := 4

		Convey("When add the two numbers", func() {
			c := a + b

			Convey("Then the result is still even", func() {
				// 期望是偶数 %2的结果 是0 就是期望 打√与打x
				So(c%2, ShouldEqual, 0)
			})
		})
	})
}
