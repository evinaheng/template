package sqlt

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/jmoiron/sqlx"
)

// Error list
var (
	ErrNoConnectionDetected = errors.New("No connection detected")
)

// DB struct wrapper for sqlx connection
type DB struct {
	sqlxdb     []*sqlx.DB
	activedb   []int
	inactivedb []int
	dsn        []string
	driverName string
	groupName  string
	length     int
	count      uint64
	// for stats
	stats     []DbStatus
	heartBeat bool
	stopBeat  chan bool
	lastBeat  string
	// only use when needed
	debug bool
}

// DbStatus for status response
type DbStatus struct {
	Name       string      `json:"name"`
	Connected  bool        `json:"connected"`
	LastActive string      `json:"last_active"`
	Error      interface{} `json:"error"`
}

type statusResponse struct {
	Dbs       interface{} `json:"db_list"`
	Heartbeat bool        `json:"heartbeat"`
	Lastbeat  string      `json:"last_beat"`
}

const defaultGroupName = "sqlt_open"

var dbLengthMutex = &sync.Mutex{}

func openConnection(driverName, sources string, groupName string) (*DB, error) {
	var err error

	conns := strings.Split(sources, ";")
	connsLength := len(conns)

	// check if no source is available
	if connsLength < 1 {
		return nil, errors.New("No sources found")
	}

	db := &DB{
		sqlxdb: make([]*sqlx.DB, connsLength),
		stats:  make([]DbStatus, connsLength),
		dsn:    make([]string, connsLength),
	}
	db.length = connsLength
	db.driverName = driverName

	for i := range conns {
		db.sqlxdb[i], err = sqlx.Open(driverName, conns[i])
		if err != nil {
			db.inactivedb = append(db.inactivedb, i)
			return nil, err
		}

		constatus := true

		// set the name
		name := ""
		if i == 0 {
			name = "master"
		} else {
			name = "slave-" + strconv.Itoa(i)
		}

		status := DbStatus{
			Name:       name,
			Connected:  constatus,
			LastActive: time.Now().String(),
		}

		db.stats[i] = status
		db.activedb = append(db.activedb, i)
		db.dsn[i] = conns[i]
	}

	// set the default group name
	db.groupName = defaultGroupName
	if groupName != "" {
		db.groupName = groupName
	}

	// ping database to retrieve error
	err = db.Ping()
	return db, err
}

// Open connection to database
func Open(driverName, sources string) (*DB, error) {
	return openConnection(driverName, sources, "")
}

// OpenWithName open the connection and set connection group name
func OpenWithName(driverName, sources string, name string) (*DB, error) {
	return openConnection(driverName, sources, name)
}

// SetDebug for sqlt
func (db *DB) SetDebug(v bool) {
	db.debug = v
}

// GetStatus return database status
func (db *DB) GetStatus() ([]DbStatus, error) {
	if len(db.stats) == 0 {
		return db.stats, ErrNoConnectionDetected
	}

	// if heartbeat is not enabled, ping to get status before send status
	if !db.heartBeat {
		db.Ping()
	}
	return db.stats, nil
}

// DoHeartBeat will automatically spawn a goroutines to ping your database every one second, use this carefully
func (db *DB) DoHeartBeat() {
	if !db.heartBeat {
		ticker := time.NewTicker(time.Second * 2)
		db.stopBeat = make(chan bool)
		go func() {
			for {
				select {
				case <-ticker.C:
					db.Ping()
					db.lastBeat = time.Now().Format(time.RFC1123)
				case <-db.stopBeat:
					return
				}
			}
		}()
	}
	db.heartBeat = true
}

// StopBeat will stop heartbeat, exit from goroutines
func (db *DB) StopBeat() {
	if !db.heartBeat {
		return
	}
	db.stopBeat <- true
}

// Ping database
func (db *DB) Ping() error {
	var err error

	if !db.heartBeat {
		for _, val := range db.sqlxdb {
			err = val.Ping()
			if err != nil {
				return err
			}
		}
		return err
	}

	for i := 0; i < len(db.activedb); i++ {
		val := db.activedb[i]
		err = db.sqlxdb[val].Ping()
		name := db.stats[val].Name

		if err != nil {
			if db.length <= 1 {
				return err
			}

			db.stats[val].Connected = false
			db.activedb = append(db.activedb[:i], db.activedb[i+1:]...)
			i--
			db.inactivedb = append(db.inactivedb, val)
			db.stats[val].Error = errors.New(name + ": " + err.Error())
			dbLengthMutex.Lock()
			db.length--
			dbLengthMutex.Unlock()
		} else {
			db.stats[val].Connected = true
			db.stats[val].LastActive = time.Now().Format(time.RFC1123)
			db.stats[val].Error = nil
		}
	}

	for i := 0; i < len(db.inactivedb); i++ {
		val := db.inactivedb[i]
		err = db.sqlxdb[val].Ping()
		name := db.stats[val].Name

		if err != nil {
			db.stats[val].Connected = false
			db.stats[val].Error = errors.New(name + ": " + err.Error())
		} else {
			db.stats[val].Connected = true
			db.inactivedb = append(db.inactivedb[:i], db.inactivedb[i+1:]...)
			i--
			db.activedb = append(db.activedb, val)
			db.stats[val].LastActive = time.Now().Format(time.RFC1123)
			db.stats[val].Error = nil
			dbLengthMutex.Lock()
			db.length++
			dbLengthMutex.Unlock()
		}
	}
	return err
}

// Prepare return sql stmt
func (db *DB) Prepare(query string) (Stmt, error) {
	var err error
	stmt := Stmt{}
	stmts := make([]*sql.Stmt, len(db.sqlxdb))

	for i := range db.sqlxdb {
		stmts[i], err = db.sqlxdb[i].Prepare(query)

		if err != nil {
			return stmt, err
		}
	}
	stmt.db = db
	stmt.stmts = stmts
	return stmt, nil
}

// Preparex sqlx stmt
func (db *DB) Preparex(query string) (*Stmtx, error) {
	var err error
	stmts := make([]*sqlx.Stmt, len(db.sqlxdb))

	for i := range db.sqlxdb {
		stmts[i], err = db.sqlxdb[i].Preparex(query)

		if err != nil {
			return nil, err
		}
	}

	return &Stmtx{db: db, stmts: stmts}, nil
}

// SetMaxOpenConnections to set max connections
func (db *DB) SetMaxOpenConnections(max int) {
	for i := range db.sqlxdb {
		db.sqlxdb[i].SetMaxOpenConns(max)
	}
}

// SetMaxIdleConnections to set max idle connections
func (db *DB) SetMaxIdleConnections(max int) {
	for i := range db.sqlxdb {
		db.sqlxdb[i].SetMaxIdleConns(max)
	}
}

// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
// Expired connections may be closed lazily before reuse.
// If d <= 0, connections are reused forever.
func (db *DB) SetConnMaxLifetime(d time.Duration) {
	for i := range db.sqlxdb {
		db.sqlxdb[i].SetConnMaxLifetime(d)
	}
}

// Slave return slave database
func (db *DB) Slave() *sqlx.DB {
	return db.sqlxdb[db.slave()]
}

// Master return master database
func (db *DB) Master() *sqlx.DB {
	return db.sqlxdb[0]
}

// Query queries the database and returns an *sql.Rows.
func (db *DB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	r, err := db.sqlxdb[db.slave()].Query(query, args...)
	return r, err
}

// QueryRow queries the database and returns an *sqlx.Row.
func (db *DB) QueryRow(query string, args ...interface{}) *sql.Row {
	rows := db.sqlxdb[db.slave()].QueryRow(query, args...)
	return rows
}

// Queryx queries the database and returns an *sqlx.Rows.
func (db *DB) Queryx(query string, args ...interface{}) (*sqlx.Rows, error) {
	r, err := db.sqlxdb[db.slave()].Queryx(query, args...)
	return r, err
}

// QueryRowx queries the database and returns an *sqlx.Row.
func (db *DB) QueryRowx(query string, args ...interface{}) *sqlx.Row {
	rows := db.sqlxdb[db.slave()].QueryRowx(query, args...)
	return rows
}

// Exec using master db
func (db *DB) Exec(query string, args ...interface{}) (sql.Result, error) {
	return db.sqlxdb[0].Exec(query, args...)
}

// MustExec (panic) runs MustExec using master database.
func (db *DB) MustExec(query string, args ...interface{}) sql.Result {
	return db.sqlxdb[0].MustExec(query, args...)
}

// Select using slave db.
func (db *DB) Select(dest interface{}, query string, args ...interface{}) error {
	return db.sqlxdb[db.slave()].Select(dest, query, args...)
}

// SelectMaster using master db.
func (db *DB) SelectMaster(dest interface{}, query string, args ...interface{}) error {
	return db.sqlxdb[0].Select(dest, query, args...)
}

// Get using slave.
func (db *DB) Get(dest interface{}, query string, args ...interface{}) error {
	return db.sqlxdb[db.slave()].Get(dest, query, args...)
}

// GetMaster using master.
func (db *DB) GetMaster(dest interface{}, query string, args ...interface{}) error {
	return db.sqlxdb[0].Get(dest, query, args...)
}

// NamedExec using master db.
func (db *DB) NamedExec(query string, arg interface{}) (sql.Result, error) {
	return db.sqlxdb[0].NamedExec(query, arg)
}

// Begin sql transaction
func (db *DB) Begin() (*sql.Tx, error) {
	return db.sqlxdb[0].Begin()
}

// Beginx sqlx transaction
func (db *DB) Beginx() (*sqlx.Tx, error) {
	return db.sqlxdb[0].Beginx()
}

// MustBegin starts a transaction, and panics on error. Returns an *sqlx.Tx instead
// of an *sql.Tx.
func (db *DB) MustBegin() *sqlx.Tx {
	tx, err := db.sqlxdb[0].Beginx()
	if err != nil {
		panic(err)
	}
	return tx
}

// Rebind query
func (db *DB) Rebind(query string) string {
	return db.sqlxdb[db.slave()].Rebind(query)
}

// RebindMaster will rebind query for master
func (db *DB) RebindMaster(query string) string {
	return db.sqlxdb[0].Rebind(query)
}

// Stmt implement sql stmt
type Stmt struct {
	db    *DB
	stmts []*sql.Stmt
}

// Exec will always go to production
func (st *Stmt) Exec(args ...interface{}) (sql.Result, error) {
	return st.stmts[0].Exec(args...)
}

// Query will always go to slave
func (st *Stmt) Query(args ...interface{}) (*sql.Rows, error) {
	return st.stmts[st.db.slave()].Query(args...)
}

// QueryMaster will use master db
func (st *Stmt) QueryMaster(args ...interface{}) (*sql.Rows, error) {
	return st.stmts[0].Query(args...)
}

// QueryRow will always go to slave
func (st *Stmt) QueryRow(args ...interface{}) *sql.Row {
	return st.stmts[st.db.slave()].QueryRow(args...)
}

// QueryRowMaster will use master db
func (st *Stmt) QueryRowMaster(args ...interface{}) *sql.Row {
	return st.stmts[0].QueryRow(args...)
}

// Close stmt
func (st *Stmt) Close() error {
	for i := range st.stmts {
		err := st.stmts[i].Close()

		if err != nil {
			return err
		}
	}
	return nil
}

// Stmtx implement sqlx stmt
type Stmtx struct {
	db    *DB
	stmts []*sqlx.Stmt
}

// Close all dbs connection
func (st *Stmtx) Close() error {
	for i := range st.stmts {
		err := st.stmts[i].Close()

		if err != nil {
			return err
		}
	}
	return nil
}

// Exec will always go to production
func (st *Stmtx) Exec(args ...interface{}) (sql.Result, error) {
	return st.stmts[0].Exec(args...)

}

// Query will always go to slave
func (st *Stmtx) Query(args ...interface{}) (*sql.Rows, error) {
	return st.stmts[st.db.slave()].Query(args...)
}

// QueryMaster will use master db
func (st *Stmtx) QueryMaster(args ...interface{}) (*sql.Rows, error) {
	return st.stmts[0].Query(args...)
}

// QueryRow will always go to slave
func (st *Stmtx) QueryRow(args ...interface{}) *sql.Row {
	return st.stmts[st.db.slave()].QueryRow(args...)
}

// QueryRowMaster will use master db
func (st *Stmtx) QueryRowMaster(args ...interface{}) *sql.Row {
	return st.stmts[0].QueryRow(args...)
}

// MustExec using master database
func (st *Stmtx) MustExec(args ...interface{}) sql.Result {
	return st.stmts[0].MustExec(args...)
}

// Queryx will always go to slave
func (st *Stmtx) Queryx(args ...interface{}) (*sqlx.Rows, error) {
	return st.stmts[st.db.slave()].Queryx(args...)
}

// QueryRowx will always go to slave
func (st *Stmtx) QueryRowx(args ...interface{}) *sqlx.Row {
	return st.stmts[st.db.slave()].QueryRowx(args...)
}

// QueryRowxMaster will always go to master
func (st *Stmtx) QueryRowxMaster(args ...interface{}) *sqlx.Row {
	return st.stmts[0].QueryRowx(args...)
}

// Get will always go to slave
func (st *Stmtx) Get(dest interface{}, args ...interface{}) error {
	return st.stmts[st.db.slave()].Get(dest, args...)
}

// GetMaster will always go to master
func (st *Stmtx) GetMaster(dest interface{}, args ...interface{}) error {
	return st.stmts[0].Get(dest, args...)
}

// Select will always go to slave
func (st *Stmtx) Select(dest interface{}, args ...interface{}) error {
	return st.stmts[st.db.slave()].Select(dest, args...)
}

// SelectMaster will always go to master
func (st *Stmtx) SelectMaster(dest interface{}, args ...interface{}) error {
	return st.stmts[0].Select(dest, args...)
}

// slave
func (db *DB) slave() int {
	dbLengthMutex.Lock()
	defer dbLengthMutex.Unlock()
	if db.length <= 1 {
		if db.debug {
			fmt.Print("selecting master, slave is not exists")
		}
		return 0
	}

	slave := int(1 + (atomic.AddUint64(&db.count, 1) % uint64(db.length-1)))
	active := db.activedb[slave]
	if db.debug {
		fmt.Printf("slave: %d. dsn: %s", active, db.dsn[active])
	}
	return active
}

//InitMocking initialize the dbconnection mocking
func InitMocking(dbConn *sql.DB, slaveAmount int) *DB {

	db := &DB{
		sqlxdb: make([]*sqlx.DB, slaveAmount+1),
		stats:  make([]DbStatus, slaveAmount+1),
	}

	for i := 0; i <= slaveAmount; i++ {
		db.sqlxdb[i] = sqlx.NewDb(dbConn, "postgres")
		name := fmt.Sprintf("slave-%d", i)
		if i == 0 {
			name = "master"
		}

		db.stats[i] = DbStatus{
			Name:       name,
			Connected:  true,
			LastActive: time.Now().String(),
		}
		db.activedb = append(db.activedb, i)
	}

	db.driverName = "postgres"
	db.groupName = "sqlt-open"
	db.length = len(db.sqlxdb)
	return db
}
