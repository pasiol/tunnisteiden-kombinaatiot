package internal

import (
	"database/sql"
	"time"
)

type AccountingIdentifiers struct {
	CostCentre int
	Identifier int
	Location   int
	Timestamp  time.Time
}

type ReporRow struct {
	EmployeeId  int
	Email       string
	CostCentre  int
	Identifier1 int
	Location    int
	Date        sql.NullTime
	Type        string
	Start       sql.NullTime
	End         sql.NullTime
	Month       string
}
