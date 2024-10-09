package svelte

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	htmlCommentLeft  = "<!-- "
	htmlCommentRight = " -->"
	pageSuffix       = "/+page.svelte"
)

type Link struct {
	Addr string `json:"addr"`
	Name string `json:"name"`
}

type Dir struct {
	// TODO: Flesh out the nested structure -- for now this isn't needed.
	Dirs map[string]*Dir `json:"dirs,omitempty"`
	Link *Link           `json:"link,omitempty"`
}

// From a given path, get a valid "addressable" string we can use as an href in an <a> tag.
func getAddressable(base, path string) string {
	return strings.TrimPrefix(strings.TrimSuffix(strings.TrimPrefix(path, base), pageSuffix), "/")
}

func trimHTMLComment(str string) string {
	return strings.TrimSuffix(strings.TrimPrefix(str, htmlCommentLeft), htmlCommentRight)
}

// TODO: This probably belongs in the "db" module, maybe?
func Walk(root, lang string) (*Dir, error) {
	dir := &Dir{Dirs: make(map[string]*Dir)}
	base := fmt.Sprintf("%s/%s", root, lang)
	err := filepath.Walk(base, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Only register pages, skip auxilary files since they don't count as a "page"
		if !strings.HasSuffix(path, pageSuffix) {
			return nil
		}

		addr := getAddressable(base, path)
		// Skip base path, since we'll always have a link back to it.
		// Also skip slugs, since it doesn't make sense to hyperlink to them w/o the value.
		if addr == "" || strings.Contains(addr, "[") {
			return nil
		}

		bytes, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		title := trimHTMLComment(strings.Split(string(bytes), "\n")[0])

		tempDir := dir
		for _, part := range strings.Split(addr, "/") {
			if tempDir.Dirs[part] == nil {
				tempDir.Dirs[part] = &Dir{Dirs: make(map[string]*Dir)}
			}
			tempDir = tempDir.Dirs[part]
		}
		tempDir.Link = &Link{Name: title, Addr: addr}

		return nil
	})

	return dir, err
}
