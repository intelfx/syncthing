// Copyright (C) 2025 The Syncthing Authors.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this file,
// You can obtain one at https://mozilla.org/MPL/2.0/.

//go:build cgo

package sqlite

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/mattn/go-sqlite3" // also registers sqlite3 database driver
)

const (
	dbDriver      = "sqlite3_ex"
	commonOptions = "_fk=true&_rt=true&_sync=1&_txlock=immediate"
)

var connPragmas []baseDBPragma

func init() {
	connPragmas = []baseDBPragma{
		{"temp_store", "MEMORY"},
	}

	sql.Register(dbDriver, &sqlite3.SQLiteDriver{
		ConnectHook: func(conn *sqlite3.SQLiteConn) error {
			// Connection-scoped pragmas: must be applied on every sqlite3* handle.
			for _, pragma := range connPragmas {
				stmt := fmt.Sprintf("PRAGMA %s = %s", pragma.K, pragma.V)
				if _, err := conn.Exec(stmt, nil); err != nil {
					return wrap(err, stmt)
				}
			}

			return nil
		},
	})
}
