package daemon

const (
	creatorTask = iota + 1
	postsTask
	fileTask
	sharedFilesTask
)

type task struct {
	ttype int
	data  interface{}
}
