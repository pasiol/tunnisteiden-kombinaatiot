package internal

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/labstack/gommon/log"
	_ "github.com/lib/pq"
)

const (
	createIdentifiers = `CREATE TABLE IF NOT EXISTS identifiers (kustannuspaikka INT NOT NULL, toiminto1 INT NOT NULL, paikkakunta INT NOT NULL);`
	createBookings    = `CREATE TABLE IF NOT EXISTS bookings (tid INT NOT NULL, email VARCHAR(128), kustannuspaikka INT, toiminto1 INT, paikkakunta INT, pvm DATE NOT NULL, kirjaustyyppi VARCHAR(256), alkaa TIME WITH TIME ZONE NOT NULL, päättyy TIME WITH TIME ZONE NOT NULL, raportti VARCHAR(7) NOT NULL)`
)

func formatSQLDate(t sql.NullTime) string {
	date := t.Time
	return date.Format("2006-02-01 15:04:05")
}

func exportIdentifiers2SQL(accountingIdentifiers []AccountingIdentifiers) error {
	db, err := sql.Open("postgres", getConnectionString())
	CheckError(err)
	defer db.Close()

	// check db
	err = db.Ping()
	CheckError(err)

	log.Info("SQL connected")

	insertDynStmt := `insert into "identifiers"("kustannuspaikka", "toiminto1", "paikkakunta") values($1, $2, $3)`
	count := 0
	for _, ai := range accountingIdentifiers {
		_, e := db.Exec(insertDynStmt, ai.CostCentre, ai.Identifier, ai.Location)
		CheckError(e)
		count = count + 1
	}
	log.Infof("%d accountingIdentifiers exported", count)
	return nil
}

func TruncateIdentifiersTable() error {
	db, err := sql.Open("postgres", getConnectionString())
	CheckError(err)
	defer db.Close()

	// check db
	err = db.Ping()
	CheckError(err)

	log.Info("SQL connected")

	truncateStatement := `TRUNCATE identifiers;`

	_, err = db.Exec(truncateStatement)
	CheckError(err)
	log.Info("identifiers table succesfully truncated")
	return nil

}

func TruncateBookingsTable() error {
	db, err := sql.Open("postgres", getConnectionString())
	CheckError(err)
	defer db.Close()

	// check db
	err = db.Ping()
	CheckError(err)

	log.Info("SQL connected")

	truncateStatement := `TRUNCATE bookings;`

	_, err = db.Exec(truncateStatement)
	CheckError(err)
	log.Info("bookings table succesfully truncated")
	return nil

}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

func getPgConnection() (*sql.DB, error) {
	db, err := sql.Open("postgres", getConnectionString())
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	return db, err
}

func getConnectionString() string {
	host, port, user, passsword, db := ReadPGSecrets()
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, passsword, db)
}

func InitializeDb() error {
	db, err := sql.Open("postgres", getConnectionString())
	CheckError(err)
	defer db.Close()

	_, e := db.Exec(createIdentifiers)
	CheckError(e)

	_, e = db.Exec(createBookings)
	CheckError(e)
	log.Info("db succesfully initialized")
	return nil
}

func insertBooking(db *sql.DB, b Booking, employees map[int64]int64, emails map[int64]string) error {
	costCentreIdentifiers := strings.TrimSpace(TransformString(256, b.CostCentreID))
	costCentre, function1, function2, _, _ := getCostCentreData(costCentreIdentifiers)
	email := emails[int64(b.PrimusID)]
	employeeId := employees[int64(b.PrimusID)]
	date, _ := transformDate(b.Date)
	shortCaption := strings.Replace(TransformString(256, b.ShortCaption), "'", "", -1)
	startTime := getStartTime(b)
	endTime := getEndTime(b)
	month := ""
	if len(b.Date) > 6 {
		month = string(b.Date[len(b.Date)-7:])
	}

	insertDynStmt := `insert into "bookings"("tid","email","kustannuspaikka", "toiminto1", "paikkakunta", "pvm", "kirjaustyyppi", "alkaa", "päättyy", "raportti") values($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	_, e := db.Exec(insertDynStmt, employeeId, email, costCentre, function1, function2, date, shortCaption, startTime, endTime, month)
	CheckError(e)

	return nil
}

func GetReport(month string) error {
	db, err := getPgConnection()
	if err != nil {
		return err
	}
	selecStmt := fmt.Sprintf("SELECT * FROM bookings b WHERE NOT EXISTS(SELECT 1 FROM identifiers i WHERE b.kustannuspaikka = i.kustannuspaikka AND b.toiminto1 = i.toiminto1 AND b.paikkakunta = i.paikkakunta) AND (b.kirjaustyyppi LIKE '%%Sääntelemättömät tehtävät%%' OR b.kirjaustyyppi LIKE '%%Opetus ja ohjaus%%') AND raportti='%s';", month)
	log.Infof("%s", selecStmt)
	rows, err := db.Query(selecStmt)
	if err != nil {
		return err
	}

	data := "\"työntekijäid\";\"email\";\"kustannuspaikka\";\"toiminto 1\";\"paikkakunta\";\"pvm\";\"kirjaustyyppi\";\"alkaa\";\"päättyy\";\"kuukausi\"\r\n"

	counter := 0
	for rows.Next() {
		var r = ReporRow{}

		err = rows.Scan(&r.EmployeeId, &r.Email, &r.CostCentre, &r.Identifier1, &r.Location, &r.Date, &r.Type, &r.Start, &r.End, &r.Month)
		CheckError(err)

		rowString := fmt.Sprintf("\"%d\";\"%s\";\"%d\";\"%d\";\"%d\";%s;\"%s\";%s;%s;\"%s\"\r\n", r.EmployeeId, r.Email, r.CostCentre, r.Identifier1, r.Location, formatSQLDate(r.Date), r.Type, string(formatSQLDate(r.Start)[len(formatSQLDate(r.Start))-10:]), string(formatSQLDate(r.End)[len(formatSQLDate(r.End))-10:]), r.Month)
		data = data + rowString
		rowString = ""
		counter++
		r = ReporRow{}

	}
	data = strings.Replace(data, "\"-1\"", "", -1)
	createFile(month+".csv", data)

	return nil
}
