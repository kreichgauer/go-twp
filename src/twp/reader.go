package twp

import (
    "encoding/binary"
    "bufio"
    "fmt"
    "errors"
    "io"
)

var (
    ErrReservedTag = errors.New("Reserved tag read")
)

type Reader struct {
    *bufio.Reader
}

func NewReader(rd io.Reader) (*Reader) {
    return &Reader{bufio.NewReader(rd)}
}

func (rd *Reader) ReadFull(buf []byte) (err error) {
    var length, i, n int
    length = len(buf)
    for i < length {
        if n, err = rd.Read(buf[i:length]); err != nil {
            return err
        }
        i += n
    }
    return nil
}

func (rd *Reader) ReadValue() (val string, err error) {
    var tag byte
    if tag, err = rd.ReadByte(); err != nil {
        return "", err
    }
    fmt.Printf("Read tag %d\n", tag)
    switch {
    case tag == 0:
        panic("End-Of-Content not implemented.")
    case 1 == tag:
        return "", nil
    case 2 == tag:
        panic("struct not implemented.")
    case 3 == tag:
        panic("sequence not implemented.")
    case 4 <= tag && tag <= 11:
        panic("Message/Union not implemented")
    case 12 == tag:
        panic("Extension not implemnted")
    case 13 <= tag:
        var v uint8
        v, err = rd.ReadShortInt()
        val = fmt.Sprintf("%d", v)
    case 14 == tag:
        var v uint32
        v, err = rd.ReadLongInt()
        val = fmt.Sprintf("%d", v)
    case 15 == tag:
        var v []byte
        v, err = rd.ReadShortBinary()
        val = fmt.Sprintf("%s", v)
    case 16 == tag:
        var v []byte
         v, err = rd.ReadLongBinary()
         val = fmt.Sprintf("%s", v)
    case 17 <= tag && tag <= 126:
        length := int(tag - 17)
        val, err = rd.ReadShortString(length)
    case 127 == tag:
        // val rd.ReadLongString()
    case 128 <= tag && tag <= 159:
        err = ErrReservedTag

    }
    return val, err
}

func (rd *Reader) ReadShortInt() (val uint8, err error) {
    err = binary.Read(rd, binary.BigEndian, &val)
    return val, err
}

func (rd *Reader) ReadLongInt() (val uint32, err error) {
    err = binary.Read(rd, binary.BigEndian, &val)
    return val, err
}

func (rd *Reader) ReadShortBinary() (val []byte, err error) {
    var l uint8
    if err = binary.Read(rd, binary.BigEndian, &l); err != nil {
        return nil, err
    }
    val = make([]byte, l)
    if err = rd.ReadFull(val); err != nil {
        return nil, err
    }
    return val, nil
}

func (rd *Reader) ReadLongBinary() (val []byte, err error) {
    var l uint32
    if err = binary.Read(rd, binary.BigEndian, &l); err != nil {
        return nil, err
    }
    val = make([]byte, l)
    if err = rd.ReadFull(val); err != nil {
        return nil, err
    }
    return val, nil
}

func (rd *Reader) ReadShortString(length int) (val string, err error) {
    if !(0 <= length && length <= 109) {
        return "", fmt.Errorf("Invalid short string length %d", length)
    }
    buf := make([]byte, length)
    if err = rd.ReadFull(buf); err != nil {
        return "", err
    }
    return string(buf), nil
}

func (rd *Reader) ReadLongString() (val string, err error) {
    var length uint32
    if err = binary.Read(rd, binary.BigEndian, &length); err != nil {
        return "", err
    }
    buf := make([]byte, length)
    if err = rd.ReadFull(buf); err != nil {
        return "", err
    }
    return string(buf), nil
}
