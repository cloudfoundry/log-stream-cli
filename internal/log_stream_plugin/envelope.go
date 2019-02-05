package log_stream_plugin

import "code.cloudfoundry.org/go-loggregator/rpc/loggregator_v2"

type Envelope__Base64DecodedLog struct {
	Log *Base64DecodedLog `json:"log"`
}

type Base64DecodedEnvelope struct {
	Envelope__Base64DecodedLog
	Timestamp      int64                            `json:"timestamp,omitempty"`
	SourceId       string                           `json:"source_id,omitempty"`
	InstanceId     string                           `json:"instance_id,omitempty"`
	DeprecatedTags map[string]*loggregator_v2.Value `json:"deprecated_tags,omitempty"`
	Tags           map[string]string                `json:"tags,omitempty"`
}

type Base64DecodedLog struct {
	Payload string                  `json:"payload,omitempty"`
	Type    loggregator_v2.Log_Type `json:"type,omitempty"`
}

func BuildBase64DecodedLog(l *loggregator_v2.Envelope) *Base64DecodedEnvelope {
	return &Base64DecodedEnvelope{
		Envelope__Base64DecodedLog{
			Log: &Base64DecodedLog{
				Payload: string(l.Message.(*loggregator_v2.Envelope_Log).Log.Payload),
				Type:    l.Message.(*loggregator_v2.Envelope_Log).Log.Type,
			},
		},
		l.Timestamp,
		l.SourceId,
		l.InstanceId,
		l.DeprecatedTags,
		l.Tags,
	}
}
