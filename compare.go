package compare_msg

import (
	"fmt"
	"log"
	"os"

	"github.com/kpawlik/om"
)

func readMessageFile(path string) (*om.OrderedMap, error) {
	var (
		content []byte
		err     error
	)
	if content, err = os.ReadFile(path); err != nil {
		err = fmt.Errorf("error open %s, %v", path, err)
		return nil, err
	}
	messageFile := om.NewOrderedMap()
	if err = messageFile.UnmarshalJSON(content); err != nil {
		err = fmt.Errorf("error unmarshal %s, %v", path, err)
		return nil, err
	}
	return messageFile, err
}

func writeOut(destMap *om.OrderedMap, outFile string, overwrite bool) (err error) {
	if len(outFile) > 0 {
		if _, statErr := os.Stat(outFile); statErr == nil {
			if !overwrite {
				fmt.Printf("File %s already exists\n", outFile)
				return
			}
		}
		var out = []byte{}
		if out, err = destMap.MarshalIndent("  "); err != nil {
			log.Panicf("Error marshal output %v", err)
		}
		if err = os.WriteFile(outFile, out, 0644); err != nil {
			log.Panicf("Error write file %s. %v", outFile, err)
		}
	}
	return
}

func CompareUpdate(baseFile, messageFile, translationFilePath, outFile string, overwrite bool) (err error) {
	var (
		translationFile          Translation
		sourceSubMap, destSubMap *om.OrderedMap
		sourceMap, destMap       *om.OrderedMap
		translation              string
	)

	if sourceMap, err = readMessageFile(baseFile); err != nil {
		return
	}
	if destMap, err = readMessageFile(messageFile); err != nil {
		return
	}
	if len(translationFilePath) > 0 {
		translationFile = NewCSV(translationFilePath)
		if err = translationFile.Read(); err != nil {
			err = fmt.Errorf("error reading translation file  %s, %w", translationFilePath, err)
			return
		}
	}
	for _, namespace := range sourceMap.Keys {
		sourceTranslation := sourceMap.Map[namespace]
		destTranslation := destMap.Map[namespace]
		if destTranslation != nil {
			destSubMap, _ = destTranslation.(*om.OrderedMap)
		} else {
			destSubMap = destMap.CreateChild(namespace)
		}
		sourceSubMap, _ = sourceTranslation.(*om.OrderedMap)
		for _, messageId := range sourceSubMap.Keys {
			currentTranslation := destSubMap.Map[messageId]
			if currentTranslation == nil && translationFile == nil {
				fmt.Printf("missing translation for %s.%s\n", namespace, messageId)
				continue
			}
			if currentTranslation != nil || translationFile == nil {
				continue
			}
			if translation, err = translationFile.GetTranslation(namespace, messageId, 2); err != nil {
				fmt.Println(err)
			}
			destSubMap.Set(messageId, translation)
		}
	}
	if err = writeOut(destMap, outFile, overwrite); err != nil {
		return
	}
	return
}

func Compare(baseFile, messageFile string) (err error) {
	var (
		sourceSubMap, destSubMap *om.OrderedMap
		sourceMap, destMap       *om.OrderedMap
	)
	if sourceMap, err = readMessageFile(baseFile); err != nil {
		return
	}
	if destMap, err = readMessageFile(messageFile); err != nil {
		return
	}
	for _, namespace := range sourceMap.Keys {
		sourceTranslation := sourceMap.Map[namespace]
		destTranslation := destMap.Map[namespace]
		if destTranslation != nil {
			destSubMap, _ = destTranslation.(*om.OrderedMap)
		} else {
			destSubMap = destMap.CreateChild(namespace)
		}
		sourceSubMap, _ = sourceTranslation.(*om.OrderedMap)
		for _, messageId := range sourceSubMap.Keys {
			currentTranslation := destSubMap.Map[messageId]
			if currentTranslation == nil {
				fmt.Printf("missing translation for %s.%s\n", namespace, messageId)
				continue
			}
		}
	}
	return
}

func Update(baseFile, translationFilePath, outFile string, overwrite bool, force bool) (err error) {
	var (
		translationFile Translation
		sourceSubMap    *om.OrderedMap
		sourceMap       *om.OrderedMap
	)
	if sourceMap, err = readMessageFile(baseFile); err != nil {
		return
	}
	translationFile = NewCSV(translationFilePath)
	if err = translationFile.Read(); err != nil {
		err = fmt.Errorf("error reading translation file  %s, %w", translationFilePath, err)
		return
	}
	translations := translationFile.GetTranslations()
	for _, translation := range translations {
		namespace := translation[0]
		messageId := translation[1]
		message := translation[2]
		sourceTranslation, ok := sourceMap.Map[namespace]
		if ok {
			sourceSubMap = sourceTranslation.(*om.OrderedMap)
		} else {
			sourceSubMap = sourceMap.CreateChild(namespace)
		}
		_, ok = sourceSubMap.Map[messageId]
		if !ok {
			fmt.Printf("%s.%s added\n", namespace, messageId)
			sourceSubMap.Set(messageId, message)
		} else {
			if force {
				fmt.Printf("%s.%s updates\n", namespace, messageId)
				sourceSubMap.Set(messageId, message)
			}
		}
	}
	if err = writeOut(sourceMap, outFile, overwrite); err != nil {
		return
	}
	return

}
