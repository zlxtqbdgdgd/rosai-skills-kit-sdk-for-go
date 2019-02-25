package speechlet

// Verifier for validating the received SpeechletRequestEnvelopes from the devices.
type RequestVerifier interface {
	// Verifies a Request within the context of the Session in which it was
	// received. Returns true if the verify succeeded, false otherwise.
	// RequestEnvelope: the speechlet request envelope to verify
	// return true if the verify succeeded, false otherwise
	Verify(re *RequestEnvelope) bool
}

type AppIdRequestVerifier struct {
	AppIds []string
}

func (verifier *AppIdRequestVerifier) Verify(re *RequestEnvelope) bool {
	return true
}
