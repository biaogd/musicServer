package main

import "fmt"

type music struct {
	ID         int
	SongName   string
	SongAuthor string
	AllTime    int
	SongSize   int
	URL        string
}

func (m music) String() string {
	return fmt.Sprintf("ID=%d,SongName=%s,SongAuthor=%s,AllTime=%d,SongSize=%d,URL=%s", m.ID, m.SongName, m.SongAuthor, m.AllTime, m.SongSize, m.URL)
}
