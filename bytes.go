package rock

func (B *Bytes) makeW() {
	B.p.w.c = make(chan []byte, B.Len)
	go postIfClient(B.p.w.c, Tbytes, B.Name)
}

func (B *Bytes) makeR() {
	B.p.r.c = make(chan []byte, B.Len)
	go getIfClient(B.p.r.c, Tbytes, B.Name)
}

func (B *Bytes) makeN() {
	B.p.n.c = make(chan int)
}

func (B *Bytes) SelSend(b []byte) chan<- interface{} {
	send := make(chan interface{})
	go func() { B.To(b); <-send }()
	return send
}

func (B *Bytes) SelRecv() <-chan []byte {
	recv := make(chan []byte)
	go func() { recv <- B.From() }()
	return recv
}

func (B *Bytes) To(b []byte) {
	go started.Do(getAndOrPostIfServer)

	bytesDict.Lock()
	if bytesDict.m == nil {
		bytesDict.m = map[string]*Bytes{}
	}
	if _, found := bytesDict.m[B.Name]; !found {
		bytesDict.m[B.Name] = B
	}
	bytesDict.Unlock()

	B.p.w.Do(B.makeW)
	if IsClient {
		B.p.w.c <- b
		return
	}

	B.p.n.Do(B.makeN)
	for {
		<-B.p.n.c
		B.p.w.c <- b
		if len(B.p.n.c) == 0 {
			break
		}
	}
}

func (B *Bytes) From() []byte {
	go started.Do(getAndOrPostIfServer)

	bytesDict.Lock()
	if bytesDict.m == nil {
		bytesDict.m = map[string]*Bytes{}
	}
	if _, found := bytesDict.m[B.Name]; !found {
		bytesDict.m[B.Name] = B
	}
	bytesDict.Unlock()

	B.p.r.Do(B.makeR)
	return <-B.p.r.c
}
