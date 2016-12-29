package core

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/couchbase/moss"
)

// The EventRecord is a wrapper around a KV store that is useful for storing
// raw events from the app.  For example, storing all of the SQS messages received.
type EventRecord interface {
	StoreSQSMessage(sqsMessage *sqs.Message) error
	GetStoredSQSMessages() (sqsMessages []sqs.Message, err error)
	Close() error
}

type NoOpEventRecord struct{}

func (n NoOpEventRecord) StoreSQSMessage(sqsMessage *sqs.Message) error {
	return nil
}

func (n NoOpEventRecord) GetStoredSQSMessages() (sqsMessages []sqs.Message, err error) {
	return []sqs.Message{}, nil
}

func (n NoOpEventRecord) Close() error {
	return nil
}

// And EventRecord that uses the Moss KV store as a backend
type MossEventRecord struct {
	collection moss.Collection
	store      *moss.Store
}

// Create a new Moss EventRecord impl
func NewMossEventRecord(persistToDisk bool, storageDir string) (*MossEventRecord, error) {

	mossEventRecord := &MossEventRecord{}

	if persistToDisk {
		// Open moss in persistent mode
		store, collection, err := moss.OpenStoreCollection(
			storageDir,
			moss.StoreOptions{},
			moss.StorePersistOptions{},
		)
		if err != nil {
			return nil, fmt.Errorf("Error setting up persistent event record: %v", err)
		}
		mossEventRecord.store = store
		mossEventRecord.collection = collection

		// Apparently .Start() shouldn't be called on persisted collections
		// https://github.com/couchbase/moss/issues/5
		// mossEventRecord.collection.Start()

	} else {
		// Open moss in-memory store only
		collection, err := moss.NewCollection(moss.CollectionOptions{})
		if err != nil {
			return nil, fmt.Errorf("Error setting up event record: %v", err)
		}
		mossEventRecord.collection = collection

		// Call Start() or else it will panic when trying to close it
		// https://github.com/couchbase/moss/issues/4
		mossEventRecord.collection.Start()
	}

	return mossEventRecord, nil

}

// Store an SQS message in Moss
func (mer *MossEventRecord) StoreSQSMessage(sqsMessage *sqs.Message) error {

	if sqsMessage.MessageId == nil {
		return fmt.Errorf("Cannot store SQS message since MessageId is nil")
	}

	// serialize to JSON and store in Moss KV store
	sqsMessageBytes, err := json.Marshal(sqsMessage)
	if err != nil {
		return err
	}

	batch, err := mer.collection.NewBatch(0, 0)
	if err != nil {
		return err
	}

	defer batch.Close()

	batch.Set([]byte(*sqsMessage.MessageId), sqsMessageBytes)

	return mer.collection.ExecuteBatch(
		batch,
		moss.WriteOptions{},
	)

}

func (mer *MossEventRecord) GetStoredSQSMessages() (sqsMessages []sqs.Message, err error) {

	result := []sqs.Message{}

	snapshot, err := mer.collection.Snapshot()
	if err != nil {
		return result, err
	}
	if snapshot == nil {
		return result, fmt.Errorf("Unable to take moss collection snapshot")
	}
	defer snapshot.Close()

	iter, err := snapshot.StartIterator(nil, nil, moss.IteratorOptions{})
	if err != nil {
		return result, err
	}
	if iter == nil {
		return result, fmt.Errorf("Unable to get moss collection iterator")
	}
	defer iter.Close()

	for {
		_, v, err := iter.Current()

		if err == moss.ErrIteratorDone {
			return result, nil
		}

		sqsMessage := sqs.Message{}
		err = json.Unmarshal(v, &sqsMessage)
		if err != nil {
			return result, err
		}
		result = append(result, sqsMessage)

		err = iter.Next()
		if err == moss.ErrIteratorDone {
			return result, nil
		}

	}

	return result, nil

}

func (mer *MossEventRecord) Close() error {

	if err := mer.collection.Close(); err != nil {
		return err
	}
	if mer.store != nil {

		if err := mer.store.Close(); err != nil {
			return err
		}
	}

	return nil
}
