aws-sdk-go-wrapper
====

[![Build Status](https://travis-ci.org/evalphobia/aws-sdk-go-wrapper.svg?branch=master)](https://travis-ci.org/evalphobia/aws-sdk-go-wrapper) [![Coverage Status](https://coveralls.io/repos/evalphobia/logrus_kinesis/badge.svg?branch=master&service=github)](https://coveralls.io/github/evalphobia/logrus_kinesis?branch=master) [![Coverage Status](https://coveralls.io/repos/evalphobia/aws-sdk-go-wrapper/badge.svg?branch=master)](https://coveralls.io/r/evalphobia/aws-sdk-go-wrapper?branch=master) [![GoDoc](https://godoc.org/github.com/evalphobia/aws-sdk-go-wrapper?status.svg)](https://godoc.org/github.com/evalphobia/aws-sdk-go-wrapper) [![Join the chat at https://gitter.im/evalphobia/aws-sdk-go-wrapper](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/evalphobia/aws-sdk-go-wrapper?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

(checked SDK version [aws-sdk-go](https://github.com/awslabs/aws-sdk-go/) :: [v1.4.10](https://github.com/awslabs/aws-sdk-go/tree/v1.4.10)

Simple wrapper for aws-sdk-go
At this time, it suports services below,

- `DynamoDB`
- `S3`
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
- `SNS`

# Quick Usage

### DynamoDB

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

### S3

```go

import(
    "os"

    // import this
    "github.com/evalphobia/aws-sdk-go-wrapper/s3"
)

func main(){
    cli := s3.NewClient()
    bucket := cli.GetBucket("MyBucket")

    // upload file
    var file *os.File
    file = getFile() // dummy code. this expects return data of "*os.File", like from POST form.
    s3obj := s3.NewS3Object(file)
    bucket.AddObject(s3obj, "/foo/bar/new_file")
    bucket.Put()

    // upload file from text data
    text := "Lorem ipsum"
    s3obj2 := s3.NewS3ObjectString(text)
    bucket.AddObject(s3obj2, "/log/new_text.txt")

    // upload file of ACL authenticated-read
    bucket.AddSecretObject(s3obj2, "/secret_path/new_secret_file.txt")


    // put all added objects.
    bucket.Put() // upload "/log/new_text.txt" & "/secret_path/new_secret_file.txt"
}
```


### SQS

```go

import(
    "fmt"

    // import this
    "github.com/evalphobia/aws-sdk-go-wrapper/config"
    "github.com/evalphobia/aws-sdk-go-wrapper/sqs"
)

func main(){
    svc, err := sqs.New(config.Config{
        AccessKey: "access key",
        SecretKey: "access key",
    })
    if err != nil {
        panic("error on create client")
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
        // shoe message content
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
