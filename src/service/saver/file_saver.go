package saver

import (
	"bufio"
	"github.com/zored/edit/src/service/formatters"
	"github.com/zored/edit/src/service/navigation"
	"github.com/zored/edit/src/service/replacers"
	"github.com/zored/edit/src/service/tokenizers"
	"github.com/zored/edit/src/service/tokens"
	"io"
	"os"
)

type fileSaver struct {
	tokenizer tokenizers.ITokenizer
	formatter formatters.IFormatter
	replacer  replacers.IReplacer
}

func NewFileSaver() *fileSaver {
	return &fileSaver{
		tokenizer: tokenizers.NewTokenizer(),
		formatter: formatters.NewFormatter(),
		replacer:  replacers.NewReplacer(),
	}
}

func (f *fileSaver) Save(options *FileOptions) error {
	path := options.file
	molecule, err := f.getMolecule(path, options)
	if err != nil {
		return err
	}

	moleculeString := f.formatter.Format(molecule.Tokens, options.GetFormatterOptions())
	content, err := f.replaceText(path, molecule.Interval, moleculeString)
	if err != nil {
		return err
	}

	return f.save(path, content)
}

func (f *fileSaver) replaceText(path string, interval *navigation.Interval, formatted string) (content string, err error) {
	err = f.withFile(path, func(file io.Reader) (err error) {
		content, err = f.replacer.Replace(file, interval, formatted)
		return
	})
	return
}

func (f *fileSaver) getMolecule(path string, config *FileOptions) (mol *tokens.Molecule, err error) {
	err = f.withFile(path, func(file io.Reader) (err error) {
		mol, err = f.tokenizer.Tokenize(file, config.cursor, config.wrapper, config.separator)
		return
	})
	return
}

func (f *fileSaver) save(path string, content string) error {
	tmpFile := path + ".new"
	newFile, err := os.Create(tmpFile)
	if err != nil {
		return err
	}
	defer func() { _ = newFile.Close() }()
	if _, err := newFile.WriteString(content); err != nil {
		return err
	}
	if err := os.Rename(tmpFile, path); err != nil {
		return err
	}
	return nil
}

func (f fileSaver) withFile(path string, h func(reader io.Reader) error) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()
	reader := bufio.NewReader(file)
	return h(reader)
}
