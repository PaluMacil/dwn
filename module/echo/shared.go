package echo

import (
	"github.com/PaluMacil/dwn/webserver/handler"
	"github.com/gorilla/mux"
)

var history = make(History, 0, 102)

func NewEchoHandler() EchoModule {
	return EchoModule{&history}
}

func NewHistoryHandler() HistoryModule {
	return HistoryModule{&history}
}

func RegisterRoutes(rt *mux.Router, factory handler.Factory) {
	historyPath := "/s/echo-history/"
	if factory.Config.Prod {
		historyPath = "/"
	}
	rt.Handle(historyPath, factory.Handler(NewHistoryHandler().historyHandler))
}
