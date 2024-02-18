package time_test

import (
	"fmt"
	"testing"
	"time"
)

func TestParse(t *testing.T) {
	dtTime, err := time.Parse("2006-01-02", "2024/02/03")
	fmt.Println(dtTime, err)

	// dtTime, err := time.Parse("2006-01-02", "2024年02月03日")
	// fmt.Println(dtTime, err)

}
