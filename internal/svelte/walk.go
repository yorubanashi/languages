package svelte

import (
	"fmt"
	"os"
	"path/filepath"
)

func Walk() {
	// TODO: Use config.yaml / function arguments to grab values
	root := "../yorubanashi.github.io/src/routes"
	lang := "cn"

	actualRoot := fmt.Sprintf("%s/%s", root, lang)
	err := filepath.Walk(actualRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		fmt.Println(path) // Print the path of each file/directory

		return nil
	})

	if err != nil {
		fmt.Println("Error:", err)
	}
}
