package main

type Message struct {
	Code  int
	Title string
	Value string
}

type Slalomer struct {
	ID       string   `json:"id"`
	Name     string   `json:"user" binding:"required"`
	Email    string   `json:"email" binding:"required"`
	Password string   `json:"pass"`
	Photo    string   `json:"photo"`
	Location Location `json:"_" `
}

type Location struct {
	ID   string
	Lat  string
	Lon  string
	Name string
}
