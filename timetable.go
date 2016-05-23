package main

type TimetableInfo []TimetableInfoRow

type TimetableInfoRow struct {
	Start   string
	End     string
	Name    string
	Entries []string
}
