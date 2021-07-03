package scaffold

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)
func TestScaffold(t *testing.T) {

	tempDir, err := filepath.Abs("./out")
	if err!=nil{
		panic(err)
	}

	fmt.Printf("tempDir:%s\n", tempDir)
	assert.NoError(t, New(true).Generate(tempDir))
}
