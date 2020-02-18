package files

import (
	"bufio"
	"fmt"
	"github.com/zored/edit/src/service/formatters"
	"github.com/zored/edit/src/service/parsers"
	"github.com/zored/edit/src/service/replacers"
	"os"
)

type fileFormatter struct {
	parser    parsers.IParser
	formatter formatters.IFormatter
	replacer  replacers.IReplacer
}

func NewFileFormatter() *fileFormatter {
	return &fileFormatter{
		parser:    parsers.NewParser(),
		formatter: formatters.NewFormatter(),
		replacer:  replacers.NewReplacer(),
	}
}

func (f *fileFormatter) Format(config *FileFormatConfig) error {
	file, err := os.Open(config.file)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()
	replaceInterval, tokens_, err := f.parser.Parse(bufio.NewReader(file), config.cursor, config.wrapper, config.separator)
	if err != nil {
		return err
	}

	if err := file.Close(); err != nil {
		return err
	}
	file, err = os.Open(config.file)
	if err != nil {
		return err
	}

	newCode := f.formatter.Format(tokens_.All, config.formatRule, config.indent, config.separator)
	newContents, err := f.replacer.Replace(bufio.NewReader(file), replaceInterval, newCode)
	if err != nil {
		return err
	}

	if err := file.Close(); err != nil {
		return err
	}

	newFile, err := os.Create(config.file + ".new")
	if err != nil {
		return err
	}
	defer func() { _ = newFile.Close() }()
	if _, err := newFile.WriteString(newContents); err != nil {
		return err
	}

	if err := os.Rename(config.file, config.file+".old"); err != nil {
		return err
	}

	if err := os.Rename(config.file+".new", config.file); err != nil {
		return err
	}

	fmt.Println(newCode, replaceInterval)
	return nil
}
