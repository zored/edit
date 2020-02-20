package saver

type IFileSaver interface {
	Save(options *FileOptions) error
}
