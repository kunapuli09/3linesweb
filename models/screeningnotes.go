package models

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"time"
)

func NewScreeningNotes(db *sqlx.DB) *ScreeningNotes {
	is := &ScreeningNotes{}
	is.db = db
	is.table = "screeningnotes"
	is.hasID = true
	return is
}

type ScreeningNotes struct {
	Base
}

type ScreeningNotesRow struct {
	ID              int64     `db:"id"`
	Application_ID  int64     `db:"application_id"`
	ScreenerEmail   string    `db:"ScreenerEmail"`
	ScreeningDate   time.Time `db:"ScreeningDate"`
	Need            string    `db:"Need"`
	Status          string    `db:"Status"`
	TeamRisk        int8      `db:"TeamRisk"`
	BarrierToEntry  int8      `db:"BarrierToEntry"`
	TechRisk        int8      `db:"TechRisk"`
	CompetitionRisk int8      `db:"CompetitionRisk"`
	PoliticalRisk   int8      `db:"PoliticalRisk"`
	SupplierRisk    int8      `db:"SupplierRisk"`
	ExecutionRisk   int8      `db:"ExecutionRisk"`
	MarketRisk      int8      `db:"MarketRisk"`
	ScalingRisk     int8      `db:"ScalingRisk"`
	ExitRisk        int8      `db:"ExitRisk"`
	Comments        string    `db:"Comments"`
}

func (i *ScreeningNotesRow) FormattedScreeningDate() string {
	return i.ScreeningDate.Format("01/02/2006")
}

func (i *ScreeningNotes) userRowFromSqlResult(tx *sqlx.Tx, sqlResult sql.Result) (*ScreeningNotesRow, error) {
	isId, err := sqlResult.LastInsertId()
	if err != nil {
		return nil, err
	}

	return i.GetById(tx, isId)
}

// AllUsers returns all user rows.
func (i *ScreeningNotes) AllScreeningNotes(tx *sqlx.Tx) ([]*ScreeningNotesRow, error) {
	isrs := []*ScreeningNotesRow{}
	query := fmt.Sprintf("SELECT * FROM %v", i.table)
	err := i.db.Select(&isrs, query)

	return isrs, err
}

// AllUsers returns all user rows.
func (i *ScreeningNotes) AllScreeningNotesByApplicationId(tx *sqlx.Tx, Application_ID int64) ([]*ScreeningNotesRow, error) {
	isrs := []*ScreeningNotesRow{}
	query := fmt.Sprintf("SELECT * FROM %v WHERE application_id=%v", i.table, Application_ID)
	err := i.db.Select(&isrs, query)

	return isrs, err
}

// GetById returns record by id.
func (i *ScreeningNotes) GetById(tx *sqlx.Tx, id int64) (*ScreeningNotesRow, error) {
	isr := &ScreeningNotesRow{}
	query := fmt.Sprintf("SELECT * FROM %v WHERE id=?", i.table)
	err := i.db.Get(isr, query, id)

	return isr, err
}

// GetByName returns record by name.
func (i *ScreeningNotes) GetByName(tx *sqlx.Tx, name string) (*ScreeningNotesRow, error) {
	isr := &ScreeningNotesRow{}
	query := fmt.Sprintf("SELECT * FROM %v WHERE name=?", i.table)
	err := i.db.Get(isr, query, name)

	return isr, err
}

// create a new record of user.
func (i *ScreeningNotes) Create(tx *sqlx.Tx, m map[string]interface{}) (*ScreeningNotesRow, error) {
	sqlResult, err := i.InsertIntoTable(tx, m)
	if err != nil {
		return nil, err
	}
	return i.userRowFromSqlResult(tx, sqlResult)
}

// UpdateEmailAndPasswordById updates user email and password.
func (i *ScreeningNotes) UpdateById(tx *sqlx.Tx, isId int64, data map[string]interface{}) (*ScreeningNotesRow, error) {
	if len(data) > 0 {
		//calling base.go function
		_, err := i.UpdateByID(tx, data, isId)
		if err != nil {
			return nil, err
		}
	}

	return i.GetById(tx, isId)
}

// Get All by Application ID.
func (i *ScreeningNotes) GetByApplicationIdAndScreener(tx *sqlx.Tx, ScreeningNotes_ID int64, Application_ID int64, ScreenerEmail string) (*ScreeningNotesRow, error) {
	isr := &ScreeningNotesRow{}
	var query string
	if ScreeningNotes_ID == 0 {
		query = fmt.Sprintf("SELECT * FROM %v WHERE Application_ID=%v AND ScreenerEmail='%v'", i.table, Application_ID, ScreenerEmail)
		//fmt.Printf("Executing Query %v", query)
	} else {
		query = fmt.Sprintf("SELECT * FROM %v WHERE ID=%v AND Application_ID=%v AND ScreenerEmail='%v'", i.table, ScreeningNotes_ID, Application_ID, ScreenerEmail)
		//fmt.Printf("Executing Query for Existing Notes %v", query)
	}
	err := i.db.Get(isr, query)
	return isr, err
}

// UpdateEmailAndPasswordById updates user email and password.
func (i *ScreeningNotes) DeleteByID(tx *sqlx.Tx, csId int64) (sql.Result, error) {

	//calling base.go function
	sqlResult, err := i.DeleteById(tx, csId)
	if err != nil {
		return nil, err
	}

	return sqlResult, nil
}
