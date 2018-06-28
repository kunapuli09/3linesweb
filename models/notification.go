package models

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"time"
)

func NewNotification(db *sqlx.DB) *Notification {
	n := &Notification{}
	n.db = db
	n.table = "notifications"
	n.hasID = true
	return n
}

type Notification struct {
	Base
}

type NotificationRow struct {
	ID               int64     `db:"id"`
	Investment_ID    int64     `db:"investment_id"`
	News_ID          int64     `db:"news_id"`
	NotificationDate time.Time `db:"NotificationDate"`
	NewsDate         time.Time `db:"NewsDate"`
	StartupName      string    `db:"StartupName"`
	Title            string    `db:"Title"`
	Status           string    `db:"Status"`
	Email            string    `db:"Email"`
}

func (n *NotificationRow) FormattedNewsDate() string {
	return n.NewsDate.Format("01/02/2006")
}

func (i *Notification) userRowFromSqlResult(tx *sqlx.Tx, sqlResult sql.Result) (*NotificationRow, error) {
	nId, err := sqlResult.LastInsertId()
	if err != nil {
		return nil, err
	}

	return i.GetById(tx, nId)
}

func (n *Notification) AllNotifications(tx *sqlx.Tx, email string) ([]*NotificationRow, error) {
	notifications := []*NotificationRow{}
	query := fmt.Sprintf("SELECT * FROM %v where Email='%v' and Status='%v'", n.table, email, "UNREAD")
	err := n.db.Select(&notifications, query)
	return notifications, err
}

// GetById returns record by id.
func (i *Notification) GetById(tx *sqlx.Tx, id int64) (*NotificationRow, error) {
	n := &NotificationRow{}
	query := fmt.Sprintf("SELECT * FROM %v WHERE id=?", i.table)
	err := i.db.Get(n, query, id)

	return n, err
}

// GetByName returns record by name.
func (i *Notification) CountByEmail(tx *sqlx.Tx, email string) (int, error) {
	var count int
	query := fmt.Sprintf("SELECT count(*) FROM %v where Email=? and Status=?", i.table)
	err := i.db.Get(&count, query, email, "UNREAD")
	return count, err
}

// GetByName returns record by name.
func (i *Notification) GetByName(tx *sqlx.Tx, name string) (*NotificationRow, error) {
	n := &NotificationRow{}
	query := fmt.Sprintf("SELECT * FROM %v WHERE name=?", i.table)
	err := i.db.Get(n, query, name)

	return n, err
}

// GetByName returns record by name.
func (i *Notification) ReadByUser(tx *sqlx.Tx, email string) ([]*NotificationRow, error) {
	notifications := []*NotificationRow{}
	query := fmt.Sprintf("SELECT * FROM %v WHERE Status='%v' and Email='%v'", i.table, "READ", email)
	err := i.db.Select(&notifications, query)
	return notifications, err
}

// create a new record of user.
func (i *Notification) Create(tx *sqlx.Tx, m map[string]interface{}) (*NotificationRow, error) {
	sqlResult, err := i.InsertIntoTable(tx, m)
	if err != nil {
		return nil, err
	}
	return i.userRowFromSqlResult(tx, sqlResult)
}

// UpdateEmailAndPasswordById updates user email and password.
func (i *Notification) UpdateById(tx *sqlx.Tx, nId int64, data map[string]interface{}) (*NotificationRow, error) {
	if len(data) > 0 {
		//calling base.go function
		_, err := i.UpdateByID(tx, data, nId)
		if err != nil {
			return nil, err
		}
	}

	return i.GetById(tx, nId)
}

// UpdateStatusById updates user email and password.
func (n *Notification) UpdateStatusById(tx *sqlx.Tx, NotificationId int64) (*NotificationRow, error) {
	data := make(map[string]interface{})
	data["Status"] = "READ"

	if len(data) > 0 {
		_, err := n.UpdateByID(tx, data, NotificationId)
		if err != nil {
			return nil, err
		}
	}

	return n.GetById(nil, NotificationId)
}

// Get All by Investment ID.
func (i *Notification) GetAllByInvestmentId(tx *sqlx.Tx, Investment_ID int64) ([]*NotificationRow, error) {
	css := []*NotificationRow{}
	query := fmt.Sprintf("SELECT * FROM %v WHERE Investment_ID=%v", i.table, Investment_ID)
	err := i.db.Select(&css, query)
	return css, err
}

// UpdateEmailAndPasswordById updates user email and password.
func (i *Notification) DeleteByID(tx *sqlx.Tx, csId int64) (sql.Result, error) {

	//calling base.go function
	sqlResult, err := i.DeleteById(tx, csId)
	if err != nil {
		return nil, err
	}

	return sqlResult, nil
}

func (i *Notification) BatchPublish(tx *sqlx.Tx, emails []string, StartupName string, news *NewsRow) (sql.Result, error) {
	sqlStr := "INSERT INTO notifications(investment_id, news_id, NotificationDate, StartupName, Status, Email, Title, NewsDate) VALUES "
	vals := []interface{}{}
	for _, email := range emails {
	    sqlStr += "(?, ?, ?, ?, ?, ?, ?, ?),"
	    vals = append(vals, news.Investment_ID, news.ID, time.Now(), StartupName, "UNREAD", email, news.Title, news.NewsDate)
	}
	//trim the last ,
	sqlStr = sqlStr[0:len(sqlStr)-1]
	//prepare the statement
	stmt, err := i.db.Prepare(sqlStr)
	if err != nil {
		fmt.Println(err)
	}

	//format all vals at once
	sqlResult, err := stmt.Exec(vals...)

	return sqlResult, err
}
