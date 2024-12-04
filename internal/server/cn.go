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

func (s *Server) cnRoutes() map[string]HandlerFunc {
	return map[string]HandlerFunc{
		"/songs":   s.cnSongsHandler,
		"/artists": s.cnArtistsHandler,
	}
}

type SongRequest struct {
	Title  string `yaml:"title" json:"title"`   // Song title, in the primary language
	Artist string `yaml:"artist" json:"artist"` // Main artist
}
type SongResponse struct {
	Songs []db.Song `json:"songs,omitempty"`
}

// If one or both of title, artist is omitted, the endpoint will fallback to returning everything.
func (s *Server) cnSongs(_ context.Context, req *SongRequest) (*SongResponse, error) {
	var songs []db.Song
	var err error
	if len(req.Title) > 0 && len(req.Artist) > 0 {
		path := fmt.Sprintf("%s/%s/%s.yaml", s.config.DBPaths.Songs.Base, req.Artist, req.Title)
		var song db.Song
		err = db.FetchYAML(path, &song)
		songs = []db.Song{song}
	} else {
		err = db.FetchYAML(s.config.DBPaths.Songs.Indexed, &songs)
	}
	return &SongResponse{Songs: songs}, err
}

type ArtistRequest struct{}
type ArtistResponse struct {
	Artists []db.Line `json:"artists,omitempty"`
}

func (s *Server) cnArtists(_ context.Context, _ *ArtistRequest) (*ArtistResponse, error) {
	var artists []db.Line
	err := db.FetchYAML(s.config.DBPaths.Artists, &artists)
	return &ArtistResponse{Artists: artists}, err
}

// TODO: Some of this logic probably belongs in the internal/db module?
//
// index walks through the data/songs directory and saves an appendix in another file.
func (s *Server) index() {
	s.logger.Println("Indexing CN songs...")
	defer s.logger.Println("Indexing CN songs complete!")

	out := []db.Song{}
	err := filepath.Walk(s.config.DBPaths.Songs.Base, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		parts := strings.Split(path, "/")
		if len(parts) == 4 && strings.HasSuffix(path, ".yaml") {
			artist := parts[2]
			title := strings.TrimSuffix(parts[3], ".yaml")
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

	file, err := os.OpenFile(s.config.DBPaths.Songs.Indexed, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
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
