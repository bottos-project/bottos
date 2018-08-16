package main

import (
	//"encoding/binary"
	"fmt"
	"time"
	"reflect"
	//"io/ioutil"
	//"encoding/json"
	//"path/filepath"
	//"github.com/bottos-project/core/vm/testcase"
	"github.com/bottos-project/bottos/vm/wasm/exec"
	//"github.com/bottos-project/core/vm/wasm"
	//"reflect"
	//"unsafe"
	"os"
	//"bytes"
	//"github.com/bottos-project/core/contract"
	//"github.com/bottos-project/core/contract"
	//"github.com/bottos-project/core/db"
	//"github.com/bottos-project/core/role"
	//"github.com/bottos-project/core/config"
	//"golang.org/x/text/internal/number"
	"github.com/bottos-project/bottos/contract"
	"github.com/bottos-project/bottos/common/types"
	log "github.com/cihub/seelog"
	"github.com/bottos-project/bottos/contract/msgpack"
)

//var service = NewInteropService()

/*
func add() {
	engine := exec.NewVMEngine(nil, "test")

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
	engine := exec.NewVMEngine(nil, "test")
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

func TestContract1() {
	engine := exec.NewVMEngine(nil, "product")
	//test
	code, err := ioutil.ReadFile("C:\\Users\\stewa\\go\\src\\github.com\\bottos-project\\core\\vm\\testcase\\test_data2\\contract.wasm")
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
		fmt.Printf(err.Error())
	}
	fmt.Println(string(bytes))

	input := make([]interface{}, 3)
	input[0] = "invoke"
	input[1] = "add"
	input[2] = string(bytes)

	fmt.Printf("<><><> input is %v\n", input)

	//code是读取的wasm文件
	//input是需要执行的方法和参数,例如[invoke add {"Params":[{"type":"int","value":"20"},{"type":"int","value":"30"}]}]
	res, err := engine.CallInf(nil, code, input, nil)
	if err != nil {
		fmt.Println("call error!", err.Error())
	}
	fmt.Printf("res:%v\n", res)

	retbytes, err := engine.GetVM().GetPtrMemory(uint64(binary.LittleEndian.Uint32(res)))
	if err != nil {
		fmt.Println(err)
		fmt.Printf("errors:" + err.Error())
	}

	fmt.Println("retbytes is " + string(retbytes))

	result := &exec.Result{}
	json.Unmarshal(retbytes, result)

	fmt.Println(engine.GetVM().GetMemory()[:20])
	fmt.Println(engine.GetVM().GetMemory()[16384:])

	fmt.Println(string(engine.GetVM().GetMemory()[7:50]))

	if result.Pval != "50" {
		fmt.Println("result should be 50")
	}
}

func bytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
*/

/*
func TestContract2() {
	engine := exec.NewVMEngine(nil, "product")
	//test
	code, err := ioutil.ReadFile("C:\\Users\\stewa\\go\\src\\github.com\\bottos-project\\core\\vm\\testcase\\test_data2\\contract.wasm")
	if err != nil {
		fmt.Println("error in read file", err.Error())
		return
	}

	par := make([]exec.Param, 2)
	par[0] = exec.Param{Ptype: "int", Pval: "20"}
	par[1] = exec.Param{Ptype: "int", Pval: "30"}

	p := exec.Args{Params: par}
	jbytes, err := json.Marshal(p)
	if err != nil {
		fmt.Println(err)
		fmt.Println(err.Error())
	}
	fmt.Println("jbytes = ",string(jbytes))

	bf := bytes.NewBufferString("add")
	bf.WriteString("|")
	bf.Write(jbytes)

	fmt.Printf("input is %v\n", bf.Bytes())
	//fmt.Printf("<<<<<<<<<<<<<<<<<<<<<<<<<< input is %v\n", bytesToString(bf.Bytes()))

	res, err := engine.Call(nil, code, bf.Bytes())
	if err != nil {
		fmt.Println("<<<<<<<< call error!", err.Error())
	}
	fmt.Printf("------------ >res:%v\n", res)

	retbytes, err := engine.GetVM().GetData(uint64(binary.LittleEndian.Uint32(res)))
	if err != nil {
		fmt.Println(err)
		fmt.Println("errors:" + err.Error())
	}

	fmt.Println("retbytes is " + string(retbytes))

	result := &exec.Result{}
	json.Unmarshal(retbytes, result)

	fmt.Println(engine.GetVM().GetMemory()[:20])
	fmt.Println(engine.GetVM().GetMemory()[16384:])

	fmt.Println(string(engine.GetVM().GetMemory()[7:50]))

	fmt.Println("result.Pval = ",result.Pval)
	if result.Pval != "50" {
		fmt.Println("result should be 50")
	}
}
*/

/*
func TestContract22() {
	engine := exec.NewVMEngine(nil, "product")
	//test
	/*code, err := ioutil.ReadFile("C:\\Users\\stewa\\go\\src\\github.com\\bottos-project\\core\\vm\\testcase\\test_data2\\contract.wasm")


	par := make([]exec.Param, 2)
	par[0] = exec.Param{Ptype: "int", Pval: "20"}
	par[1] = exec.Param{Ptype: "int", Pval: "30"}

	p := exec.Args{Params: par}
	engine.Execute(nil , "123" , "add" , p)
}
*/

/*
func TestContract3() {
	engine := exec.NewVMEngine(nil, "product")
	//test
	code, err := ioutil.ReadFile("C:\\Users\\stewa\\go\\src\\github.com\\bottos-project\\core\\vm\\testcase\\test_data2\\contract.wasm")
	if err != nil {
		fmt.Println("error in read file", err.Error())
		return
	}

	par := make([]exec.Param, 2)
	par[0] = exec.Param{Ptype: "string", Pval: "hello "}
	par[1] = exec.Param{Ptype: "string", Pval: "world!"}

	p := exec.Args{Params: par}
	jbytes, err := json.Marshal(p)
	if err != nil {
		fmt.Println(err)
		fmt.Println(err.Error())
	}
	fmt.Println(string(jbytes))

	bf := bytes.NewBufferString("concat")
	bf.WriteString("|")
	bf.Write(jbytes)

	fmt.Printf("input is %v\n", bf.Bytes())

	res, err := engine.Call(nil, code, bf.Bytes())
	if err != nil {
		fmt.Println("call error!", err.Error())
	}
	fmt.Printf("res:%v\n", res)

	retbytes, err := engine.GetVM().GetPtrMemory(uint64(binary.LittleEndian.Uint32(res)))
	if err != nil {
		fmt.Println(err)
		fmt.Println("errors:" + err.Error())
	}

	fmt.Println("retbytes is " + string(retbytes))

	result := &exec.Result{}
	json.Unmarshal(retbytes, result)

	if result.Pval != "hello world!" {
		fmt.Println("the res should be 'hello world!'")
	}

}

func TestContract4() {
	engine := exec.NewVMEngine(nil, "product")
	//test
	code, err := ioutil.ReadFile("C:\\Users\\stewa\\go\\src\\github.com\\bottos-project\\core\\vm\\testcase\\test_data2\\contract.wasm")
	if err != nil {
		fmt.Println("error in read file", err.Error())
		return
	}

	par := make([]exec.Param, 2)
	par[0] = exec.Param{Ptype: "int_array", Pval: "1,2,3,4,5,6"}
	par[1] = exec.Param{Ptype: "int_array", Pval: "7,8,9,10"}

	p := exec.Args{Params: par}
	jbytes, err := json.Marshal(p)
	if err != nil {
		fmt.Println(err)
		fmt.Println(err.Error())
	}
	fmt.Println(string(jbytes))

	bf := bytes.NewBufferString("sumArray")
	bf.WriteString("|")
	bf.Write(jbytes)

	fmt.Printf("input is %v\n", bf.Bytes())

	res, err := engine.Call(nil, code, bf.Bytes())
	if err != nil {
		fmt.Println("call error!", err.Error())
	}
	fmt.Printf("res:%v\n", res)

	retbytes, err := engine.GetVM().GetPtrMemory(uint64(binary.LittleEndian.Uint32(res)))
	if err != nil {
		fmt.Println(err)
		fmt.Println("errors:" + err.Error())
	}

	fmt.Println("retbytes is " + string(retbytes))

	result := &exec.Result{}
	json.Unmarshal(retbytes, result)

	if result.Pval != "55" {
		fmt.Println("the res should be '55'")
	}

}

func TestRawContract() {
	engine := exec.NewVMEngine(nil, "product")
	//test
	code, err := ioutil.ReadFile("C:\\Users\\stewa\\go\\src\\github.com\\bottos-project\\core\\vm\\testcase\\test_data2\\rawcontract.wasm")
	if err != nil {
		fmt.Println("error in read file", err.Error())
		return
	}
	bf := bytes.NewBufferString("add")
	bf.WriteString("|")

	tmp := make([]byte, 8)
	binary.LittleEndian.PutUint32(tmp[:4], uint32(10))
	binary.LittleEndian.PutUint32(tmp[4:], uint32(20))
	bf.Write(tmp)

	fmt.Printf("input is %v\n", bf.Bytes())

	res, err := engine.Call(nil, code, bf.Bytes())
	if err != nil {
		fmt.Println("call error!", err.Error())
	}
	fmt.Printf("res:%v\n", res)

	retbytes, err := engine.GetVM().GetPtrMemory(uint64(binary.LittleEndian.Uint32(res)))
	if err != nil {
		fmt.Println(err)
		fmt.Println("errors:" + err.Error())
	}

	fmt.Println("retbytes is " + string(retbytes))

	result := &exec.Result{}
	json.Unmarshal(retbytes, result)

	if result.Pval != "30" {
		fmt.Println("the res should be '30'")
	}

}

func TestRawContract4() {
	engine := exec.NewVMEngine(nil, "product")
	//test
	code, err := ioutil.ReadFile("./test_data2/rawcontract2.wasm")
	if err != nil {
		fmt.Println("error in read file", err.Error())
		return
	}
	bf := bytes.NewBufferString("add")
	bf.WriteString("|")

	tmp := make([]byte, 8)
	binary.LittleEndian.PutUint32(tmp[:4], uint32(10))
	binary.LittleEndian.PutUint32(tmp[4:], uint32(20))
	bf.Write(tmp)

	fmt.Printf("input is %v\n", bf.Bytes())

	res, err := engine.Call(nil, code, bf.Bytes())
	if err != nil {
		fmt.Println("call error!", err.Error())
	}
	fmt.Printf("res:%v\n", res)

	retbytes, err := engine.GetVM().GetPtrMemory(uint64(binary.LittleEndian.Uint32(res)))
	if err != nil {
		fmt.Println(err)
		fmt.Println("errors:" + err.Error())
	}

	fmt.Println("retbytes is " + string(retbytes))

	result := &exec.Result{}
	json.Unmarshal(retbytes, result)

	if result.Pval != "30" {
		fmt.Println("the res should be '30'")
	}

}
*/

/*
func NewAc1() *exec.Apply_context{

	par := make([]exec.ParamInfo, 2)
	par[0] = exec.ParamInfo{Type: "int", Val: "20"}
	par[1] = exec.ParamInfo{Type: "int", Val: "30"}

	p := exec.ParamList{Params: par}

	jbytes, err := json.Marshal(p)
	if err != nil {
		fmt.Println(err.Error())
	}


	auth := exec.Authorization{
		Accout: "iwojima",
	}

	msg := exec.Message{
		Wasm_name: "CTX_WASM_FILE",
		Method_name: "add",
		Auth: auth,
		Method_param: jbytes, //这里将参数的byte[]传入
	}

	Ac := &exec.Apply_context{
		Msg: msg,
	}

	return Ac
}

func Test1() {
	ac := NewAc1()
	fmt.Println("ac = ",ac)

	exec.GetInstance().Apply(ac , 0 , true)
}

func Test2() {
	par := make([]exec.ParamInfo, 2)
	par[0] = exec.ParamInfo{Type: "int_array", Val: "1,2,3,4,5,6"}
	par[1] = exec.ParamInfo{Type: "int_array", Val: "7,8,9,10"}

	p := exec.ParamList{Params: par}

	jbytes, err := json.Marshal(p)
	if err != nil {
		fmt.Println(err.Error())
	}


	auth := exec.Authorization{
		Accout: "iwojima",
	}

	msg := exec.Message{
		Wasm_name: "CTX_WASM_FILE",
		Method_name: "sumArray",
		Auth: auth,
		Method_param: jbytes, //这里将参数的byte[]传入
	}

	ac := &exec.Apply_context{
		Msg: msg,
	}

	exec.GetInstance().Apply(ac , 0 , true)
}

func Test3() {
	par := make([]exec.ParamInfo, 2)
	par[0] = exec.ParamInfo{Type: "string", Val: "hello "}
	par[1] = exec.ParamInfo{Type: "string", Val: "world!"}

	p := exec.ParamList{Params: par}

	jbytes, err := json.Marshal(p)
	if err != nil {
		fmt.Println(err.Error())
	}


	auth := exec.Authorization{
		Accout: "iwojima",
	}

	msg := exec.Message {
		Wasm_name: "iwojima_wasm",
		Method_name: "concat",
		Auth: auth,
		Method_param: jbytes, //这里将参数的byte[]传入
	}

	ac := &exec.Apply_context{
		Msg: msg,
	}

	exec.GetInstance().Apply(ac , 0 , true)

	//vm.wasm.exec.GetInstance()
}

func TestTX() {
	par := make([]exec.ParamInfo, 2)
	par[0] = exec.ParamInfo{Type: "string", Val: "hello "}
	par[1] = exec.ParamInfo{Type: "string", Val: "world!"}

	p := exec.ParamList{Params: par}

	jbytes, err := json.Marshal(p)
	if err != nil {
		fmt.Println(err.Error())
	}


	auth := exec.Authorization{
		Accout: "iwojima",
	}

	msg := exec.Message{
		Wasm_name: "iwojima_wasm",
		Method_name: "_start",
		Auth: auth,
		Method_param: jbytes, //这里将参数的byte[]传入
	}

	ac := &exec.Apply_context{
		Msg: msg,
	}

	exec.GetInstance().Apply(ac , 0 , true)
}

func Test4() {
	par := make([]exec.ParamInfo, 1)
	//par[0] = exec.ParamInfo{Type: "int", Val: "1"}
	par[0] = exec.ParamInfo{Type: "string", Val: "Kamikaze"}

	p := exec.ParamList{Params: par}

	jbytes, err := json.Marshal(p)
	if err != nil {
		fmt.Println(err.Error())
	}


	auth := exec.Authorization{
		Accout: "iwojima",
	}

	msg := exec.Message{
		Wasm_name: "iwojima_wasm",
		Method_name: "apply",
		Auth: auth,
		Method_param: jbytes, //这里将参数的byte[]传入
	}

	ac := &exec.Apply_context{
		Msg: msg,
	}

	exec.GetInstance().Apply2 (ac , 0 , true)
}
*/

func init() {
	defer log.Flush()
	logger, err := log.LoggerFromConfigAsFile("C:\\Users\\stewa\\go\\src\\github.com\\bottos-project\\config\\log.xml")
	if err != nil {
		log.Critical("err parsing config log file", err)
		os.Exit(1)
		return
	}
	log.ReplaceLogger(logger)
}

func Test5() {

	type transferparam struct {
		To			string
		Amount		uint32
	}

	param := transferparam {
		To     : "stewart",
		Amount : 1233,
	}

	bf , err :=  msgpack.Marshal(param)
	fmt.Println(" Test5() bf = ",bf," , err = ",err)

	trx := &types.Transaction{
		Version:1,
		CursorNum:1,
		CursorLabel:1,
		Lifetime:1,
		Sender:"bottos",
		Contract:"usermng",
		Method:  "reguser",
		Param: bf,
		SigAlg:1,
		Signature:[]byte{},
	}
	ctx := &contract.Context{ Trx:trx}
	//go exec.GetInstance().StartHandler()

	res , err := exec.GetInstance().Start(ctx, 1, false)
	if err != nil {
		fmt.Println("*ERROR* fail to execute start !!!")
		return
	}

	fmt.Println("*SUCCESS* res = ",res)

	time.Sleep(time.Second * 3)

	//exec.GetInstance().StopHandler()

	//select {}
	for {}
}

func Test6() {

	type transferparam struct {
		To			string
		Amount		uint32
	}

	param := transferparam {
		To     : "stewart",
		Amount : 1233,
	}

	bf , err :=  msgpack.Marshal(param)
	fmt.Println(" Test6() bf = ",bf," , err = ",err)

	/*
	var tf transferparam
	msgpack.Unmarshal(bf , &tf)
	fmt.Println(" Test6() tf = ",tf)
	*/

	trx := &types.Transaction{
		Version:1,
		CursorNum:1,
		CursorLabel:1,
		Lifetime:1,
		Sender:"bottos",
		Contract: "usermng",
		Method:  "r",
		Param: bf,
		SigAlg:1,
		Signature:[]byte{},
	}

	ctx := &contract.Context{ Trx:trx}


	res , err := exec.GetInstance().Start(ctx, 1, false)
	if err != nil {
		fmt.Println("*ERROR* fail to execute start !!!")
		fmt.Println("err = ",err)
		return
	}

	fmt.Println("handled_trx = ",res)
	var tf transferparam
	for _ , sub_trx := range res {
		//var tf transferparam
		msgpack.Unmarshal(sub_trx.Param , &tf)
		fmt.Println("Test6 sub_trx = ",sub_trx.Param ," , tf = ",tf)
	}

	fmt.Println("=========================== *SUCCESS* res = ",res, " , err = ",err)

	/*
	handled_trx , err := exec.GetInstance().Start(ctx, 1, false)
	if err != nil {
		fmt.Println("*ERROR* fail to execute start !!!")
		fmt.Println("err = ",err)
		return
	}

	fmt.Println("handled_trx = ",handled_trx)

	for _ , sub_trx  := range handled_trx {
		fmt.Println("sub_trx = ",sub_trx.Contract) //reflect.TypeOf(sub_trx))
	}


	fmt.Println("*SUCCESS* res = ",res, " , err = ",err)
	*/
}

func Test7() {
	//
	type transferparam struct {
		To			string
		Amount		uint32
	}

	type teststrcut struct {
		valueA uint32
		valueB uint32
	}

	param := transferparam {
		To     : "stewart",
		Amount : 1233,
	}

	_ , err :=  msgpack.Marshal(param)
	//fmt.Println(" Test7() bf = ",bf," , err = ",err)

	//var p      string = "dc0004da00087465737466726f6dda000b646174616465616c6d6e67da000344544fcf0000000000000064"
	var p string      = "12345"
	var data []byte   = []byte(p)

	//fmt.Println("data = ",data)

	trx := &types.Transaction{
		Version:     1,
		CursorNum:   1,
		CursorLabel: 1,
		Lifetime:    1,
		Sender:      "bottos",
		Contract:    "usermng",
		Method:      "test_method",
		Param:       data,
		SigAlg:      1,
		Signature:   []byte{},
	}

	ctx := &contract.Context{Trx:trx}

	res , err := exec.GetInstance().Start(ctx, 1, false)
	if err != nil {
		fmt.Println("*ERROR* fail to execute start !!!")
		fmt.Println("err = ",err)
		return
	}

	//fmt.Println("err: ",err," ,handled_trx: ",res)
	var tf transferparam
	//var tt teststrcut
	for _ , sub_trx := range res {
		//var tf transferparam
		msgpack.Unmarshal(sub_trx.Param , &tf)
		fmt.Println("Test7 sub_trx = ",sub_trx.Param ," , tf = ",tf)
	}
	fmt.Println("end of testcase")


	time.Sleep(time.Second * 3)

	vmi := exec.GetInstance().GetWasteVM()
	if vmi == nil {
		fmt.Println("777")
		return
	}

	fmt.Println("After exec.GetInstance().Start(): ",reflect.TypeOf(vmi),",vmi: ",vmi)
	//vmi
	return
}

func Test8() {
	//
}
/*
func TestStart() {
	wasm_engine = &exec.WASM_ENGINE{
		vm_map: make(map[string]*VM),
	}

	module, err := wasm.ReadModule(bytes.NewBuffer(wasm_code), importer)
	if err != nil {
		fmt.Println("*ERROR* Failed to parse the wasm module !!! " + err.Error())
		return
	}

	if module.Export == nil  {
		fmt.Println("*ERROR* Failed to find export method from wasm module !!!")
		return nil
	}
	//fmt.Println("<<<<<<<<<<<<<<<< NewWASM")
	vm , err := NewVM(module)
	if err != nil {
		return nil
	}

	code, err := ioutil.ReadFile("C:\\Users\\stewa\\Desktop\\t.wasm")
	if err != nil {
		fmt.Println("error in read file", err.Error())
		return
	}

	method2 := "_start"
	input2 := make([]byte, 9)
	input2[0] = byte(len(method2))
	copy(input2[1:len(method2)+1], []byte(method2))

	input2[len(method2)+1] = byte(2) //param count
	input2[len(method2)+2] = byte(1) //param1 length
	input2[len(method2)+3] = byte(1) //param2 length
	input2[len(method2)+4] = byte(5) //param1
	input2[len(method2)+5] = byte(9) //param2

	fmt.Println(input2)

	res2, err := wasm_engine.Call(nil, code, input2)
	if err != nil {
		fmt.Println("call error!", err.Error())
	}

	methodName, err := exec.getCallMethodName(input)
	if err != nil {
		return nil, err
	}

	fmt.Printf("res:%v\n", res2)
	if binary.LittleEndian.Uint32(res2) != uint32(14) {
		t.Error("the result should be 14")
	}

}

*/

func main() {
	/*
	defer log.Flush()
	logger, err := log.LoggerFromConfigAsFile("C:\\Users\\stewa\\go\\src\\github.com\\bottos-project\\core\\log.xml")
	if err != nil {
		log.Critical("err parsing config log file", err)
		os.Exit(1)
		return
	}
	log.ReplaceLogger(logger)
	*/
	/*
	dbInst := db.NewDbService(config.Param.DataDir, filepath.Join(config.Param.DataDir, "blockchain"))
	if dbInst == nil {
		fmt.Println("Create DB service fail")
		os.Exit(1)
	}
	*/
	//role.Init(dbInst)

	//TestContract1()
	fmt.Println("==================================================")
	//TestContract2()
	//TestTX()
	//TestAdd()
	Test7()
	fmt.Println("==================================================")
	//TestContract3()
	//fmt.Println("==================================================")
	//TestContract4()
	//fmt.Println("==================================================")
	//fmt.Println("helloworld")
	//exec.Query_contract()
	//fmt.Println("Query_contract = ",exec.Query_contract())
}
