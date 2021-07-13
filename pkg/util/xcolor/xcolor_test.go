package xcolor

import (
	"fmt"
	"testing"
)

func TestBlue(t *testing.T) {
	fmt.Println(Blue("hello "))
	fmt.Println(Blue("hello ", "world"))
}
