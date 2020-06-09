package r

// UseQuit returns a single function that when invoked
// will exit the application
func UseQuit() func() {
	return func() {
		close(quitChan)
	}
}
