package ziface

/*
	encasuplate Message
*/
type IMessage interface {
	// obtain message ID
	GetMsgId() uint32
	// obtain message length
	GetMsgLen() uint32
	// obtain message content
	GetData() []byte

	// set message ID
	SetMsgId(uint32)
	// set message content
	SetData([]byte)
	// set message length
	SetDataLen(uint32)
}
