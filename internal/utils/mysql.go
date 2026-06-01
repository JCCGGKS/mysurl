package utils

import (
	"errors"
	"strings"

	mysqlDriver "github.com/go-sql-driver/mysql"
)

const mysqlDuplicateEntry = 1062

func IsDuplicateEntryError(err error, indexName string) bool {
	var mysqlErr *mysqlDriver.MySQLError
	if !errors.As(err, &mysqlErr) {
		return false
	}

	if mysqlErr.Number != mysqlDuplicateEntry {
		return false
	}

	if indexName == "" {
		return true
	}

	return strings.Contains(mysqlErr.Message, indexName)
}
