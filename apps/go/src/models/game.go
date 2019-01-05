package models

// ID represents a globally unique identifier for this game
type GameID int

// Game tracks all the information we know about a game
type Game struct {
	ID   GameID
	Slug string
	Name string
}
