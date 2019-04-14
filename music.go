package main

import "fmt"

type music struct {
	ID         int
	SongName   string
	SongAuthor string
	AllTime    int
	SongSize   int64
	URL        string
	Count      int
}

func (m music) String() string {
	return fmt.Sprintf("ID=%d,SongName=%s,SongAuthor=%s,AllTime=%d,SongSize=%d,URL=%s", m.ID, m.SongName, m.SongAuthor, m.AllTime, m.SongSize, m.URL)
}

type myApp struct {
	Content string
	Name    string
	Status  string
}