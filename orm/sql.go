package orm

import (
	"database/sql"
	"errors"
	"fmt"
)

type DB struct {
	sql.DB
}

func Open(dialect string, args ...interface{}) (db *DB, err error) {
	if len(args) == 0 {
		err = errors.New("invalid database source")
		return 
	}
	buf := make([]byte, 0, len(dialect) * 2)
	num := 0
	for i := 0;i < len(dialect);i++ {
		if dialect[i] == '?' {
			num++
			if num > len(args) {
				// 返回参数 ? 过多错误
				return
			}
			buf = append(buf, []byte(fmt.Sprintf("%v", args[num-1]))...)
		} else {
			buf = append(buf, dialect[i])
		}
	}
	if num != len(args) {
		// 返回参数 ? 过少错误
		return
	}
	return
}

func (db *DB) Close() {
	db.DB.Close()
}
