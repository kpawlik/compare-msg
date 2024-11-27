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
	GetTranslation(string, string, int) (string, error)
	GetTranslations() [][]string
}

type CSV struct {
	file string
	rows [][]string
}

func NewCSV(file string) *CSV {
	return &CSV{
		file: file,
		rows: make([][]string, 0),
	}
}

func (c *CSV) Read() (err error) {
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

func (c *CSV) GetTranslation(namespace, key string, col int) (translation string, err error) {
	for _, row := range c.rows {
		col1 := row[0]
		parts := strings.Split(col1, ".")
		csvNS := parts[0]
		csvKey := parts[1]
		if namespace == csvNS && key == csvKey {
			translation = row[col]
			return
		}
	}
	err = fmt.Errorf("missing translation for %s.%s", namespace, key)
	return
}


func (c *CSV) GetTranslations() (res [][]string){
	for _, row := range c.rows {
		col1 := row[0]
		parts := strings.Split(col1, ".")
		namespace := parts[0]
		messageId := parts[1]
		message := row[3]
		res = append(res, []string{namespace, messageId, message})
	}
	return
}