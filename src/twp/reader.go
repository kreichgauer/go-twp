package twp

import (
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
    rd io.Reader
}

func NewReader(rd io.Reader) (*Reader) {
    return &Reader{rd}
}

func (r *Reader) ReadSequence(seq *[]Value) (err error) {
    var tag Tag
    var v Value
    for {
        if err := binary.Read(r.rd, binary.BigEndian, &tag); err != nil {
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
            if v, err = r.ReadShortInt(); err != nil {
                return err
            }
            *seq = append(*seq, v)
        case LongInteger == tag:
            if v, err = r.ReadLongInt(); err != nil {
                return err
            }
            *seq = append(*seq, v)
        case ShortBinary == tag:
            if v, err = r.ReadShortBinary(); err != nil {
                return err
            }
            *seq = append(*seq, v)
        case LongBinary == tag:
             if v, err = r.ReadLongBinary(); err != nil {
                return err
             }
             *seq = append(*seq, v)
        case ShortString <= tag && tag < LongString:
            length := int(tag - 17)
            if v, err = r.ReadShortString(length); err != nil {
               return err
            }
            *seq = append(*seq, v)
        case LongString == tag:
            if v, err = r.ReadLongString(); err != nil {
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

func (r *Reader) ReadMessage() (val *RawMessage, err error) {
    var tag uint8
    if err = binary.Read(r.rd, binary.BigEndian, &tag); err != nil {
        return nil, err
    }
    id := tag - 4
    fmt.Printf("Message ID %d\n", id)
    val = new(RawMessage)
    val.Id = id
    if err = r.ReadSequence(&val.Fields); err != nil {
        return nil, err
    }
    return val, nil
}

func (r *Reader) ReadShortInt() (val uint8, err error) {
    err = binary.Read(r.rd, binary.BigEndian, &val)
    return val, err
}

func (r *Reader) ReadLongInt() (val uint32, err error) {
    err = binary.Read(r.rd, binary.BigEndian, &val)
    return val, err
}

func (r *Reader) ReadShortBinary() (val []byte, err error) {
    var l uint8
    if err = binary.Read(r.rd, binary.BigEndian, &l); err != nil {
        return nil, err
    }
    val = make([]byte, l)
    if _, err = io.ReadFull(r.rd, val); err != nil {
        return nil, err
    }
    return val, nil
}

func (r *Reader) ReadLongBinary() (val []byte, err error) {
    var l uint32
    if err = binary.Read(r.rd, binary.BigEndian, &l); err != nil {
        return nil, err
    }
    val = make([]byte, l)
    if _, err = io.ReadFull(r.rd, val); err != nil {
        return nil, err
    }
    return val, nil
}

func (r *Reader) ReadShortString(length int) (val string, err error) {
    if !(0 <= length && length <= 109) {
        return "", fmt.Errorf("Invalid short string length %d", length)
    }
    buf := make([]byte, length)
    if _, err = io.ReadFull(r.rd, buf); err != nil {
        return "", err
    }
    return string(buf), nil
}

func (r *Reader) ReadLongString() (val string, err error) {
    var length uint32
    if err = binary.Read(r.rd, binary.BigEndian, &length); err != nil {
        return "", err
    }
    buf := make([]byte, length)
    if _, err = io.ReadFull(r.rd, buf); err != nil {
        return "", err
    }
    return string(buf), nil
}
