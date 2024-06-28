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

// DefaultQueueGroupName is the default queue group name.
const DefaultQueueGroupName = "streamline"

// NATSEventStream is an event stream that uses NATS JetStream as the underlying transport.
type NATSEventStream struct {
	js            nats.JetStream
	subjectPrefix string
	opts          *streamOpts
}

// NewNATSEventStream creates a new NATS event stream.
//
// The subject prefix is used to prefix the event subject.
//
// Example:
//
//	`sales` as the prefix will result in the subject format of `sales.<aggregate-name>.<aggregate-id>.<event-name>`.
func NewNATSEventStream(js nats.JetStream, subjectPrefix string, opts ...StreamOpt) (*NATSEventStream, error) {
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

	return &NATSEventStream{
		js:            js,
		subjectPrefix: subjectPrefix,
		opts:          strmOpts,
	}, nil
}

// Publish the event to the NATS JetStream.
func (es *NATSEventStream) Publish(_ context.Context, event Event) error {
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

	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}

	header := make(nats.Header)
	header.Set("Content-Type", "application/json")

	_, err = es.js.PublishMsg(&nats.Msg{
		Subject: subject,
		Data:    payload,
		Header:  header,
	})
	if err != nil {
		return err
	}

	return nil
}

// StreamTo streams the event from the NATS JetStream to the recv.
func (es *NATSEventStream) StreamTo(ctx context.Context, recv Receiver) error {
	queueGroupName := DefaultQueueGroupName
	if es.opts.queueGroupName != "" {
		queueGroupName = es.opts.queueGroupName
	}

	subOpts := []nats.SubOpt{nats.AckExplicit()}
	if es.opts.durableName != "" {
		subOpts = append(subOpts, nats.Durable(es.opts.durableName))
	}

	msgCh := make(chan *nats.Msg, 100)

	_, err := es.js.ChanQueueSubscribe(
		es.subjectPrefix+".>", queueGroupName,
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
	queueGroupName string
	durableName    string
}

// StreamOpt is an option to configure the NATS event stream.
type StreamOpt interface {
	configureStream(*streamOpts)
}

type streamOptFunc func(*streamOpts)

func (f streamOptFunc) configureStream(opts *streamOpts) {
	f(opts)
}

// QueueGroup sets the queue group name.
func QueueGroup(name string) StreamOpt {
	return streamOptFunc(func(opts *streamOpts) {
		opts.queueGroupName = name
	})
}

// Durable sets the durable name.
func Durable(name string) StreamOpt {
	return streamOptFunc(func(opts *streamOpts) {
		opts.durableName = name
	})
}
