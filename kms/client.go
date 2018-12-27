package kms

import (
	"encoding/base64"

	SDK "github.com/aws/aws-sdk-go/service/kms"

	"github.com/evalphobia/aws-sdk-go-wrapper/config"
	"github.com/evalphobia/aws-sdk-go-wrapper/log"
	"github.com/evalphobia/aws-sdk-go-wrapper/private/pointers"
)

const (
	serviceName = "KMS"
)

// KMS has KMS client.
type KMS struct {
	client *SDK.KMS

	logger log.Logger
	prefix string
}

// New returns initialized *Rekognition.
func New(conf config.Config) (*KMS, error) {
	sess, err := conf.Session()
	if err != nil {
		return nil, err
	}

	svc := &KMS{
		client: SDK.New(sess),
		logger: log.DefaultLogger,
		prefix: conf.DefaultPrefix,
	}
	return svc, nil
}

// SetLogger sets logger.
func (svc *KMS) SetLogger(logger log.Logger) {
	svc.logger = logger
}

// CreateAlias executes CreateAlias operation.
func (svc *KMS) CreateAlias(keyID, aliasName string) error {
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
	metaData, err := svc.CreateKey(tags...)
	if err != nil {
		return nil, err
	}

	err = svc.CreateAlias(*metaData.KeyId, aliasName)
	return metaData, err
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

// Infof logging information.
func (svc *KMS) Infof(format string, v ...interface{}) {
	svc.logger.Infof(serviceName, format, v...)
}

// Errorf logging error information.
func (svc *KMS) Errorf(format string, v ...interface{}) {
	svc.logger.Errorf(serviceName, format, v...)
}
