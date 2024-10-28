package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/yorubanashi/languages/internal/db"
	"gopkg.in/yaml.v3"
)

func processText(text string, secondary bool) ([]byte, error) {
	song := db.Song{Verses: []db.Verse{}}
	for _, verseText := range strings.Split(text, "\n\n") {
		if len(verseText) == 0 {
			continue
		}

		verse := db.Verse{Lines: make([]db.Line, 1)}
		verseLine := db.Line{}
		lines := strings.Split(verseText, "\n")
		for i := 0; i < len(lines)/3; i++ {
			pri := lines[i*3+0]
			rom := lines[i*3+1]
			eng := lines[i*3+2]

			verseLine.Pri += pri + "\n"
			verseLine.Rom += rom + "\n"
			verseLine.Eng += eng + "\n"
		}
		if secondary {
			verseLine.Sec = verseLine.Pri
			verseLine.Pri = ""
		}

		verse.Lines[0] = verseLine
		song.Verses = append(song.Verses, verse)
	}

	var b bytes.Buffer
	encoder := yaml.NewEncoder(&b)
	encoder.SetIndent(2)
	err := encoder.Encode(song)
	return b.Bytes(), err
}

func main() {
	filename := flag.String("filename", "", "Path to the file to parse")
	isTraditional := flag.Bool("traditional", false, "Whether the Chinese input is traditional")
	flag.Parse()

	if *filename == "" {
		log.Fatalln("Filename is empty")
	}

	data, err := os.ReadFile(*filename)
	if err != nil {
		log.Fatalln(err)
	}

	out, err := processText(string(data), *isTraditional)
	if err != nil {
		log.Fatalln(err)
	}

	// This assumes the filename has only one dot -- the file extension
	outpath := fmt.Sprintf("%s.yml", strings.Split(*filename, ".")[0])
	err = os.WriteFile(outpath, out, os.ModePerm)
	if err != nil {
		log.Fatalln(err)
	}
}
