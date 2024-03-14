package streamline

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/nats-io/nats.go"
)

type NatsEventStream struct {
	js            nats.JetStream
	subjectPrefix string
	opts          *streamOpts
}

func NewNatsEventStream(js nats.JetStream, subjectPrefix string, opts ...StreamOpt) (*NatsEventStream, error) {
	if js == nil {
		return nil, errors.New("streamline: nil js")
	}

	if subjectPrefix == "" {
		return nil, errors.New("streamline: empty subject prefix")
	}

	strmOpts := &streamOpts{}
	for _, opt := range opts {
		opt.configureStream(strmOpts)
	}

	return &NatsEventStream{
		js:            js,
		subjectPrefix: subjectPrefix,
		opts:          strmOpts,
	}, nil
}

func (es *NatsEventStream) Publish(ctx context.Context, event Event) error {
	eventType := reflect.TypeOf(event)
	eventName, ok := TagValue(eventType)
	if !ok {
		return errors.New("streamline: missing streamline tag")
	}

	objectID, ok := TagFieldValue(event)
	if !ok {
		return errors.New("streamline: missing object id")
	}

	if objectID == "" {
		return errors.New("streamline: empty object id")
	}

	subject, err := natsSubject(es.subjectPrefix, eventName, objectID)
	if err != nil {
		return err
	}

	// TODO: doesn't have to be json
	payload, err := json.Marshal(event)
	if err != nil {
		panic(err)
	}

	_, err = es.js.Publish(subject, payload)
	if err != nil {
		return err
	}

	return nil
}

func (es *NatsEventStream) StreamTo(ctx context.Context, recv Receiver) error {
	subOpts := []nats.SubOpt{nats.AckExplicit()}
	if es.opts.durableName != "" {
		subOpts = append(subOpts, nats.Durable(es.opts.durableName))
	}

	msgCh := make(chan *nats.Msg, 100)
	_, err := es.js.ChanQueueSubscribe(
		es.subjectPrefix+".>", "streamline",
		msgCh,
		subOpts...)
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case msg := <-msgCh:
			eventName, err := eventNameFromSubject(es.subjectPrefix, msg.Subject)
			if err != nil {
				return err
			}

			err = recv.Receive(eventName, msg.Data)
			if err != nil {
				return err
			}

			msg.Ack()
		}
	}
}

func natsSubject(prefix, eventName string, objectID string) (string, error) {
	// eventName format: <aggregate-name>.<event-name>
	segments := strings.Split(eventName, ".")
	if len(segments) != 2 {
		return "", errors.New("streamline: invalid event name")
	}

	if objectID == "" {
		return "", errors.New("streamline: empty object id")
	}

	// format: <prefix>.<aggregate-name>.<aggregate-id>.<event-name>
	return fmt.Sprintf("%s.%s.%s.%s", prefix, segments[0], objectID, segments[1]), nil
}

func eventNameFromSubject(prefix, subject string) (string, error) {
	// format: <prefix>.<aggregate-name>.<aggregate-id>.<event-name>
	if !strings.HasPrefix(subject, prefix+".") {
		return "", errors.New("streamline: invalid subject")
	}

	// format: <aggregate-name>.<aggregate-id>.<event-name>
	fullyQualifiedEventName := subject[len(prefix+"."):]
	segments := strings.Split(fullyQualifiedEventName, ".")
	if len(segments) != 3 {
		return "", errors.New("streamline: invalid subject")
	}

	// format: <aggregate-name>.<event-name>
	return segments[0] + "." + segments[2], nil
}

type streamOpts struct {
	durableName string
}
type StreamOpt interface {
	configureStream(*streamOpts)
}

type streamOptFunc func(*streamOpts)

func (f streamOptFunc) configureStream(opts *streamOpts) {
	f(opts)
}

func Durable(name string) StreamOpt {
	return streamOptFunc(func(opts *streamOpts) {
		opts.durableName = name
	})
}
