package twp

import (
    "bytes"
    "encoding/binary"
    "strings"
    "testing"
)

func TestReadFull(t *testing.T) {
    buf := []byte{1, 2, 3, 4, 5, 6}
    src := bytes.NewReader(buf)
    r := NewReader(src)
    dst := make([]byte, 6)
    if err := r.ReadFull(dst); err != nil {
        t.Fatalf("Error: %s", err)
    }
    if !bytes.Equal(buf, dst) {
        t.Fatalf("Expected %x to be %x", dst, buf)
    }
}

func TestReadShortInt(t *testing.T) {
    var (
        val uint8
        err error
    )
    buf := bytes.NewReader([]byte{42})
    r := NewReader(buf)
    if val, err = r.ReadShortInt(); err != nil {
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
    r := NewReader(buf)
    if val, err = r.ReadLongInt(); err != nil {
        t.Errorf("%s", err)
    }
    if val != 16777216 {
        t.Errorf("Expected %d to be 16777216\n", val)
    }
}

func TestReadShortBinary(t *testing.T) {
    buf := []byte{5, 1, 2, 3, 4, 5}
    src := bytes.NewReader(buf)
    r := NewReader(src)
    dst, err := r.ReadShortBinary()
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
    r := NewReader(src)
    dst, err := r.ReadLongBinary()
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
    r := NewReader(buf)
    s, err := r.ReadShortString(length)
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
    r := NewReader(bytes.NewReader(buf))
    s, err := r.ReadLongString()
    if err != nil {
        pError(err, t)
    }
    verify(s, str, t)
}

func pError(err error, t *testing.T) {
    t.Errorf("Error: %s\n", err)
}

func pFatal(err error, t *testing.T) {
    t.Fatalf("Fatal: %s\n", err)
}

func verify(a, b interface{}, t *testing.T) {
    if a != b {
        t.Fatalf("Expected %s to be %s.\n", a, b)
    }
}
