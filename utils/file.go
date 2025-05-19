package utils

import (
	"github.com/samber/lo"
	"os"
)

func CreateAndWriteToFile(path string, content string) {
	f := lo.Must(os.Create(path))
	defer f.Close()
	f.WriteString(content)
}
