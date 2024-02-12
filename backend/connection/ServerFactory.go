package connection

type Server interface {
	newServer() Server
	StartConnection(string)
}

func CreateServer() Server {
	return createNativeWSS()
}

func createNativeWSS() Server {
	var wsServer NativeServer

	return wsServer.newServer()
}