package tools

import (
	"encoding/json"

	"errors"

	aws_events "github.com/aws/aws-lambda-go/events"

	"github.com/aws/aws-sdk-go/aws/arn"
)

type Record struct {
	EventSource    string
	EventSourceArn string
	AWSRegion      string
	S3             aws_events.S3Entity
	SQS            aws_events.SQSMessage
	SNS            aws_events.SNSEntity
}

// Event incoming event
type Event struct {
	Records []Record
}

type eventType int

const (
	unknownEventType eventType = iota
	s3EventType
	snsEventType
	sqsEventType
)

func (event *Event) UnmarshalJSON(data []byte) error {
	var err error

	switch event.getEventType(data) {
	case s3EventType:
		s3Event := &aws_events.S3Event{}
		err = json.Unmarshal(data, s3Event)

		if err == nil {
			return event.mapS3EventRecords(s3Event)
		}

	case snsEventType:
		snsEvent := &aws_events.SNSEvent{}
		err = json.Unmarshal(data, snsEvent)

		if err == nil {
			return event.mapSNSEventRecords(snsEvent)
		}

	case sqsEventType:
		sqsEvent := &aws_events.SQSEvent{}
		err = json.Unmarshal(data, sqsEvent)

		if err == nil {
			return event.mapSQSEventRecords(sqsEvent)
		}
	}

	return err
}

func (event *Event) mapS3EventRecords(s3Event *aws_events.S3Event) error {
	event.Records = make([]Record, 0)

	for _, s3Record := range s3Event.Records {
		event.Records = append(event.Records, Record{
			EventSource:    s3Record.EventSource,
			EventSourceArn: s3Record.S3.Bucket.Arn,
			AWSRegion:      s3Record.AWSRegion,
			S3:             s3Record.S3,
		})
	}

	return nil
}

func (event *Event) mapSNSEventRecords(snsEvent *aws_events.SNSEvent) error {
	event.Records = make([]Record, 0)

	for _, snsRecord := range snsEvent.Records {
		// decode sns message to s3 event
		s3Event := &aws_events.S3Event{}
		err := json.Unmarshal([]byte(snsRecord.SNS.Message), s3Event)
		if err != nil {
			//return errors.Wrap(err, "Failed to decode sns message to an S3 event")
			return err
		}

		if len(s3Event.Records) == 0 {
			return errors.New("SNS s3 event records is empty")
		}

		for _, s3Record := range s3Event.Records {
			topicArn, err := arn.Parse(snsRecord.SNS.TopicArn)
			if err != nil {
				return err
			}

			event.Records = append(event.Records, Record{
				EventSource:    snsRecord.EventSource,
				EventSourceArn: snsRecord.SNS.TopicArn,
				AWSRegion:      topicArn.Region,
				SNS:            snsRecord.SNS,
				S3:             s3Record.S3,
			})
		}
	}

	return nil
}

func (event *Event) mapSQSEventRecords(sqsEvent *aws_events.SQSEvent) error {
	event.Records = make([]Record, 0)

	for _, sqsRecord := range sqsEvent.Records {

		// decode sqs body to s3 event
		s3Event := &aws_events.S3Event{}
		err := json.Unmarshal([]byte(sqsRecord.Body), s3Event)
		if err != nil {
			return errors.New("Failed to decode sqs body to an S3 event")
		}

		if len(s3Event.Records) == 0 {
			return errors.New("SQS s3 event records is empty")
		}

		for _, s3Record := range s3Event.Records {
			event.Records = append(event.Records, Record{
				EventSource:    sqsRecord.EventSource,
				EventSourceArn: sqsRecord.EventSourceARN,
				AWSRegion:      sqsRecord.AWSRegion,
				SQS:            sqsRecord,
				S3:             s3Record.S3,
			})
		}
	}

	return nil
}

func (event *Event) getEventType(data []byte) eventType {
	temp := make(map[string]interface{})
	json.Unmarshal(data, &temp)

	recordsList, _ := temp["Records"].([]interface{})
	record, _ := recordsList[0].(map[string]interface{})

	var eventSource string

	if es, ok := record["EventSource"]; ok {
		eventSource = es.(string)

	} else if es, ok := record["eventSource"]; ok {
		eventSource = es.(string)
	}

	switch eventSource {
	case "aws:s3":
		return s3EventType
	case "aws:sns":
		return snsEventType
	case "aws:sqs":
		return sqsEventType
	}

	return unknownEventType
}
