package repositories

import (
	"database/sql"
	"errors"
	"time"
)

func PingWithTimeout(db *sql.DB, timeout time.Duration) error {
	ch := make(chan error)
	go func() {
		ch <- db.Ping()
	}()
	select {
	case err := <-ch:
		return err
	case <-time.After(time.Second * 5):
		return errors.New("db ping timeout")
	}
}
