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

const tsFile = "../OpenCC/data/dictionary/TSCharacters.txt"

var ts map[string]string

func createTS() error {
	data, err := os.ReadFile(tsFile)
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
	ts = mapping
	return nil
}

func convertSimplified(text string) string {
	out := ""
	for _, char := range text {
		s := ts[string(char)]
		if s == "" {
			out += string(char)
		} else {
			out += s
		}
	}
	return out
}

func processText(artist, songname, text string, secondary bool) ([]byte, error) {
	song := db.Song{Artist: artist, Title: songname, Verses: []db.Verse{}}
	for _, verseText := range strings.Split(text, "\n\n") {
		if len(verseText) == 0 {
			continue
		}

		lines := strings.Split(verseText, "\n")
		verse := db.Verse{Lines: make([]db.Line, len(lines)/3)}
		for i := 0; i < len(lines)/3; i++ {
			verseLine := db.Line{}
			if secondary {
				verseLine.Sec = lines[i*3+0]
				verseLine.Pri = convertSimplified(lines[i*3+0])
			} else {
				verseLine.Pri = lines[i*3+0]
			}
			verseLine.Rom = lines[i*3+1]
			verseLine.Eng = lines[i*3+2]

			verse.Lines[i] = verseLine
		}
		song.Verses = append(song.Verses, verse)
	}

	var b bytes.Buffer
	encoder := yaml.NewEncoder(&b)
	encoder.SetIndent(2)
	err := encoder.Encode(song)
	return b.Bytes(), err
}

func parseFilename(filename string) (string, string) {
	parts := strings.Split(filename, "/")
	artist := parts[len(parts)-2]
	song := strings.Split(parts[len(parts)-1], ".")[0]
	return artist, song
}

func main() {
	filename := flag.String("filename", "", "Path to the file to parse")
	isTraditional := flag.Bool("traditional", false, "Whether the Chinese input is traditional")
	flag.Parse()

	if *filename == "" {
		log.Fatalln("Filename is empty")
	}

	err := createTS()
	if err != nil {
		log.Fatalln(err)
	}

	data, err := os.ReadFile(*filename)
	if err != nil {
		log.Fatalln(err)
	}

	artist, song := parseFilename(*filename)
	out, err := processText(artist, song, string(data), *isTraditional)
	if err != nil {
		log.Fatalln(err)
	}

	// This assumes the filename has only one dot -- the file extension
	outpath := fmt.Sprintf("%s.yaml", strings.Split(*filename, ".")[0])
	err = os.WriteFile(outpath, out, os.ModePerm)
	if err != nil {
		log.Fatalln(err)
	}
}
