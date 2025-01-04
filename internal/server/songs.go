package server

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/yorubanashi/languages/internal/db"
	"gopkg.in/yaml.v3"
)

func (s *Server) songRoutes() map[string]HandlerFunc {
	return map[string]HandlerFunc{
		"/songs":   s.songsHandler,
		"/artists": s.artistsHandler,
	}
}

type SongRequest struct {
	Language string `yaml:"language" json:"language"` // CN, JP
	Title    string `yaml:"title" json:"title"`       // Song title, in the primary language
	Artist   string `yaml:"artist" json:"artist"`     // Main artist
}
type SongResponse struct {
	Songs []db.Song `json:"songs,omitempty"`
}

// If one or both of title, artist is omitted, the endpoint will fallback to returning everything.
func (s *Server) songs(_ context.Context, req *SongRequest) (*SongResponse, error) {
	key := fmt.Sprintf("songs+%s+%s+%s", req.Language, req.Artist, req.Title)
	if val, ok := s.cache[key]; ok {
		songs, ok := val.([]db.Song)
		if ok {
			return &SongResponse{Songs: songs}, nil
		}
	}

	var songs []db.Song
	var err error
	if len(req.Title) > 0 && len(req.Artist) > 0 {
		var song db.Song
		err = db.FetchYAML(s.config.SongPath(req.Language, req.Artist, req.Title), &song)
		songs = []db.Song{song}
	} else {
		err = db.FetchYAML(s.config.IndexedSongsPath(req.Language), &songs)
	}

	if err == nil {
		s.cache[key] = songs
	}
	return &SongResponse{Songs: songs}, err
}

type ArtistRequest struct{}
type ArtistResponse struct {
	Artists []db.Line `json:"artists,omitempty"`
}

func (s *Server) artists(_ context.Context, _ *ArtistRequest) (*ArtistResponse, error) {
	var artists []db.Line
	err := db.FetchYAML(s.config.DBPaths.Artists, &artists)
	return &ArtistResponse{Artists: artists}, err
}

// TODO: Some of this logic probably belongs in the internal/db module?
//
// index walks through the data/songs directory and saves an appendix in another file.
func (s *Server) indexAll() {
	for _, lang := range []string{"cn"} {
		s.index(lang)
	}
}

func (s *Server) index(lang string) {
	s.logger.Printf("Indexing %s songs...\n", strings.ToUpper(lang))
	defer s.logger.Printf("Indexing %s songs complete!\n", strings.ToUpper(lang))

	out := []db.Song{}
	err := filepath.Walk(s.config.SongBasePath(lang), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		parts := strings.Split(path, "/")
		if len(parts) == 5 && strings.HasSuffix(path, ".yaml") {
			artist := parts[3]
			title := strings.TrimSuffix(parts[4], ".yaml")
			out = append(out, db.Song{Title: title, Artist: artist})
		}
		return nil
	})
	if err != nil {
		s.logger.Fatal(err)
	}

	// Marshal the struct into YAML
	data, err := yaml.Marshal(out)
	if err != nil {
		s.logger.Fatal(err)
	}

	file, err := os.OpenFile(s.config.IndexedSongsPath(lang), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		s.logger.Fatal(err)
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		s.logger.Fatal(err)
	}

	err = file.Sync()
	if err != nil {
		s.logger.Fatal(err)
	}
}
