package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/kpawlik/om"
)


var(
	file1, file2, translationFile string
	outFile string
	overwrite bool 
)

func init() {
	
	flag.StringVar(&file1, "f1", "", "Source file to compare")
	flag.StringVar(&file2, "f2", "", "File to compare")
	flag.StringVar(&translationFile, "translation-file", "", "CSV translation file")
	flag.StringVar(&outFile, "out", "", "JSON out file")
	flag.BoolVar(&overwrite, "overwrite", false, "Overwrite out")
	
	flag.Parse()
	
}
func main() {
	if len(file1) == 0 || len(file2) == 0{
		flag.Usage()
		return;
	}
	var (
		err error
		f1, f2 []byte
		csvTr *CSV
		omtr1, omtr2 *om.OrderedMap
		translation string
	)
	if f1, err = os.ReadFile(file1); err != nil{
		log.Panicf("Error open %s, %v", file1, err)
	}
	if f2, err = os.ReadFile(file2); err != nil{
		log.Panicf("Error open %s, %v", file2, err)
	}
	map1 := om.NewOrderedMap()
	map2 := om.NewOrderedMap()
	if err = map1.UnmarshalJSON(f1); err != nil{
		log.Panicf("Error unmarshal %s, %v", file1, err)
	}
	if err = map2.UnmarshalJSON(f2); err != nil{
		log.Panicf("Error unmarshal %s, %v", file2, err)
	}
	if len(translationFile)>0{
		csvTr = NewCSV(translationFile)
		if err = csvTr.Read(); err != nil{
			log.Panicf("Error reading translation file  %s, %v", translationFile, err)
		}
	}
	for _, key := range map1.Keys{
		// fmt.Println(key)
		tr1 := map1.Map[key]
		tr2 := map2.Map[key]
		if (tr2 != nil){
			omtr2, _ = tr2.(*om.OrderedMap)	
		}else{
			omtr2 = om.NewOrderedMap()
		}
		omtr1, _ = tr1.(*om.OrderedMap)
		
		for _, key2 := range omtr1.Keys{
			v2 := omtr2.Map[key2]
			if v2 != nil{
				continue
				fmt.Printf("Misisng %s.%s in %s\n", key, key2, file2)
			}
			if csvTr == nil{
				continue
			}
			if translation, err = csvTr.getTranslation(key, key2, 2); err != nil{
				fmt.Println(err)
			}
			omtr2.Map[key2] = translation
			omtr2.Keys = append(omtr2.Keys, key2)
		}
	}
	if len(outFile) > 0 {
		if _, statErr := os.Stat(outFile); statErr == nil{
			if !overwrite{
				fmt.Printf("File %s already exists\n", outFile)
				return
			}
		}
		var out = []byte{}
		if out, err = map2.MarshalIndent("  "); err != nil{
			log.Panicf("Error marshal output %v", err)
		}
		if err = os.WriteFile(outFile, out, 0644); err != nil{
			log.Panicf("Error write file %s. %v", outFile, err)
		}
	}
}

type CSV struct {
	file string
	rows  [][]string
}

func NewCSV(file string) *CSV {
	return &CSV{
		file: file,
		rows:  make([][]string, 0),
	}
}

func (c *CSV) Read() (err error) {
	var (
		f *os.File
		row []string
	)
	if f, err = os.Open(c.file); err != nil{
		return fmt.Errorf("error opening file %s. %w", c.file, err)
	}
	reader := csv.NewReader(f)
	for {
		row, err = reader.Read()
		if err == io.EOF{
			err = nil
			break
		}
		c.rows = append(c.rows, row)
	}
	return
}

func (c *CSV) getTranslation(namespace, key string, col int) (translation string, err error){
	for _, row := range c.rows {
		col1:=row[0]
		parts := strings.Split(col1, ".")
		csvNS := parts[0]
		csvKey := parts[1]
		if namespace == csvNS && key == csvKey{
			if len(row) -1 < col{
				continue
			}
			translation = row[col]
			return
		}
	}
	err = fmt.Errorf("missing translation for %s.%s", namespace, key)
	return
}