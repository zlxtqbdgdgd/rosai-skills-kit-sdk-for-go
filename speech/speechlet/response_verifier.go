package speechlet

// Verifier for validating the ResponseEnvelope received from the Speechlets.
type ResponseVerifier interface {
	// Verifies a ResponseEnvelope within the context of the link Session in
	// which it was received. Returns true if the verify succeeded, false otherwise.
	// responseEnvelope: ResponseEnvelope to verify
	// session: Session context within which to verify the call
	// return true if the verify succeeded, false otherwise
	Verify(re *ResponseEnvelope, session *Session) bool
}
