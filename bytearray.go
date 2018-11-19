package network

import (
	"errors"
	"math"
)

const (
	PACKET_MAX_LEN = 1024
)

type ByteArray struct {
	pos  int
	data []byte
}

func (p *ByteArray) Data() []byte {
	return p.data
}

func (p *ByteArray) Length() int {
	return len(p.data)
}

func CreateByteArray() *ByteArray {
	return &ByteArray{data: make([]byte, 0, PACKET_MAX_LEN)}
}

//=============================================== Readers
func (p *ByteArray) ReadBool() (ret bool, err error) {
	b, _err := p.ReadByte()

	if b != byte(1) {
		return false, _err
	}

	return true, _err
}

func (p *ByteArray) ReadByte() (ret byte, err error) {
	if p.pos >= len(p.data) {
		err = errors.New("read byte failed")
		return
	}

	ret = p.data[p.pos]
	p.pos++
	return
}

func (p *ByteArray) ReadBytes() (ret []byte, err error) {
	if p.pos+2 > len(p.data) {
		err = errors.New("read bytes header failed")
		return
	}
	size, _ := p.ReadU16()
	if p.pos+int(size) > len(p.data) {
		err = errors.New("read bytes data failed")
		return
	}

	ret = p.data[p.pos : p.pos+int(size)]
	p.pos += int(size)
	return
}

func (p *ByteArray) ReadString() (ret string, err error) {
	if p.pos+2 > len(p.data) {
		err = errors.New("read string header failed")
		return
	}

	size, _ := p.ReadU16()
	if p.pos+int(size) > len(p.data) {
		err = errors.New("read string data failed")
		return
	}

	bytes := p.data[p.pos : p.pos+int(size)]
	p.pos += int(size)
	ret = string(bytes)
	return
}

func (p *ByteArray) ReadU16() (ret uint16, err error) {
	if p.pos+2 > len(p.data) {
		err = errors.New("read uint16 failed")
		return
	}

	buf := p.data[p.pos : p.pos+2]
	ret = uint16(buf[0])<<8 | uint16(buf[1])
	p.pos += 2
	return
}

func (p *ByteArray) ReadS16() (ret int16, err error) {
	_ret, _err := p.ReadU16()
	ret = int16(_ret)
	err = _err
	return
}

func (p *ByteArray) ReadU24() (ret uint32, err error) {
	if p.pos+3 > len(p.data) {
		err = errors.New("read uint24 failed")
		return
	}

	buf := p.data[p.pos : p.pos+3]
	ret = uint32(buf[0])<<16 | uint32(buf[1])<<8 | uint32(buf[2])
	p.pos += 3
	return
}

func (p *ByteArray) ReadS24() (ret int32, err error) {
	_ret, _err := p.ReadU24()
	ret = int32(_ret)
	err = _err
	return
}

func (p *ByteArray) ReadU32() (ret uint32, err error) {
	if p.pos+4 > len(p.data) {
		err = errors.New("read uint32 failed")
		return
	}

	buf := p.data[p.pos : p.pos+4]
	ret = uint32(buf[0])<<24 | uint32(buf[1])<<16 | uint32(buf[2])<<8 | uint32(buf[3])
	p.pos += 4
	return
}

func (p *ByteArray) ReadS32() (ret int32, err error) {
	_ret, _err := p.ReadU32()
	ret = int32(_ret)
	err = _err
	return
}

func (p *ByteArray) ReadU64() (ret uint64, err error) {
	if p.pos+8 > len(p.data) {
		err = errors.New("read uint64 failed")
		return
	}

	ret = 0
	buf := p.data[p.pos : p.pos+8]
	for i, v := range buf {
		ret |= uint64(v) << uint((7-i)*8)
	}
	p.pos += 8
	return
}

func (p *ByteArray) ReadS64() (ret int64, err error) {
	_ret, _err := p.ReadU64()
	ret = int64(_ret)
	err = _err
	return
}

func (p *ByteArray) ReadFloat32() (ret float32, err error) {
	bits, _err := p.ReadU32()
	if _err != nil {
		return float32(0), _err
	}

	ret = math.Float32frombits(bits)
	if math.IsNaN(float64(ret)) || math.IsInf(float64(ret), 0) {
		return 0, nil
	}

	return ret, nil
}

func (p *ByteArray) ReadFloat64() (ret float64, err error) {
	bits, _err := p.ReadU64()
	if _err != nil {
		return float64(0), _err
	}

	ret = math.Float64frombits(bits)
	if math.IsNaN(ret) || math.IsInf(ret, 0) {
		return 0, nil
	}

	return ret, nil
}

//================================================ Writers
func (p *ByteArray) WriteZeros(n int) {
	for i := 0; i < n; i++ {
		p.data = append(p.data, byte(0))
	}
}

func (p *ByteArray) WriteBool(v bool) {
	if v {
		p.data = append(p.data, byte(1))
	} else {
		p.data = append(p.data, byte(0))
	}
}

func (p *ByteArray) WriteByte(v byte) {
	p.data = append(p.data, v)
}

func (p *ByteArray) WriteBytes(v []byte) {
	p.WriteU16(uint16(len(v)))
	p.data = append(p.data, v...)
}

func (p *ByteArray) WriteRawBytes(v []byte) {
	p.data = append(p.data, v...)
}

func (p *ByteArray) WriteString(v string) {
	bytes := []byte(v)
	p.WriteU16(uint16(len(bytes)))
	p.data = append(p.data, bytes...)
}

func (p *ByteArray) WriteU16(v uint16) {
	p.data = append(p.data, byte(v>>8), byte(v))
}

func (p *ByteArray) WriteS16(v int16) {
	p.WriteU16(uint16(v))
}

func (p *ByteArray) WriteU24(v uint32) {
	p.data = append(p.data, byte(v>>16), byte(v>>8), byte(v))
}

func (p *ByteArray) WriteU32(v uint32) {
	p.data = append(p.data, byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
}

func (p *ByteArray) WriteS32(v int32) {
	p.WriteU32(uint32(v))
}

func (p *ByteArray) WriteU64(v uint64) {
	p.data = append(p.data, byte(v>>56), byte(v>>48), byte(v>>40), byte(v>>32), byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
}

func (p *ByteArray) WriteS64(v int64) {
	p.WriteU64(uint64(v))
}

func (p *ByteArray) WriteFloat32(f float32) {
	v := math.Float32bits(f)
	p.WriteU32(v)
}

func (p *ByteArray) WriteFloat64(f float64) {
	v := math.Float64bits(f)
	p.WriteU64(v)
}
