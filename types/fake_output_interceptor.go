package types

type FakeOutputInterceptor struct {
	DidStartInterceptingOutput bool
	DidStopInterceptingOutput  bool
	InterceptedOutput          string
}

func (interceptor *FakeOutputInterceptor) StartInterceptingOutput() error {
	interceptor.DidStartInterceptingOutput = true
	return nil
}

func (interceptor *FakeOutputInterceptor) StopInterceptingAndReturnOutput() (string, error) {
	interceptor.DidStopInterceptingOutput = true
	return interceptor.InterceptedOutput, nil
}
