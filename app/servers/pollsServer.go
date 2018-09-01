package servers

type PollsServer struct {
	*DefaultServer
}

func NewPollsServer() *PollsServer {
	defaultSrv := NewDefaultServer()
	defaultSrv.name = "pollsServer"
	return &PollsServer{defaultSrv}
}
