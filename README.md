# aws-sdk-go-wrapper
Simple wrapper for aws-sdk-go
At this time, it suports services below,
- `S3` 
- `DynamoDB`
- `SQS`
- `SNS`

# Quick Usage

for DynamoDB

```go
import (
    "github.com/evalphobia/aws-sdk-go-wrapper/dynamodb"
)

func main() {
    // Create connection client
    cli := ddb.NewClient()
    
    // Get dynamodb table
    table, err := cli.GetTable("MyDynamoTable")
    if err != nil {
        panic("error on loading dynamodb table")
    }
    
    // Create new dynamodb item (row on RDBMS)
    item := dynamodb.NewItem()
    item.AddAttribute("user_id", 999)
    item.AddAttribute("status", 1)
    
    // Add item to the wait list.
    table.AddItem(item)
    
    item2 := dynamodb.NewItem()
    item.AddAttribute("user_id", 1000)
    item.AddAttribute("status", 2)
    item.AddConditionEQ("status", 3) // Add condition for write 
    table.AddItem(item2)
    
    // write all
    cli.PutAll()
}
```

# Configure

```sh
vim aws.json

{
    "access_key": "XXXXXXXXXXXXXXXXXXXX",
    "secret_key": "abcdefg",
}

```


# License

Apache License, Version 2.0
