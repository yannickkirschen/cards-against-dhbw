package communication

// Sender is an interface that is used to send messages to a Receiver.
type Sender interface {
	// Send sends a message to the Receiver.
	Send(title string, message ...any)
}

// Receiver is an interface that is used to receive messages from a Sender.
type Receiver interface {
	// Receive is called when a message is received.
	Receive(player string, title string, message any)
}
