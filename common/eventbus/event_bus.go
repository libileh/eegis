package eventbus

import (
	"errors"
	"sync"
	"time"

	"go.uber.org/zap"
)

// Event represents a generic event with a type, occurrence time, and associated data.
type Event struct {
	Type       string
	OccurredAt time.Time
	Data       interface{}
}

// EventBus manages subscribers for various event types.
type EventBus struct {
	logger      *zap.SugaredLogger
	subscribers map[string][]chan<- Event
	mutex       sync.RWMutex
}

// New creates a new instance of EventBus with the provided logger.
func New(logger *zap.SugaredLogger) *EventBus {
	return &EventBus{
		logger:      logger,
		subscribers: make(map[string][]chan<- Event),
	}
}

// Subscribe adds a subscriber channel for a specific event type.
// Returns an error if eventType is empty, the subscriber channel is nil, or if subscription fails.
func (eb *EventBus) Subscribe(eventType string, subscriber chan<- Event) error {
	if eventType == "" {
		return errors.New("event type cannot be empty")
	}
	if subscriber == nil {
		return errors.New("subscriber channel cannot be nil")
	}

	eb.mutex.Lock()
	defer eb.mutex.Unlock()

	eb.subscribers[eventType] = append(eb.subscribers[eventType], subscriber)
	eb.logger.Infof("subscriber added for event type: %s", eventType)
	return nil
}

// Publish sends an event to all subscriber channels registered for the event's type.
// Returns an error if no subscribers exist for the given event type.
func (eb *EventBus) Publish(event Event) error {
	eb.mutex.RLock()
	subs, exists := eb.subscribers[event.Type]
	eb.mutex.RUnlock()

	if !exists || len(subs) == 0 {
		eb.logger.Warnf("publish failed: no subscribers for event type: %s", event.Type)
		return errors.New("no subscribers for event type")
	}

	// Log the publishing event
	eb.logger.Infof("publishing event of type: %s at %v", event.Type, event.OccurredAt)

	// Deliver events asynchronously to avoid blocking publisher,
	// though in high load scenarios a buffered channel approach might be necessary.
	for _, sub := range subs {
		go func(ch chan<- Event) {
			ch <- event
		}(sub)
	}
	return nil
}
