package pipeline

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/mitchellh/mapstructure"
)

type (
	BswtanResult struct {
		URL          string `mapstructure:"_url"`
		Title        string `mapstructure:"title"`
		Author       string `mapstructure:"author"`
		ChapterTitle string `mapstructure:"chapter_title"`
		Content      string `mapstructure:"content"`
	}
	Bswtan struct {
	}
)

func (b *BswtanResult) GetDir() string {
	return b.GetTitle() + "-" + b.GetAuthor()
}

func (b *BswtanResult) GetTitle() string {
	return strings.Trim(b.Title, " ")
}

func (b *BswtanResult) GetFilename() string {
	return strings.TrimSuffix(filepath.Base(b.URL), filepath.Ext(b.URL))
}

func (b *BswtanResult) GetAuthor() string {
	return strings.Trim(strings.ReplaceAll(b.Author, "作    者：", ""), " ")
}

func (b *BswtanResult) GetChapterTitle() string {
	return strings.Trim(b.ChapterTitle, " ")
}

func (b *BswtanResult) GetContent() string {
	c := strings.Trim(b.Content, "    ")
	return strings.ReplaceAll(c, "    ", "\n") + "\n"
}

func (b *Bswtan) Name() string {
	return "bswtan"
}

func (b *Bswtan) Run(path string, v interface{}) error {
	result := new(BswtanResult)
	if err := mapstructure.Decode(v, result); err != nil {
		return err
	}
	fileDir := filepath.Join(path, result.GetDir())
	if !isDir(fileDir) {
		if err := os.MkdirAll(fileDir, 0755); err != nil {
			return err
		}
	}
	buf := &bytes.Buffer{}
	buf.WriteString(result.GetChapterTitle())
	buf.WriteString("\n")
	buf.WriteString(result.GetContent())
	return ioutil.WriteFile(filepath.Join(fileDir, result.GetFilename()), buf.Bytes(), 0664)
}
