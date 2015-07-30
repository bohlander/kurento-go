package kurento

import "fmt"

type IDispatcherOneToMany interface {
	SetSource(source HubPort) error
	RemoveSource() error
}

// A `Hub` that sends a given source to all the connected sinks
type DispatcherOneToMany struct {
	Hub
}

// Return contructor params to be called by "Create".
func (elem *DispatcherOneToMany) getConstructorParams(from IMediaObject, options map[string]interface{}) map[string]interface{} {

	// Create basic constructor params
	ret := map[string]interface{}{
		"mediaPipeline": fmt.Sprintf("%s", from),
	}

	// then merge options
	mergeOptions(ret, options)

	return ret

}

// Sets the source port that will be connected to the sinks of every `HubPort` of
// the dispatcher
func (elem *DispatcherOneToMany) SetSource(source HubPort) error {
	req := elem.getInvokeRequest()

	params := make(map[string]interface{})

	setIfNotEmpty(params, "source", source)

	reqparams := map[string]interface{}{
		"operation":       "setSource",
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

// Remove the source port and stop the media pipeline.
func (elem *DispatcherOneToMany) RemoveSource() error {
	req := elem.getInvokeRequest()

	reqparams := map[string]interface{}{
		"operation": "removeSource",
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
