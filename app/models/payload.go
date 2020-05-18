package models

type LogFile struct {
	Name string  `json:"name"`
	Size float64 `json:"size"`
	Path string  `json:"path"`
	Date string  `json:"date"`
}