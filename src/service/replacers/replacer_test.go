package replacers

import (
	"github.com/stretchr/testify/assert"
	"github.com/zored/edit/src/service/navigation"
	"strings"
	"testing"
)

func Test_replacer_Replace(t *testing.T) {
	result, err := NewReplacer().Replace(strings.NewReader(`
Hello!
It's time to replace text here.
Isn't it?
`), navigation.NewInterval(navigation.NewPosition(3, 6), navigation.NewPosition(4, 6)), "XXX")
	assert.Nil(t, err)
	assert.Equal(t, `
Hello!
It's XXX it?
`, result)
}
