package sms

var tokenSubscribers = make(map[string]TokenCallback)

func SubscribeToTokenEvents(clientID string, callback TokenCallback) {
	tokenSubscribers[clientID] = callback
}

func UnsubscribeFromTokenEvents(clientID string) {
	delete(tokenSubscribers, clientID)
}

func NotifyTokenEvent(clientID, token string) {
	if callback, ok := tokenSubscribers[clientID]; ok {
		callback(clientID, token)
	}
}
