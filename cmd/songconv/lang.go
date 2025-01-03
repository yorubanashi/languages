package main

import (
	"log"
	"os"
	"strings"
)

const TS_PATH = "../../OpenCC/data/dictionary/TSCharacters.txt"

type LangConv struct {
	ts map[string]string // Mapping from traditional Chinese to simplified
}

func NewLC() *LangConv {
	return &LangConv{}
}

func (lc *LangConv) createTS() error {
	data, err := os.ReadFile(TS_PATH)
	if err != nil {
		return err
	}

	mapping := make(map[string]string)
	for _, line := range strings.Split(string(data), "\n") {
		if len(line) == 0 {
			continue
		}

		ts := strings.Split(line, "\t")
		t := ts[0]
		// Default to the first character for now
		s := strings.Split(ts[1], " ")[0]
		mapping[t] = s
	}

	lc.ts = mapping
	return nil
}

// Convert TS (traditional -> simplified)
func (lc *LangConv) ConvertTS(text string) string {
	if lc.ts == nil {
		log.Println("Creating TS table...")
		lc.createTS()
		log.Println("TS table creation complete!")
	}

	out := ""
	for _, char := range text {
		s := lc.ts[string(char)]
		if s == "" {
			out += string(char)
		} else {
			out += s
		}
	}
	return out
}
