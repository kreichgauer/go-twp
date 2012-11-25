package twp

import (
    "bytes"
    "encoding/binary"
    "strings"
    "testing"
)

func TestReadShortInt(t *testing.T) {
    var (
        val uint8
        err error
    )
    buf := bytes.NewReader([]byte{42})
    d := NewDecoder(buf)
    if val, err = d.DecodeShortInt(); err != nil {
        t.Errorf("%s", err)
    }
    if val != 42 {
        t.Errorf("Expected %d to be 42\n", val)
    }
}

func TestReadLongInt(t *testing.T) {
    var (
        val uint32
        err error
    )
    buf := bytes.NewReader([]byte{1, 0, 0, 0})
    d := NewDecoder(buf)
    if val, err = d.DecodeLongInt(); err != nil {
        t.Errorf("%s", err)
    }
    if val != 16777216 {
        t.Errorf("Expected %d to be 16777216\n", val)
    }
}

func TestReadShortBinary(t *testing.T) {
    buf := []byte{5, 1, 2, 3, 4, 5}
    src := bytes.NewReader(buf)
    d := NewDecoder(src)
    dst, err := d.DecodeShortBinary()
    if err != nil {
        t.Fatalf("Error: %s", err)
    }
    if !bytes.Equal(buf[1:], dst) {
        t.Fatalf("Expected %x to be %x", dst, buf)
    }
}


func TestReadLongBinary(t *testing.T) {
    // 4 bytes length, 2**24 bytes content
    buf := make([]byte, 4 + 1<<24)
    // set kength to 2**4
    copy(buf[:4], []byte{1, 0, 0, 0})
    // fill data with 42s
    data := buf[4:]
    for i, _ := range data {
        data[i] = 42
    }
    src := bytes.NewReader(buf[:])
    d := NewDecoder(src)
    dst, err := d.DecodeLongBinary()
    if err != nil {
        t.Fatalf("Error: %s", err)
    }
    if !bytes.Equal(data, dst) {
        t.Fatalf("Buffers don't match")
    }
}

func TestReadShortString(t *testing.T) {
    length := 109
    str := strings.Repeat("a", length)
    buf := bytes.NewBufferString(str)
    d := NewDecoder(buf)
    s, err := d.DecodeShortString(length)
    if err != nil {
        pError(err, t)
    }
    verify(s, str, t)
}

func TestReadLongString(t *testing.T) {
    var length uint32 = 1 << 24
    buf := make([]byte, 4 + length)
    binary.BigEndian.PutUint32(buf[:4], length)
    str := strings.Repeat("a", int(length))
    data := buf[4:]
    copy(data, []byte(str))
    d := NewDecoder(bytes.NewReader(buf))
    s, err := d.DecodeLongString()
    if err != nil {
        pError(err, t)
    }
    verify(s, str, t)
}
