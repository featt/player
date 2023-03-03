package models

import "time"

type Song struct {
	Name     string
	Duration time.Duration
	FilePath string
}

type LinkedListNode struct {
	Song *Song
	Prev *LinkedListNode
	Next *LinkedListNode
}

type LinkedList struct {
	Head *LinkedListNode
	Tail *LinkedListNode
}

func NewLinkedList() *LinkedList {
	return &LinkedList{}
}

type Track struct {
	Name     string
	Duration time.Duration
	Prev     *Track
	Next     *Track
}

type PlaylistController struct {
	playlist *LinkedList
	current  *LinkedListNode
	playing  bool
	pausedAt time.Duration
}

type Playlist struct {
	List       *LinkedList
	Current    *LinkedListNode
	IsPlaying  bool
	PlaybackCh chan bool
}