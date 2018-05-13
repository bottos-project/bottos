# msgpack

msgpack use a subset of MessagePack protocol, which support types:

```
uint8
uint16
uint32
uint64
str
bin
array
```

# encode

```
func Marshal(v interface{}) ([]byte, error)
```

Sample code:

```
type TestSubStruct struct{
    V1 string
    V2 uint32
}

type TestStruct struct{
    V1 string
    V2 uint32
    V3 TestSubStruct
}

ts := TestStruct {
    V1: "testuser",
    V2: 99,
    V3: TestSubStruct{V1:"123", V2:3},
}
b, err := Marshal(ts)

// BytesToHex(b)
// dc0003da00087465737475736572ce00000063dc0002da0003313233ce00000003
```

# decode

```
func Unmarshal(data []byte, dst interface{}) error
```

Sample code:

```
type TestSubStruct struct{
    V1 string
    V2 uint32
}

type TestStruct struct{
    V1 string
    V2 uint32
    V3 TestSubStruct
}

b, err := HexToBytes("dc0003da00087465737475736572ce00000063dc0002da0003313233ce00000003")

ts := TestStruct {}
err = Unmarshal(b, &ts)
fmt.Println("ts", ts)
// ts  {testuser 99 {123 3}}
```
