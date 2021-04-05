// +build !pgonline

// This code is available on the terms of the project LICENSE.md file,
// also available online at https://blueoakcouncil.org/license/1.0.0.

package pg

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"os"
	"regexp"
	"testing"

	"decred.org/dcrdex/dex"
)

func TestMain(m *testing.M) {
	startLogger()

	mktInfo, err := dex.NewMarketInfoFromSymbols("dcr", "btc", LotSize, EpochDuration, MarketBuyBuffer)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	AssetDCR = mktInfo.Base
	AssetBTC = mktInfo.Quote

	sql.Register("stub", &dbStub{})

	os.Exit(m.Run())
}

// driver.Driver
type dbStub struct{}

var _ driver.Driver = (*dbStub)(nil)

func (db *dbStub) Open(name string) (driver.Conn, error) {
	return &dbStubConn{}, nil
}

// driver.Conn
type dbStubConn struct{}

var _ driver.Conn = (*dbStubConn)(nil)

func (dbc *dbStubConn) Prepare(query string) (driver.Stmt, error) {
	re := regexp.MustCompile(`\$\d+`)
	matches := re.FindAllStringIndex(query, -1)
	return &dbStubStmt{len(matches)}, nil
}
func (dbc *dbStubConn) Close() error              { return nil }
func (dbc *dbStubConn) Begin() (driver.Tx, error) { return &dbStubTx{}, nil }

// driver.Tx
type dbStubTx struct{}

var _ driver.Tx = (*dbStubTx)(nil)

func (dbt *dbStubTx) Commit() error   { return nil }
func (dbt *dbStubTx) Rollback() error { return nil }

// driver.Stmt
type dbStubStmt struct {
	numPlaceholders int
}

var _ driver.Stmt = (*dbStubStmt)(nil)

func (dbs *dbStubStmt) Close() error  { return nil }
func (dbs *dbStubStmt) NumInput() int { return dbs.numPlaceholders }
func (dbs *dbStubStmt) Exec(args []driver.Value) (driver.Result, error) {
	return &dbStubResult{args}, nil
}
func (dbs *dbStubStmt) Query(args []driver.Value) (driver.Rows, error) {
	return &dbStubRows{}, nil
}

// driver.Rows
type dbStubRows struct{}

var _ driver.Rows = (*dbStubRows)(nil)

func (dbr *dbStubRows) Columns() []string {
	return nil
}
func (dbr *dbStubRows) Close() error                   { return nil }
func (dbr *dbStubRows) Next(dest []driver.Value) error { return nil }

// driver.Result
type dbStubResult struct {
	values []driver.Value
}

var _ driver.Result = (*dbStubResult)(nil)

func (dbs *dbStubResult) LastInsertId() (int64, error) {
	return 0, nil
}
func (dbs *dbStubResult) RowsAffected() (int64, error) {
	return 1, nil
}

// Values is not required for driver.Result
func (dbs *dbStubResult) Values() []driver.Value {
	return dbs.values
}
