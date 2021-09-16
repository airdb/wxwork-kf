package schema

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestModelPk_create(t *testing.T) {
	db, mock, clean := prepareDB(t)
	defer clean()

	mock.ExpectExec("^INSERT INTO (.+)").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), nil, "", "").
		WillReturnResult(sqlmock.NewResult(0, 1))

	talk := Talk{}
	db.Create(&talk)
	assert.NotEmpty(t, talk.ID)

	err := mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestModelPk_update(t *testing.T) {
	db, mock, clean := prepareDB(t)
	defer clean()

	talk := Talk{}
	talk.ID.renew()
	talk.OpenKFID = "001"

	rows := sqlmock.NewRows([]string{"id", "open_kfid"}).
		AddRow(talk.ID, talk.OpenKFID)
	mock.ExpectQuery("^SELECT").WithArgs(talk.ID).WillReturnRows(rows)

	mock.ExpectExec("^UPDATE").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(0, 1))

	assert.NoError(t, db.First(&talk).Error)

	db.Model(&talk).Update("open_kfid", "002")

	err := mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func prepareDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, func()) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	mock.ExpectExec(`^CREATE TABLE (.+)`).
		WillReturnResult(sqlmock.NewResult(0, 0))

	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      mockDB,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{
		SkipDefaultTransaction: true, // for sqlmock
	})
	if err != nil {
		t.Fatal("failed to connect database")
	}
	err = db.Migrator().AutoMigrate(&Talk{})
	assert.NoError(t, err)

	return db, mock, func() {
		mockDB.Close()
	}
}
