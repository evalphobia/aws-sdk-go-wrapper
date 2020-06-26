package rekognition

import (
	"io/ioutil"

	SDK "github.com/aws/aws-sdk-go/service/rekognition"

	"github.com/evalphobia/aws-sdk-go-wrapper/private/pointers"
)

const emptyThreshold = -1.0

// Collection has Collection client.
type Collection struct {
	service *Rekognition

	name           string
	nameWithPrefix string
	limit          int64
}

// NewCollection returns initialized *Collection.
func NewCollection(svc *Rekognition, name string) *Collection {
	collectionName := svc.prefix + name
	return &Collection{
		service:        svc,
		name:           name,
		nameWithPrefix: collectionName,
	}
}

// SetLimit sets limit, which is used for getting result list.
//  (e.g. MaxResults, MaxFaces etc...)
func (c *Collection) SetLimit(limit int64) {
	c.limit = limit
}

// ListAllFaces returns all of faces in a collection.
func (c *Collection) ListAllFaces() ([]*Face, error) {
	const defaultMaxResult = 1000 // hard limit is 4096
	limit := c.limit
	if limit == 0 {
		limit = defaultMaxResult
	}

	var nextToken *string
	faces := make([]*Face, 0, limit)
	hasNextToken := true
	for hasNextToken {
		list, token, err := c.listFaces(&SDK.ListFacesInput{
			CollectionId: pointers.String(c.nameWithPrefix),
			MaxResults:   pointers.Long64(limit),
			NextToken:    nextToken,
		})
		if err != nil {
			return nil, err
		}

		faces = append(faces, list...)
		hasNextToken = (token != "")
		if hasNextToken {
			nextToken = &token
		}
	}

	return faces, nil
}

// listFaces executes `ListFaces` operation and gets faces in a collection.
func (c *Collection) listFaces(input *SDK.ListFacesInput) (faces []*Face, nextToken string, err error) {
	svc := c.service
	op, err := svc.client.ListFaces(input)
	if err != nil {
		svc.Errorf("error on `ListFaces` operation; collection=%s; error=%s;", c.nameWithPrefix, err.Error())
		return nil, "", err
	}

	faces = make([]*Face, len(op.Faces))
	for i, f := range op.Faces {
		faces[i] = NewFaceFromAWSFace(f)
	}

	if op.NextToken != nil {
		nextToken = *op.NextToken
	}
	return faces, nextToken, nil
}

// DeleteFaces executes `DeleteFaces` operation and delete faces in a collection.
func (c *Collection) DeleteFaces(faceIDs []string) ([]string, error) {
	svc := c.service

	ids := make([]*string, len(faceIDs))
	for i, id := range faceIDs {
		id := id
		ids[i] = &id
	}

	op, err := svc.client.DeleteFaces(&SDK.DeleteFacesInput{
		CollectionId: pointers.String(c.nameWithPrefix),
		FaceIds:      ids,
	})
	if err != nil {
		svc.Errorf("error on `DeleteFaces` operation; collection=%s; error=%s;", c.nameWithPrefix, err.Error())
		return nil, err
	}

	result := make([]string, 0, len(op.DeletedFaces))
	for _, f := range op.DeletedFaces {
		if f != nil {
			result = append(result, *f)
		}
	}
	return result, nil
}

// IndexFacesFromLocalFile saves image into a collection from local image file.
func (c *Collection) IndexFacesFromLocalFile(filepath string) (*IndexFacesResponse, error) {
	return c.IndexFacesFromLocalFileWithExternalImageID("", filepath)
}

// IndexFacesFromLocalFileWithExternalImageID saves image into a collection from local image file.
func (c *Collection) IndexFacesFromLocalFileWithExternalImageID(externalID string, filepath string) (*IndexFacesResponse, error) {
	byt, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	return c.IndexFacesByBytesWithExternalImageID(externalID, byt)
}

// IndexFacesFromURL saves image into a collection from URL.
func (c *Collection) IndexFacesFromURL(url string) (*IndexFacesResponse, error) {
	return c.IndexFacesFromURLWithExternalImageID("", url)
}

// IndexFacesFromURLWithExternalImageID saves image into a collection from URL.
func (c *Collection) IndexFacesFromURLWithExternalImageID(externalID string, url string) (*IndexFacesResponse, error) {
	byt, err := c.service.getImageFromURL(url)
	if err != nil {
		return nil, err
	}
	return c.IndexFacesByBytesWithExternalImageID(externalID, byt)
}

// IndexFacesByBytes saves image of bytes into a collection.
func (c *Collection) IndexFacesByBytes(byt []byte) (*IndexFacesResponse, error) {
	return c.IndexFacesByBytesWithExternalImageID("", byt)
}

// IndexFacesByBytesWithExternalImageID saves image of bytes into a collection.
func (c *Collection) IndexFacesByBytesWithExternalImageID(externalID string, byt []byte) (*IndexFacesResponse, error) {
	input := &SDK.IndexFacesInput{
		CollectionId: pointers.String(c.nameWithPrefix),
		Image: &SDK.Image{
			Bytes: byt,
		},
	}
	if externalID != "" {
		input.ExternalImageId = pointers.String(externalID)
	}
	return c.indexFaces(input)
}

// IndexFacesByS3Object saves image of S3 object into a collection.
func (c *Collection) IndexFacesByS3Object(bucket, name string, version ...string) (*IndexFacesResponse, error) {
	return c.IndexFacesByS3ObjectWithExternalImageID("", bucket, name, version...)
}

// IndexFacesByS3ObjectWithExternalImageID saves image of S3 object into a collection.
func (c *Collection) IndexFacesByS3ObjectWithExternalImageID(externalID string, bucket, name string, version ...string) (*IndexFacesResponse, error) {
	input := &SDK.IndexFacesInput{
		CollectionId: pointers.String(c.nameWithPrefix),
		Image: &SDK.Image{
			S3Object: createS3Object(bucket, name, version...),
		},
	}
	if externalID != "" {
		input.ExternalImageId = pointers.String(externalID)
	}
	return c.indexFaces(input)
}

// indexFaces executes `IndexFaces` operation and gets faces in a collection.
func (c *Collection) indexFaces(input *SDK.IndexFacesInput) (*IndexFacesResponse, error) {
	svc := c.service
	op, err := svc.client.IndexFaces(input)
	if err != nil {
		svc.Errorf("error on `IndexFaces` operation; collection=%s; error=%s;", c.nameWithPrefix, err.Error())
		return nil, err
	}

	return NewIndexFacesResponseFromAWSOutput(op), nil
}

// SearchFacesByFaceID searches similar faces by face id in a collection.
func (c *Collection) SearchFacesByFaceID(faceID string, threshold ...float64) (*SearchFacesResponse, error) {
	input := &SDK.SearchFacesInput{
		CollectionId: pointers.String(c.nameWithPrefix),
		FaceId:       pointers.String(faceID),
	}
	if len(threshold) == 1 {
		input.FaceMatchThreshold = &threshold[0]
	}
	if c.limit != 0 {
		input.MaxFaces = pointers.Long64(c.limit)
	}
	return c.searchFaces(input)
}

// searchFaces executes `SearchFaces` operation and gets faces in a collection.
func (c *Collection) searchFaces(input *SDK.SearchFacesInput) (*SearchFacesResponse, error) {
	svc := c.service
	op, err := svc.client.SearchFaces(input)
	if err != nil {
		svc.Errorf("error on `SearchFaces` operation; collection=%s; error=%s;", c.nameWithPrefix, err.Error())
		return nil, err
	}

	return NewSearchFacesResponseFromAWSOutput(op), nil
}

// SearchFacesByBytes searches similar faces by image of bytes in a collection.
func (c *Collection) SearchFacesByBytes(byt []byte) (*SearchFacesResponse, error) {
	return c.SearchFacesByBytesWithThreshold(emptyThreshold, byt)
}

// SearchFacesByBytesWithThreshold searches similar faces by image of bytes in a collection.
func (c *Collection) SearchFacesByBytesWithThreshold(threshold float64, byt []byte) (*SearchFacesResponse, error) {
	input := &SDK.SearchFacesByImageInput{
		CollectionId:       pointers.String(c.nameWithPrefix),
		FaceMatchThreshold: &threshold,
		Image: &SDK.Image{
			Bytes: byt,
		},
	}
	if threshold != emptyThreshold {
		input.FaceMatchThreshold = &threshold
	}
	if c.limit != 0 {
		input.MaxFaces = pointers.Long64(c.limit)
	}
	return c.searchFacesByImage(input)
}

// SearchFacesByS3Object searches similar faces by image of S3 object in a collection.
func (c *Collection) SearchFacesByS3Object(bucket, name string, version ...string) (*SearchFacesResponse, error) {
	return c.SearchFacesByS3ObjectWithThreshold(emptyThreshold, bucket, name, version...)
}

// SearchFacesByS3ObjectWithThreshold searches similar faces by image of S3 object in a collection.
func (c *Collection) SearchFacesByS3ObjectWithThreshold(threshold float64, bucket, name string, version ...string) (*SearchFacesResponse, error) {
	input := &SDK.SearchFacesByImageInput{
		CollectionId:       pointers.String(c.nameWithPrefix),
		FaceMatchThreshold: &threshold,
		Image: &SDK.Image{
			S3Object: createS3Object(bucket, name, version...),
		},
	}
	if threshold != emptyThreshold {
		input.FaceMatchThreshold = &threshold
	}
	if c.limit != 0 {
		input.MaxFaces = pointers.Long64(c.limit)
	}
	return c.searchFacesByImage(input)
}

// searchFacesByImage executes `SearchFacesByImage` operation and gets faces in a collection.
func (c *Collection) searchFacesByImage(input *SDK.SearchFacesByImageInput) (*SearchFacesResponse, error) {
	svc := c.service
	op, err := svc.client.SearchFacesByImage(input)
	if err != nil {
		svc.Errorf("error on `SearchFacesByImage` operation; collection=%s; error=%s;", c.nameWithPrefix, err.Error())
		return nil, err
	}

	return NewSearchFacesResponseFromAWSOutputByImage(op), nil
}
