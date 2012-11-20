package twp

import (
    "bufio"
    "bytes"
    "encoding/binary"
    "errors"
    "fmt"
    "io"
)

var (
    ErrReservedTag = errors.New("Reserved tag read")
)

type Reader struct {
    bufio.Reader
}

func NewReader(rd io.Reader) (*Reader) {
    return &Reader{*bufio.NewReader(rd)}
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
    var tag Tag
    if err = binary.Read(rd, binary.BigEndian, &tag); err != nil {
        return "", err
    }
    fmt.Printf("Read tag %d\n", tag)
    switch {
    case EndOfContent == tag:
        val = ""
    case NoValue == tag:
        val = ""
    case Struct == tag:
        panic("struct not implemented.")
    case Sequence == tag:
        panic("sequence not implemented.")
    case MessageOrUnion <= tag && tag <= MessageOrUnionEnd:
        id := int(tag - 4)
        val, err = rd.ReadMessage(id)
    case RegisteredExtension == tag:
        panic("Extension not implemnted")
    case ShortInteger == tag:
        var v uint8
        v, err = rd.ReadShortInt()
        val = fmt.Sprintf("%d", v)
    case LongInteger == tag:
        var v uint32
        v, err = rd.ReadLongInt()
        val = fmt.Sprintf("%d", v)
    case ShortBinary == tag:
        var v []byte
        v, err = rd.ReadShortBinary()
        val = fmt.Sprintf("%s", v)
    case LongBinary == tag:
        var v []byte
         v, err = rd.ReadLongBinary()
         val = fmt.Sprintf("%s", v)
    case ShortString <= tag && tag < LongString:
        length := int(tag - 17)
        val, err = rd.ReadShortString(length)
    case LongString == tag:
        val, err = rd.ReadLongString()
    case Reserved <= tag && tag <= ReservedEnd:
        err = ErrReservedTag
    }
    return val, err
}

func (rd *Reader) ReadMessage(tag int) (val string, err error) {
    fmt.Printf("Message ID %d\n", tag)
    var buffer bytes.Buffer
    for {
        v, err := rd.ReadValue()
        if err != nil {
            return "", err
        }
        if v == "" {
            fmt.Println("Message End")
            break
        }
        buffer.WriteString(v)
    }
    return buffer.String(), nil
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
