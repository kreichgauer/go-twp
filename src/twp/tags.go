package twp

type Tag uint8

const (
    EndOfContentTag        = 0
    NoValueTag             = 1
    StructTag              = 2
    SequenceTag            = 3
    MessageOrUnionTag      = 4
    MessageOrUnionEndTag   = 11
    RegisteredExtensionTag = 12
    ShortIntegerTag        = 13
    LongIntegerTag         = 14
    ShortBinaryTag         = 15
    LongBinaryTag          = 16
    ShortStringTag         = 17
    LongStringTag          = 127
    ReservedTag            = 128
    ReservedEndTag         = 159
    ApplicationTypeTag     = 160
    ApplicationTypeEndTag  = 255
)
