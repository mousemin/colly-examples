package pipeline

import (
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strconv"
)

type Pipeline interface {
	Name() string
	Run(path string, v interface{}) error
}

type Factory struct {
	path      string
	pipelines map[string]Pipeline
}

func New(path string) (*Factory, error) {
	f := &Factory{
		path:      path,
		pipelines: make(map[string]Pipeline),
	}
	if err := f.init(); err != nil {
		return nil, err
	}
	return f, nil
}

func (f *Factory) init() error {
	if isDir(f.path) {
		return nil
	}
	return os.MkdirAll(f.path, 0755)
}

func (f *Factory) SetPath(path string) error {
	f.path = path
	return f.init()
}

func RegisterPipeline(f *Factory) {
	f.RegisterPipeline(&Bswtan{})
}

func (f *Factory) RegisterPipeline(p Pipeline) *Factory {
	f.pipelines[p.Name()] = p
	return f
}

func (f *Factory) Pipeline(name string, v interface{}) error {
	p, ok := f.pipelines[name]
	if !ok {
		return errors.New("not found " + name + " pipeline")
	}
	return p.Run(f.path, v)
}

func (f *Factory) Merge() error {
	fp, err := os.Open(f.path)
	if err != nil {
		return err
	}
	names, err := fp.Readdirnames(-1)
	fp.Close()
	if err != nil {
		return err
	}
	sort.Strings(names)
	for _, name := range names {
		if err := f.merge(filepath.Join(f.path, name)); err != nil {
			return err
		}
	}
	return nil
}

func (f *Factory) merge(path string) error {
	if !isDir(path) {
		return nil
	}
	chapterList := make(sort.IntSlice, 0, 1024)
	err := filepath.Walk(path, func(p string, info fs.FileInfo, err error) error {
		if p == path {
			return nil
		}
		a, err := strconv.Atoi(info.Name())
		if err != nil {
			return err
		}
		chapterList = append(chapterList, a)
		return nil
	})

	if err != nil {
		return err
	}
	sort.Sort(chapterList)
	fp, err := os.Create(path + ".txt")
	if err != nil {
		return err
	}
	defer fp.Close()
	for _, val := range chapterList {
		f1, err := os.Open(filepath.Join(path, strconv.Itoa(val)))
		if err != nil {
			return err
		}
		_, err = io.Copy(fp, f1)
		if err != nil {
			return err
		}
		if err := f1.Close(); err != nil {
			return err
		}
	}
	return os.RemoveAll(path)
}
