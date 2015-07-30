package kurento

// Media Profile.
// Currently WEBM and MP4 are supported.
type MediaProfileSpecType string

// Implement fmt.Stringer interface
func (t MediaProfileSpecType) String() string {
	return string(t)
}

const (
	MEDIAPROFILESPECTYPE_WEBM            MediaProfileSpecType = "WEBM"
	MEDIAPROFILESPECTYPE_MP4             MediaProfileSpecType = "MP4"
	MEDIAPROFILESPECTYPE_WEBM_VIDEO_ONLY MediaProfileSpecType = "WEBM_VIDEO_ONLY"
	MEDIAPROFILESPECTYPE_WEBM_AUDIO_ONLY MediaProfileSpecType = "WEBM_AUDIO_ONLY"
	MEDIAPROFILESPECTYPE_MP4_VIDEO_ONLY  MediaProfileSpecType = "MP4_VIDEO_ONLY"
	MEDIAPROFILESPECTYPE_MP4_AUDIO_ONLY  MediaProfileSpecType = "MP4_AUDIO_ONLY"
)

type IceCandidate struct {
	Candidate     string
	SdpMid        string
	SdpMLineIndex int
}

func (t IceCandidate) CustomSerialize() map[string]interface{} {
	ret := make(map[string]interface{})

	ret["candidate"] = t.Candidate

	ret["sdpMid"] = t.SdpMid

	ret["sdpMLineIndex"] = t.SdpMLineIndex

	return ret
}

// States of an ICE component.
type IceComponentState string

// Implement fmt.Stringer interface
func (t IceComponentState) String() string {
	return string(t)
}

const (
	ICECOMPONENTSTATE_DISCONNECTED IceComponentState = "DISCONNECTED"
	ICECOMPONENTSTATE_GATHERING    IceComponentState = "GATHERING"
	ICECOMPONENTSTATE_CONNECTING   IceComponentState = "CONNECTING"
	ICECOMPONENTSTATE_CONNECTED    IceComponentState = "CONNECTED"
	ICECOMPONENTSTATE_READY        IceComponentState = "READY"
	ICECOMPONENTSTATE_FAILED       IceComponentState = "FAILED"
)

type ServerInfo struct {
	Version      string
	Modules      []ModuleInfo
	Type         ServerType
	Capabilities []string
}

// Indicates if the server is a real media server or a proxy
type ServerType string

// Implement fmt.Stringer interface
func (t ServerType) String() string {
	return string(t)
}

const (
	SERVERTYPE_KMS ServerType = "KMS"
	SERVERTYPE_KCS ServerType = "KCS"
)

// Details of gstreamer dot graphs
type GstreamerDotDetails string

// Implement fmt.Stringer interface
func (t GstreamerDotDetails) String() string {
	return string(t)
}

const (
	GSTREAMERDOTDETAILS_SHOW_MEDIA_TYPE         GstreamerDotDetails = "SHOW_MEDIA_TYPE"
	GSTREAMERDOTDETAILS_SHOW_CAPS_DETAILS       GstreamerDotDetails = "SHOW_CAPS_DETAILS"
	GSTREAMERDOTDETAILS_SHOW_NON_DEFAULT_PARAMS GstreamerDotDetails = "SHOW_NON_DEFAULT_PARAMS"
	GSTREAMERDOTDETAILS_SHOW_STATES             GstreamerDotDetails = "SHOW_STATES"
	GSTREAMERDOTDETAILS_SHOW_ALL                GstreamerDotDetails = "SHOW_ALL"
)

type ModuleInfo struct {
	Version   string
	Name      string
	Factories []string
}

// State of the media.
type MediaState string

// Implement fmt.Stringer interface
func (t MediaState) String() string {
	return string(t)
}

const (
	MEDIASTATE_DISCONNECTED MediaState = "DISCONNECTED"
	MEDIASTATE_CONNECTED    MediaState = "CONNECTED"
)

// State of the connection.
type ConnectionState string

// Implement fmt.Stringer interface
func (t ConnectionState) String() string {
	return string(t)
}

const (
	CONNECTIONSTATE_DISCONNECTED ConnectionState = "DISCONNECTED"
	CONNECTIONSTATE_CONNECTED    ConnectionState = "CONNECTED"
)

// Type of media stream to be exchanged.
// Can take the values AUDIO, DATA or VIDEO.
type MediaType string

// Implement fmt.Stringer interface
func (t MediaType) String() string {
	return string(t)
}

const (
	MEDIATYPE_AUDIO MediaType = "AUDIO"
	MEDIATYPE_DATA  MediaType = "DATA"
	MEDIATYPE_VIDEO MediaType = "VIDEO"
)

// Type of filter to be created.
// Can take the values AUDIO, VIDEO or AUTODETECT.
type FilterType string

// Implement fmt.Stringer interface
func (t FilterType) String() string {
	return string(t)
}

const (
	FILTERTYPE_AUDIO      FilterType = "AUDIO"
	FILTERTYPE_AUTODETECT FilterType = "AUTODETECT"
	FILTERTYPE_VIDEO      FilterType = "VIDEO"
)

// Codec used for transmission of video.
type VideoCodec string

// Implement fmt.Stringer interface
func (t VideoCodec) String() string {
	return string(t)
}

const (
	VIDEOCODEC_VP8  VideoCodec = "VP8"
	VIDEOCODEC_H264 VideoCodec = "H264"
	VIDEOCODEC_RAW  VideoCodec = "RAW"
)

// Codec used for transmission of audio.
type AudioCodec string

// Implement fmt.Stringer interface
func (t AudioCodec) String() string {
	return string(t)
}

const (
	AUDIOCODEC_OPUS AudioCodec = "OPUS"
	AUDIOCODEC_PCMU AudioCodec = "PCMU"
	AUDIOCODEC_RAW  AudioCodec = "RAW"
)

type Fraction struct {
	Numerator   int
	Denominator int
}

type AudioCaps struct {
	Codec   AudioCodec
	Bitrate int
}

type VideoCaps struct {
	Codec     VideoCodec
	Framerate Fraction
}

type ElementConnectionData struct {
	Source            MediaElement
	Sink              MediaElement
	Type              MediaType
	SourceDescription string
	SinkDescription   string
}

type Tag struct {
	Key   string
	Value string
}

// The type of the object.
type StatsType string

// Implement fmt.Stringer interface
func (t StatsType) String() string {
	return string(t)
}

const (
	STATSTYPE_inboundrtp      StatsType = "inboundrtp"
	STATSTYPE_outboundrtp     StatsType = "outboundrtp"
	STATSTYPE_session         StatsType = "session"
	STATSTYPE_datachannel     StatsType = "datachannel"
	STATSTYPE_track           StatsType = "track"
	STATSTYPE_transport       StatsType = "transport"
	STATSTYPE_candidatepair   StatsType = "candidatepair"
	STATSTYPE_localcandidate  StatsType = "localcandidate"
	STATSTYPE_remotecandidate StatsType = "remotecandidate"
)

type Stats struct {
	Id        string
	Type      StatsType
	Timestamp float64
}

type RTCStats struct {
}

type RTCRTPStreamStats struct {
	Ssrc             string
	AssociateStatsId string
	IsRemote         bool
	MediaTrackId     string
	TransportId      string
	CodecId          string
	FirCount         int64
	PliCount         int64
	NackCount        int64
	SliCount         int64
	Remb             int64
	PacketsLost      int64
	FractionLost     float64
}

type RTCCodec struct {
	PayloadType int64
	Codec       string
	ClockRate   int64
	Channels    int64
	Parameters  string
}

type RTCInboundRTPStreamStats struct {
	PacketsReceived int64
	BytesReceived   int64
	Jitter          float64
}

type RTCOutboundRTPStreamStats struct {
	PacketsSent   int64
	BytesSent     int64
	TargetBitrate float64
	RoundTripTime float64
}

type RTCPeerConnectionStats struct {
	DataChannelsOpened int64
	DataChannelsClosed int64
}

type RTCMediaStreamStats struct {
	StreamIdentifier string
	TrackIds         []string
}

type RTCMediaStreamTrackStats struct {
	TrackIdentifier           string
	RemoteSource              bool
	SsrcIds                   []string
	FrameWidth                int64
	FrameHeight               int64
	FramesPerSecond           float64
	FramesSent                int64
	FramesReceived            int64
	FramesDecoded             int64
	FramesDropped             int64
	FramesCorrupted           int64
	AudioLevel                float64
	EchoReturnLoss            float64
	EchoReturnLossEnhancement float64
}

// Represents the state of the RTCDataChannel
type RTCDataChannelState string

// Implement fmt.Stringer interface
func (t RTCDataChannelState) String() string {
	return string(t)
}

const (
	RTCDATACHANNELSTATE_connecting RTCDataChannelState = "connecting"
	RTCDATACHANNELSTATE_open       RTCDataChannelState = "open"
	RTCDATACHANNELSTATE_closing    RTCDataChannelState = "closing"
	RTCDATACHANNELSTATE_closed     RTCDataChannelState = "closed"
)

type RTCDataChannelStats struct {
	Label            string
	Protocol         string
	Datachannelid    int64
	State            RTCDataChannelState
	MessagesSent     int64
	BytesSent        int64
	MessagesReceived int64
	BytesReceived    int64
}

type RTCTransportStats struct {
	BytesSent               int64
	BytesReceived           int64
	RtcpTransportStatsId    string
	ActiveConnection        bool
	SelectedCandidatePairId string
	LocalCertificateId      string
	RemoteCertificateId     string
}

// Types of candidates
type RTCStatsIceCandidateType string

// Implement fmt.Stringer interface
func (t RTCStatsIceCandidateType) String() string {
	return string(t)
}

const (
	RTCSTATSICECANDIDATETYPE_host            RTCStatsIceCandidateType = "host"
	RTCSTATSICECANDIDATETYPE_serverreflexive RTCStatsIceCandidateType = "serverreflexive"
	RTCSTATSICECANDIDATETYPE_peerreflexive   RTCStatsIceCandidateType = "peerreflexive"
	RTCSTATSICECANDIDATETYPE_relayed         RTCStatsIceCandidateType = "relayed"
)

type RTCIceCandidateAttributes struct {
	IpAddress        string
	PortNumber       int64
	Transport        string
	CandidateType    RTCStatsIceCandidateType
	Priority         int64
	AddressSourceUrl string
}

// Represents the state of the checklist for the local and remote candidates in a
// pair.
type RTCStatsIceCandidatePairState string

// Implement fmt.Stringer interface
func (t RTCStatsIceCandidatePairState) String() string {
	return string(t)
}

const (
	RTCSTATSICECANDIDATEPAIRSTATE_frozen     RTCStatsIceCandidatePairState = "frozen"
	RTCSTATSICECANDIDATEPAIRSTATE_waiting    RTCStatsIceCandidatePairState = "waiting"
	RTCSTATSICECANDIDATEPAIRSTATE_inprogress RTCStatsIceCandidatePairState = "inprogress"
	RTCSTATSICECANDIDATEPAIRSTATE_failed     RTCStatsIceCandidatePairState = "failed"
	RTCSTATSICECANDIDATEPAIRSTATE_succeeded  RTCStatsIceCandidatePairState = "succeeded"
	RTCSTATSICECANDIDATEPAIRSTATE_cancelled  RTCStatsIceCandidatePairState = "cancelled"
)

type RTCIceCandidatePairStats struct {
	TransportId              string
	LocalCandidateId         string
	RemoteCandidateId        string
	State                    RTCStatsIceCandidatePairState
	Priority                 int64
	Nominated                bool
	Writable                 bool
	Readable                 bool
	BytesSent                int64
	BytesReceived            int64
	RoundTripTime            float64
	AvailableOutgoingBitrate float64
	AvailableIncomingBitrate float64
}

type RTCCertificateStats struct {
	Fingerprint          string
	FingerprintAlgorithm string
	Base64Certificate    string
	IssuerCertificateId  string
}

type CodecConfiguration struct {
	Name       string
	Properties map[string]interface{}
}

type RembParams struct {
	PacketsRecvIntervalTop int
	ExponentialFactor      float64
	LinealFactorMin        int
	LinealFactorGrade      float64
	DecrementFactor        float64
	ThresholdFactor        float64
	UpLosses               int
	RembOnConnect          int
}
