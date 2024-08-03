package entity

import "time"

type Rate struct {
	Param1 int
	Param2 int
	Param3 int
	Param4 int
}

type Rating struct {
	ID        int
	UserID    int
	MusicID   int
	Rating    int
	Comment   string
	CreatedAt time.Time
}
