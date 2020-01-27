package daemon

type Taskif interface {
	Domain() string
	Description() string
	Do() []Taskif
}
