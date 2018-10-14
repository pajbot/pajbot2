package common

import "github.com/go-sql-driver/mysql"

func IsDuplicateKey(err error) bool {
	me, ok := err.(*mysql.MySQLError)
	if !ok {
		return false
	}

	return me.Number == 1062
}
