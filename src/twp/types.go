package twp

import (
    "bytes"
    "fmt"
    "reflect"
)

type Value interface {}

type Sequence interface {
    GetValue(i int) (Value, bool)
    SetValue(i int, v Value) (error)
}

type Message struct {
    Id uint8
    Fields []Value
}

func (msg *Message) GetValue(i int) (Value, bool) {
    if i >= len(msg.Fields) {
        return nil, false
    }
    return msg.Fields[i], true
}

func (msg *Message) SetValue(i int, v Value) (error) {
    if i >= len(msg.Fields) {
        // Reslice or grow msg.Fields
        if i < cap(msg.Fields) {
            msg.Fields = msg.Fields[:i+1]
        } else {
            fields := make([]Value, i+1, 2*(i+1))
            copy(fields, msg.Fields)
            msg.Fields = fields
        }
    }
    msg.Fields[i] = v
    return nil
}

func (msg *Message) String() (string) {
    var b bytes.Buffer
    b.WriteString(fmt.Sprintf("Message %d\n", msg.Id))
    for _, val := range msg.Fields {
        b.WriteString(fmt.Sprintf("%s: %v\n", reflect.TypeOf(val).Name(), val))
    }
    return b.String()
}
