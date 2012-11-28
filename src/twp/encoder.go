package twp

import (
    "bytes"
    "encoding/binary"
    "io"
)

type Encoder struct {
    wr io.Writer
}

const twpMagic = "TWP3\n"

func NewEncoder(wr io.Writer) (*Encoder) {
    return &Encoder{wr}
}

func (en *Encoder) InitWithProtocol(id int) (err error) {
    if err = en.EncodeMagic(); err != nil {
        return err
    }
    if err = en.EncodeProtocolId(id); err != nil {
        return err
    }
    return nil
}

// FIXME Remove
func (en *Encoder) Write(buf []byte) (n int, err error) {
    return en.wr.Write(buf)
}

func (en *Encoder) EncodeMagic() (err error) {
    _, err = en.wr.Write([]byte(twpMagic))
    return err
}

func (en *Encoder) EncodeProtocolId(id int) (err error) {
    return en.EncodeInteger(id)
}

func (en *Encoder) EncodeInteger(val int) (err error) {
    return binary.Write(en.wr, binary.BigEndian, []byte{13, byte(val)})
}

func (en *Encoder) EncodeString(val string) (err error) {
    var buf bytes.Buffer
    if length := len(val); length <= 109 {
        buf.WriteByte(byte(ShortStringTag + length))
    } else {
        buf.WriteByte(LongStringTag)
        binary.Write(&buf, binary.BigEndian, length)
    }
    buf.WriteString(val)
    _, err = en.wr.Write(buf.Bytes())
    return err
}
