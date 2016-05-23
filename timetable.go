package main

type TimetableInfo map[string]TimetableInfoRow

type TimetableInfoRow struct {
	Start   string
	End     string
	Name    string
	Entries []string
}
