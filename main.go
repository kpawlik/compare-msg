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
			translation = row[col]
			return
		}
	}
	err = fmt.Errorf("missing translation for %s.%s", namespace, key)
	return
}

func readMessageFile(path string) (*om.OrderedMap, error){
	var (
		content []byte
		err error
	)
	if content, err = os.ReadFile(path); err != nil{
		err = fmt.Errorf("error open %s, %v", path, err)
		return nil, err
	}
	oMap := om.NewOrderedMap()
	if err = oMap.UnmarshalJSON(content); err != nil{
		err = fmt.Errorf("error unmarshal %s, %v", path, err)
		return nil, err
	}
	return oMap, err
}

func writeOut(destMap *om.OrderedMap) (err error) {
	if len(outFile) > 0 {
		if _, statErr := os.Stat(outFile); statErr == nil{
			if !overwrite{
				fmt.Printf("File %s already exists\n", outFile)
				return
			}
		}
		var out = []byte{}
		if out, err = destMap.MarshalIndent("  "); err != nil{
			log.Panicf("Error marshal output %v", err)
		}
		if err = os.WriteFile(outFile, out, 0644); err != nil{
			log.Panicf("Error write file %s. %v", outFile, err)
		}
	}
	return
}

func main() {
	var (
		err error
		csvTr *CSV
		sourceSubMap, destSubMap *om.OrderedMap
		sourceMap, destMap *om.OrderedMap
		translation string
	)
	if len(file1) == 0 || len(file2) == 0{
		flag.Usage()
		return;
	}
	if sourceMap, err = readMessageFile(file1); err != nil{
		log.Fatal(err)
	}
	if destMap, err = readMessageFile(file2); err != nil{
		log.Fatal(err)
	}
	if len(translationFile)>0{
		csvTr = NewCSV(translationFile)
		if err = csvTr.Read(); err != nil{
			log.Panicf("Error reading translation file  %s, %v", translationFile, err)
		}
	}
	for _, namespace := range sourceMap.Keys{
		sourceTranslation := sourceMap.Map[namespace]	
		destTranslation := destMap.Map[namespace]
		if (destTranslation != nil){
			destSubMap, _ = destTranslation.(*om.OrderedMap)	
		}else{
			destSubMap = destMap.CreateChild(namespace)
		}
		sourceSubMap, _ = sourceTranslation.(*om.OrderedMap)
		for _, messageId := range sourceSubMap.Keys{
			v2 := destSubMap.Map[messageId]
			if v2 != nil || csvTr == nil{
				continue
			}
			if translation, err = csvTr.getTranslation(namespace, messageId, 2); err != nil{
				fmt.Println(err)
			}
			destSubMap.Map[messageId] = translation
			destSubMap.Keys = append(destSubMap.Keys, messageId)
		}
	}
	if err = writeOut(destMap); err != nil{
		log.Panic(err)
	}

}