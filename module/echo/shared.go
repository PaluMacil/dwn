package echo

var history = make(History, 0, 102)

func NewEchoHandler() EchoModule {
	return EchoModule{&history}
}

func NewHistoryHandler() HistoryModule {
	return HistoryModule{&history}
}
