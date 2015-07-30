package kurento

import "fmt"

type IWebRtcEndpoint interface {
	GatherCandidates() error
	AddIceCandidate(candidate IceCandidate) error
}

// WebRtcEndpoint interface. This type of "Endpoint" offers media streaming using
// WebRTC.
type WebRtcEndpoint struct {
	BaseRtpEndpoint

	// Address of the STUN server (Only IP address are supported)
	StunServerAddress string

	// Port of the STUN server
	StunServerPort int

	// TURN server URL with this format:
	// 'user:password@address:port(?transport=[udp|tcp|tls])'.
	// 'address' must be an IP (not a domain).
	// 'transport' is optional (UDP by default).
	TurnUrl string
}

// Return contructor params to be called by "Create".
func (elem *WebRtcEndpoint) getConstructorParams(from IMediaObject, options map[string]interface{}) map[string]interface{} {

	// Create basic constructor params
	ret := map[string]interface{}{
		"mediaPipeline":   fmt.Sprintf("%s", from),
		"useDataChannels": false,
	}

	// then merge options
	mergeOptions(ret, options)

	return ret

}

// Init the gathering of ICE candidates.
// It must be called after SdpEndpoint::generateOffer or SdpEndpoint::processOffer
func (elem *WebRtcEndpoint) GatherCandidates() error {
	req := elem.getInvokeRequest()

	reqparams := map[string]interface{}{
		"operation": "gatherCandidates",
		"object":    elem.Id,
	}
	if elem.connection.SessionId != "" {
		reqparams["sessionId"] = elem.connection.SessionId
	}
	req["params"] = reqparams

	// Call server and wait response
	response := <-elem.connection.Request(req)

	// Returns error or nil
	if response.Error != nil {
		return response.Error
	} else {
		return nil
	}

}

// Provide a remote ICE candidate
func (elem *WebRtcEndpoint) AddIceCandidate(candidate IceCandidate) error {
	req := elem.getInvokeRequest()

	params := make(map[string]interface{})

	setIfNotEmpty(params, "candidate", candidate)

	reqparams := map[string]interface{}{
		"operation":       "addIceCandidate",
		"object":          elem.Id,
		"operationParams": params,
	}
	if elem.connection.SessionId != "" {
		reqparams["sessionId"] = elem.connection.SessionId
	}
	req["params"] = reqparams

	// Call server and wait response
	response := <-elem.connection.Request(req)

	// Returns error or nil
	if response.Error != nil {
		return response.Error
	} else {
		return nil
	}

}
