package models

import "time"

type Film struct {
	ID          string
	Name        string
	Description string
	ShowDate    time.Time
	Location    Location
	PosterURL   string
	IsOpen      bool
	PlacesNum   uint64
}

type Location struct {
	Lat         float64
	Long        float64
	Title       string
	Description string
	VideoURL    string
}

type User struct {
	TgID      uint64
	Username  string
	Name      string
	Surname   string
	Group     string
	WhereFind string
}

type State struct {
	Val  int
	Info map[string]string
}
