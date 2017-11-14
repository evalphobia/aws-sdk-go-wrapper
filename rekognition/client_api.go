package rekognition

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"

	SDK "github.com/aws/aws-sdk-go/service/rekognition"

	"github.com/evalphobia/aws-sdk-go-wrapper/private/pointers"
)

// DetectFacesFromLocalFile gets label info from local image file.
func (svc *Rekognition) DetectFacesFromLocalFile(filepath string) (*FaceDetailResponse, error) {
	byt, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	return svc.DetectFacesByBytes(byt)
}

// DetectFacesFromURL gets label info from URL.
func (svc *Rekognition) DetectFacesFromURL(url string) (*FaceDetailResponse, error) {
	byt, err := svc.getImageFromURL(url)
	if err != nil {
		return nil, err
	}
	return svc.DetectFacesByBytes(byt)
}

// DetectFacesByBytes gets label info from image of bytes.
func (svc *Rekognition) DetectFacesByBytes(byt []byte) (*FaceDetailResponse, error) {
	input := &SDK.DetectFacesInput{
		Image: &SDK.Image{
			Bytes: byt,
		},
	}
	return svc.detectFaces(input)
}

// DetectFacesByS3Object gets label info from image of S3 object.
func (svc *Rekognition) DetectFacesByS3Object(bucket, name string, version ...string) (*FaceDetailResponse, error) {
	input := &SDK.DetectFacesInput{
		Image: &SDK.Image{
			S3Object: createS3Object(bucket, name, version...),
		},
	}
	return svc.detectFaces(input)
}

// detectFaces executes `detectFace` operation and gets face info.
func (svc *Rekognition) detectFaces(input *SDK.DetectFacesInput) (*FaceDetailResponse, error) {
	op, err := svc.client.DetectFaces(input)
	if err != nil {
		svc.Errorf("error on `DetectFaces` operation; error=%s;", err.Error())
		return nil, err
	}

	list := make([]*FaceDetail, len(op.FaceDetails))
	for i, f := range op.FaceDetails {
		list[i] = NewFaceDetailFromAWSFaceDetail(f)
	}

	result := &FaceDetailResponse{
		List: list,
	}
	if op.OrientationCorrection != nil {
		result.Orientation = *op.OrientationCorrection
	}
	return result, nil
}

// DetectLabelsFromLocalFile gets label info from local image file.
func (svc *Rekognition) DetectLabelsFromLocalFile(filepath string) (*LabelResponse, error) {
	byt, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	return svc.DetectLabelsByBytes(byt)
}

// DetectLabelsFromURL gets label info from URL.
func (svc *Rekognition) DetectLabelsFromURL(url string) (*LabelResponse, error) {
	byt, err := svc.getImageFromURL(url)
	if err != nil {
		return nil, err
	}
	return svc.DetectLabelsByBytes(byt)
}

// DetectLabelsByBytes gets label info from image of bytes.
func (svc *Rekognition) DetectLabelsByBytes(byt []byte) (*LabelResponse, error) {
	input := &SDK.DetectLabelsInput{
		Image: &SDK.Image{
			Bytes: byt,
		},
	}
	return svc.detectLabels(input)
}

// DetectLabelsByS3Object gets label info from image of S3 object.
func (svc *Rekognition) DetectLabelsByS3Object(bucket, name string, version ...string) (*LabelResponse, error) {
	input := &SDK.DetectLabelsInput{
		Image: &SDK.Image{
			S3Object: createS3Object(bucket, name, version...),
		},
	}
	return svc.detectLabels(input)
}

// detectLabels executes `DetectLabels` operation and gets face info.
func (svc *Rekognition) detectLabels(input *SDK.DetectLabelsInput) (*LabelResponse, error) {
	op, err := svc.client.DetectLabels(input)
	if err != nil {
		svc.Errorf("error on `DetectLabels` operation; error=%s;", err.Error())
		return nil, err
	}

	list := make([]*Label, len(op.Labels))
	for i, l := range op.Labels {
		list[i] = NewLabelFromAWSLabel(l)
	}

	result := &LabelResponse{
		List: list,
	}
	if op.OrientationCorrection != nil {
		result.Orientation = *op.OrientationCorrection
	}
	return result, nil
}

// DetectModerationLabelsFromLocalFile gets moderation info from local image file.
func (svc *Rekognition) DetectModerationLabelsFromLocalFile(filepath string) (*ModerationLabelResponse, error) {
	byt, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	return svc.DetectModerationLabelsByBytes(byt)
}

// DetectModerationLabelsFromURL gets moderation info from URL.
func (svc *Rekognition) DetectModerationLabelsFromURL(url string) (*ModerationLabelResponse, error) {
	byt, err := svc.getImageFromURL(url)
	if err != nil {
		return nil, err
	}
	return svc.DetectModerationLabelsByBytes(byt)
}

// DetectModerationLabelsByBytes gets moderation info from image of bytes.
func (svc *Rekognition) DetectModerationLabelsByBytes(byt []byte) (*ModerationLabelResponse, error) {
	input := &SDK.DetectModerationLabelsInput{
		Image: &SDK.Image{
			Bytes: byt,
		},
	}
	return svc.detectModerationLabels(input)
}

// DetectModerationLabelsByS3Object gets moderation info from image of S3 object.
func (svc *Rekognition) DetectModerationLabelsByS3Object(bucket, name string, version ...string) (*ModerationLabelResponse, error) {
	input := &SDK.DetectModerationLabelsInput{
		Image: &SDK.Image{
			S3Object: createS3Object(bucket, name, version...),
		},
	}
	return svc.detectModerationLabels(input)
}

// detectModerationLabels executes `DetectModerationLabels` operation and gets face info.
func (svc *Rekognition) detectModerationLabels(input *SDK.DetectModerationLabelsInput) (*ModerationLabelResponse, error) {
	op, err := svc.client.DetectModerationLabels(input)
	if err != nil {
		svc.Errorf("error on `DetectModerationLabels` operation; error=%s;", err.Error())
		return nil, err
	}

	list := make([]*ModerationLabel, len(op.ModerationLabels))
	for i, l := range op.ModerationLabels {
		list[i] = NewModerationLabelFromAWSModerationLabel(l)
	}

	return &ModerationLabelResponse{
		List: list,
	}, nil
}

// RecognizeCelebritiesFromLocalFile gets celebrities info from local image file.
func (svc *Rekognition) RecognizeCelebritiesFromLocalFile(filepath string) (*CelebrityResponse, error) {
	byt, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	return svc.RecognizeCelebritiesByBytes(byt)
}

// RecognizeCelebritiesFromURL gets celebrities info from URL.
func (svc *Rekognition) RecognizeCelebritiesFromURL(url string) (*CelebrityResponse, error) {
	byt, err := svc.getImageFromURL(url)
	if err != nil {
		return nil, err
	}
	return svc.RecognizeCelebritiesByBytes(byt)
}

// RecognizeCelebritiesByBytes gets celebrities info from image of bytes.
func (svc *Rekognition) RecognizeCelebritiesByBytes(byt []byte) (*CelebrityResponse, error) {
	input := &SDK.RecognizeCelebritiesInput{
		Image: &SDK.Image{
			Bytes: byt,
		},
	}
	return svc.recognizeCelebrities(input)
}

// RecognizeCelebritiesByS3Object gets celebrities info from image of S3 object.
func (svc *Rekognition) RecognizeCelebritiesByS3Object(bucket, name string, version ...string) (*CelebrityResponse, error) {
	input := &SDK.RecognizeCelebritiesInput{
		Image: &SDK.Image{
			S3Object: createS3Object(bucket, name, version...),
		},
	}
	return svc.recognizeCelebrities(input)
}

// recognizeCelebrities executes `RecognizeCelebrities` operation and gets celebrities info.
func (svc *Rekognition) recognizeCelebrities(input *SDK.RecognizeCelebritiesInput) (*CelebrityResponse, error) {
	op, err := svc.client.RecognizeCelebrities(input)
	if err != nil {
		svc.Errorf("error on `RecognizeCelebrities` operation; error=%s;", err.Error())
		return &CelebrityResponse{}, err
	}

	list := make([]*Celebrity, len(op.CelebrityFaces))
	for i, c := range op.CelebrityFaces {
		var urls []string
		if len(c.Urls) != 0 {
			urls := make([]string, len(c.Urls))
			for j, url := range c.Urls {
				urls[j] = *url
			}
		}

		list[i] = &Celebrity{
			Name:            *c.Name,
			MatchConfidence: *c.MatchConfidence,
			Urls:            urls,
		}
	}

	return &CelebrityResponse{
		List: list,
	}, nil
}

// GetCelebrityInfoByID gets celebrities info by id.
func (svc *Rekognition) GetCelebrityInfoByID(id string) (*CelebrityInfoResponse, error) {
	return svc.getCelebrityInfo(&SDK.GetCelebrityInfoInput{
		Id: pointers.String(id),
	})
}

// getCelebrityInfo executes `GetCelebrityInfo` operation and gets celebrities info.
func (svc *Rekognition) getCelebrityInfo(input *SDK.GetCelebrityInfoInput) (*CelebrityInfoResponse, error) {
	op, err := svc.client.GetCelebrityInfo(input)
	if err != nil {
		svc.Errorf("error on `GetCelebrityInfo` operation; error=%s;", err.Error())
		return &CelebrityInfoResponse{}, err
	}

	return NewCelebrityInfoResponseFromAWSOutput(op), nil
}

// CompareFacesFromLocalFile compares faces from local image files.
func (svc *Rekognition) CompareFacesFromLocalFile(source, target string) (*CompareFaceResponse, error) {
	sourceByte, err := ioutil.ReadFile(source)
	if err != nil {
		return nil, err
	}
	targetByte, err := ioutil.ReadFile(target)
	if err != nil {
		return nil, err
	}

	return svc.CompareFacesByBytes(sourceByte, targetByte)
}

// CompareFacesFromURL compares faces from URL.
func (svc *Rekognition) CompareFacesFromURL(source, target string) (*CompareFaceResponse, error) {
	sourceByte, err := svc.getImageFromURL(source)
	if err != nil {
		return nil, err
	}
	targetByte, err := svc.getImageFromURL(target)
	if err != nil {
		return nil, err
	}

	return svc.CompareFacesByBytes(sourceByte, targetByte)
}

// CompareFacesByBytes compares faces from image of bytes.
func (svc *Rekognition) CompareFacesByBytes(source, target []byte) (*CompareFaceResponse, error) {
	sourceImage := &SDK.Image{
		Bytes: source,
	}
	targetImage := &SDK.Image{
		Bytes: target,
	}

	input := &SDK.CompareFacesInput{
		SourceImage: sourceImage,
		TargetImage: targetImage,
	}
	return svc.compareFaces(input)
}

// CompareFaces executes `CompareFaces` operation and gets celebrities info.
func (svc *Rekognition) CompareFaces(input *SDK.CompareFacesInput) (*CompareFaceResponse, error) {
	return svc.CompareFaces(input)
}

// compareFaces executes `CompareFaces` operation and gets celebrities info.
func (svc *Rekognition) compareFaces(input *SDK.CompareFacesInput) (*CompareFaceResponse, error) {
	op, err := svc.client.CompareFaces(input)
	if err != nil {
		svc.Errorf("error on `CompareFaces` operation; error=%s;", err.Error())
		return &CompareFaceResponse{}, err
	}

	fmt.Printf("op=%#v\n", op)
	return NewCompareFaceResponseFromAWSOutput(op), nil
}

// getImageFromURL gets image data from url.
func (svc *Rekognition) getImageFromURL(url string) ([]byte, error) {
	cli := svc.httpClient
	if cli == nil {
		err := errors.New("error on `getImageFromURL`; error=`svc.httpClient is nil`;")
		svc.Errorf(err.Error())
		return nil, err
	}

	resp, err := cli.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, resp.Body)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// createS3Object creates *SDK.S3Object from bucket and name.
func createS3Object(bucket, name string, version ...string) *SDK.S3Object {
	s3object := &SDK.S3Object{
		Bucket: pointers.String(bucket),
		Name:   pointers.String(name),
	}
	if len(version) == 1 {
		s3object.Version = pointers.String(version[0])
	}
	return s3object
}
