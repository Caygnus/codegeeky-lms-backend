package types

// PubSubType defines the type of pubsub implementation
type PubSubType string

const (
	// MemoryPubSub uses in-memory implementation
	MemoryPubSub PubSubType = "memory"

	// KafkaPubSub uses Kafka implementation
	KafkaPubSub PubSubType = "kafka"
)
