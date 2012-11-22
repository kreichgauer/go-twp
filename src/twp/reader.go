package twp

import (
    "bufio"
    "encoding/binary"
    "errors"
    "fmt"
    "io"
)

var (
    ErrReservedTag = errors.New("Reserved tag read")
    ErrInvalidTag = errors.New("Invalid tag read")
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

func (rd *Reader) ReadSequence(seq *[]Value) (err error) {
    var tag Tag
    var v Value
    for {
        if err := binary.Read(rd, binary.BigEndian, &tag); err != nil {
            return err
        }
        fmt.Printf("Read tag %d\n", tag)
        switch {
        case EndOfContent == tag:
            return nil
        case NoValue == tag:
            *seq = append(*seq, nil)
        case Struct == tag:
            panic("struct not implemented.")
        case Sequence == tag:
            panic("sequence not implemented.")
        case MessageOrUnion <= tag && tag <= MessageOrUnionEnd:
            // id := int(tag - 4)
            // val, err = rd.ReadMessage(id)
            panic("union not implemented.")
        case RegisteredExtension == tag:
            panic("extension not implemneted.")
        case ShortInteger == tag:
            if v, err = rd.ReadShortInt(); err != nil {
                return err
            }
            *seq = append(*seq, v)
        case LongInteger == tag:
            if v, err = rd.ReadLongInt(); err != nil {
                return err
            }
            *seq = append(*seq, v)
        case ShortBinary == tag:
            if v, err = rd.ReadShortBinary(); err != nil {
                return err
            }
            *seq = append(*seq, v)
        case LongBinary == tag:
             if v, err = rd.ReadLongBinary(); err != nil {
                return err
             }
             *seq = append(*seq, v)
        case ShortString <= tag && tag < LongString:
            length := int(tag - 17)
            if v, err = rd.ReadShortString(length); err != nil {
               return err
            }
            *seq = append(*seq, v)
        case LongString == tag:
            if v, err = rd.ReadLongString(); err != nil {
               return err
            }
            *seq = append(*seq, v)
        case Reserved <= tag && tag <= ReservedEnd:
            return ErrReservedTag
        default:
            return ErrInvalidTag
        }
    }
    return nil
}

func (rd *Reader) ReadMessage() (val *RawMessage, err error) {
    var tag uint8
    if err = binary.Read(rd, binary.BigEndian, &tag); err != nil {
        return nil, err
    }
    id := tag - 4
    fmt.Printf("Message ID %d\n", id)
    val = new(RawMessage)
    val.Id = id
    if err = rd.ReadSequence(&val.Fields); err != nil {
        return nil, err
    }
    return val, nil
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
