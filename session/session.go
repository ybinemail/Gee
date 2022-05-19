package session

import (
	"database/sql"
	"geego/log"
	"strings"
)

type Session struct {
	db *sql.DB

	sql strings.Builder

	sqlVars []interface{}
}

func New(db *sql.DB) *Session {

	return &Session{
		db: db,
	}
}

func (s *Session) Clear() {
	s.sql.Reset()

	s.sqlVars = nil
}

func (s *Session) DB() *sql.DB {
	return s.db
}

// 用来拼接 SQL 语句和 SQL 语句中占位符的对应值
func (s *Session) Raw(sql string, value ...interface{}) *Session {
	s.sql.WriteString(sql)
	s.sql.WriteString(" ")
	s.sqlVars = append(s.sqlVars, value...)
	return s
}

func (s *Session) Exec() (result sql.Result, err error) {

	defer s.Clear()

	log.Info(s.sql.String(), s.sqlVars)
	if result, err = s.DB().Exec(s.sql.String(), s.sqlVars...); err != nil {
		log.Error(err)
	}
	return
}

func (s *Session) QueryRow() *sql.Row {

	defer s.Clear()

	log.Info(s.sql.String(), s.sqlVars)
	return s.DB().QueryRow(s.sql.String(), s.sqlVars...)
}

func (s *Session) Query() (rows *sql.Rows, err error) {

	defer s.Clear()

	log.Info(s.sql.String(), s.sqlVars)
	if rows, err = s.DB().Query(s.sql.String(), s.sqlVars...); err != nil {
		log.Error(err)
	}
	return
}
