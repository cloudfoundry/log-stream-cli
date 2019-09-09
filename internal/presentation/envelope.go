package presentation

import (
	"code.cloudfoundry.org/go-loggregator/rpc/loggregator_v2"
	"errors"
)

func Envelope(l *loggregator_v2.Envelope) (EnvelopeWrapper, error) {
	switch l.Message.(type) {
	case *loggregator_v2.Envelope_Log:
		return buildLog(l), nil
	case *loggregator_v2.Envelope_Counter:
		return buildCounter(l), nil
	case *loggregator_v2.Envelope_Gauge:
		return buildGauge(l), nil
	case *loggregator_v2.Envelope_Timer:
		return buildTimer(l), nil
	case *loggregator_v2.Envelope_Event:
		return buildEvent(l), nil
	default:
		return nil, errors.New("Unknown envelope type")
	}
}

type baseEnvelope struct {
	Timestamp      int64                            `json:"timestamp,omitempty"`
	SourceId       string                           `json:"source_id,omitempty"`
	InstanceId     string                           `json:"instance_id,omitempty"`
	DeprecatedTags map[string]*loggregator_v2.Value `json:"deprecated_tags,omitempty"`
	Tags           map[string]string                `json:"tags,omitempty"`
}

func (b *baseEnvelope) presentationEnvelope() {}

type EnvelopeWrapper interface {
	presentationEnvelope()
}

type counterEnvelope struct {
	baseEnvelope
	Counter *counter `json:"counter"`
}

type counter struct {
	Name  string `json:"name"`
	Delta uint64 `json:"delta"`
	Total uint64 `json:"total"`
}

func buildCounter(l *loggregator_v2.Envelope) EnvelopeWrapper {
	c := l.GetCounter()
	return &counterEnvelope{
		baseEnvelope{
			l.Timestamp,
			l.SourceId,
			l.InstanceId,
			l.DeprecatedTags,
			l.Tags,
		},
		&counter{
			Name:  c.GetName(),
			Delta: c.GetDelta(),
			Total: c.GetTotal(),
		},
	}
}

type logEnvelope struct {
	baseEnvelope
	Log *log `json:"log"`
}

type log struct {
	Payload string                  `json:"payload,omitempty"`
	Type    loggregator_v2.Log_Type `json:"type,omitempty"`
}

func buildLog(l *loggregator_v2.Envelope) EnvelopeWrapper {
	return &logEnvelope{
		baseEnvelope{
			l.Timestamp,
			l.SourceId,
			l.InstanceId,
			l.DeprecatedTags,
			l.Tags,
		},
		&log{
			Payload: string(l.Message.(*loggregator_v2.Envelope_Log).Log.Payload),
			Type:    l.Message.(*loggregator_v2.Envelope_Log).Log.Type,
		},
	}
}

type gaugeEnvelope struct {
	baseEnvelope
	Gauge *gauge `json:"gauge"`
}

type gauge struct {
	Metrics map[string]*gaugeValue `json:"metrics"`
}

type gaugeValue struct {
	Unit  string  `json:"unit"`
	Value float64 `json:"value"`
}

func buildGauge(l *loggregator_v2.Envelope) EnvelopeWrapper {
	metrics := make(map[string]*gaugeValue)
	for k, v := range l.GetGauge().GetMetrics() {
		metrics[k] = &gaugeValue{
			Unit:  v.GetUnit(),
			Value: v.GetValue(),
		}
	}

	return &gaugeEnvelope{
		baseEnvelope{
			l.Timestamp,
			l.SourceId,
			l.InstanceId,
			l.DeprecatedTags,
			l.Tags,
		},
		&gauge{
			Metrics: metrics,
		},
	}
}

type timerEnvelope struct {
	baseEnvelope
	Timer *timer `json:"timer"`
}

type timer struct {
	Name  string `json:"name"`
	Start int64  `json:"start"`
	Stop  int64  `json:"stop"`
}

func buildTimer(l *loggregator_v2.Envelope) EnvelopeWrapper {
	t := l.GetTimer()
	return &timerEnvelope{
		baseEnvelope{
			l.Timestamp,
			l.SourceId,
			l.InstanceId,
			l.DeprecatedTags,
			l.Tags,
		},
		&timer{
			Name:  t.GetName(),
			Start: t.GetStart(),
			Stop:  t.GetStop(),
		},
	}
}

type eventEnvelope struct {
	baseEnvelope
	Event *event `json:"event"`
}

type event struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

func buildEvent(l *loggregator_v2.Envelope) EnvelopeWrapper {
	e := l.GetEvent()
	return &eventEnvelope{
		baseEnvelope{
			l.Timestamp,
			l.SourceId,
			l.InstanceId,
			l.DeprecatedTags,
			l.Tags,
		},
		&event{
			Title: e.GetTitle(),
			Body:  e.GetBody(),
		},
	}
}
