package twp

import (
    "bytes"
    "encoding/binary"
    "io"
)

type Writer struct {
    io.Writer
}

const twpMagic = "TWP3\n"

func NewWriter(wr io.Writer) (*Writer) {
    return &Writer{wr}
}

func (wr *Writer) InitWithProtocol(id int) (err error) {
    if err = wr.WriteMagic(); err != nil {
        return err
    }
    if err = wr.WriteProtocolId(id); err != nil {
        return err
    }
    return nil
}

func (wr *Writer) WriteMagic() (err error) {
    _, err = wr.Write([]byte(twpMagic))
    return err
}

func (wr *Writer) WriteProtocolId(id int) (err error) {
    return wr.WriteInteger(id)
}

func (wr *Writer) WriteInteger(val int) (err error) {
    return binary.Write(wr, binary.BigEndian, []byte{13, byte(val)})
}

func (wr *Writer) WriteString(val string) (err error) {
    var buf bytes.Buffer
    if length := len(val); length <= 109 {
        buf.WriteByte(byte(ShortString + length))
    } else {
        buf.WriteByte(LongString)
        binary.Write(&buf, binary.BigEndian, length)
    }
    buf.WriteString(val)
    _, err = wr.Write(buf.Bytes())
    return err
}
