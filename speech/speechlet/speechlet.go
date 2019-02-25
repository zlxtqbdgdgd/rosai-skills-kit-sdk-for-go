package speechlet

// A Speechlet is a speech-enabled web service that runs in the cloud.
// A Speechlet receives and responds to speech initiated requests.
// The methods in the Speechlet interface define:
// * life-cycle events for the skill as experienced by the user
// * a way to service speech requests from the user
// * call-backs from events happening on the device the user is interacting with

// Because a Speechlet is a cloud-based service, the life-cycle of the actual object
// implementing the Speechlet interface is unrelated to the life-cycle of the Alexa skill
/// as experienced by the user interacting with the speech enabled device.
// A single Speechlet object typically handles requests from many users, from many devices,
// for many sessions, using many concurrent threads. Also, multiple requests
// that are part of a single session coming from a single user and a single device
// may actually be serviced by different Speechlet objects if your service is deployed
// on multiple machines and load-balanced. The Session object is what defines a single
// run of a skill as experienced by a user.

type Speechlet interface {
	// Used to notify that a new session started as a result of a user interacting with the device.
	// This method enables Speechlets to perform initialization logic and allows for session
	// attributes to be stored for subsequent requests.
	// requestEnvelope: the session started request envelope
	OnSessionStarted(requestEnvelope *RequestEnvelope) error

	// without providing an Intent.<br>
	// This method is only invoked when {@link Session#isNew()} is true.
	// requestEnvelope: the launch request envelope
	// return the response, spoken and visual, to the request
	OnLaunch(requestEnvelope *RequestEnvelope) (*Response, error)

	// Entry point for handling speech initiated requests.
	// This is where the bulk of the Speechlet logic lives. Intent requests are handled by
	// this method and return responses to render to the user.
	// If this is the initial request of a new Speechlet session, Session#isNew
	// returns true. Otherwise, this is a subsequent request within an existing session.
	// requestEnvelope: the intent request envelope to handle
	// return the response, spoken and visual, to the request
	OnIntent(requestEnvelope *RequestEnvelope) (*Response, *Context, error)

	// Callback used to notify that the session ended as a result of the user interacting,
	// or not interacting with the device. This method is not invoked if the
	// itself ended the session using Response#setNullableShouldEndSession(bool).
	// requestEnvelope: the end of session request envelope
	OnSessionEnded(requestEnvelope *RequestEnvelope) error
}
