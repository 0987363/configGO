package daemon

type Daemon struct {
	CloseHook []func()
}

var daemon *Daemon
var daemonClose chan struct{} = make(chan struct{})

func init() {
	daemon = &Daemon{}
}

func AddCloseHook(f func()) {
	daemon.CloseHook = append(daemon.CloseHook, f)
}

func Exit() {
	for _, f := range daemon.CloseHook {
		f()
	}
	close(daemonClose)
}

func Wait() {
	<-daemonClose
}
