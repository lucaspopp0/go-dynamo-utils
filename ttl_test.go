package dynamoutils

import (
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/require"
)

// Type checks
var (
	_ attributevalue.Marshaler   = (*TTL)(nil)
	_ attributevalue.Unmarshaler = (*TTL)(nil)
)

func TestTTL_Conversions(t *testing.T) {
	t.Run("success-direct-calls", func(t *testing.T) {
		// Create a TTL
		now := time.Now().Truncate(time.Second)
		ttl := NewTTL(now)

		// Marshal it to an attribute value
		marshaled, err := attributevalue.Marshal(ttl)
		require.NoError(t, err)

		// Validate the marshaled result is correct
		require.IsType(t, &types.AttributeValueMemberN{}, marshaled)
		require.Equal(t, fmt.Sprint(now.Unix()), marshaled.(*types.AttributeValueMemberN).Value)

		// Unmarshal back to a TTL object
		var unmarshaled TTL
		err = attributevalue.Unmarshal(marshaled, &unmarshaled)
		require.NoError(t, err)

		// Validate the unmarshaled result is correct
		require.WithinDuration(t, ttl.Time(), unmarshaled.Time(), 10*time.Millisecond)
	})

	t.Run("success-as-field", func(t *testing.T) {
		type ArbitraryStruct struct {
			TTL TTL `dynamodbav:"ttl"`
		}

		// Create a TTL
		now := time.Now().Truncate(time.Second)
		original := ArbitraryStruct{TTL: NewTTL(now)}

		// Marshal it to an attribute value
		marshaled, err := attributevalue.Marshal(original)
		require.NoError(t, err)

		// Validate the marshaled result is correct
		require.IsType(t, &types.AttributeValueMemberM{}, marshaled)
		ttlAttr, ok := marshaled.(*types.AttributeValueMemberM).Value["ttl"]
		require.True(t, ok)
		require.IsType(t, &types.AttributeValueMemberN{}, ttlAttr)
		require.Equal(t, fmt.Sprint(now.Unix()), ttlAttr.(*types.AttributeValueMemberN).Value)

		// Unmarshal back to an ArbitraryStruct
		var unmarshaled ArbitraryStruct
		err = attributevalue.Unmarshal(marshaled, &unmarshaled)
		require.NoError(t, err)

		// Validate the unmarshaled result is correct
		require.WithinDuration(t, original.TTL.Time(), unmarshaled.TTL.Time(), 10*time.Millisecond)
	})

	t.Run("unmarshal-err-non-number", func(t *testing.T) {
		av := &types.AttributeValueMemberS{Value: "1"}

		var unmarshaled TTL
		err := attributevalue.Unmarshal(av, &unmarshaled)

		require.ErrorContains(t, err, "AttributeValueMemberN")
	})

	t.Run("unmarshal-err-parseint-failed", func(t *testing.T) {
		av := &types.AttributeValueMemberN{Value: "hello"}

		var unmarshaled TTL
		err := attributevalue.Unmarshal(av, &unmarshaled)

		require.ErrorContains(t, err, "error parsing value")
	})
}
