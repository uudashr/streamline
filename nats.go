package streamline

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/nats-io/nats.go"
)

type NatsEventStream struct {
	nc            *nats.Conn
	subjectPrefix string
}

func NewNatsEventStream(nc *nats.Conn, subjectPrefix string) (*NatsEventStream, error) {
	if nc == nil {
		return nil, errors.New("streamline: nil nc")
	}

	if subjectPrefix == "" {
		return nil, errors.New("streamline: empty subject prefix")
	}

	return &NatsEventStream{
		nc:            nc,
		subjectPrefix: subjectPrefix,
	}, nil
}

func (ep *NatsEventStream) PublishEvents(ctx context.Context, events []Event) error {
	js, err := ep.nc.JetStream(nats.PublishAsyncMaxPending(100))
	if err != nil {
		return err
	}

	for _, event := range events {
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

		subject, err := natsSubject(ep.subjectPrefix, eventName, objectID)
		if err != nil {
			return err
		}

		// TODO: doesn't have to be json
		payload, err := json.Marshal(event)
		if err != nil {
			panic(err)
		}

		log.Printf("Publishing to %s: %s\n", subject, payload)
		_, err = js.Publish(subject, payload)
		if err != nil {
			return err
		}
	}

	return nil
}

func (ep *NatsEventStream) StreamTo(ctx context.Context, d Dispatcher) error {
	js, err := ep.nc.JetStream()
	if err != nil {
		return nil
	}

	msgCh := make(chan *nats.Msg, 100)
	_, err = js.ChanQueueSubscribe(ep.subjectPrefix+".>", "streamline", msgCh, nats.AckExplicit())
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case msg := <-msgCh:
			eventName, err := eventNameFromSubject(ep.subjectPrefix, msg.Subject)
			if err != nil {
				return err
			}

			err = d.Dispatch(eventName, msg.Data)
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
