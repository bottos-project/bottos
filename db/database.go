package db

import (
	"fmt"
	"sync"
	"time"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/errors"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/opt"

	gometrics "github.com/rcrowley/go-metrics"
)

var OpenFileLimit = 64

type KVDatabase struct {
	fn string      // filename for reporting
	db *leveldb.DB // LevelDB instance

	getTimer       gometrics.Timer // Timer for measuring the database get request counts and latencies
	putTimer       gometrics.Timer // Timer for measuring the database put request counts and latencies
	delTimer       gometrics.Timer // Timer for measuring the database delete request counts and latencies
	missMeter      gometrics.Meter // Meter for measuring the missed database get requests
	readMeter      gometrics.Meter // Meter for measuring the database get request data usage
	writeMeter     gometrics.Meter // Meter for measuring the database put request data usage
	compTimeMeter  gometrics.Meter // Meter for measuring the total time spent in database compaction
	compReadMeter  gometrics.Meter // Meter for measuring the data read during compaction
	compWriteMeter gometrics.Meter // Meter for measuring the data written during compaction

	quitLock sync.Mutex      // Mutex protecting the quit channel access
	quitChan chan chan error // Quit channel to stop the metrics collection before closing the database
}

// NewLDBDatabase returns a LevelDB wrapped object. LDBDatabase does not persist data by
// it self but requires a background poller which syncs every X. `Flush` should be called
// when data needs to be stored and written to disk.
func NewKVDatabase(file string) (*KVDatabase, error) {
	// Open the db
	db, err := leveldb.OpenFile(file, &opt.Options{OpenFilesCacheCapacity: OpenFileLimit})
	// check for corruption and attempt to recover
	if _, iscorrupted := err.(*errors.ErrCorrupted); iscorrupted {
		db, err = leveldb.RecoverFile(file, nil)
	}
	// (re) check for errors and abort if opening of the db failed
	if err != nil {
		return nil, err
	}
	return &KVDatabase{
		fn: file,
		db: db,
	}, nil
}

// Put puts the given key / value to the queue
func (self *KVDatabase) Put(key []byte, value []byte) error {
	// Measure the database put latency, if requested
	if self.putTimer != nil {
		defer self.putTimer.UpdateSince(time.Now())
	}
	// Generate the data to write to disk, update the meter and write
	//value = rle.Compress(value)

	if self.writeMeter != nil {
		self.writeMeter.Mark(int64(len(value)))
	}
	return self.db.Put(key, value, nil)
}

// Get returns the given key if it's present.
func (self *KVDatabase) Get(key []byte) ([]byte, error) {
	// Measure the database get latency, if requested
	if self.getTimer != nil {
		defer self.getTimer.UpdateSince(time.Now())
	}
	// Retrieve the key and increment the miss counter if not found
	dat, err := self.db.Get(key, nil)
	if err != nil {
		if self.missMeter != nil {
			self.missMeter.Mark(1)
		}
		return nil, err
	}
	// Otherwise update the actually retrieved amount of data
	if self.readMeter != nil {
		self.readMeter.Mark(int64(len(dat)))
	}
	return dat, nil
	//return rle.Decompress(dat)
}

// Delete deletes the key from the queue and database
func (self *KVDatabase) Delete(key []byte) error {
	// Measure the database delete latency, if requested
	if self.delTimer != nil {
		defer self.delTimer.UpdateSince(time.Now())
	}
	// Execute the actual operation
	return self.db.Delete(key, nil)
}

func (self *KVDatabase) NewIterator() iterator.Iterator {
	return self.db.NewIterator(nil, nil)
}

// Flush flushes out the queue to leveldb
func (self *KVDatabase) Flush() error {
	return nil
}

func (self *KVDatabase) Close() {
	// Stop the metrics collection to avoid internal database races
	self.quitLock.Lock()
	defer self.quitLock.Unlock()

	if self.quitChan != nil {
		errc := make(chan error)
		self.quitChan <- errc
		if err := <-errc; err != nil {
			//glog.V(logger.Error).Infof("metrics failure in '%s': %v\n", self.fn, err)
			fmt.Println("err", err)
		}
	}
	// Flush and close the database
	if err := self.Flush(); err != nil {
		//glog.V(logger.Error).Infof("flushing '%s' failed: %v\n", self.fn, err)
		fmt.Println("err", err)
	}
	self.db.Close()
	//glog.V(logger.Error).Infoln("flushed and closed db:", self.fn)
	fmt.Println("flushed and closed db:", self.fn)
}

func (self *KVDatabase) LDB() *leveldb.DB {
	return self.db
}
