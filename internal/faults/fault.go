package faults

type Action int

const (
    Pass Action = iota
    Drop
)

type Decision struct {
    Action Action
}
