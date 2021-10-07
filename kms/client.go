package kms

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	SDK "github.com/aws/aws-sdk-go/service/kms"

	"github.com/evalphobia/aws-sdk-go-wrapper/config"
	"github.com/evalphobia/aws-sdk-go-wrapper/log"
	"github.com/evalphobia/aws-sdk-go-wrapper/private/pointers"
)

const (
	serviceName = "KMS"
	aliasPrefix = "alias/"
)

// KMS has KMS client.
type KMS struct {
	client *SDK.KMS

	logger log.Logger
}

// New returns initialized *Rekognition.
func New(conf config.Config) (*KMS, error) {
	sess, err := conf.Session()
	if err != nil {
		return nil, err
	}

	return NewFromSession(sess), nil
}

// NewFromSession returns initialized *KMS from aws.Session.
func NewFromSession(sess *session.Session) *KMS {
	return &KMS{
		client: SDK.New(sess),
		logger: log.DefaultLogger,
	}
}

// GetClient gets aws client.
func (svc *KMS) GetClient() *SDK.KMS {
	return svc.client
}

// SetLogger sets logger.
func (svc *KMS) SetLogger(logger log.Logger) {
	svc.logger = logger
}

// CreateAlias executes CreateAlias operation.
func (svc *KMS) CreateAlias(keyID, aliasName string) error {
	if !strings.HasPrefix(aliasName, aliasPrefix) {
		aliasName = aliasPrefix + aliasName
	}

	_, err := svc.client.CreateAlias(&SDK.CreateAliasInput{
		TargetKeyId: pointers.String(keyID),
		AliasName:   pointers.String(aliasName),
	})
	if err != nil {
		svc.Errorf("error on `CreateAlias` operation; keyID=%s; aliasName=%s; error=%s;", keyID, aliasName, err.Error())
		return err
	}

	svc.Infof("success on `CreateAlias` operation; keyID=%s; aliasName=%s;", keyID, aliasName)
	return nil
}

// CreateKey executes CreateKey operation.
func (svc *KMS) CreateKey(tags ...Tag) (*SDK.KeyMetadata, error) {
	var kmsTags []*SDK.Tag
	for _, tag := range tags {
		kmsTags = append(kmsTags, tag.Tag())
	}

	output, err := svc.client.CreateKey(&SDK.CreateKeyInput{
		Tags: kmsTags,
	})
	if err != nil {
		svc.Errorf("error on `CreateKey` operation; error=%s;", err.Error())
		return nil, err
	}

	metaData := output.KeyMetadata
	svc.Infof("success on `CreateKey` operation; keyID=%s; arn=%s;", *metaData.KeyId, *metaData.Arn)
	return metaData, nil
}

// CreateKeyWithAlias creates a key and sets alias name.
func (svc *KMS) CreateKeyWithAlias(aliasName string, tags ...Tag) (*SDK.KeyMetadata, error) {
	if !strings.HasPrefix(aliasName, aliasPrefix) {
		aliasName = aliasPrefix + aliasName
	}

	_, err := svc.DescribeKey(aliasName)
	if err == nil {
		return nil, fmt.Errorf("error: aliasName=[%s] is already exists", aliasName)
	}
	aerr, ok := err.(awserr.Error)
	if !ok {
		return nil, err
	}
	switch aerr.Code() {
	case SDK.ErrCodeNotFoundException:
		// error must be NotFoundException.
	default:
		return nil, err
	}

	metaData, err := svc.CreateKey(tags...)
	if err != nil {
		return nil, err
	}

	err = svc.CreateAlias(*metaData.KeyId, aliasName)
	return metaData, err
}

// DescribeKey executes DescribeKey operation.
func (svc *KMS) DescribeKey(keyName string) (metaData *SDK.KeyMetadata, err error) {
	output, err := svc.client.DescribeKey(&SDK.DescribeKeyInput{
		KeyId: pointers.String(keyName),
	})
	if err != nil {
		svc.Errorf("error on `DescribeKey` operation; keyName=%s; error=%s;", keyName, err.Error())
		return nil, err
	}

	return output.KeyMetadata, nil
}

// DeleteKey executes ScheduleKeyDeletion operation from Key name(key id, arn or alias).
func (svc *KMS) DeleteKey(keyName string, day int) error {
	metaData, err := svc.DescribeKey(keyName)
	if err != nil {
		return err
	}

	return svc.ScheduleKeyDeletion(*metaData.KeyId, day)
}

// ScheduleKeyDeletion executes ScheduleKeyDeletion operation.
func (svc *KMS) ScheduleKeyDeletion(keyID string, day int) error {
	_, err := svc.client.ScheduleKeyDeletion(&SDK.ScheduleKeyDeletionInput{
		KeyId:               pointers.String(keyID),
		PendingWindowInDays: pointers.Long(day),
	})
	if err != nil {
		svc.Errorf("error on `ScheduleKeyDeletion` operation; keyID=%s; error=%s;", keyID, err.Error())
	}
	return err
}

// Encrypt executes Encrypt operation.
func (svc *KMS) Encrypt(keyName string, plainData []byte) (encryptedData []byte, err error) {
	output, err := svc.client.Encrypt(&SDK.EncryptInput{
		KeyId:     pointers.String(keyName),
		Plaintext: plainData,
	})
	if err != nil {
		svc.Errorf("error on `Encrypt` operation; keyName=%s; error=%s;", keyName, err.Error())
		return nil, err
	}

	return output.CiphertextBlob, nil
}

// EncryptString executes Encrypt operation with base64 string.
func (svc *KMS) EncryptString(keyName, plainText string) (base64Text string, err error) {
	encryptedData, err := svc.Encrypt(keyName, []byte(plainText))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(encryptedData), nil
}

// Decrypt executes Decrypt operation.
func (svc *KMS) Decrypt(encryptedData []byte) (plainData []byte, err error) {
	output, err := svc.client.Decrypt(&SDK.DecryptInput{
		CiphertextBlob: encryptedData,
	})
	if err != nil {
		svc.Errorf("error on `Decrypt` operation; error=%s;", err.Error())
		return nil, err
	}

	return output.Plaintext, nil
}

// DecryptString executes Decrypt operation with base64 string.
func (svc *KMS) DecryptString(base64Text string) (plainText string, err error) {
	byt, err := base64.StdEncoding.DecodeString(base64Text)
	if err != nil {
		return "", err
	}

	plainData, err := svc.Decrypt(byt)
	if err != nil {
		return "", err
	}
	return string(plainData), nil
}

// ReEncrypt executes ReEncrypt operation.
func (svc *KMS) ReEncrypt(destinationKey string, encryptedData []byte) (resultEncryptedData []byte, err error) {
	output, err := svc.client.ReEncrypt(&SDK.ReEncryptInput{
		DestinationKeyId: pointers.String(destinationKey),
		CiphertextBlob:   encryptedData,
	})
	if err != nil {
		svc.Errorf("error on `ReEncrypt` operation; destinationKey=%s; error=%s;", destinationKey, err.Error())
		return nil, err
	}

	return output.CiphertextBlob, nil
}

// ReEncryptString executes ReEncrypt operation with base64 string.
func (svc *KMS) ReEncryptString(destinationKey, base64Text string) (resultBase64Text string, err error) {
	byt, err := base64.StdEncoding.DecodeString(base64Text)
	if err != nil {
		return "", err
	}

	encryptedData, err := svc.ReEncrypt(destinationKey, byt)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(encryptedData), nil
}

// Infof logging information.
func (svc *KMS) Infof(format string, v ...interface{}) {
	svc.logger.Infof(serviceName, format, v...)
}

// Errorf logging error information.
func (svc *KMS) Errorf(format string, v ...interface{}) {
	svc.logger.Errorf(serviceName, format, v...)
}
