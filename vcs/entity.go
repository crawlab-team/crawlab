package vcs

import (
	"time"
)

type GitOptions struct {
	checkout []GitCheckoutOption
}

type GitRef struct {
	Type        string    `json:"type"`
	Name        string    `json:"name"`
	FullName    string    `json:"full_name"`
	Hash        string    `json:"hash"`
	Timestamp   time.Time `json:"timestamp"`
	RemoteTrack string    `json:"remote_track"`
}

type GitLog struct {
	Hash        string    `json:"hash"`
	Msg         string    `json:"msg"`
	AuthorName  string    `json:"author_name"`
	AuthorEmail string    `json:"author_email"`
	Timestamp   time.Time `json:"timestamp"`
	Refs        []GitRef  `json:"refs"`
}

type GitFileStatus struct {
	Path     string          `json:"path"`
	Name     string          `json:"name"`
	IsDir    bool            `json:"is_dir"`
	Staging  string          `json:"staging"`
	Worktree string          `json:"worktree"`
	Extra    string          `json:"extra"`
	Children []GitFileStatus `json:"children"`
}
