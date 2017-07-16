package app

import "time"

const (
	TypeText               = "text"
	TypeTextStatus         = "text-status"
	TypeTextStatus2        = "text-status-2"
	TypeTextStatusProgress = "text-status-bar"
)

const (
	Size1x = "1x"
	Size2x = "2x"
	Size4x = "4x"
)

const (
	StateDefault       = "default"
	StateSuccess       = "success"
	StateIndeterminate = "indeterminate"
	StateWarning       = "warning"
	StateError         = "error"
)

type tileJSON struct {
	ID              string    `json:"id"`
	LastChangeTime  time.Time `json:"updated"`
	Source          string    `json:"source"`
	Type            string    `json:"type"`
	Size            string    `json:"size"`
	State           string    `json:"state"`
	TitleText       string    `json:"titleText,omitempty"`
	DescriptionText string    `json:"descrText,omitempty"`
	StatusValue     *int      `json:"statusValue,omitempty"`
}

type dataJSON struct {
	Version string     `json:"version"`
	Tiles   []tileJSON `json:"tiles"`
}
