package saver

import (
	"encoding/csv"
	"fmt"
	"os"

	"lab4/internal/core"
	"lab4/pkg/cfg"

	"github.com/xuri/excelize/v2"
)

func SaveResultsToExcel(results []core.Ship) error {
	f := excelize.NewFile()

	headers := []string{"Название", "IMO", "MMSI", "Тип", "Ссылка", "Ошибка"}
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue("Sheet1", cell, header)
	}

	for i, vessel := range results {
		row := i + 2
		f.SetCellValue("Sheet1", fmt.Sprintf("A%d", row), vessel.Name)
		f.SetCellValue("Sheet1", fmt.Sprintf("B%d", row), vessel.IMO)
		f.SetCellValue("Sheet1", fmt.Sprintf("C%d", row), vessel.MMSI)
		f.SetCellValue("Sheet1", fmt.Sprintf("D%d", row), vessel.Type)
		f.SetCellValue("Sheet1", fmt.Sprintf("E%d", row), vessel.URL)
		f.SetCellValue("Sheet1", fmt.Sprintf("F%d", row), vessel.Error)
	}

	if err := f.SaveAs(cfg.ResultExcelFile); err != nil {
		return fmt.Errorf("ошибка сохранения Excel файла '%s': %w", cfg.ResultExcelFile, err)
	}

	return nil
}

func SaveResultsToCSV(results []core.Ship) error {
	file, err := os.Create(cfg.ResultCSVFile)
	if err != nil {
		return fmt.Errorf("ошибка создания CSV файла '%s': %w", cfg.ResultCSVFile, err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	headers := []string{"Название", "IMO", "MMSI", "Тип", "Ссылка", "Ошибка"}
	if err := writer.Write(headers); err != nil {
		return fmt.Errorf("ошибка записи заголовков CSV: %w", err)
	}

	for _, vessel := range results {
		record := []string{
			vessel.Name,
			vessel.IMO,
			vessel.MMSI,
			vessel.Type,
			vessel.URL,
			vessel.Error,
		}
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("ошибка записи строки CSV: %w", err)
		}
	}

	if err := writer.Error(); err != nil {
		return fmt.Errorf("ошибка записи CSV: %w", err)
	}

	return nil
}
