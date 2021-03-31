package paxos

import (
	"errors"
	"log"
	"strconv"
	"strings"

	"github.com/fudute/GoPaxos/sm"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

var ErrorBadLogFormat = errors.New("bad log format")

func init() {
	DB = NewLevelDB("../db")

	DB.Restore(&proposer, &acceptor, sm.GetKVStatMachineInstance())
}

// [key = index, set key value ]
// [key = index, delete key ]
// [key = index, nop ]

type LevelDB struct {
	db  *leveldb.DB
	opt *opt.Options
	ro  *opt.ReadOptions
	wo  *opt.WriteOptions
}

func NewLevelDB(path string) *LevelDB {

	db := &LevelDB{}
	db.opt = &opt.Options{}
	db.ro = &opt.ReadOptions{}
	db.wo = &opt.WriteOptions{}

	var err error
	db.db, err = leveldb.OpenFile(path, db.opt)

	if err != nil {
		log.Fatal(err)
	}
	return db
}

func (db *LevelDB) WriteLog(index int, entry *LogEntry) error {

	str := entry.String()
	key := strconv.Itoa(index)
	return db.db.Put([]byte(key), []byte(str), db.wo)
}

func (db *LevelDB) ReadLog(index int) (*LogEntry, error) {

	key := strconv.Itoa(index)

	value, err := db.db.Get([]byte(key), db.ro)

	if err != nil {
		if err == leveldb.ErrNotFound {
			return nil, ErrorNotFound
		}
		return nil, err
	}

	le, err := parseLog(string(value))

	if err != nil {
		return nil, err
	}

	return le, nil
}

func (db *LevelDB) Restore(p *Proposer, a *Acceptor, sm sm.StatMachine) error {
	iter := db.db.NewIterator(nil, nil)

	log.Println("start restore...")
	for iter.Next() {
		key := iter.Key()
		value := iter.Value()

		le, err := parseLog(string(value))

		log.Println("restore: log entry ", *le)

		if err != nil {
			return err
		}

		index, err := strconv.Atoi(string(key))
		if err != nil {
			log.Fatal("db format is incorrect")
		}
		if p.LogIndex == index && le.isCommited {
			p.LogIndex++
			sm.Execute(le.AcceptedValue)
		}
	}
	return nil
}

func (db *LevelDB) Close() error {
	return db.db.Close()
}

func parseLog(str string) (*LogEntry, error) {
	var err error

	le := &LogEntry{}
	tokens := strings.Split(str, ":")

	if len(tokens) < 3 {
		return nil, ErrorBadLogFormat
	}
	le.MinProposal, err = strconv.Atoi(tokens[0])
	if err != nil {
		return nil, ErrorBadLogFormat
	}

	le.AcceptedProposal, err = strconv.Atoi(tokens[1])
	if err != nil {
		return nil, ErrorBadLogFormat
	}

	le.AcceptedValue = tokens[2]

	le.isCommited, err = strconv.ParseBool(tokens[3])
	if err != nil {
		return nil, ErrorBadLogFormat
	}

	return le, nil
}
