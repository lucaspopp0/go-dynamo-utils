package dynamoutils

import (
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type TTL time.Time

// MarshalDynamoDBAttributeValue implements [attributevalue.Marshaler].
func (t TTL) MarshalDynamoDBAttributeValue() (types.AttributeValue, error) {
	// Convert the TTL timestamp to unix seconds
	unixSeconds := time.Time(t).Unix()

	// Format it as a string, for storage
	asString := strconv.FormatInt(unixSeconds, 10)

	return &types.AttributeValueMemberN{
		Value: asString,
	}, nil
}

// UnmarshalDynamoDBAttributeValue implements [attributevalue.Unmarshaler].
func (t *TTL) UnmarshalDynamoDBAttributeValue(av types.AttributeValue) error {
	// Check if the provided value is a number
	asNumber, isNumber := av.(*types.AttributeValueMemberN)
	if !isNumber {
		return fmt.Errorf("expected *types.AttributeValueMemberN, got %v", reflect.TypeOf(av))
	}

	// Parse the value to unix seconds
	unixSeconds, err := strconv.ParseInt(asNumber.Value, 10, 0)
	if err != nil {
		return fmt.Errorf("error parsing value: %w", err)
	}

	// Construct a timestamp from the unix seconds
	*t = TTL(time.Unix(unixSeconds, 0))

	return nil
}
