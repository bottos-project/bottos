package main

import (
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"encoding/json"
	//"math"
	//"testing"
	//"github.com/ontio/ontology-wasm/exec"
	//"github.com/ontio/ontology-wasm/memory"
	//"github.com/ontio/ontology-wasm/util"
	//"github.com/bottos-project/core/vm"
	"github.com/bottos-project/core/vm/exec"
	"reflect"
	"os"
	//"golang.org/x/text/currency"
)

//var service = NewInteropService()

func add() {
	engine := exec.NewExecutionEngine(nil, "test")

	code, err := ioutil.ReadFile("C:\\Users\\stewa\\go\\src\\github.com\\ontio\\ontology-wasm\\exec\\test_data\\math.wasm")
	//code, err := ioutil.ReadFile("/home/iwojima/iWork/go/src/github.com/ontio/ontology-wasm/exec/test_data/math.wasm")
	//code, err := ioutil.ReadFile("C:\\Users\\stewa\\go\\src\\github.com\\ontio\\ontology-wasm\\exec\\test_data2\\contract.wasm")
	if err != nil {
		fmt.Println("error in read file", err.Error())
		return
	}

	fmt.Println("main.add() code's type = ",reflect.TypeOf(code))

	method2 := "add"
	input2 := make([]byte, 9)
	input2[0] = byte(len(method2))
	copy(input2[1:len(method2)+1], []byte(method2))

	fmt.Println("[]byte(method2) = ",[]byte(method2))

	input2[len(method2)+1] = byte(2) //param count
	input2[len(method2)+2] = byte(1) //param1 length
	input2[len(method2)+3] = byte(1) //param2 length
	input2[len(method2)+4] = byte(5) //param1
	input2[len(method2)+5] = byte(9) //param2

	fmt.Println(input2)
	res2, err := engine.Call(nil, code, input2)
	if err != nil {
		fmt.Println("call error!", err.Error())
	}
	fmt.Printf("res:%v\n", res2)
	fmt.Println("binary.LittleEndian.Uint32(res) = ",binary.LittleEndian.Uint32(res2))
	if binary.LittleEndian.Uint32(res2) != uint32(14) {
		fmt.Println("the result should be 14")
	}

}

func square() {
	engine := exec.NewExecutionEngine(nil, "test")
	code, err := ioutil.ReadFile("C:\\Users\\stewa\\go\\src\\github.com\\ontio\\ontology-wasm\\exec\\test_data\\math.wasm")
	//code, err := ioutil.ReadFile("/home/iwojima/iWork/go/src/github.com/ontio/ontology-wasm/exec/test_data/math.wasm")
	if err != nil {
		fmt.Println("error in read file", err.Error())
		return
	}
	method := "square"

	input := make([]byte, 10)
	input[0] = byte(len(method))
	copy(input[1:len(method)+1], []byte(method))
	input[len(method)+1] = byte(1) //param count
	input[len(method)+2] = byte(1) //param1 length
	input[len(method)+3] = byte(5) //param1

	fmt.Println(input)
	res, err := engine.Call(nil, code, input)
	if err != nil {
		fmt.Println("call error!", err.Error())
	}
	fmt.Printf("res:%v\n", res)
	//fmt.Println("binary.LittleEndian.Uint32(res) = ",binary.LittleEndian.Uint32(res))
	//fmt.Println("uint32(25) = ",vm.ctx.stack)
	if binary.LittleEndian.Uint32(res) != uint32(25) {
		fmt.Println("the result should be 14")
	}


}

/*
func TestContract1() {
	engine := exec.NewExecutionEngine(nil, "product")
	//test
	code, err := ioutil.ReadFile("C:\\Users\\stewa\\go\\src\\github.com\\ontio\\ontology-wasm\\exec\\test_data2\\contract.wasm")
	//code, err := ioutil.ReadFile("C:\\Users\\stewa\\go\\src\\github.com\\ontio\\documentation\\smart-contract-tutorial\\examples\\contract.wasm")
	//
	if err != nil {
		fmt.Println("error in read file", err.Error())
		return
	}

	//code 是十进制标识的WASM
	//fmt.Println("code = ",code)

	par := make([]exec.Param, 2)
	par[0] = exec.Param{Ptype: "int", Pval: "20"}
	par[1] = exec.Param{Ptype: "int", Pval: "30"}

	p := exec.Args{Params: par}
	bytes, err := json.Marshal(p)
	if err != nil {
		fmt.Println(err)
		//t.Fatal(err.Error())
		fmt.Println(err.Error())
	}

	fmt.Println("<------------------------>")
	fmt.Println("string(bytes) = ", string(bytes))
	fmt.Println("<------------------------>")

	input := make([]interface{}, 3)
	input[0] = "invoke"
	input[1] = "add"
	input[2] = string(bytes)

	////code 是十进制标识的WASM
	fmt.Println("code = ", code)
	//input是导入命令加操作命令，后面跟j'son
	fmt.Println("input = ", input)

	res, err := engine.CallInf(nil, code, input, nil)
	if err != nil {
		fmt.Println("call error!", err.Error())
	}
	fmt.Printf("res:%v\n", res)

	//retbytes, err := engine.vm.GetPointerMemory(uint64(binary.LittleEndian.Uint32(res)))
	retbytes, err := engine.GetVM().GetPointerMemory(uint64(binary.LittleEndian.Uint32(res)))
	if err != nil {
		fmt.Println(err)
		//t.Fatal("errors:" + err.Error())
		fmt.Println("errors:", err.Error())
	}

	fmt.Println("retbytes is " + string(retbytes))

	result := &exec.Result{}
	json.Unmarshal(retbytes, result)

	//fmt.Println(engine.GetVM().GetMemory().Memory[:20])

	//fmt.Println(engine.GetVM().GetMemory().Memory[16384:])

	//fmt.Println(string(engine.GetVM().GetMemory().Memory[7:50]))

	if result.Pval != "50" {
		//t.Fatal("result should be 50")
		fmt.Println("result should be 50")
	}
}
*/

func TestContract2() error {
	engine := exec.NewExecutionEngine(nil, "product")
	//test
	code, err := ioutil.ReadFile("C:\\Users\\stewa\\Desktop\\wasm_case\\currency.wasm")
	//code, err := ioutil.ReadFile("/home/iwojima/iWork/go/src/github.com/ontio/ontology-wasm/exec/test_data2/contract.wasm")
	if err != nil {
		fmt.Println("error in read file", err.Error())
		return err
	}

	//code 是十进制标识的WASM
	//fmt.Println("code = ",code)
	var bytes []byte
	debug := false
	if debug { //make a easy json file
		par := make([]exec.Param, 2)
		par[0] = exec.Param{Ptype: "int", Pval: "20"}
		par[1] = exec.Param{Ptype: "int", Pval: "30"}

		p := exec.Args{Params: par}
		bytes, err = json.Marshal(p)
		if err != nil {
			fmt.Println(err)
			//t.Fatal(err.Error())
			fmt.Println(err.Error())
		}
	}else{ //read a json of EOS
		filename := "C:\\Users\\stewa\\Desktop\\wasm_case\\currency.abi"
		bytes, err = ioutil.ReadFile(filename)
		if err != nil {
			return err
		}
	}
	//

	fmt.Println("<------------------------>")
	fmt.Println("bytes = ",bytes)
	//fmt.Println("string(bytes) = ",string(bytes))  //代表需要导入的json
	fmt.Println("<------------------------>")

	input := make([]interface{}, 3)
	input[0] = "init"
	input[1] = "apply"
	input[2] = string(bytes) //代表需要导入的json文件

	////code 是十进制标识的WASM
	fmt.Println("code = ",code)
	//input是导入命令加操作命令，后面跟json
	fmt.Println("input = ",input)

	res, err := engine.CallInf(nil, code, input, nil)
	if err != nil {
		fmt.Println("call error!", err.Error())
	}
	fmt.Printf("res:%v\n", res)

	//retbytes, err := engine.vm.GetPointerMemory(uint64(binary.LittleEndian.Uint32(res)))

	/*
	retbytes, err := engine.GetVM().GetPointerMemory(uint64(binary.LittleEndian.Uint32(res)))
	if err != nil {
		fmt.Println(err)
		//t.Fatal("errors:" + err.Error())
		fmt.Println("errors:" , err.Error())
	}

	fmt.Println("retbytes is " + string(retbytes))

	result := &exec.Result{}
	json.Unmarshal(retbytes, result)

	//fmt.Println(engine.GetVM().GetMemory().Memory[:20])

	//fmt.Println(engine.GetVM().GetMemory().Memory[16384:])

	//fmt.Println(string(engine.GetVM().GetMemory().Memory[7:50]))

	if result.Pval != "50" {
		//t.Fatal("result should be 50")
		fmt.Println("result should be 50")
	}
	*/
	return nil
}

func TestContract1() {
	engine := exec.NewExecutionEngine(nil, "product")//这个engine是提供给客户端执行的句并
	//test
	//code, err := ioutil.ReadFile("C:\\Users\\stewa\\go\\src\\github.com\\ontio\\ontology-wasm\\exec\\test_data2\\contract.wasm")
	code, err := ioutil.ReadFile("C:\\Users\\stewa\\go\\src\\github.com\\bottos-project\\core\\vm\\exec\\bi_contracts1\\contract.wasm")
	if err != nil {
		fmt.Println("error in read file", err.Error())
		return
	}

	par := make([]exec.Param, 2)
	par[0] = exec.Param{Ptype: "int", Pval: "20"}
	par[1] = exec.Param{Ptype: "int", Pval: "30"}

	p := exec.Args{Params: par}
	bytes, err := json.Marshal(p)
	if err != nil {
		fmt.Println(err)
		//t.Fatal(err.Error())
	}
	fmt.Println(string(bytes))

	input := make([]interface{}, 3)
	input[0] = "invoke"
	input[1] = "add"
	input[2] = string(bytes)

	fmt.Printf("input is %v\n", input)

	res, err := engine.CallInf(nil, code, input, nil)
	if err != nil {
		fmt.Println("call error!", err.Error())
	}
	fmt.Printf("res:%v\n", res)

	retbytes, err := engine.GetVM().GetPointerMemory(uint64(binary.LittleEndian.Uint32(res)))
	if err != nil {
		fmt.Println(err)
		//t.Fatal("errors:" + err.Error())
	}

	fmt.Println("======================= retbytes is " + string(retbytes))

	result := &exec.Result{}
	json.Unmarshal(retbytes, result)

	fmt.Println(engine.GetVM().GetMemory().Memory[:20])
	fmt.Println(engine.GetVM().GetMemory().Memory[16384:])

	fmt.Println(string(engine.GetVM().GetMemory().Memory[7:50]))

	if result.Pval != "50" {
		fmt.Println("result should be 50")
	}
}


type stlHead struct {
	Name [80]byte
	FaceNum uint32
}


func read_bin() {
	file, err := os.Open("C:\\Windows\\Boot\\EFI\\boot.stl")
	if err != nil {
		fmt.Print(err)
		return
	}

	head := new(stlHead)

	if err := binary.Read(file, binary.LittleEndian, head); err != nil {
		fmt.Print(err)
		return
	}

	fmt.Println("main.read_bin()")

	fmt.Printf("name: %s\r\n", head.Name)
	fmt.Println("FaceNum: ", head.FaceNum)
}

func main() {
	//add()
	//fmt.Println("<----------------------------->")
	//square()
	TestContract1()
	//TestContract2()

	//read_bin()
}
