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

func (d *Decoder) DecodeSequence(seq Sequence) (err error) {
    var tag Tag
    var v Value
    for i := 0; ; i++ {
        if err := binary.Read(d.rd, binary.BigEndian, &tag); err != nil {
            return err
        }
        fmt.Printf("Read tag %d\n", tag)
        switch {
        case EndOfContentTag == tag:
            return nil
        case NoValueTag == tag:
            seq.SetValue(i, nil)
        case StructTag == tag:
            panic("struct not implemented.")
        case SequenceTag == tag:
            panic("sequence not implemented.")
        case MessageOrUnionTag <= tag && tag <= MessageOrUnionEndTag:
            // id := int(tag - 4)
            // val, err = rd.ReadMessage(id)
            panic("union not implemented.")
        case RegisteredExtensionTag == tag:
            panic("extension not implemneted.")
        case ShortIntegerTag == tag:
            if v, err = d.DecodeShortInt(); err != nil {
                return err
            }
            seq.SetValue(i, v)
        case LongIntegerTag == tag:
            if v, err = d.DecodeLongInt(); err != nil {
                return err
            }
            seq.SetValue(i, v)
        case ShortBinaryTag == tag:
            if v, err = d.DecodeShortBinary(); err != nil {
                return err
            }
            seq.SetValue(i, v)
        case LongBinaryTag == tag:
             if v, err = d.DecodeLongBinary(); err != nil {
                return err
             }
             seq.SetValue(i, v)
        case ShortStringTag <= tag && tag < LongStringTag:
            length := int(tag - 17)
            if v, err = d.DecodeShortString(length); err != nil {
               return err
            }
            seq.SetValue(i, v)
        case LongStringTag == tag:
            if v, err = d.DecodeLongString(); err != nil {
               return err
            }
            seq.SetValue(i, v)
        case ReservedTag <= tag && tag <= ReservedEndTag:
            return ErrReservedTag
        default:
            return ErrInvalidTag
        }
    }
    return nil
}

func (d *Decoder) DecodeMessage() (msg *Message, err error) {
    var tag uint8
    if err = binary.Read(d.rd, binary.BigEndian, &tag); err != nil {
        return nil, err
    }
    id := tag - 4
    fmt.Printf("Message ID %d\n", id)
    msg = new(Message)
    msg.Id = id
    if err = d.DecodeSequence(msg); err != nil {
        return nil, err
    }
    return msg, nil
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
