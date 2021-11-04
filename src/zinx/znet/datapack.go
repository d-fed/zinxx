package znet

import (
	"awesomeProject/src/zinx/utils"
	"awesomeProject/src/zinx/ziface"
	"bytes"
	"encoding/binary"
	"errors"
)

//Pack, Unpack for Type-Length-Value
// TLV-encoded data stream contain code of the record type, followed by record value length, and finally the value itself

type DataPack struct{}

// instance for pack/unpack

func NewDataPack() *DataPack {
	return &DataPack{}
}

// obtain length of package head
func (dp *DataPack) GetHeadLen() uint32 {
	// Datalen uint32(4 bytes) + ID uint32(4 bytes)
	return 8
}

// pack, Message -> TLV (Type-length-value)
// datalen | msgID | data
func (dp *DataPack) Pack(msg ziface.IMessage) ([]byte, error) {
	// create a bytes stream buffer
	dataBuff := bytes.NewBuffer([]byte{})
	// !!!! Sequence Matters

	// Write dataLen into dataBuff
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetDataLen()); err != nil {
		return nil, err
	}

	// Write MessageID into dataBuff
	//if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetData()); err != nil {
	//	return nil, err
	//}

	// Write data into dataBuff
	return dataBuff.Bytes(), nil
	// 返回该二进制文件

}

// unpack
// Extract head
// read again based on data length in head
func (dp *DataPack) Unpack(binaryData []byte) (ziface.IMessage, error) {
	// create a ioReader from binary data
	dataBuff := bytes.NewReader(binaryData)

	// unencrypt head, to obtain datalen and MsgID
	msg := &Message{}

	// read dataLen
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}

	// read MsgID
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}

	// decide whether dataLen has exceed max packet length

	if utils.GlobalObject.MaxPacketSize > 0 && msg.DataLen > utils.GlobalObject.MaxPacketSize {
		return nil, errors.New("too Large msg data recv!")
	}

	return msg, nil

}
