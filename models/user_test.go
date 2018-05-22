package models

import (
	//_ "github.com/go-sql-driver/mysql"
	"fmt"
	"github.com/shopspring/decimal"
	"testing"
)

// func newUserForTest(t *testing.T) *User {
// 	return NewUser(newDbForTest(t))
// }

// func TestUserCRUD(t *testing.T) {
// 	u := newUserForTest(t)

// 	// Signup
// 	userRow, err := u.Signup(nil, newEmailForTest(), "abc123", "abc123")
// 	if err != nil {
// 		t.Errorf("Signing up user should work. Error: %v", err)
// 	}
// 	if userRow == nil {
// 		t.Fatal("Signing up user should work.")
// 	}
// 	if userRow.ID <= 0 {
// 		t.Fatal("Signing up user should work.")
// 	}

// 	// DELETE FROM users WHERE id=...
// 	_, err = u.DeleteById(nil, userRow.ID)
// 	if err != nil {
// 		t.Fatalf("Deleting user by id should not fail. Error: %v", err)
// 	}

// }

func TestFloatToDecimal(t *testing.T) {
	d1 := decimal.NewFromFloat(10000024.00)
	fmt.Printf("%v", d1)
}
