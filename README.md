aws-sdk-go-wrapper
====

[![Build Status](https://travis-ci.org/evalphobia/aws-sdk-go-wrapper.svg?branch=master)](https://travis-ci.org/evalphobia/aws-sdk-go-wrapper) [![Coverage Status](https://coveralls.io/repos/evalphobia/logrus_kinesis/badge.svg?branch=master&service=github)](https://coveralls.io/github/evalphobia/logrus_kinesis?branch=master) [![Coverage Status](https://coveralls.io/repos/evalphobia/aws-sdk-go-wrapper/badge.svg?branch=master)](https://coveralls.io/r/evalphobia/aws-sdk-go-wrapper?branch=master) [![GoDoc](https://godoc.org/github.com/evalphobia/aws-sdk-go-wrapper?status.svg)](https://godoc.org/github.com/evalphobia/aws-sdk-go-wrapper) [![Join the chat at https://gitter.im/evalphobia/aws-sdk-go-wrapper](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/evalphobia/aws-sdk-go-wrapper?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

(checked SDK version [aws-sdk-go](https://github.com/awslabs/aws-sdk-go/) :: [v1.4.10](https://github.com/awslabs/aws-sdk-go/tree/v1.4.10)

Simple wrapper for aws-sdk-go
At this time, it suports services below,

- [`DynamoDB`](https://github.com/evalphobia/aws-sdk-go-wrapper/tree/master/dynamodb)
    - ListTables
    - DescribeTable
    - CreateTable
    - UpdateTable
    - DeleteTable
    - Scan
    - Query
    - GetItem
    - PutItem
    - DeleteItem
- [`S3`](https://github.com/evalphobia/aws-sdk-go-wrapper/tree/master/s3)
    - GetObject
    - PutObject
    - DeleteObject
- [`SNS`](https://github.com/evalphobia/aws-sdk-go-wrapper/tree/master/sns)
    - CreatePlatformEndpoint
    - CreateTopic
    - DeleteTopic
    - Subscribe
    - Publish
    - GetEndpointAttributes
    - SetEndpointAttributes
- [`SQS`](https://github.com/evalphobia/aws-sdk-go-wrapper/tree/master/sqs)
    - GetQueueUrl
    - CreateQueue
    - DeleteQueue
    - PurgeQueue
    - SendMessageBatch
    - ReceiveMessage
    - DeleteMessage
    - DeleteMessageBatch
    - GetQueueAttributes


# Quick Usage

### DynamoDB

```go
import (
    "github.com/evalphobia/aws-sdk-go-wrapper/config"
    "github.com/evalphobia/aws-sdk-go-wrapper/dynamodb"
)

func main() {
    // Create DynamoDB service
    svc, err := dynamodb.New(config.Config{
        AccessKey: "access",
        SecretKey: "secret",
        Region: "ap-north-east1",
        Endpoint:  "http://localhost:8000", // option for DynamoDB Local
    })
    if err != nil {
        panic("error to create client")
    }

    // Get DynamoDB table
    table, err := svc.GetTable("MyDynamoTable")
    if err != nil {
        panic("error to get table")
    }

    // Create new DynamoDB item (row on RDBMS)
    item := dynamodb.NewPutItem()
    item.AddAttribute("user_id", 999)
    item.AddAttribute("status", 1)

    // Add item to the put spool
    table.AddItem(item)

    item2 := dynamodb.NewItem()
    item.AddAttribute("user_id", 1000)
    item.AddAttribute("status", 2)
    item.AddConditionEQ("status", 3) // Add condition for write
    table.AddItem(item2)

    // Put all items in the put spool
    err = table.Put()

    // Use svc.PutAll() to put all of the tables,
    // `err = svc.PutAll()`

    // Scan items
    cond = table.NewConditionList()
    cond.SetLimit(1000)
    cond.FilterEQ("status", 2)
    result, err = table.ScanWithCondition(cond)
    data := result.ToSliceMap() // `result.ToSliceMap()` returns []map[string]interface{}

    //Scan from last key
    cond.SetStartKey(result.LastEvaluatedKey)
    result, err = table.ScanWithCondition(cond)
    data = append(data, result.ToSliceMap())

    // Query items
    cond := table.NewConditionList()
    cond.AndEQ("user_id", 999)
    cond.FilterLT("age", 20)
    cond.SetLimit(100)
    result, err := table.Query(cond)
    if err != nil {
        panic("error to query")
    }

    // mapping result data to the struct
    type User struct {
        ID int64 `dynamodb:"user_id"`
        Age int `dynamodb:"age"`
        Status int `dynamodb:"status"`
    }
    var list []*User
    err = result.Unmarshal(&list)
    if err != nil {
        panic("error to unmarshal")
    }

    if len(list) == int(result.Count) {
        fmt.Println("success to get items")
    }
}
```

### S3

```go

import(
    "os"

    "github.com/evalphobia/aws-sdk-go-wrapper/config"
    "github.com/evalphobia/aws-sdk-go-wrapper/s3"
)

func main(){
    // Create S3 service
    svc, err := s3.New(config.Config{
        AccessKey: "access",
        SecretKey: "secret",
        Region: "ap-north-east1",
        S3ForcePathStyle: true,
        Endpoint:  "http://localhost:4567", // option for FakeS3
    })
    if err != nil {
        panic("error to create client")
    }

    bucket := svc.GetBucket("MyBucket")

    // upload file
    var file *os.File
    file = getFile() // dummy code. this expects return data of "*os.File". e.g. from POST form.
    s3obj := s3.NewPutObject(file)
    bucket.AddObject(s3obj, "/foo/bar/new_file")
    err = bucket.PutAll()
    if err != nil {
       panic("error to put file")
    }

    // upload file from text data
    text := "Lorem ipsum"
    s3obj2 := s3.NewPutObjectString(text)
    bucket.AddObject(s3obj2, "/log/new_text.txt")

    // upload file of ACL authenticated-read
    bucket.AddSecretObject(s3obj2, "/secret_path/new_secret_file.txt")

    // put all added objects.
    err = bucket.PutAll() // upload "/log/new_text.txt" & "/secret_path/new_secret_file.txt"
    if err != nil {
       panic("error to put files")
    }

    byt, err := bucket.GetObjectByte("/log/new_text.txt")
    if err != nil {
       panic("error to get file")
    }

    fmt.Println(string(byt)) // => Lorem ipsum
}
```


### SNS

```go

import(
    "fmt"

    "github.com/evalphobia/aws-sdk-go-wrapper/config"
    "github.com/evalphobia/aws-sdk-go-wrapper/sns"
)

func main(){
    svc, err := sns.New(config.Config{
        AccessKey: "access key",
        SecretKey: "access key",
        Region: "ap-north-east1",
    }, Platforms{
        Production: false, // flag for APNS/APNS sandbox.
        Apple:      "arn:aws:sns:us-east-1:0000000000:app/APNS/foo_apns", // Endpoint ARN for APNS
        Google:     "arn:aws:sns:us-east-1:0000000000:app/GCM/foo_gcm", // Endpoint ARN for GCM
    })
    if err != nil {
        panic("error to create client")
    }

    // send message to iOS devices.
    tokenListForIOS := []string{"fooEndpoint"}
    err = svc.BulkPublishByDevice("ios", tokenListForIOS, "push message!")
    if err != nil {
        panic("error to publish")
    }

    // send message to multiple devices.
    tokenList := map[string][]string{
        "android": {"token1", "token2"},
        "ios": {"token3", "token4"},
    }
    err = svc.BulkPublish(tokenList, "push message!")
    if err != nil {
        panic("error to publish")
    }
}
```


### SQS

```go

import(
    "fmt"

    "github.com/evalphobia/aws-sdk-go-wrapper/config"
    "github.com/evalphobia/aws-sdk-go-wrapper/sqs"
)

func main(){
    svc, err := sqs.New(config.Config{
        AccessKey: "access key",
        SecretKey: "access key",
        Region: "ap-north-east1",
    })
    if err != nil {
        panic("error to create client")
    }

    queue := svc.GetQueue("my-queue")

    // add message to spool
    queue.AddMessage("my message")

    // send messages in spool
    err := queue.Send()
    if err != nil {
        panic("error on sending sqs message")
    }

    // count message in SQS Queue
    num, _, _ := queue.CountMessage()
    if num > 0 {
        panic("message count must be sent")
    }

    // fetch messages from SQS Queue
    // maximum 10 message
    messageList, err := queue.Fetch(10)
    if err != nil {
        panic("error on getting sqs message")
    }

    for _, msg := messageList {
        // print message content
        fmt.Println(msg.Body())

        // delete message manually
        // if set queue.AutoDelete(true), messages are delete on fetching process
        queue.DeleteMessage(msg)
    }

    // purge queue
    queue.Purge()
}
```

# License

Apache License, Version 2.0
