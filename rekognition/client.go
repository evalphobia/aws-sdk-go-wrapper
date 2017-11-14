package rekognition

import (
	"net/http"
	"sync"
	"time"

	SDK "github.com/aws/aws-sdk-go/service/rekognition"

	"github.com/evalphobia/aws-sdk-go-wrapper/config"
	"github.com/evalphobia/aws-sdk-go-wrapper/log"
	"github.com/evalphobia/aws-sdk-go-wrapper/private/errors"
	"github.com/evalphobia/aws-sdk-go-wrapper/private/pointers"
)

const (
	serviceName    = "Rekognition"
	defaultTimeout = 30 * time.Second
)

// Rekognition has Rekognition client.
type Rekognition struct {
	client *SDK.Rekognition

	logger log.Logger
	prefix string

	collectionsMu sync.RWMutex
	collections   map[string]*Collection

	httpClient HTTPClient
}

// New returns initialized *Rekognition.
func New(conf config.Config) (*Rekognition, error) {
	sess, err := conf.Session()
	if err != nil {
		return nil, err
	}

	svc := &Rekognition{
		client:      SDK.New(sess),
		logger:      log.DefaultLogger,
		prefix:      conf.DefaultPrefix,
		collections: make(map[string]*Collection),
		httpClient: &http.Client{
			Timeout: defaultTimeout,
		},
	}
	return svc, nil
}

// SetLogger sets logger.
func (svc *Rekognition) SetLogger(logger log.Logger) {
	svc.logger = logger
}

// SetHTTPClient sets httpClient.
func (svc *Rekognition) SetHTTPClient(httpClient HTTPClient) {
	svc.httpClient = httpClient
}

// GetCollection gets Rekognition collection.
func (svc *Rekognition) GetCollection(name string) (*Collection, error) {
	collectionName := svc.prefix + name

	// get the stream from cache
	svc.collectionsMu.RLock()
	c, ok := svc.collections[collectionName]
	svc.collectionsMu.RUnlock()
	if ok {
		return c, nil
	}

	c = NewCollection(svc, name)
	svc.collectionsMu.Lock()
	svc.collections[collectionName] = c
	svc.collectionsMu.Unlock()
	return c, nil
}

// CreateCollection creates new Rekognition collection.
func (svc *Rekognition) CreateCollection(name string) error {
	collectionName := svc.prefix + name
	_, err := svc.client.CreateCollection(&SDK.CreateCollectionInput{
		CollectionId: pointers.String(collectionName),
	})
	if err != nil {
		svc.Errorf("error on `CreateCollection` operation; collection=%s; error=%s;", collectionName, err.Error())
		return err
	}

	svc.Infof("success on `CreateCollection` operation; collection=%s;", collectionName)
	return nil
}

// IsExistCollection checks if the Collection already exists or not.
func (svc *Rekognition) IsExistCollection(name string) (bool, error) {
	collectionName := svc.prefix + name

	collections, err := svc.ListCollections()
	if err != nil {
		svc.Errorf("error on `ListCollections` operation; error=%s", err.Error())
		return false, err
	}
	for _, collection := range collections {
		if collection == collectionName {
			return true, nil
		}
	}

	return false, nil
}

// ListCollections returns list of collections.
func (svc *Rekognition) ListCollections() ([]string, error) {
	const maxResult = 1000 // hard limit is 4096

	var nextToken *string
	collections := make([]string, 0, maxResult)
	hasNextToken := true
	for hasNextToken {
		list, token, err := svc.listCollections(&SDK.ListCollectionsInput{
			MaxResults: pointers.Long64(maxResult),
			NextToken:  nextToken,
		})
		if err != nil {
			return nil, err
		}

		collections = append(collections, list...)
		hasNextToken = (token != "")
		if hasNextToken {
			nextToken = &token
		}
	}

	return collections, nil
}

// listCollections returns list of collections.
func (svc *Rekognition) listCollections(input *SDK.ListCollectionsInput) (collections []string, nextToken string, err error) {
	op, err := svc.client.ListCollections(input)
	if err != nil {
		svc.Errorf("error on `ListCollections` operation; error=%s;", err.Error())
		return nil, "", err
	}

	collections = make([]string, len(op.CollectionIds))
	for i, id := range op.CollectionIds {
		collections[i] = *id
	}

	if op.NextToken != nil {
		nextToken = *op.NextToken
	}
	return collections, nextToken, nil
}

// ForceDeleteCollection deletes Rekognition collection by given name with prefix.
func (svc *Rekognition) ForceDeleteCollection(name string) error {
	collectionName := svc.prefix + name
	_, err := svc.client.DeleteCollection(&SDK.DeleteCollectionInput{
		CollectionId: pointers.String(collectionName),
	})
	if err != nil {
		svc.Errorf("error on `DeleteCollection` operation; collection=%s; error=%s;", collectionName, err.Error())
		return err
	}

	svc.Infof("success on `DeleteCollection` operation; collection=%s;", collectionName)
	return nil
}

// Infof logging information.
func (svc *Rekognition) Infof(format string, v ...interface{}) {
	svc.logger.Infof(serviceName, format, v...)
}

// Errorf logging error information.
func (svc *Rekognition) Errorf(format string, v ...interface{}) {
	svc.logger.Errorf(serviceName, format, v...)
}

func newErrors() *errors.Errors {
	return errors.NewErrors(serviceName)
}

// HTTPClient is used for fetching image data from URL.
type HTTPClient interface {
	Get(url string) (*http.Response, error)
}
