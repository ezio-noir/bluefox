package runner

type Done <-chan struct{}

func Run(funcToRun func()) Done {
	done := make(chan struct{})
	go func() {
		defer close(done)
		funcToRun()
	}()
	return done
}

func WaitAll(dones ...Done) {
	for _, done := range dones {
		<-done
	}
}
