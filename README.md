# go-dynamo-utils

Small Go utilities for working with [AWS DynamoDB](https://aws.amazon.com/dynamodb/).

## Installation

```bash
go get github.com/lucaspopp0/go-dynamo-utils
```

Requires Go 1.25 or later.

## TTL

DynamoDB [TTL](https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/TTL.html) attributes store an expiration time as a **Number** of Unix seconds.
The `TTL` type wraps `time.Time` and implements the AWS SDK v2 `attributevalue` marshaler interfaces, so it round-trips correctly with `attributevalue.Marshal` and `attributevalue.Unmarshal`.

| Function / method | Description |
|---|---|
| `NewTTL(time.Time) TTL` | Construct a `TTL` from a `time.Time` |
| `(TTL) Time() time.Time` | Convert a `TTL` back to `time.Time` |

```go
import (
    "time"

    dynamoutils "github.com/lucaspopp0/go-dynamo-utils"
    "github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
)

type Item struct {
    ID  string          `dynamodbav:"id"`
    TTL dynamoutils.TTL `dynamodbav:"ttl"`
}

item := Item{
    ID:  "abc-123",
    TTL: dynamoutils.NewTTL(time.Now().Add(24 * time.Hour)),
}

av, err := attributevalue.Marshal(item)
// av is a map with "ttl" set to a numeric Unix timestamp

var decoded Item
err = attributevalue.Unmarshal(av, &decoded)

expiresAt := decoded.TTL.Time()
```

`TTL` can also be marshaled and unmarshaled on its own, without wrapping it in a struct.

## Development

```bash
make lint       # go mod tidy, vet, and fmt
make unit-test  # go test ./...
```

## License

MIT — see [LICENSE](LICENSE).
