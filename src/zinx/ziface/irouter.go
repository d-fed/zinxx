package ziface

/*
	router abstract interface
	rotuter's data is IRequest
*/

type IRouter interface {
	//  Hook before conn
	PreHandle(request IRequest)
	//  Hook during conn
	Handle(request IRequest)
	// Hook after conn
	PostHandle(request IRequest)
}
