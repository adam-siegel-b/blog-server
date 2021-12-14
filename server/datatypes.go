package main

type Message struct {
	Code  int
	Title string
	Value string
}

type Slalomer struct {
	ID       string
	Name     string
	Email    string
	Photo    string
	Location Location
}

type Location struct {
	ID   string
	Lat  string
	Lon  string
	Name string
}
