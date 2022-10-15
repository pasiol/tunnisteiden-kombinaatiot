package internal

import (
	"context"
	"database/sql"
	b64 "encoding/base64"
	"errors"
	"fmt"
	"regexp"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/gommon/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	mssqlutils "github.com/pasiol/go-mssql-utils"
	pq "github.com/pasiol/gopq"
	"github.com/pasiol/mongoutils"
	"github.com/pasiol/serviceLog"

	_ "github.com/denisenkom/go-mssqldb"
)

// Booking struct
type Booking struct {
	TransformID           interface{} `json:"_id,omitempty" bson:"_id,omitempty"`
	PrimusID              int16
	PrimusName            string
	CostCentreID          string
	CostCentrePercent     int16
	ShortCaption          string
	Caption               string
	CourseTypeID          string
	CourseTypeName        string
	CourseTypeAccountCode string
	Info                  string `json:"Info,omitempty" bson:"Info,omitempty"`
	KurreID               int32
	Date                  string
	StartHour             int
	StartMinutes          int
	EndHour               int
	EndMinutes            int
	ScheduleSize          int
	ArrayIndex            int
	JobTimestamp          time.Time
	ExtractionID          interface{} `json:"ExtractionID,omitempty" bson:"ExtractionID,omitempty"`
}

func transformDate(Date string) (string, error) {
	if len(Date) == 10 {
		return Date[6:10] + "-" + Date[3:5] + "-" + Date[:2], nil
	}
	return "", fmt.Errorf("trasforming date failed, string %s len != 10", Date)
}

func decodeValue(s string) string {
	decodedS, err := b64.StdEncoding.DecodeString(s)
	if err != nil {
		log.Print(serviceLog.GetLogMessage(serviceLog.Service{}, serviceLog.Event{ShortMessage: "Decoding error nothing to do.", FullMessage: err.Error(), Succesful: false, Severity: "critical"}))
	}
	return string(decodedS)
}

func getStartTime(b Booking) string {
	dString, _ := transformDate(b.Date)
	hour := strconv.Itoa(b.StartHour)
	minutes := strconv.Itoa(b.StartMinutes)
	sqlDate := fmt.Sprintf("%s %s:%s:00.000", dString, hour, minutes)
	if strings.Contains(sqlDate, " 24:") {
		sqlDateTransformed, err := mssqlutils.SQLMidnight24To00(sqlDate)
		if err == nil {
			return sqlDateTransformed
		}
	}
	return sqlDate
}

func getEndTime(b Booking) string {
	dString, _ := transformDate(b.Date)
	hour := strconv.Itoa(b.EndHour)
	minutes := strconv.Itoa(b.EndMinutes)
	sqlDate := fmt.Sprintf("%s %s:%s:00.000", dString, hour, minutes)
	if strings.Contains(sqlDate, " 24:") {
		sqlDateTransformed, err := mssqlutils.SQLMidnight24To00(sqlDate)
		if err == nil {
			return sqlDateTransformed
		}
	}
	return sqlDate
}

func mongoDate2SQLDate(d time.Time) string {
	layout := "2006-01-02 15:04:05"
	return d.Format(layout) + ".000"
}

func getCostCentreData(data string) (int16, int16, int16, int16, error) {
	re := regexp.MustCompile(`\d[\d,]*[\.]?[\d{2}]*`)
	numbers := re.FindAllString(data, -1)
	count := len(numbers)
	if count > 4 || count < 1 {
		return -1, -1, -1, -1, errors.New("count of arguments is illegal on the cost centres data")
	}
	integer1, err := strconv.ParseInt(numbers[0], 10, 16)
	if err != nil {
		return -1, -1, -1, -1, errors.New("type conversion error on the cost centres data")
	}
	if count == 1 {
		return int16(integer1), -1, -1, -1, nil
	}
	integer2, err := strconv.ParseInt(numbers[1], 10, 16)
	if err != nil {
		return int16(integer1), -1, -1, -1, errors.New("type conversion error on the cost centres data")
	}
	if count == 2 {
		return int16(integer1), int16(integer2), -1, -1, nil
	}

	integer3, err := strconv.ParseInt(numbers[2], 10, 16)
	if err != nil {
		return int16(integer1), int16(integer2), -1, -1, errors.New("type conversion error on the cost centres data")
	}
	if count == 3 {
		return int16(integer1), int16(integer2), int16(integer3), -1, nil
	}
	integer4, err := strconv.ParseInt(numbers[3], 10, 16)
	if err != nil {
		return int16(integer1), int16(integer2), int16(integer3), -1, errors.New("type conversion error on the cost centres data")
	}
	return int16(integer1), int16(integer2), int16(integer3), int16(integer4), nil
}

func truncateTable(db *sql.DB) error {
	ctx := context.Background()
	err := db.PingContext(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}
	tsql := fmt.Sprintf("TRUNCATE TABLE [%s]", "table")
	c, err := db.QueryContext(ctx, tsql)

	if err != nil {
		defer c.Close()
		log.Printf("Truncating table %s failed %s", "table", err)

		return err
	}

	return nil
}

func stringToNULL(s string, tsql string) string {
	stringPattern, err := regexp.Compile(s)
	if err == nil {
		match := stringPattern.ReplaceAll([]byte(tsql), []byte("NULL"))
		converted := string(match)
		if len(converted) > 0 {
			return converted
		}

		return tsql
	}

	return tsql
}

func removeBooking(noSQLDb *mongo.Database, b Booking) error {

	coll := noSQLDb.Collection("teachersTransformed")
	_, err := coll.DeleteOne(context.TODO(), bson.M{"_id": b.TransformID})

	if err != nil {
		return err
	}

	return nil
}

func getEmployeesAndEmails(emails map[int64]string, employees map[int64]int64) (int, error) {
	query := Teachers()
	pq.Debug = false
	host, port, user, password := ReadPrimusSecrets()
	query.Host = host
	query.Port = port
	query.User = user
	query.Pass = password
	query.Output = "testi.pq"

	output, err := pq.ExecuteAndRead(query, 30)
	if err != nil {
		return 0, err
	}
	if output == "" {
		return 0, nil
	}
	rows := strings.Fields(output)
	for _, row := range rows {
		fields := strings.Split(row, ";")
		id, err := strconv.ParseInt(fields[0], 10, 64)
		if err != nil {
			return 0, err
		}
		employeeId, err := strconv.ParseInt(fields[1], 10, 64)
		if err != nil {
			return 0, err
		}
		email := fields[2]
		if err != nil {
			return 0, err
		}
		employees[id] = employeeId
		emails[id] = email
	}
	return len(rows), nil
}

func GetBookins() {
	debug.SetGCPercent(100)

	var employees = map[int64]int64{}
	var emails = map[int64]string{}
	teacherCount, err := getEmployeesAndEmails(emails, employees)
	if err != nil {
		teacherCount = 0
		log.Fatalf("no teachers data check pq and network connections")
	} else {
		log.Infof("teacher count: %d", teacherCount)
	}

	if teacherCount > 0 {
		var succeed, failed = 0, 0
		m := ReadMongoSecrets()

		noSQLDb, session := mongoutils.ConnectToMongo(m, false)
		findOptions := options.Find()
		findOptions.SetLimit(100000)

		pg, err := getPgConnection()
		if err != nil {
			log.Fatalf("db connection failed: %s", err)
		}

		dataToMove := true
		for dataToMove {
			coll := noSQLDb.Collection("teachersTransformed")
			documentCount, err := coll.CountDocuments(context.TODO(), bson.D{{}})
			log.Infof("Founded %d transformed documents to be processed.", documentCount)
			if err != nil {
				log.Fatal("getting transformed document count failed")
			}
			if documentCount > 0 {
				cur, err := coll.Find(context.TODO(), bson.D{{}}, findOptions)
				if err != nil {
					log.Fatalf("getting transformed documents failed: %s", err.Error())
				}
				for cur.Next(context.TODO()) {
					var document Booking
					err := cur.Decode(&document)
					if err != nil {
						log.Fatalf("decoding transformed document failed: %s", err.Error())
					}
					err = insertBooking(pg, document, employees, emails)
					if err == nil {
						err := removeBooking(noSQLDb, document)
						if err != nil {
							log.Infof("removing transformed document failed: %s", err.Error())
							failed++
						} else {
							succeed++
						}
					} else {
						log.Infof("inserting transformed document failed: %s", err.Error())
					}
					document = Booking{}
				}
			} else {
				dataToMove = false
			}
		}

		err = session.Disconnect(context.TODO())
		if err != nil {
			log.Errorf("disconnecting mongo connection failed: %s", err)
		}
		log.Infof("importing bookings ended controlled manner, succeed: %d, failed: %d", succeed, failed)
	}
}
