
// Package db handles database connection, migrations, and queries
package db

import (
    "database/sql"
)

// DB is the global database connection pointer used by all handlers
var DB *sql.DB