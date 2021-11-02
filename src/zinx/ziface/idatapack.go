package ziface

/*
Wrap & Unwrap

Data stream for TCP Connection
*/
type IDataPack interface {
	// obtain length of package head
	GetHeadLen() uint32
	// pack
	Pack(msg IMessage) ([]byte, error)

	// unpack
	Unpack([]byte) (IMessage, error)
}
