package ziface
// conn + data
// conn+data encapsulate to Request
type IRequest interface {
	// conn which has established connection with client-side
	GetConnection() IConnection

	// request's Message Data
	GetData() []byte

	// obtain ID of request Message
	GetMsgID() uint32
}


