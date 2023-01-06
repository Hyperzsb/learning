package task

import "fmt"

type Task struct {
	Name   string
	Desc   string
	Status Status
	// TODO: add a property describing time points
}

type Status int

const (
	Todo Status = iota + 1
	Doing
	Done
)

func (t *Task) ToString() string {
	switch t.Status {
	case Todo:
		return fmt.Sprintf("(TODO) %s: %s", t.Name, t.Desc)
	case Doing:
		return fmt.Sprintf("(DOING) %s: %s", t.Name, t.Desc)
	case Done:
		return fmt.Sprintf("(DONE) %s: %s", t.Name, t.Desc)
	default:
		return fmt.Sprintf("(UNKNOWN) %s: %s", t.Name, t.Desc)
	}
}

type DuplicateErr struct {
	Name string
}

func (de DuplicateErr) Error() string {
	return fmt.Sprintf("duplicate task: '%s'; task name should be unique", de.Name)
}
