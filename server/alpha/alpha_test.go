package alpha

import (
	"fmt"
	"testing"
)

func TestGetAlpha(t *testing.T) {
	tmp, _ := GetAlphaStock("PDD")
	fmt.Println(tmp)

}
