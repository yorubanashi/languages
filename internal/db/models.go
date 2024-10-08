package db

type Line struct {
	// The original (primary) line. Currently either Japanese or Chinese (CN).
	Pri string `yaml:"pri" json:"pri"`
	// An alternate (secondary) version of the line. Currently either Japanese (Furigana) or Chinese (TW).
	Sec string `yaml:"sec" json:"sec"`
	// Romanized version of the line, for easy pronounciation.
	Rom string `yaml:"rom" json:"rom"`
	// The translation of the line in English.
	Eng string `yaml:"eng" json:"eng"`
}

type Verse struct {
	Name  string `yaml:"name" json:"name"`
	Lines []Line `yaml:"lines" json:"lines"`
}

type Song struct {
	Title      string   `yaml:"title" json:"title"`                           // Song title, in the primary language
	Artist     string   `yaml:"artist" json:"artist"`                         // Main artist
	Featured   []string `yaml:"featured,omitempty" json:"featured,omitempty"` // Featured artist(s)
	VerseOrder []string `yaml:"order" json:"order"`                           // Order of the verses
	Verses     []Verse  `yaml:"verses" json:"verses"`                         // Verse definitions
}
