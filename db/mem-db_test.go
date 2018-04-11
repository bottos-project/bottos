package db_test
import (

"testing"
	"fmt"
	"github.com/bottos-project/core/db"
)

func TestNewDatabase(t *testing.T) {

	db,err := db.NewMemDatabase()
	fmt.Println(db,err)
}
