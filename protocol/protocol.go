package protocol

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

type Command byte

const (
	CmdSet Command = iota
	CmdGet
)

type CommandSet struct {
	Key   []byte
	Value []byte
	TTL   int32
}

type CommandGet struct {
	Key []byte
}

func (c *CommandSet) Bytes() []byte {
	buf := new(bytes.Buffer)

	binary.Write(buf, binary.LittleEndian, CmdSet)

	binary.Write(buf, binary.LittleEndian, int32(len(c.Key)))
	binary.Write(buf, binary.LittleEndian, c.Key)

	binary.Write(buf, binary.LittleEndian, int32(len(c.Value)))
	binary.Write(buf, binary.LittleEndian, c.Value)

	binary.Write(buf, binary.LittleEndian, c.TTL)

	return buf.Bytes()
}

func (c *CommandGet) Bytes() []byte {
	buf := new(bytes.Buffer)

	binary.Write(buf, binary.LittleEndian, CmdGet)
	binary.Write(buf, binary.LittleEndian, int32(len(c.Key)))
	binary.Write(buf, binary.LittleEndian, c.Key)

	return buf.Bytes()
}

func ParseCommand(r io.Reader) (any, error) {
	var cmd Command

	err := binary.Read(r, binary.LittleEndian, &cmd)
	if err != nil {
		return nil, err
	}

	switch cmd {
	case CmdSet:
		return parseSet(r)
	case CmdGet:
		return parseGet(r)
	default:
		return nil, fmt.Errorf("Unknown Command")
	}
}

func parseSet(r io.Reader) (*CommandSet, error) {
	c := &CommandSet{}

	var keyLen int32
	binary.Read(r, binary.LittleEndian, &keyLen)

	c.Key = make([]byte, keyLen)
	io.ReadFull(r, c.Key)

	var valLen int32
	binary.Read(r, binary.LittleEndian, valLen)

	c.Value = make([]byte, valLen)
	io.ReadFull(r, c.Value)

	binary.Read(r, binary.LittleEndian, &c.TTL)

	return c, nil
}

func parseGet(r io.Reader) (*CommandGet, error) {
	c := &CommandGet{}

	var keyLen int32
	binary.Read(r, binary.LittleEndian, keyLen)

	c.Key = make([]byte, keyLen)
	io.ReadFull(r, c.Key)

	return c, nil
}
