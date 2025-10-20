package ds

type Writer interface {
	Path() string
	Write(spec Spec, rows [][]string) error
}
