package internal

import (
	"strconv"
	"time"

	"github.com/labstack/gommon/log"
)

func ReadAccountingIdentifiers(f string) ([]AccountingIdentifiers, error) {
	rows, err := ReadCsv(f)
	if err != nil {
		return []AccountingIdentifiers{}, err
	}
	ts := time.Now()
	identifiers := []AccountingIdentifiers{}
	count := 0
	for i, row := range rows[1:] {
		f1, err := strconv.Atoi(row[0])
		if err != nil {
			log.Infof("skipping line %d, check file %s integrity", i, f)
			continue
		}
		f2, err := strconv.Atoi(row[1])
		if err != nil {
			log.Infof("skipping line %d, check file %s integrity", i, f)
			continue
		}
		f3, err := strconv.Atoi(row[2])
		if err != nil {
			log.Infof("skipping line %d, check file %s integrity", i, f)
			continue
		}
		ai := AccountingIdentifiers{
			CostCentre: f1,
			Identifier: f2,
			Location:   f3,
			Timestamp:  ts,
		}
		identifiers = append(identifiers, ai)
		count = count + 1
	}
	log.Infof("readed %d identifiers from file %s", count, f)
	return identifiers, nil
}

func UpdateAccountingIdentifiers(newData []AccountingIdentifiers) error {
	err := exportIdentifiers2SQL(newData)
	if err != nil {
		return err
	}
	return nil
}
