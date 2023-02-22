package testutil

import (
	"context"
	"crypto/rand"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/kelseyhightower/envconfig"
	"hello-circleci/pkg/db"
	"math/big"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

type Env struct {
	MySQL MySQL
}

type MySQL struct {
	Host     string `envconfig:"MYSQL_HOST" required:"true"`
	Port     string `envconfig:"MYSQL_PORT" required:"true"`
	User     string `envconfig:"MYSQL_USER" required:"true"`
	Password string `envconfig:"MYSQL_PASSWORD" required:"true"`
}

func TearUp(t *testing.T) (*db.DBContext, func()) {
	t.Helper()

	var c Env
	if err := envconfig.Process("", &c); err != nil {
		t.Fatalf("failed to envconfig.Process() :%v", err)
	}

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to os.Getwd() :%v", err)
	}

	targetRoot, err := filepath.Abs(filepath.Join(wd, "../../hack/db/scripts"))
	if err != nil {
		t.Fatalf("failed to filepath.Abs :%v", err)
	}

	rootDB := connectDB(MySQL{
		Host:     c.MySQL.Host,
		Port:     c.MySQL.Port,
		User:     "root",
		Password: "",
	}, "")

	dbName := GenRandomString(10)
	for _, v := range []string{
		fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s;", dbName),
		fmt.Sprintf("GRANT ALL ON %s.* TO %s@'%%';", dbName, c.MySQL.User),
		fmt.Sprintf("UPDATE mysql.user SET Super_priv='Y' WHERE user='%s'", c.MySQL.User),
		"FLUSH PRIVILEGES;",
	} {
		if _, err := rootDB.Exec(context.TODO(), v); err != nil {
			t.Logf("query=%s", v)
			t.Fatalf("failed to db.Exec(create database) err=%v", err)
		}
	}

	db := connectDB(c.MySQL, dbName)
	if _, err := db.Exec(context.TODO(), fmt.Sprintf("USE %s;", dbName)); err != nil {
		t.Fatalf("failed to db.Exec(use database) err=%v", err)
	}

	err = filepath.Walk(targetRoot, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		// exclude not .sql file
		// datalake's file (.sql_) do not reload (cause drone test error)
		if filepath.Ext(info.Name()) != ".sql" {
			return nil
		}

		text, e := os.ReadFile(path)
		if e != nil {
			t.Fatalf("failed to ioutil.ReadFile(path) [path=%v, err=%v]", path, e)
		}

		query := []string{}
		delimiter := ";"
		newDelimiter := regexp.MustCompile(`^DELIMITER\s+(\S*)\s*$`)
		for _, line := range strings.Split(string(text), "\n") {
			if strings.TrimSpace(line) == "" {
				continue
			}

			if newDelimiter.MatchString(line) {
				l := newDelimiter.FindStringSubmatch(line)
				delimiter = l[1]
				continue
			}

			query = append(query, line)

			if strings.HasSuffix(line, delimiter) {
				q := strings.TrimRight(strings.Join(query, " "), delimiter)

				if _, err := db.Exec(context.TODO(), q); err != nil {
					t.Logf("query=%s", q)
					t.Fatalf("file: %s, err: %+v", path, err)
				}
				query = []string{}
			}
		}

		return nil
	})
	if err != nil {
		t.Fatalf("failed to filepath.Walk :%v", err)
	}

	return db, func() {
		_, err := db.Exec(context.TODO(), fmt.Sprintf("DROP SCHEMA IF EXISTS %s;", dbName))
		if err != nil {
			t.Fatal(err)
		}
	}
}

func connectDB(c MySQL, dbName string) *db.DBContext {
	d, err := sqlx.Connect("mysql", fmt.Sprintf("%s:%s@(%s:%s)/%s?parseTime=true&loc=Asia%%2FTokyo", c.User, c.Password, c.Host, c.Port, dbName))
	if err != nil {
		panic(err)
	}

	return db.NewDBWithDB(d)
}

const rs2Letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// https://qiita.com/srtkkou/items/ccbddc881d6f3549baf1#2-const%E3%82%92%E4%BD%BF%E3%81%A3%E3%81%A6%E3%81%BF%E3%82%8B
func GenRandomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(rs2Letters))))
		b[i] = rs2Letters[n.Int64()]

	}
	return string(b)
}
