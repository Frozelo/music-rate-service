package entity

import "time"

type Music struct {
	Id          int
	Name        string
	Artist      string
	Genre       string
	Duration    time.Duration
	ReleaseDate time.Time
}

type MusicRate struct {
	Param1 int
	Param2 int
	Param3 int
	Param4 int
}
