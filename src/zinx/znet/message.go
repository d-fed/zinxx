package znet

/*
encapsulate Message to TLV
1. load const length of head
2. read Messaage content from conn


Message Length,
MessageID
Message Content

*/
type Message struct {
	Id      uint32 // message ID
	DataLen uint32 // message length
	Data    []byte // message content

}

// Create a Message Package
func NewMsgPackage(id uint32, data []byte) *Message {
	return &Message{
		Id:      id,
		DataLen: uint32(len(data)),
		Data:    data,
	}
}

// obtain message ID
func (m *Message) GetMsgId() uint32 {
	return m.Id
}

// obtain message length
func (m *Message) GetDataLen() uint32 {
	return m.DataLen
}

// obtain message content
func (m *Message) GetData() []byte {
	return m.Data
}

// set message ID
func (m *Message) SetMsgId(id uint32) {
	m.Id = id
}

// set message content
func (m *Message) SetData(data []byte) {
	m.Data = data
}

// set message length
func (m *Message) SetDataLen(len uint32) {
	m.DataLen = len
}
