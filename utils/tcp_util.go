package utils

func Decode(bytes []byte) (code uint16,data []byte) {
	code = uint16(bytes[0]) << 8 | uint16(bytes[1])
	data = bytes[2:]
	return
}

func Encode(code uint16, data []byte) (bytes []byte) {
	length := len(data)+4
	bytes = make([]byte,length)
	bytes = append(bytes,byte(length >> 8),byte(length))
	bytes = append(bytes,byte(code >> 8),byte(code))
	bytes = append(bytes,data...)
	return
}

