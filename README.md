aws-sdk-go-wrapper
----

[![GoDoc][1]][2] [![License: MIT][3]][4] [![Release][5]][6] [![Build Status][7]][8] [![Coveralls Coverage][9]][10] [![Codecov Coverage][11]][12] [![Go Report Card][13]][14] [![Code Climate][19]][20] [![BCH compliance][21]][22]

[1]: https://godoc.org/github.com/evalphobia/aws-sdk-go-wrapper?status.svg
[2]: https://godoc.org/github.com/evalphobia/aws-sdk-go-wrapper
[3]: https://img.shields.io/badge/License-MIT-blue.svg
[4]: LICENSE.md
[5]: https://img.shields.io/github/release/evalphobia/aws-sdk-go-wrapper.svg
[6]: https://github.com/evalphobia/aws-sdk-go-wrapper/releases/latest
[7]: https://github.com/evalphobia/aws-sdk-go-wrapper/workflows/test/badge.svg
[8]: https://github.com/evalphobia/aws-sdk-go-wrapper/actions?query=workflow%3Atest
[9]: https://coveralls.io/repos/evalphobia/aws-sdk-go-wrapper/badge.svg?branch=master&service=github
[10]: https://coveralls.io/github/evalphobia/aws-sdk-go-wrapper?branch=master
[11]:https://codecov.io/gh/evalphobia/aws-sdk-go-wrapper/branch/master/graph/badge.svg
[12]: https://codecov.io/gh/evalphobia/aws-sdk-go-wrapper
[13]: https://goreportcard.com/badge/github.com/evalphobia/aws-sdk-go-wrapper
[14]: https://goreportcard.com/report/github.com/evalphobia/aws-sdk-go-wrapper
[15]: https://img.shields.io/github/downloads/evalphobia/aws-sdk-go-wrapper/total.svg?maxAge=1800
[16]: https://github.com/evalphobia/aws-sdk-go-wrapper/releases
[17]: https://img.shields.io/github/stars/evalphobia/aws-sdk-go-wrapper.svg
[18]: https://github.com/evalphobia/aws-sdk-go-wrapper/stargazers
[19]: https://codeclimate.com/github/evalphobia/aws-sdk-go-wrapper/badges/gpa.svg
[20]: https://codeclimate.com/github/evalphobia/aws-sdk-go-wrapper
[21]: https://bettercodehub.com/edge/badge/evalphobia/aws-sdk-go-wrapper?branch=master
[22]: https://bettercodehub.com/
[23]: https://img.shields.io/badge/License-Apache%202.0-blue.svg
[24]: LICENSE.md


Simple wrapper for aws-sdk-go
At this time, this library suports these AWS services below,

| Service | API |
| :--- | :-- |
| [`CloudTrail`](/cloudtrail) | LookupEvents |
| [`CloudWatch`](/cloudwatch) | GetMetricStatistics |
| [`CostExplorer`](/costexplorer) | GetCostAndUsage |
| [`DynamoDB`](/dynamodb) | BatchWriteItem |
|  | CreateTable |
|  | DeleteItem |
|  | DeleteTable |
|  | DescribeTable |
|  | GetItem |
|  | ListTables |
|  | PutItem |
|  | Query |
|  | UpdateTable |
|  | Scan |
| [`IAM`](/iam) | GetGroup |
|  | GetGroupPolicy |
|  | GetPolicyVersion |
|  | GetRolePolicy |
|  | GetUserPolicy |
|  | ListEntitiesForPolicy |
|  | ListGroups |
|  | ListGroupPolicies |
|  | ListPolicies |
|  | ListUsers |
|  | ListUserPolicies |
|  | ListRoles |
|  | ListRolePolicies |
| [`Kinesis`](/kinesis) | CreateStream |
|  | DeleteStream |
|  | DescribeStream |
|  | GetRecords |
|  | GetShardIterator |
|  | PutRecord |
| [`KMS`](/kms) | CreateAlias |
|  | CreateKey |
|  | Decrypt |
|  | DescribeKey |
|  | Encrypt |
|  | ReEncrypt |
|  | ReEncrypt |
|  | ScheduleKeyDeletion |
| [`Pinpoint`](/pinpoint) | SendEmail |
| [`Rekognition`](/rekognition) | CompareFaces |
|  | CreateCollection |
|  | DeleteCollection |
|  | DeleteFaces |
|  | DetectFaces |
|  | DetectLabels |
|  | DetectModerationLabels |
|  | GetCelebrityInfo |
|  | IndexFaces |
|  | ListCollections |
|  | ListFaces |
|  | RecognizeCelebrities |
|  | SearchFaces |
|  | SearchFacesByImage |
| [`S3`](/s3) | CreateBucket |
|  | CopyObject |
|  | DeleteBucket |
|  | DeleteObject |
|  | GetObject |
|  | HeadObject |
|  | ListObjectsV2 |
|  | PutObject |
| [`SNS`](/sns) | CreatePlatformEndpoint |
|  | CreateTopic |
|  | DeleteTopic |
|  | GetEndpointAttributes |
|  | GetPlatformApplicationAttributes |
|  | Publish |
|  | SetEndpointAttributes |
|  | Subscribe |
| [`SQS`](/sqs) | ChangeMessageVisibility |
|  | CreateQueue |
|  | DeleteMessage |
|  | DeleteMessageBatch |
|  | DeleteQueue |
|  | GetQueueAttributes |
|  | GetQueueUrl |
|  | ListQueues |
|  | PurgeQueue |
|  | ReceiveMessage |
|  | SendMessageBatch |
| [`X-Ray`](/xray) | PutTraceSegments |


# Quick Usage

### CloudTrail

```go
import (
    "fmt"
    "time"

    "github.com/evalphobia/aws-sdk-go-wrapper/cloudtrail"
    "github.com/evalphobia/aws-sdk-go-wrapper/config"
)

func main() {
    // Create CloudTrail service
    svc, err := cloudtrail.New(config.Config{
        AccessKey: "access",
        SecretKey: "secret",
        Region: "ap-north-east1",
    })
    if err != nil {
        panic("error to create client")
    }

    // Get all of CloudTrail events with in the time.
    results, err := svc.LookupEventsAll(cloudtrail.LookupEventsInput{
        StartTime: time.Now(),
        EndTime:   time.Now(),
        LookupAttributes: []LookupAttribute{
            {
                Key:   "EventName",
                Value: "GetEndpointAttributes", // sns
            },
        },
    })
    if err != nil {
        panic("error to get table")
    }

    for _, v := results.Events {
        fmt.Printf("%+v\n", v)
    }
```

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

### Kinesis

```go
import(
    "encoding/json"

    "github.com/evalphobia/aws-sdk-go-wrapper/config"
    "github.com/evalphobia/aws-sdk-go-wrapper/kinesis"
)

func main(){
    // Create Kinesis service
    svc, err := kinesis.New(config.Config{
        AccessKey: "access key",
        SecretKey: "access key",
        Region: "ap-north-east1",
    })
    if err != nil {
        panic("error on creating client")
    }

    // Get Kinesis Stream
    stream, err := svc.GetStream("my-stream")
    if err != nil {
        panic("error on getting stream")
    }

    // Get ShardID list of the stream
    shardIDs, err := stream.GetShardIDs()
    if err != nil {
        panic("error on getting shard id")
    }

    // get records from all of the shards
    for _, shardID := range shardIDs {
        // get records
        result, err := stream.GetRecords(kinesis.GetCondition{
            ShardID:           shardID,
            ShardIteratorType: kinesis.IteratorTypeLatest,
        })
        if err != nil {
            panic("error on getting records")
        }

        // get next records from the last result.
        result, err = stream.GetRecords(kinesis.GetCondition{
            ShardID:           shardID,
            ShardIteratorType: kinesis.IteratorTypeLatest,
            ShardIterator:     result.NextShardIterator,
        })
    }

    data := make(map[string]interface{})
    data["foo"] = 999
    data["bar"] = "some important info"

    bytData, _ := json.Marshal(data)

    // put data into stream record
    err = stream.PutRecord(bytData)
    if err != nil {
        panic("error on putting record")
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

MIT
