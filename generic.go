package corto

import "unsafe"

func GenericGet(data interface{}, i int) interface{} {
	switch data := data.(type) {
	case []int8:
		return data[i:]
	case []uint8:
		return data[i:]
	case []int16:
		return data[i:]
	case []uint16:
		return data[i:]
	case []int32:
		return data[i:]
	case []uint32:
		return data[i:]
	case []float32:
		return data[i:]
	case []float64:
		return data[i:]
	}
	return 0
}

func GenericResize(data interface{}, si int) interface{} {
	switch cdata := data.(type) {
	case []int8:
		new := make([]int8, si)
		copy(new, cdata)
		return new
	case []uint8:
		new := make([]uint8, si)
		copy(new, cdata)
		return new
	case []int16:
		new := make([]int16, si)
		copy(new, cdata)
		return new
	case []uint16:
		new := make([]uint16, si)
		copy(new, cdata)
		return new
	case []int32:
		new := make([]int32, si)
		copy(new, cdata)
		return new
	case []uint32:
		new := make([]uint32, si)
		copy(new, cdata)
		return new
	case []float32:
		new := make([]float32, si)
		copy(new, cdata)
		return new
	case []float64:
		new := make([]float64, si)
		copy(new, cdata)
		return new
	}
	return 0
}

func GenericSize(data interface{}) int {
	switch data := data.(type) {
	case []int8:
		return len(data)
	case []uint8:
		return len(data)
	case []int16:
		return len(data)
	case []uint16:
		return len(data)
	case []int32:
		return len(data)
	case []uint32:
		return len(data)
	case []float32:
		return len(data)
	case []float64:
		return len(data)
	}
	return 0
}

func GenericGetInt(data interface{}, i int) int {
	switch data := data.(type) {
	case []int8:
		return int(data[i])
	case []uint8:
		return int(data[i])
	case []int16:
		return int(data[i])
	case []uint16:
		return int(data[i])
	case []int32:
		return int(data[i])
	case []uint32:
		return int(data[i])
	case []float32:
		return int(data[i])
	case []float64:
		return int(data[i])
	}
	return 0
}

func GenericSet(data interface{}, i int, v interface{}) {
	switch data := data.(type) {
	case []int8:
		data[i] = v.(int8)
		break
	case []uint8:
		data[i] = v.(uint8)
		break
	case []int16:
		data[i] = v.(int16)
		break
	case []uint16:
		data[i] = v.(uint16)
		break
	case []int32:
		data[i] = v.(int32)
		break
	case []uint32:
		data[i] = v.(uint32)
		break
	case []float32:
		data[i] = v.(float32)
		break
	case []float64:
		data[i] = v.(float64)
		break
	}
	return
}

func GenericLess(a, b interface{}) bool {
	switch data := a.(type) {
	case int8:
		return data < b.(int8)
	case uint8:
		return data < b.(uint8)
	case int16:
		return data < b.(int16)
	case uint16:
		return data < b.(uint16)
	case int32:
		return data < b.(int32)
	case uint32:
		return data < b.(uint32)
	case float32:
		return data < b.(float32)
	case float64:
		return data < b.(float64)
	}
	return false
}

func GenericGreater(a, b interface{}) bool {
	switch data := a.(type) {
	case int8:
		return data > b.(int8)
	case uint8:
		return data > b.(uint8)
	case int16:
		return data > b.(int16)
	case uint16:
		return data > b.(uint16)
	case int32:
		return data > b.(int32)
	case uint32:
		return data > b.(uint32)
	case float32:
		return data > b.(float32)
	case float64:
		return data > b.(float64)
	}
	return false
}

func GenericInt(a interface{}) int {
	switch data := a.(type) {
	case int8:
		return int(data)
	case uint8:
		return int(data)
	case int16:
		return int(data)
	case uint16:
		return int(data)
	case int32:
		return int(data)
	case uint32:
		return int(data)
	case float32:
		return int(data)
	case float64:
		return int(data)
	}
	return 0
}

func GenericGetPtr(data interface{}) uintptr {
	switch data := data.(type) {
	case []int8:
		return uintptr(unsafe.Pointer(&data[0]))
	case []uint8:
		return uintptr(unsafe.Pointer(&data[0]))
	case []int16:
		return uintptr(unsafe.Pointer(&data[0]))
	case []uint16:
		return uintptr(unsafe.Pointer(&data[0]))
	case []int32:
		return uintptr(unsafe.Pointer(&data[0]))
	case []uint32:
		return uintptr(unsafe.Pointer(&data[0]))
	case []float32:
		return uintptr(unsafe.Pointer(&data[0]))
	case []float64:
		return uintptr(unsafe.Pointer(&data[0]))
	}
	return 0
}
