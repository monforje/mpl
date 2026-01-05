package util

import (
	"log"
	"send/internal/model"
	"strconv"
	"time"
)

func CSVtoJSON(csvData []string) (model.Scan, error) {
	birthDate, _ := time.Parse("02.01.2006", csvData[1])

	log.Println("Преобразуем CSV в JSON для:", csvData[0])

	return model.Scan{
		FullName:  csvData[0],
		BirthDate: birthDate,
		Sex:       csvData[2],

		Hemoglobin:   parse(csvData[3]),
		Erythrocytes: parse(csvData[4]),
		Hematocrit:   parse(csvData[5]),
		MCV:          parse(csvData[6]),

		Leukocytes:  parse(csvData[7]),
		Neutrophils: parse(csvData[8]),
		Lymphocytes: parse(csvData[9]),
		Monocytes:   parse(csvData[10]),
		Eosinophils: parse(csvData[11]),
		Basophils:   parse(csvData[12]),

		Platelets: parse(csvData[13]),
		MPV:       parse(csvData[14]),
	}, nil
}

func parse(s string) float64 {
	v, _ := strconv.ParseFloat(s, 64)
	return v
}
