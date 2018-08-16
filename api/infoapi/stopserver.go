package infoapi

import "os"

func (rt *InfoRoute) handleStopServer() {
	p, err := os.FindProcess(os.Getpid())
	if err != nil {
		rt.API().ServeInternalServerError(err)
		return
	}
	p.Signal(os.Interrupt)
}
