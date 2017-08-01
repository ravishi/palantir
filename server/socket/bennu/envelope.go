package bennu

import (
	"encoding/json"
	"fmt"
)

type envelope struct {
	JoinRef string
	Ref string
	Topic string
	Event string
	Payload interface{}
}

func (e *envelope) MarshalJSON() ([]byte, error) {
	// let [join_ref, ref, topic, event, payload] = JSON.parse(rawPayload)
	var joinRef *string
	if e.JoinRef != "" {
		joinRef = &e.JoinRef
	}
	return json.Marshal([]interface{}{
		joinRef,
		e.Ref,
		e.Topic,
		e.Event,
		e.Payload,
	})
}

func (e *envelope) UnmarshalJSON(data []byte) error {
	// let payload = [
	// 	msg.join_ref, msg.ref, msg.topic, msg.event, msg.payload
	// ]
	arr := make([]json.RawMessage, 5)
	if err := json.Unmarshal(data, &arr); err != nil {
		return fmt.Errorf("Failed to unmarshal message: %s", err)
	}
	if err := json.Unmarshal(arr[0], &e.JoinRef); err != nil {
		return fmt.Errorf("Failed to unmarshal JoinRef: %s", err)
	}
	if err := json.Unmarshal(arr[1], &e.Ref); err != nil {
		return fmt.Errorf("Failed to unmarshal Ref: %s", err)
	}
	if err := json.Unmarshal(arr[2], &e.Topic); err != nil {
		return fmt.Errorf("Failed to unmarshal Topic: %s", err)
	}
	if err := json.Unmarshal(arr[3], &e.Event); err != nil {
		return fmt.Errorf("Failed to unmarshal Event: %s", err)
	}
	if err := json.Unmarshal(arr[4], &e.Payload); err != nil {
		return fmt.Errorf("Failed to unmarshal Payload: %s", err)
	}
	return nil
}