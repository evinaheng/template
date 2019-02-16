package dummydb_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	. "github.com/template/be/lib/storage/database/dummydb"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestRunQuerySelectError(t *testing.T) {

	sdb, sMock, _ := sqlmock.New()

	// SELECT
	sMock.ExpectQuery(`SELECT id FROM dummy`).WillReturnError(errors.New("Foo"))

	db := New(sdb)

	query := `SELECT id FROM dummy`
	rows, err := db.Queryx(query)

	assert.Nil(t, rows)
	assert.EqualError(t, err, "Foo")

}

func TestRunStructScan(t *testing.T) {

	sdb, sMock, _ := sqlmock.New()
	db := New(sdb)

	mockRows := sqlmock.NewRows([]string{"id"}).AddRow("GA")
	sMock.ExpectQuery(`SELECT (.+) FROM dummy`).WillReturnRows(mockRows)

	rows, err := db.Master().Queryx(`SELECT id FROM dummy`)
	data := struct {
		ID string `db:"id"`
	}{}
	rows.Next()
	rows.StructScan(&data)
	rows.Close()

	assert.Equal(t, "GA", data.ID)
	assert.Nil(t, err)
	assert.NotNil(t, rows)

	assert.Nil(t, sMock.ExpectationsWereMet())
}

func TestRunQueryWhere(t *testing.T) {

	sdb, sMock, _ := sqlmock.New()

	row := sqlmock.NewRows([]string{"id"}).AddRow(1)
	sMock.ExpectQuery(`SELECT (.+)`).WithArgs(1, 2).WillReturnRows(row)

	db := New(sdb)

	var newLastID int64
	// Update last ID
	query := `
	SELECT 
		id
	FROM a
	WHERE 
		status = ? OR status = ?
	ORDER BY id
	LIMIT 1
	`

	err := db.Get(&newLastID, query, 1, 2)

	assert.Nil(t, err)
	assert.Nil(t, sMock.ExpectationsWereMet())

}

func TestRunQueryUpdateBeginCommit(t *testing.T) {

	sdb, sMock, _ := sqlmock.New()

	sMock.ExpectBegin()
	sMock.ExpectExec("UPDATE test").WithArgs(123, 123).WillReturnResult(sqlmock.NewResult(1, 1))
	sMock.ExpectCommit()

	queryUpdate := `
		UPDATE test
		SET 
			id = ?
		WHERE
			id = ?
	`

	db := New(sdb)

	tx, errBegin := db.Begin()
	result, err := tx.Exec(queryUpdate, 123, 123)
	errCommit := tx.Commit()

	lastID, _ := result.LastInsertId()

	assert.Nil(t, errBegin)
	assert.Equal(t, int64(1), lastID)
	assert.Nil(t, err)
	assert.Nil(t, errCommit)
	assert.Nil(t, sMock.ExpectationsWereMet())

}

func TestRunQueryUpdateBeginRollback(t *testing.T) {

	sdb, sMock, _ := sqlmock.New()

	sMock.ExpectBegin()
	sMock.ExpectExec("UPDATE test").WithArgs(123, 123).WillReturnResult(sqlmock.NewResult(1, 1))
	sMock.ExpectRollback()

	queryUpdate := `
		UPDATE test
		SET 
			id = ?
		WHERE
			id = ?
	`

	db := New(sdb)

	tx, errBegin := db.Beginx()
	result, err := tx.Exec(queryUpdate, 123, 123)
	errRollback := tx.Rollback()

	lastID, _ := result.LastInsertId()

	assert.Nil(t, errBegin)
	assert.Equal(t, int64(1), lastID)
	assert.Nil(t, err)
	assert.Nil(t, errRollback)
	assert.Nil(t, sMock.ExpectationsWereMet())

}
