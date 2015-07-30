package kurento

import "fmt"

// Base for all objects that can be created in the media server.
type MediaObject struct {
	connection *Connection

	// `MediaPipeline` to which this MediaObject belong, or the pipeline itself if
	// invoked over a `MediaPipeline`
	MediaPipeline IMediaPipeline

	// parent of this media object. The type of the parent depends on the type of the
	// element. The parent of a `MediaPad` is its `MediaElement`; the parent of a
	// `Hub` or a `MediaElement` is its `MediaPipeline`. A `MediaPipeline` has no
	// parent, i.e. the property is null
	Parent IMediaObject

	// unique identifier of the mediaobject.
	Id string

	// Childs of current object, all returned objects have parent set to current
	// object
	Childs []IMediaObject

	// Object name. This is just a comodity to simplify developers life debugging, it
	// is not used internally for indexing nor idenfiying the objects. By default is
	// the object type followed by the object id.
	Name string

	// This property allows activate/deactivate sending the element tags in all its
	// events.
	SendTagsInEvents bool

	// Number of seconds since Epoch when the element was created
	CreationTime int
}

// Return contructor params to be called by "Create".
func (elem *MediaObject) getConstructorParams(from IMediaObject, options map[string]interface{}) map[string]interface{} {
	return options

}

// Request a SessionSpec offer.
// This can be used to initiate a connection.
func (elem *MediaObject) AddTag(key string, value string) error {
	req := elem.getInvokeRequest()

	params := make(map[string]interface{})

	setIfNotEmpty(params, "key", key)
	setIfNotEmpty(params, "value", value)

	reqparams := map[string]interface{}{
		"operation":       "addTag",
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

// Remove the tag (key and value) associated to a tag
func (elem *MediaObject) RemoveTag(key string) error {
	req := elem.getInvokeRequest()

	params := make(map[string]interface{})

	setIfNotEmpty(params, "key", key)

	reqparams := map[string]interface{}{
		"operation":       "removeTag",
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

// Returns the value associated to the given key.
// Returns:
// // The value associated to the given key.
func (elem *MediaObject) GetTag(key string) (string, error) {
	req := elem.getInvokeRequest()

	params := make(map[string]interface{})

	setIfNotEmpty(params, "key", key)

	reqparams := map[string]interface{}{
		"operation":       "getTag",
		"object":          elem.Id,
		"operationParams": params,
	}
	if elem.connection.SessionId != "" {
		reqparams["sessionId"] = elem.connection.SessionId
	}
	req["params"] = reqparams

	// Call server and wait response
	response := <-elem.connection.Request(req)

	// // The value associated to the given key.

	if response.Error == nil {
		return response.Result["value"], nil
	} else {
		return response.Result["value"], response.Error
	}
	return response.Result["value"], response.Error

}

// Returns all the MediaObject tags.
// Returns:
// // An array containing all pairs key-value associated to the MediaObject.
func (elem *MediaObject) GetTags() ([]Tag, error) {
	req := elem.getInvokeRequest()

	reqparams := map[string]interface{}{
		"operation": "getTags",
		"object":    elem.Id,
	}
	if elem.connection.SessionId != "" {
		reqparams["sessionId"] = elem.connection.SessionId
	}
	req["params"] = reqparams

	// Call server and wait response
	response := <-elem.connection.Request(req)

	// // An array containing all pairs key-value associated to the MediaObject.

	ret := []Tag{}
	return ret, response.Error

}

type IServerManager interface {
	GetKmd(moduleName string) (string, error)
}

// This is a standalone object for managing the MediaServer
type ServerManager struct {
	MediaObject

	// Server information, version, modules, factories, etc
	Info *ServerInfo

	// All the pipelines available in the server
	Pipelines []IMediaPipeline

	// All active sessions in the server
	Sessions []string

	// Metadata stored in the server
	Metadata string
}

// Return contructor params to be called by "Create".
func (elem *ServerManager) getConstructorParams(from IMediaObject, options map[string]interface{}) map[string]interface{} {
	return options

}

// Returns the kmd associated to a module
// Returns:
// // The kmd file
func (elem *ServerManager) GetKmd(moduleName string) (string, error) {
	req := elem.getInvokeRequest()

	params := make(map[string]interface{})

	setIfNotEmpty(params, "moduleName", moduleName)

	reqparams := map[string]interface{}{
		"operation":       "getKmd",
		"object":          elem.Id,
		"operationParams": params,
	}
	if elem.connection.SessionId != "" {
		reqparams["sessionId"] = elem.connection.SessionId
	}
	req["params"] = reqparams

	// Call server and wait response
	response := <-elem.connection.Request(req)

	// // The kmd file

	if response.Error == nil {
		return response.Result["value"], nil
	} else {
		return response.Result["value"], response.Error
	}
	return response.Result["value"], response.Error

}

type ISessionEndpoint interface {
}

// Session based endpoint. A session is considered to be started when the media
// exchange starts. On the other hand, sessions terminate when a timeout,
// defined by the developer, takes place after the connection is lost.
type SessionEndpoint struct {
	Endpoint
}

// Return contructor params to be called by "Create".
func (elem *SessionEndpoint) getConstructorParams(from IMediaObject, options map[string]interface{}) map[string]interface{} {
	return options

}

type IHub interface {
}

// A Hub is a routing `MediaObject`. It connects several `endpoints <Endpoint>`
// together
type Hub struct {
	MediaObject
}

// Return contructor params to be called by "Create".
func (elem *Hub) getConstructorParams(from IMediaObject, options map[string]interface{}) map[string]interface{} {
	return options

}

type IFilter interface {
}

// Base interface for all filters. This is a certain type of `MediaElement`, that
// processes media injected through its sinks, and delivers the outcome through
// its sources.
type Filter struct {
	MediaElement
}

// Return contructor params to be called by "Create".
func (elem *Filter) getConstructorParams(from IMediaObject, options map[string]interface{}) map[string]interface{} {
	return options

}

type IEndpoint interface {
}

// Base interface for all end points. An Endpoint is a `MediaElement`
// that allow `KMS` to interchange media contents with external systems,
// supporting different transport protocols and mechanisms, such as `RTP`,
// `WebRTC`, `HTTP`, "file:/" URLs... An "Endpoint" may
// contain both sources and sinks for different media types, to provide
// bidirectional communication.
type Endpoint struct {
	MediaElement
}

// Return contructor params to be called by "Create".
func (elem *Endpoint) getConstructorParams(from IMediaObject, options map[string]interface{}) map[string]interface{} {
	return options

}

type IHubPort interface {
}

// This `MediaElement` specifies a connection with a `Hub`
type HubPort struct {
	MediaElement
}

// Return contructor params to be called by "Create".
func (elem *HubPort) getConstructorParams(from IMediaObject, options map[string]interface{}) map[string]interface{} {

	// Create basic constructor params
	ret := map[string]interface{}{
		"hub": fmt.Sprintf("%s", from),
	}

	// then merge options
	mergeOptions(ret, options)

	return ret

}

type IPassThrough interface {
}

// This `MediaElement` that just passes media through
type PassThrough struct {
	MediaElement
}

// Return contructor params to be called by "Create".
func (elem *PassThrough) getConstructorParams(from IMediaObject, options map[string]interface{}) map[string]interface{} {

	// Create basic constructor params
	ret := map[string]interface{}{
		"mediaPipeline": fmt.Sprintf("%s", from),
	}

	// then merge options
	mergeOptions(ret, options)

	return ret

}

type IUriEndpoint interface {
	Pause() error
	Stop() error
}

// Interface for endpoints the require a URI to work. An example of this, would be
// a `PlayerEndpoint` whose URI property could be used to locate a file to stream
type UriEndpoint struct {
	Endpoint

	// The uri for this endpoint.
	Uri string
}

// Return contructor params to be called by "Create".
func (elem *UriEndpoint) getConstructorParams(from IMediaObject, options map[string]interface{}) map[string]interface{} {
	return options

}

// Pauses the feed
func (elem *UriEndpoint) Pause() error {
	req := elem.getInvokeRequest()

	reqparams := map[string]interface{}{
		"operation": "pause",
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

// Stops the feed
func (elem *UriEndpoint) Stop() error {
	req := elem.getInvokeRequest()

	reqparams := map[string]interface{}{
		"operation": "stop",
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

type IMediaPipeline interface {
	GetGstreamerDot(details GstreamerDotDetails) (string, error)
}

// A pipeline is a container for a collection of `MediaElements<MediaElement>` and
// `MediaMixers<MediaMixer>`. It offers the methods needed to control the
// creation and connection of elements inside a certain pipeline.
type MediaPipeline struct {
	MediaObject
}

// Return contructor params to be called by "Create".
func (elem *MediaPipeline) getConstructorParams(from IMediaObject, options map[string]interface{}) map[string]interface{} {
	return options

}

// Returns a string in dot (graphviz) format that represents the gstreamer
// elements inside the pipeline
// Returns:
// // The dot graph
func (elem *MediaPipeline) GetGstreamerDot(details GstreamerDotDetails) (string, error) {
	req := elem.getInvokeRequest()

	params := make(map[string]interface{})

	setIfNotEmpty(params, "details", details)

	reqparams := map[string]interface{}{
		"operation":       "getGstreamerDot",
		"object":          elem.Id,
		"operationParams": params,
	}
	if elem.connection.SessionId != "" {
		reqparams["sessionId"] = elem.connection.SessionId
	}
	req["params"] = reqparams

	// Call server and wait response
	response := <-elem.connection.Request(req)

	// // The dot graph

	if response.Error == nil {
		return response.Result["value"], nil
	} else {
		return response.Result["value"], response.Error
	}
	return response.Result["value"], response.Error

}

type ISdpEndpoint interface {
	GenerateOffer() (string, error)
	ProcessOffer(offer string) (string, error)
	ProcessAnswer(answer string) (string, error)
	GetLocalSessionDescriptor() (string, error)
	GetRemoteSessionDescriptor() (string, error)
}

// Implements an SDP negotiation endpoint able to generate and process
// offers/responses and that configures resources according to
// negotiated Session Description
type SdpEndpoint struct {
	SessionEndpoint

	// Maximum video bandwidth for receiving.
	// Unit: kbps(kilobits per second).
	// 0: unlimited.
	// Default value: 500
	MaxVideoRecvBandwidth int
}

// Return contructor params to be called by "Create".
func (elem *SdpEndpoint) getConstructorParams(from IMediaObject, options map[string]interface{}) map[string]interface{} {
	return options

}

// Request a SessionSpec offer.
// This can be used to initiate a connection.
// Returns:
// // The SDP offer.
func (elem *SdpEndpoint) GenerateOffer() (string, error) {
	req := elem.getInvokeRequest()

	reqparams := map[string]interface{}{
		"operation": "generateOffer",
		"object":    elem.Id,
	}
	if elem.connection.SessionId != "" {
		reqparams["sessionId"] = elem.connection.SessionId
	}
	req["params"] = reqparams

	// Call server and wait response
	response := <-elem.connection.Request(req)

	// // The SDP offer.

	if response.Error == nil {
		return response.Result["value"], nil
	} else {
		return response.Result["value"], response.Error
	}
	return response.Result["value"], response.Error

}

// Request the NetworkConnection to process the given SessionSpec offer (from the
// remote User Agent)
// Returns:
// // The chosen configuration from the ones stated in the SDP offer
func (elem *SdpEndpoint) ProcessOffer(offer string) (string, error) {
	req := elem.getInvokeRequest()

	params := make(map[string]interface{})

	setIfNotEmpty(params, "offer", offer)

	reqparams := map[string]interface{}{
		"operation":       "processOffer",
		"object":          elem.Id,
		"operationParams": params,
	}
	if elem.connection.SessionId != "" {
		reqparams["sessionId"] = elem.connection.SessionId
	}
	req["params"] = reqparams

	// Call server and wait response
	response := <-elem.connection.Request(req)

	// // The chosen configuration from the ones stated in the SDP offer

	if response.Error == nil {
		return response.Result["value"], nil
	} else {
		return response.Result["value"], response.Error
	}
	return response.Result["value"], response.Error

}

// Request the NetworkConnection to process the given SessionSpec answer (from the
// remote User Agent).
// Returns:
// // Updated SDP offer, based on the answer received.
func (elem *SdpEndpoint) ProcessAnswer(answer string) (string, error) {
	req := elem.getInvokeRequest()

	params := make(map[string]interface{})

	setIfNotEmpty(params, "answer", answer)

	reqparams := map[string]interface{}{
		"operation":       "processAnswer",
		"object":          elem.Id,
		"operationParams": params,
	}
	if elem.connection.SessionId != "" {
		reqparams["sessionId"] = elem.connection.SessionId
	}
	req["params"] = reqparams

	// Call server and wait response
	response := <-elem.connection.Request(req)

	// // Updated SDP offer, based on the answer received.

	if response.Error == nil {
		return response.Result["value"], nil
	} else {
		return response.Result["value"], response.Error
	}
	return response.Result["value"], response.Error

}

// This method gives access to the SessionSpec offered by this NetworkConnection.
// .. note:: This method returns the local MediaSpec, negotiated or not. If no
// offer has been generated yet, it returns null. It an offer has been
// generated it returns the offer and if an answer has been processed
// it returns the negotiated local SessionSpec.
// Returns:
// // The last agreed SessionSpec
func (elem *SdpEndpoint) GetLocalSessionDescriptor() (string, error) {
	req := elem.getInvokeRequest()

	reqparams := map[string]interface{}{
		"operation": "getLocalSessionDescriptor",
		"object":    elem.Id,
	}
	if elem.connection.SessionId != "" {
		reqparams["sessionId"] = elem.connection.SessionId
	}
	req["params"] = reqparams

	// Call server and wait response
	response := <-elem.connection.Request(req)

	// // The last agreed SessionSpec

	if response.Error == nil {
		return response.Result["value"], nil
	} else {
		return response.Result["value"], response.Error
	}
	return response.Result["value"], response.Error

}

// This method gives access to the remote session description.
// .. note:: This method returns the media previously agreed after a complete
// offer-answer exchange. If no media has been agreed yet, it returns null.
// Returns:
// // The last agreed User Agent session description
func (elem *SdpEndpoint) GetRemoteSessionDescriptor() (string, error) {
	req := elem.getInvokeRequest()

	reqparams := map[string]interface{}{
		"operation": "getRemoteSessionDescriptor",
		"object":    elem.Id,
	}
	if elem.connection.SessionId != "" {
		reqparams["sessionId"] = elem.connection.SessionId
	}
	req["params"] = reqparams

	// Call server and wait response
	response := <-elem.connection.Request(req)

	// // The last agreed User Agent session description

	if response.Error == nil {
		return response.Result["value"], nil
	} else {
		return response.Result["value"], response.Error
	}
	return response.Result["value"], response.Error

}

type IBaseRtpEndpoint interface {
	GetStats(mediaType MediaType) (map[string]Stats, error)
}

// Base class to manage common RTP features.
type BaseRtpEndpoint struct {
	SdpEndpoint

	// Minimum video bandwidth for receiving.
	// Unit: kbps(kilobits per second).
	// 0: unlimited.
	// Default value: 100
	MinVideoRecvBandwidth int

	// Minimum video bandwidth for sending.
	// Unit: kbps(kilobits per second).
	// 0: unlimited.
	// Default value: 100
	MinVideoSendBandwidth int

	// Maximum video bandwidth for sending.
	// Unit: kbps(kilobits per second).
	// 0: unlimited.
	// Default value: 500
	MaxVideoSendBandwidth int

	// State of the media
	MediaState *MediaState

	// State of the connection
	ConnectionState *ConnectionState

	// Parameters of the congestion control algorithm
	RembParams *RembParams
}

// Return contructor params to be called by "Create".
func (elem *BaseRtpEndpoint) getConstructorParams(from IMediaObject, options map[string]interface{}) map[string]interface{} {
	return options

}

// Provides statistics collected for this endpoint
// Returns:
// // Delivers a successful result in the form of a RTC stats report. A RTC stats
// // report represents a map between strings, identifying the inspected objects
// // (RTCStats.id), and their corresponding RTCStats objects.
func (elem *BaseRtpEndpoint) GetStats(mediaType MediaType) (map[string]Stats, error) {
	req := elem.getInvokeRequest()

	params := make(map[string]interface{})

	setIfNotEmpty(params, "mediaType", mediaType)

	reqparams := map[string]interface{}{
		"operation":       "getStats",
		"object":          elem.Id,
		"operationParams": params,
	}
	if elem.connection.SessionId != "" {
		reqparams["sessionId"] = elem.connection.SessionId
	}
	req["params"] = reqparams

	// Call server and wait response
	response := <-elem.connection.Request(req)

	// // Delivers a successful result in the form of a RTC stats report. A RTC stats
	// // report represents a map between strings, identifying the inspected objects
	// // (RTCStats.id), and their corresponding RTCStats objects.

	ret := map[string]Stats{}
	return ret, response.Error

}

type IMediaElement interface {
	GetSourceConnections(mediaType MediaType, description string) ([]ElementConnectionData, error)
	GetSinkConnections(mediaType MediaType, description string) ([]ElementConnectionData, error)
	Connect(sink IMediaElement, mediaType MediaType, sourceMediaDescription string, sinkMediaDescription string) error
	Disconnect(sink IMediaElement, mediaType MediaType, sourceMediaDescription string, sinkMediaDescription string) error
	SetAudioFormat(caps AudioCaps) error
	SetVideoFormat(caps VideoCaps) error
	GetGstreamerDot(details GstreamerDotDetails) (string, error)
	SetOutputBitrate(bitrate int) error
}

// Basic building blocks of the media server, that can be interconnected through
// the API. A `MediaElement` is a module that encapsulates a specific media
// capability. They can be connected to create media pipelines where those
// capabilities are applied, in sequence, to the stream going through the
// pipeline.
// `MediaElement` objects are classified by its supported media type (audio,
// video, etc.)
type MediaElement struct {
	MediaObject
}

// Return contructor params to be called by "Create".
func (elem *MediaElement) getConstructorParams(from IMediaObject, options map[string]interface{}) map[string]interface{} {
	return options

}

// Get the connections information of the elements that are sending media to this
// element `MediaElement`
// Returns:
// // A list of the connections information that are sending media to this
// element.
// // The list will be empty if no sources are found.
func (elem *MediaElement) GetSourceConnections(mediaType MediaType, description string) ([]ElementConnectionData, error) {
	req := elem.getInvokeRequest()

	params := make(map[string]interface{})

	setIfNotEmpty(params, "mediaType", mediaType)
	setIfNotEmpty(params, "description", description)

	reqparams := map[string]interface{}{
		"operation":       "getSourceConnections",
		"object":          elem.Id,
		"operationParams": params,
	}
	if elem.connection.SessionId != "" {
		reqparams["sessionId"] = elem.connection.SessionId
	}
	req["params"] = reqparams

	// Call server and wait response
	response := <-elem.connection.Request(req)

	// // A list of the connections information that are sending media to this
	// element.
	// // The list will be empty if no sources are found.

	ret := []ElementConnectionData{}
	return ret, response.Error

}

// Returns a list of the connections information of the elements that ere
// receiving media from this element.
// Returns:
// // A list of the connections information that arereceiving media from this
// // element. The list will be empty if no sinks are found.
func (elem *MediaElement) GetSinkConnections(mediaType MediaType, description string) ([]ElementConnectionData, error) {
	req := elem.getInvokeRequest()

	params := make(map[string]interface{})

	setIfNotEmpty(params, "mediaType", mediaType)
	setIfNotEmpty(params, "description", description)

	reqparams := map[string]interface{}{
		"operation":       "getSinkConnections",
		"object":          elem.Id,
		"operationParams": params,
	}
	if elem.connection.SessionId != "" {
		reqparams["sessionId"] = elem.connection.SessionId
	}
	req["params"] = reqparams

	// Call server and wait response
	response := <-elem.connection.Request(req)

	// // A list of the connections information that arereceiving media from this
	// // element. The list will be empty if no sinks are found.

	ret := []ElementConnectionData{}
	return ret, response.Error

}

// Connects two elements, with the given restrictions, current `MediaElement` will
// start emmit media to sink element. Connection could take place in the future,
// when both media element show capabilities for connecting with the given
// restrictions
func (elem *MediaElement) Connect(sink IMediaElement, mediaType MediaType, sourceMediaDescription string, sinkMediaDescription string) error {
	req := elem.getInvokeRequest()

	params := make(map[string]interface{})

	setIfNotEmpty(params, "sink", sink)
	setIfNotEmpty(params, "mediaType", mediaType)
	setIfNotEmpty(params, "sourceMediaDescription", sourceMediaDescription)
	setIfNotEmpty(params, "sinkMediaDescription", sinkMediaDescription)

	reqparams := map[string]interface{}{
		"operation":       "connect",
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

// Disconnects two elements, with the given restrictions, current `MediaElement`
// stops sending media to sink element. If the previously requested connection
// didn't took place it is also removed
func (elem *MediaElement) Disconnect(sink IMediaElement, mediaType MediaType, sourceMediaDescription string, sinkMediaDescription string) error {
	req := elem.getInvokeRequest()

	params := make(map[string]interface{})

	setIfNotEmpty(params, "sink", sink)
	setIfNotEmpty(params, "mediaType", mediaType)
	setIfNotEmpty(params, "sourceMediaDescription", sourceMediaDescription)
	setIfNotEmpty(params, "sinkMediaDescription", sinkMediaDescription)

	reqparams := map[string]interface{}{
		"operation":       "disconnect",
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

// Sets the type of data for the audio stream. MediaElements that do not support
// configuration of audio capabilities will raise an exception
func (elem *MediaElement) SetAudioFormat(caps AudioCaps) error {
	req := elem.getInvokeRequest()

	params := make(map[string]interface{})

	setIfNotEmpty(params, "caps", caps)

	reqparams := map[string]interface{}{
		"operation":       "setAudioFormat",
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

// Sets the type of data for the video stream. MediaElements that do not support
// configuration of video capabilities will raise an exception
func (elem *MediaElement) SetVideoFormat(caps VideoCaps) error {
	req := elem.getInvokeRequest()

	params := make(map[string]interface{})

	setIfNotEmpty(params, "caps", caps)

	reqparams := map[string]interface{}{
		"operation":       "setVideoFormat",
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

// Returns a string in dot (graphviz) format that represents the gstreamer
// elements inside
// Returns:
// // The dot graph
func (elem *MediaElement) GetGstreamerDot(details GstreamerDotDetails) (string, error) {
	req := elem.getInvokeRequest()

	params := make(map[string]interface{})

	setIfNotEmpty(params, "details", details)

	reqparams := map[string]interface{}{
		"operation":       "getGstreamerDot",
		"object":          elem.Id,
		"operationParams": params,
	}
	if elem.connection.SessionId != "" {
		reqparams["sessionId"] = elem.connection.SessionId
	}
	req["params"] = reqparams

	// Call server and wait response
	response := <-elem.connection.Request(req)

	// // The dot graph

	if response.Error == nil {
		return response.Result["value"], nil
	} else {
		return response.Result["value"], response.Error
	}
	return response.Result["value"], response.Error

}

// Allows change the target bitrate for the media output, if the media is encoded
// using VP8 or H264. This method only works if it is called before the media
// starts to flow.
func (elem *MediaElement) SetOutputBitrate(bitrate int) error {
	req := elem.getInvokeRequest()

	params := make(map[string]interface{})

	setIfNotEmpty(params, "bitrate", bitrate)

	reqparams := map[string]interface{}{
		"operation":       "setOutputBitrate",
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
