package compare_msg

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
)

type Translation interface {
	Read() error
	GetTranslation(string, string) (string, error)
	GetTranslations() [][]string
}

type CSVTranslation struct {
	file string
	rows [][]string
	translationColumn int
}

func NewCSV(file string) *CSVTranslation {
	return &CSVTranslation{
		file: file,
		rows: make([][]string, 0),
		translationColumn: 2,
	}
}

func (c *CSVTranslation) Read() (err error) {
	var (
		f   *os.File
		row []string
	)
	if f, err = os.Open(c.file); err != nil {
		return fmt.Errorf("error opening file %s. %w", c.file, err)
	}
	reader := csv.NewReader(f)
	for {
		row, err = reader.Read()
		if err == io.EOF {
			err = nil
			break
		}
		c.rows = append(c.rows, row)
	}
	return
}

func (c *CSVTranslation) GetTranslation(namespace, key string) (translation string, err error) {
	nsMessageId := fmt.Sprintf("%s.%s", namespace, key)
	for _, row := range c.rows {
		if row[0] == nsMessageId{
			translation = row[c.translationColumn]
			return
		}
	}
	err = fmt.Errorf("missing translation for %s.%s", namespace, key)
	return
}


func (c *CSVTranslation) GetTranslations() (res [][]string){
	res = make([][]string, len(c.rows))
	for i, row := range c.rows {
		nsMessageId := row[0]
		parts := strings.Split(nsMessageId, ".")
		namespace, messageId  := parts[0], parts[1]
		message := row[c.translationColumn]
		res[i] = []string{namespace, messageId, message}
	}
	return
}