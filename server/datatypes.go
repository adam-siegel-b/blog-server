package main

type Message struct {
	Code  int    `json:"status_code"`
	Title string `json:"title"`
	Value string `json:"value"`
}

type Slalomer struct {
	ID       string   `json:"id"`
	Name     string   `json:"user" binding:"required"`
	Email    string   `json:"email" binding:"required"`
	Password string   `json:"pass"`
	Photo    string   `json:"photo"`
	Location Location `json:"location" `
}

type Slalomers struct {
	Users []Slalomer
}

type Location struct {
	ID   string  `json:"loc-id"`
	Lat  float64 `json:"lat"`
	Lon  float64 `json:"lon"`
	Name string  `json:"loc-name"`
}
