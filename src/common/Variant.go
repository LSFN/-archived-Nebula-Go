package common

type Variant struct {
	data []byte
}

func (v *Variant) ConnectByte(byte b) {
	data = append(data, b)
}

func (v *Variant) IsComplete() bool {
	return v.data[-1]&0x80 == 0
}

func (v *Variant) ToUint64() uint64 {
	if data == nil {
		return 0
	}
	var value uint64
	var multiplier uint8
	for i := range data {
		value += (data[i] & 0x7f) << multiplier
		multiplier += 7
	}
}

func (v *Variant) FromUint64(value uint64) {
	for value > 0 {
		var nextByte byte = value & 0x7f
		value >>= 7
		if value > 0 {
			nextByte |= 0x80
		}
		Variant.data = append(Variant.data, nextByte)
	}
}

func (v *Variant) Reset() {
	v.data = nil
}

func (v *Variant) ToBytes() {
	result := make([]byte, len(data))
	copy(result, data)
	return result
}
