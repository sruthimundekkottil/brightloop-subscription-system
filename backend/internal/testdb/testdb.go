// Package testdb points g.DB() at a freshly migrated MySQL database for one test package.
package testdb

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	"github.com/gogf/gf/v2/database/gdb"
)

// TEST_MYSQL_ADDR defaults to 127.0.0.1:3306 (root/root, as in docker-compose.yaml).
func addr() string {
	if a := os.Getenv("TEST_MYSQL_ADDR"); a != "" {
		return a
	}
	return "127.0.0.1:3306"
}

// Connect makes existing database `name` the default group for g.DB().
func Connect(name string) {
	gdb.SetConfig(gdb.Config{"default": gdb.ConfigGroup{{
		Link: fmt.Sprintf("mysql:root:root@tcp(%s)/%s?loc=UTC&parseTime=true&multiStatements=true", addr(), name),
	}}})
}

// Setup drops and recreates database `name`, applies every migration, and connects to it.
func Setup(name string) error {
	host := addr()
	ctx := context.Background()
	admin, err := gdb.New(gdb.ConfigNode{Link: fmt.Sprintf("mysql:root:root@tcp(%s)/mysql", host)})
	if err != nil {
		return err
	}
	defer admin.Close(ctx)
	if _, err = admin.Exec(ctx, fmt.Sprintf("DROP DATABASE IF EXISTS `%s`", name)); err != nil {
		return err
	}
	if _, err = admin.Exec(ctx, fmt.Sprintf("CREATE DATABASE `%s`", name)); err != nil {
		return err
	}

	Connect(name)
	_, self, _, _ := runtime.Caller(0)
	migrations, err := filepath.Glob(filepath.Join(filepath.Dir(self), "../../db/migrations/*.up.sql"))
	if err != nil || len(migrations) == 0 {
		return fmt.Errorf("no migrations found: %v", err)
	}
	db, err := gdb.Instance()
	if err != nil {
		return err
	}
	for _, m := range migrations { // Glob sorts, so 000001 runs before 000002
		sql, err := os.ReadFile(m)
		if err != nil {
			return err
		}
		if _, err = db.Exec(ctx, string(sql)); err != nil {
			return fmt.Errorf("%s: %w", filepath.Base(m), err)
		}
	}
	return nil
}
