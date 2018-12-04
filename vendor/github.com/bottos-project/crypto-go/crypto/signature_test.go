package crypto

import "testing"
import "fmt"
import "encoding/hex"

func Test_GenerateKey(t *testing.T) {
	//t.Log(GenerateKey())
	x,y := GenerateKey()
	fmt.Println("public  key: ", hex.EncodeToString(x))
    fmt.Println("private key: ", hex.EncodeToString(y))
	
}



