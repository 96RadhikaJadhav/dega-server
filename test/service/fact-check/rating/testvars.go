package rating

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/factly/dega-server/test"
	"github.com/factly/dega-server/test/service/core/medium"
	"github.com/jinzhu/gorm/dialects/postgres"
)

var headers = map[string]string{
	"X-Space": "1",
	"X-User":  "1",
}

var Data = map[string]interface{}{
	"name": "True",
	"slug": "true",
	"description": postgres.Jsonb{
		RawMessage: []byte(`{"type":"description"}`),
	},
	"numeric_value": 5,
	"medium_id":     uint(1),
}

var resData = map[string]interface{}{
	"name": "True",
	"slug": "true",
	"description": map[string]interface{}{
		"type": "description",
	},
	"numeric_value": 5,
}

var defaultData = []map[string]interface{}{
	{
		"name":          "True",
		"slug":          "true",
		"description":   "True",
		"numeric_value": 5,
	},
	{
		"name":          "Partly True",
		"slug":          "partly-true",
		"description":   "Partly True",
		"numeric_value": 4,
	},
	{
		"name":          "Misleading",
		"slug":          "misleading",
		"description":   "Misleading",
		"numeric_value": 3,
	},
	{
		"name":          "Partly False",
		"slug":          "partly-false",
		"description":   "Partly False",
		"numeric_value": 2,
	},
	{
		"name":          "False",
		"slug":          "false",
		"description":   "False",
		"numeric_value": 1,
	},
}

var invalidData = map[string]interface{}{
	"name":          "a",
	"numeric_value": 0,
}

var columns = []string{"id", "created_at", "updated_at", "deleted_at", "created_by_id", "updated_by_id", "name", "slug", "medium_id", "description", "numeric_value", "space_id"}

var selectQuery = regexp.QuoteMeta(`SELECT * FROM "ratings"`)
var deleteQuery = regexp.QuoteMeta(`UPDATE "ratings" SET "deleted_at"=`)
var paginationQuery = `SELECT \* FROM "ratings" (.+) LIMIT 1 OFFSET 1`

var basePath = "/fact-check/ratings"
var defaultsPath = "/fact-check/ratings/default"
var path = "/fact-check/ratings/{rating_id}"

func slugCheckMock(mock sqlmock.Sqlmock, rating map[string]interface{}) {
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT slug, space_id FROM "ratings"`)).
		WithArgs(fmt.Sprint(rating["slug"], "%"), 1).
		WillReturnRows(sqlmock.NewRows(columns))
}

func ratingInsertMock(mock sqlmock.Sqlmock) {
	mock.ExpectBegin()
	medium.SelectWithSpace(mock)
	mock.ExpectQuery(`INSERT INTO "ratings"`).
		WithArgs(test.AnyTime{}, test.AnyTime{}, nil, 1, 1, Data["name"], Data["slug"], Data["description"], Data["numeric_value"], Data["medium_id"], 1).
		WillReturnRows(sqlmock.
			NewRows([]string{"id"}).
			AddRow(1))
}

func ratingInsertError(mock sqlmock.Sqlmock) {
	mock.ExpectBegin()
	medium.EmptyRowMock(mock)
	mock.ExpectRollback()
}

func ratingUpdateMock(mock sqlmock.Sqlmock, rating map[string]interface{}, err error) {
	mock.ExpectBegin()
	if err != nil {
		medium.EmptyRowMock(mock)
	} else {
		medium.SelectWithSpace(mock)
		mock.ExpectExec(`UPDATE \"ratings\"`).
			WithArgs(test.AnyTime{}, 1, rating["name"], rating["slug"], rating["description"], rating["numeric_value"], rating["medium_id"], 1).
			WillReturnResult(sqlmock.NewResult(1, 1))
		SelectWithSpace(mock)
		medium.SelectWithOutSpace(mock)
	}

}

func SelectWithOutSpace(mock sqlmock.Sqlmock, rating map[string]interface{}) {
	mock.ExpectQuery(selectQuery).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows(columns).
			AddRow(1, time.Now(), time.Now(), nil, 1, 1, rating["name"], rating["slug"], rating["medium_id"], rating["description"], rating["numeric_value"], 1))

	// Preload medium
	medium.SelectWithOutSpace(mock)
}

func SelectWithSpace(mock sqlmock.Sqlmock) {
	mock.ExpectQuery(selectQuery).
		WithArgs(1, 1).
		WillReturnRows(sqlmock.NewRows(columns).
			AddRow(1, time.Now(), time.Now(), nil, 1, 1, Data["name"], Data["slug"], Data["medium_id"], Data["description"], Data["numeric_value"], 1))
}

//check rating exits or not
func recordNotFoundMock(mock sqlmock.Sqlmock) {
	mock.ExpectQuery(selectQuery).
		WithArgs(1, 100).
		WillReturnRows(sqlmock.NewRows(columns))
}

func sameNameCount(mock sqlmock.Sqlmock, count int, name interface{}) {
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(1) FROM "ratings"`)).
		WithArgs(1, strings.ToLower(name.(string))).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(count))
}

// check rating associated with any claim before deleting
func ratingClaimExpect(mock sqlmock.Sqlmock, count int) {
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(1) FROM "claims"`)).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(count))
}

func ratingCountQuery(mock sqlmock.Sqlmock, count int) {
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(1) FROM "ratings"`)).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(count))
}

func EmptyRowMock(mock sqlmock.Sqlmock) {
	mock.ExpectQuery(selectQuery).
		WithArgs(1, 1).
		WillReturnRows(sqlmock.NewRows(columns))
}
