package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/yorubanashi/languages/internal/db"
	"gopkg.in/yaml.v3"
)

func parseFilename(filename string) (string, string) {
	parts := strings.Split(filename, "/")
	artist := parts[len(parts)-2]
	song := strings.Split(parts[len(parts)-1], ".")[0]
	return artist, song
}

func songYAML(song db.Song) ([]byte, error) {
	var b bytes.Buffer
	encoder := yaml.NewEncoder(&b)
	encoder.SetIndent(2)
	err := encoder.Encode(song)
	return b.Bytes(), err
}

func processJoinedLyrics(flags *Flags, lc *LangConv, text string, song *db.Song) {
	for _, verseText := range strings.Split(text, "\n\n") {
		if len(verseText) == 0 {
			continue
		}

		lines := strings.Split(verseText, "\n")
		verse := db.Verse{Lines: make([]db.Line, len(lines)/3)}
		for i := 0; i < len(lines)/3; i++ {
			verseLine := db.Line{}
			if flags.isTraditional {
				verseLine.Sec = lines[i*3+0]
				verseLine.Pri = lc.ConvertTS(lines[i*3+0])
			} else {
				verseLine.Pri = lines[i*3+0]
			}
			verseLine.Rom = lines[i*3+1]
			verseLine.Eng = lines[i*3+2]

			verse.Lines[i] = verseLine
		}
		song.Verses = append(song.Verses, verse)
	}
}

func processSeparateLyrics(flags *Flags, lc *LangConv, text string, song *db.Song) {
	iterations := 0
	verseNum := 0
	for _, verseText := range strings.Split(text, "\n\n") {
		if len(verseText) == 0 {
			continue
		}

		if strings.Contains(verseText, "----") {
			iterations += 1
			verseNum = 0
			continue
		}

		lines := strings.Split(verseText, "\n")
		switch iterations {
		case 0:
			verse := db.Verse{Lines: make([]db.Line, len(lines))}
			for i, line := range lines {
				verseLine := db.Line{}
				if flags.isTraditional {
					verseLine.Sec = line
					verseLine.Pri = lc.ConvertTS(line)
				} else {
					verseLine.Pri = line
				}
				verse.Lines[i] = verseLine
			}
			song.Verses = append(song.Verses, verse)
		case 1:
			verse := song.Verses[verseNum]
			for i, line := range lines {
				verse.Lines[i].Rom = line
			}
		case 2:
			verse := song.Verses[verseNum]
			for i, line := range lines {
				if len(line) == 0 {
					break
				}
				verse.Lines[i].Eng = line
			}
		}

		verseNum += 1
	}
}

func processText(flags *Flags, lc *LangConv, text string) ([]byte, error) {
	artist, songname := parseFilename(flags.filename)
	song := db.Song{Artist: artist, Title: songname, Verses: []db.Verse{}}

	if flags.isJoined {
		processJoinedLyrics(flags, lc, text, &song)
	} else {
		processSeparateLyrics(flags, lc, text, &song)
	}

	return songYAML(song)
}

func main() {
	flags, err := parseFlags()
	if err != nil {
		log.Fatalln(err)
	}

	data, err := os.ReadFile(flags.filename)
	if err != nil {
		log.Fatalln(err)
	}

	lc := NewLC()
	out, err := processText(flags, lc, string(data))
	if err != nil {
		log.Fatalln(err)
	}

	// This assumes the filename has only one dot -- the file extension
	outpath := fmt.Sprintf("%s.yaml", strings.Split(flags.filename, ".")[0])
	err = os.WriteFile(outpath, out, os.ModePerm)
	if err != nil {
		log.Fatalln(err)
	}
}
