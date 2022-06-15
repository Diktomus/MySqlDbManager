package utils

type ScanBuffer struct {
	Pointers []interface{}
	Values   []interface{}
}

func (sb *ScanBuffer) Prepare(bufferSize int) {
	sb.Pointers = make([]interface{}, bufferSize)
	sb.Values = make([]interface{}, bufferSize)

	for i := 0; i < len(sb.Values); i++ {
		sb.Pointers[i] = &sb.Values[i]
	}
}
