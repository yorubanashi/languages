package main

import (
	"flag"
	"fmt"
)

type Flags struct {
	filename      string
	isJoined      bool
	isTraditional bool
}

func parseFlags() (*Flags, error) {
	filename := flag.String("filename", "", "Path to the file to parse")
	isTraditional := flag.Bool("traditional", false, "Whether the Chinese input is traditional")
	isJoined := flag.Bool("joined", false, "Whether the lyric types are all joined into the same verse")
	flag.Parse()

	if *filename == "" {
		return nil, fmt.Errorf("filename is empty")
	}

	return &Flags{
		filename:      *filename,
		isJoined:      *isJoined,
		isTraditional: *isTraditional,
	}, nil
}
