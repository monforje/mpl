package reader

import (
	"errors"
	"fmt"
	"lab4/pkg/cfg"

	"github.com/xuri/excelize/v2"
)

func ReadURLsFromExcel() ([]string, error) {
	f, err := excelize.OpenFile(cfg.LinksFileName)
	if err != nil {
		return nil, errors.New("ошибка открытия файла " + cfg.LinksFileName + ": " + err.Error())
	}
	defer f.Close()

	rows, err := f.GetRows(cfg.SheetName)
	if err != nil {
		return nil, errors.New("ошибка чтения листа " + cfg.SheetName + ": " + err.Error())
	}

	var urls []string
	for i, row := range rows {
		if i == 0 || len(row) == 0 {
			continue
		}
		urls = append(urls, row[0])
	}

	fmt.Printf("Найдено %d ссылок для проверки\n", len(urls))
	return urls, nil
}
