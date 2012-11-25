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

type Decoder struct {
    rd io.Reader
}

func NewDecoder(rd io.Reader) (*Decoder) {
    return &Decoder{rd}
}

func (d *Decoder) DecodeSequence(seq *[]Value) (err error) {
    var tag Tag
    var v Value
    for {
        if err := binary.Read(d.rd, binary.BigEndian, &tag); err != nil {
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
            if v, err = d.DecodeShortInt(); err != nil {
                return err
            }
            *seq = append(*seq, v)
        case LongInteger == tag:
            if v, err = d.DecodeLongInt(); err != nil {
                return err
            }
            *seq = append(*seq, v)
        case ShortBinary == tag:
            if v, err = d.DecodeShortBinary(); err != nil {
                return err
            }
            *seq = append(*seq, v)
        case LongBinary == tag:
             if v, err = d.DecodeLongBinary(); err != nil {
                return err
             }
             *seq = append(*seq, v)
        case ShortString <= tag && tag < LongString:
            length := int(tag - 17)
            if v, err = d.DecodeShortString(length); err != nil {
               return err
            }
            *seq = append(*seq, v)
        case LongString == tag:
            if v, err = d.DecodeLongString(); err != nil {
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

func (d *Decoder) DecodeMessage() (val *RawMessage, err error) {
    var tag uint8
    if err = binary.Read(d.rd, binary.BigEndian, &tag); err != nil {
        return nil, err
    }
    id := tag - 4
    fmt.Printf("Message ID %d\n", id)
    val = new(RawMessage)
    val.Id = id
    if err = d.DecodeSequence(&val.Fields); err != nil {
        return nil, err
    }
    return val, nil
}

func (d *Decoder) DecodeShortInt() (val uint8, err error) {
    err = binary.Read(d.rd, binary.BigEndian, &val)
    return val, err
}

func (d *Decoder) DecodeLongInt() (val uint32, err error) {
    err = binary.Read(d.rd, binary.BigEndian, &val)
    return val, err
}

func (d *Decoder) DecodeShortBinary() (val []byte, err error) {
    var l uint8
    if err = binary.Read(d.rd, binary.BigEndian, &l); err != nil {
        return nil, err
    }
    val = make([]byte, l)
    if _, err = io.ReadFull(d.rd, val); err != nil {
        return nil, err
    }
    return val, nil
}

func (d *Decoder) DecodeLongBinary() (val []byte, err error) {
    var l uint32
    if err = binary.Read(d.rd, binary.BigEndian, &l); err != nil {
        return nil, err
    }
    val = make([]byte, l)
    if _, err = io.ReadFull(d.rd, val); err != nil {
        return nil, err
    }
    return val, nil
}

func (d *Decoder) DecodeShortString(length int) (val string, err error) {
    if !(0 <= length && length <= 109) {
        return "", fmt.Errorf("Invalid short string length %d", length)
    }
    buf := make([]byte, length)
    if _, err = io.ReadFull(d.rd, buf); err != nil {
        return "", err
    }
    return string(buf), nil
}

func (d *Decoder) DecodeLongString() (val string, err error) {
    var length uint32
    if err = binary.Read(d.rd, binary.BigEndian, &length); err != nil {
        return "", err
    }
    buf := make([]byte, length)
    if _, err = io.ReadFull(d.rd, buf); err != nil {
        return "", err
    }
    return string(buf), nil
}
