package vm

type VmType byte

const (
	VmTypeUnkonw VmType = 0
	VmTypeWasm VmType = 1
	VmTypeJS VmType = 2
)