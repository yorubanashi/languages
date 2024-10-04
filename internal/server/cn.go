package server

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/yorubanashi/languages/internal/db"
	"gopkg.in/yaml.v3"
)

func (s *Server) cnRoutes() map[string]HandlerFunc {
	return map[string]HandlerFunc{
		"/songs": s.cnSongsHandler,
	}
}

type Chinese struct {
	CN string `yaml:"cn" json:"cn,omitempty"` // Chinese (China)
	TW string `yaml:"tw" json:"tw,omitempty"` // Chinese (Taiwan)
	PY string `yaml:"py" json:"py"`           // Pinyin
	EN string `yaml:"en" json:"en"`           // English
}

type NamedChinese struct {
	Chinese
	Name string `yaml:"name" json:"name"`
}

type Song struct {
	Title    string         `yaml:"title" json:"title"`
	Artist   string         `yaml:"artist" json:"artist"`
	Featured []string       `yaml:"featured" json:"featured,omitempty"`
	Verses   []NamedChinese `yaml:"verses" json:"verses"`
	Order    []string       `yaml:"order" json:"order"`
}

type ArtistResponse struct {
	Artists []Chinese `json:"artists,omitempty"`
	Error   error     `json:"error,omitempty"`
}

type SongRequest struct {
	Name string `json:"name,omitempty"`
}
type SongResponse struct {
	Songs []Song `json:"songs,omitempty"`
}

func (s *Server) cnSongs(ctx context.Context, req *SongRequest) (*SongResponse, error) {
	var songs []Song
	err := db.FetchYAML(s.config.DBPaths.Songs.Indexed, &songs)
	return &SongResponse{Songs: songs}, err
}
func (s *Server) cnSongsHandler(ctx context.Context, decode func(interface{}) error) (interface{}, error) {
	in := &SongRequest{}
	err := decode(in)
	if err != nil {
		return nil, err
	}
	return s.cnSongs(ctx, in)
}

// func (s *Server) cnSongs(w http.ResponseWriter, r *http.Request) {
// 	res := SongResponse{}
// 	err := db.FetchYAML(s.config.DBPaths.Songs.Indexed, &res.Songs)
// 	if err != nil {
// 		res.Error = err
// 		writeJSON(w, 500, res)
// 		return
// 	}

// 	writeJSON(w, 200, res)
// 	return
// }

// func (s *Server) cnSongsHandler(w http.ResponseWriter, r *http.Request) {
// 	parts := strings.Split(r.URL.Path, "/")
// 	if len(parts) == 2 {
// 		s.cnSongs(w, r)
// 	} else if len(parts) == 4 {
// 	} else {
// 		w.WriteHeader(404)
// 	}
// }

func (s *Server) cnArtistsHandler(w http.ResponseWriter, r *http.Request) {
	res := ArtistResponse{}
	err := db.FetchYAML(s.config.DBPaths.Artists, &res.Artists)
	if err != nil {
		res.Error = err
		writeJSON(w, 500, res)
		return
	}

	writeJSON(w, 200, res)
	return
}

// TODO: Some of this logic probably belongs in the internal/db module?
//
// index walks through the data/songs directory and saves an appendix in another file.
func (s *Server) index() {
	out := []Song{}
	err := filepath.Walk(s.config.DBPaths.Songs.Base, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		parts := strings.Split(path, "/")
		if len(parts) == 4 && strings.HasSuffix(path, ".yaml") {
			artist := parts[2]
			title := strings.TrimSuffix(parts[3], ".yaml")
			out = append(out, Song{Artist: artist, Title: title})
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
