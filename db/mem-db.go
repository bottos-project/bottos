package db

/*
 * This is a test memory database. Do not use for any production it does not get persisted
 */
type MemDb struct {
	db map[string][]byte
}

func NewMemDatabase() (*MemDb, error) {
	db := &MemDb{db: make(map[string][]byte)}

	return db, nil
}

func (db *MemDb) Put(key []byte, value []byte) error {
	db.db[string(key)] = value

	return nil
}

func (db *MemDb) Set(key []byte, value []byte) {
	db.Put(key, value)
}

func (db *MemDb) Get(key []byte) ([]byte, error) {
	return db.db[string(key)], nil
}

/*
func (db *MemDatabase) GetKeys() []*common.Key {
	data, _ := db.Get([]byte("KeyRing"))

	return []*common.Key{common.NewKeyFromBytes(data)}
}
*/

func (db *MemDb) Delete(key []byte) error {
	delete(db.db, string(key))

	return nil
}

//func (db *MemDb) Print() {
//	for key, val := range db.db {
//		fmt.Printf("%x(%d): ", key, len(key))
//		node := common.NewValueFromBytes(val)
//		fmt.Printf("%q\n", node.Val)
//	}
//}

func (db *MemDb) Close() {
}

func (db *MemDb) LastKnownTD() []byte {
	data, _ := db.Get([]byte("LastKnownTotalDifficulty"))

	if len(data) == 0 || data == nil {
		data = []byte{0x0}
	}

	return data
}

func (db *MemDb) Flush() error {
	return nil
}
