package apptui

import ()

type FindingSource int
type Severity int

const (
	Low Severity = iota
	Medium
	High
	Critical
)

type Finding struct {
	ID          string
	Source      FindingSource
	Title       string
	Description string
	Severity    Severity
	Target      string
	Focus       bool
}

const (
	Source1 FindingSource = iota
	Source2
	Source3
)
