package rekognition

import (
	"fmt"

	SDK "github.com/aws/aws-sdk-go/service/rekognition"
)

// FaceDetailResponse is response data for `DetectFaces` operation.
type FaceDetailResponse struct {
	List        []*FaceDetail
	Orientation string
}

// FaceDetail has detailed face data. (wrapper struct for *SDK.FaceDetail)
type FaceDetail struct {
	// estimated age range.
	HasAgeRange  bool
	AgeRangeHigh int
	AgeRangeLow  int

	HasBeard        bool
	BeardConfidence float64

	HasMustache        bool
	MustacheConfidence float64

	HasEyeglasses        bool
	EyeglassesConfidence float64

	HasSunglasses        bool
	SunglassesConfidence float64

	HasEyesOpen        bool
	EyesOpenConfidence float64

	HasMouthOpen        bool
	MouthOpenConfidence float64

	HasSmile        bool
	SmileConfidence float64

	HasPose   bool
	PosePitch float64
	PoseRoll  float64
	PoseYaw   float64

	HasImageQuality        bool
	ImageQualityBrightness float64
	ImageQualitySharpness  float64

	HasGender        bool
	Gender           string
	GenderConfidence float64

	HasFace        bool
	FaceConfidence float64
	BoundingHeight float64
	BoundingWidth  float64
	BoundingTop    float64
	BoundingLeft   float64

	HasEmotion bool
	Emotions   []Emotion

	HasLandmark bool
	Landmarks   []Landmark
}

// NewFaceDetailFromAWSFaceDetail creates FaceDetail from *SDK.FaceDetail.
func NewFaceDetailFromAWSFaceDetail(face *SDK.FaceDetail) *FaceDetail {
	f := &FaceDetail{}
	if face.AgeRange != nil {
		f.HasAgeRange = true
		f.AgeRangeHigh = int(*face.AgeRange.High)
		f.AgeRangeLow = int(*face.AgeRange.Low)
	}

	if face.Beard != nil {
		f.HasBeard = *face.Beard.Value
		f.BeardConfidence = *face.Beard.Confidence
	}

	if face.Mustache != nil {
		f.HasMustache = *face.Mustache.Value
		f.MustacheConfidence = *face.Mustache.Confidence
	}

	if face.Eyeglasses != nil {
		f.HasEyeglasses = *face.Eyeglasses.Value
		f.EyeglassesConfidence = *face.Eyeglasses.Confidence
	}

	if face.Sunglasses != nil {
		f.HasSunglasses = *face.Sunglasses.Value
		f.SunglassesConfidence = *face.Sunglasses.Confidence
	}

	if face.EyesOpen != nil {
		f.HasEyesOpen = *face.EyesOpen.Value
		f.EyesOpenConfidence = *face.EyesOpen.Confidence
	}

	if face.MouthOpen != nil {
		f.HasMouthOpen = *face.MouthOpen.Value
		f.MouthOpenConfidence = *face.MouthOpen.Confidence
	}

	if face.Smile != nil {
		f.HasSmile = *face.Smile.Value
		f.SmileConfidence = *face.Smile.Confidence
	}

	if face.Pose != nil {
		f.HasPose = true
		f.PosePitch = *face.Pose.Pitch
		f.PoseRoll = *face.Pose.Roll
		f.PoseYaw = *face.Pose.Yaw
	}

	if face.Quality != nil {
		f.HasImageQuality = true
		f.ImageQualityBrightness = *face.Quality.Brightness
		f.ImageQualitySharpness = *face.Quality.Sharpness
	}

	if face.Gender != nil {
		f.HasGender = true
		f.Gender = *face.Gender.Value
		f.GenderConfidence = *face.Gender.Confidence
	}

	if face.Confidence != nil {
		f.FaceConfidence = *face.Confidence
	}
	if face.BoundingBox != nil {
		f.HasFace = true
		f.BoundingHeight = *face.BoundingBox.Height
		f.BoundingWidth = *face.BoundingBox.Width
		f.BoundingTop = *face.BoundingBox.Top
		f.BoundingLeft = *face.BoundingBox.Left
	}

	if len(face.Emotions) != 0 {
		f.HasEmotion = true
		f.Emotions = make([]Emotion, len(face.Emotions))
		for i, emotion := range face.Emotions {
			f.Emotions[i] = Emotion{
				Type:       *emotion.Type,
				Confidence: *emotion.Confidence,
			}
		}
	}

	if len(face.Landmarks) != 0 {
		f.HasLandmark = true
		f.Landmarks = make([]Landmark, len(face.Landmarks))
		for i, landmark := range face.Landmarks {
			f.Landmarks[i] = Landmark{
				Type: *landmark.Type,
				X:    *landmark.X,
				Y:    *landmark.Y,
			}
		}
	}

	return f
}

// NewFaceDetailFromAWSComparedFace creates FaceDetail from *SDK.ComparedFace.
func NewFaceDetailFromAWSComparedFace(face *SDK.ComparedFace) *FaceDetail {
	if face == nil {
		return nil
	}

	f := &FaceDetail{}
	if face.Pose != nil {
		f.HasPose = true
		f.PosePitch = *face.Pose.Pitch
		f.PoseRoll = *face.Pose.Roll
		f.PoseYaw = *face.Pose.Yaw
	}

	if face.Quality != nil {
		f.HasImageQuality = true
		f.ImageQualityBrightness = *face.Quality.Brightness
		f.ImageQualitySharpness = *face.Quality.Sharpness
	}

	if face.Confidence != nil {
		f.FaceConfidence = *face.Confidence
	}
	if face.BoundingBox != nil {
		f.HasFace = true
		f.BoundingHeight = *face.BoundingBox.Height
		f.BoundingWidth = *face.BoundingBox.Width
		f.BoundingTop = *face.BoundingBox.Top
		f.BoundingLeft = *face.BoundingBox.Left
	}

	if len(face.Landmarks) != 0 {
		f.HasLandmark = true
		f.Landmarks = make([]Landmark, len(face.Landmarks))
		for i, landmark := range face.Landmarks {
			f.Landmarks[i] = Landmark{
				Type: *landmark.Type,
				X:    *landmark.X,
				Y:    *landmark.Y,
			}
		}
	}
	return f
}

// Emotion has emotion data.
type Emotion struct {
	Type       string
	Confidence float64
}

// Landmark has landmark data.
type Landmark struct {
	Type string
	X    float64
	Y    float64
}

// CompareFaceResponse is response data for `CompareFaces` operation.
type CompareFaceResponse struct {
	MatchedList            []*FaceDetail
	UnmatchedList          []*FaceDetail
	SourceImage            *FaceDetail
	SourceImageOrientation string
	TargetImageOrientation string
}

// CompareFacesMatch has detailed face data and similarity. (wrapper struct for *SDK.ComparedFace)
type CompareFacesMatch struct {
	*FaceDetail
	Similarity float64
}

// NewCompareFaceResponseFromAWSOutput creates CompareFaceResponse from *SDK.CompareFacesOutput.
func NewCompareFaceResponseFromAWSOutput(output *SDK.CompareFacesOutput) *CompareFaceResponse {
	resp := &CompareFaceResponse{}
	if output.SourceImageOrientationCorrection != nil {
		resp.SourceImageOrientation = *output.SourceImageOrientationCorrection
	}
	if output.TargetImageOrientationCorrection != nil {
		resp.TargetImageOrientation = *output.TargetImageOrientationCorrection
	}

	if output.SourceImageFace != nil {
		f := &FaceDetail{}
		src := output.SourceImageFace
		if src.Confidence != nil {
			f.HasFace = true
			f.FaceConfidence = *src.Confidence
		}
		if src.BoundingBox != nil {
			f.BoundingHeight = *src.BoundingBox.Height
			f.BoundingWidth = *src.BoundingBox.Width
			f.BoundingTop = *src.BoundingBox.Top
			f.BoundingLeft = *src.BoundingBox.Left
		}

		resp.SourceImage = f
	}

	if len(output.FaceMatches) != 0 {
		matchedFaces := make([]*CompareFacesMatch, len(output.FaceMatches))
		for i, match := range output.FaceMatches {
			c := &CompareFacesMatch{}
			if match.Similarity != nil {
				c.Similarity = *match.Similarity
			}
			if match.Face != nil {
				c.FaceDetail = NewFaceDetailFromAWSComparedFace(match.Face)
			}
			matchedFaces[i] = c
		}
	}

	if len(output.UnmatchedFaces) != 0 {
		unmatchedFaces := make([]*FaceDetail, len(output.UnmatchedFaces))
		for i, unmatch := range output.UnmatchedFaces {
			fmt.Printf("unmatch=%#v\n", unmatch)
			unmatchedFaces[i] = NewFaceDetailFromAWSComparedFace(unmatch)
		}
	}

	return resp
}

// LabelResponse is response data for `DetectLabels` operation.
type LabelResponse struct {
	List        []*Label
	Orientation string
}

// Label has label data.
type Label struct {
	Name       string
	Confidence float64
}

// NewLabelFromAWSLabel creates Label from *SDK.Label.
func NewLabelFromAWSLabel(label *SDK.Label) *Label {
	l := &Label{}
	if label.Name != nil {
		l.Name = *label.Name
	}
	if label.Confidence != nil {
		l.Confidence = *label.Confidence
	}
	return l
}

// ModerationLabelResponse is response data for `DetectModerationLabels` operation.
type ModerationLabelResponse struct {
	List []*ModerationLabel
}

// ModerationLabel has moderation label data.
type ModerationLabel struct {
	Name       string
	ParentName string
	Confidence float64
}

// NewModerationLabelFromAWSModerationLabel creates ModerationLabel from *SDK.ModerationLabel.
func NewModerationLabelFromAWSModerationLabel(label *SDK.ModerationLabel) *ModerationLabel {
	l := &ModerationLabel{}
	if label.Name != nil {
		l.Name = *label.Name
	}
	if label.ParentName != nil {
		l.ParentName = *label.ParentName
	}
	if label.Confidence != nil {
		l.Confidence = *label.Confidence
	}
	return l
}

// CelebrityResponse is response data for `RecognizeCelebrities` operation.
type CelebrityResponse struct {
	List []*Celebrity
}

// Celebrity has celebrity data.
type Celebrity struct {
	Name            string
	MatchConfidence float64
	Urls            []string
}

// CelebrityInfoResponse is response data for `GetCelebrityInfo` operation.
type CelebrityInfoResponse struct {
	Name string
	URLs []string
}

// NewCelebrityInfoResponseFromAWSOutput creates CelebrityInfoResponse from *SDK.GetCelebrityInfoOutput.
func NewCelebrityInfoResponseFromAWSOutput(info *SDK.GetCelebrityInfoOutput) *CelebrityInfoResponse {
	if info == nil {
		return nil
	}

	c := &CelebrityInfoResponse{}
	if info.Name != nil {
		c.Name = *info.Name
	}

	if len(info.Urls) != 0 {
		urls := make([]string, len(info.Urls))
		for i, url := range info.Urls {
			urls[i] = *url
		}
		c.URLs = urls
	}
	return c
}

// Face has face data in collection. (wrapper struct for *SDK.Face)
type Face struct {
	HasFace        bool
	FaceConfidence float64
	BoundingHeight float64
	BoundingWidth  float64
	BoundingTop    float64
	BoundingLeft   float64

	ExternalImageID string
	FaceID          string
	ImageID         string
}

// NewFaceFromAWSFace creates Face from *SDK.Face.
func NewFaceFromAWSFace(face *SDK.Face) *Face {
	f := &Face{}
	if face.Confidence != nil {
		f.FaceConfidence = *face.Confidence
	}
	if face.BoundingBox != nil {
		f.HasFace = true
		f.BoundingHeight = *face.BoundingBox.Height
		f.BoundingWidth = *face.BoundingBox.Width
		f.BoundingTop = *face.BoundingBox.Top
		f.BoundingLeft = *face.BoundingBox.Left
	}

	if face.ExternalImageId != nil {
		f.ExternalImageID = *face.ExternalImageId
	}

	if face.FaceId != nil {
		f.FaceID = *face.FaceId
	}

	if face.ImageId != nil {
		f.ImageID = *face.ImageId
	}
	return f
}

// IndexFacesResponse is response data for `IndexFaces` operation.
type IndexFacesResponse struct {
	List        []*FaceRecord
	Orientation string
}

// FaceRecord has face and detail data. (wrapper struct for *SDK.FaceRecord)
type FaceRecord struct {
	*Face
	*FaceDetail
}

// NewIndexFacesResponseFromAWSOutput creates IndexFacesResponse from *SDK.IndexFacesOutput.
func NewIndexFacesResponseFromAWSOutput(output *SDK.IndexFacesOutput) *IndexFacesResponse {
	resp := &IndexFacesResponse{}
	if output.OrientationCorrection != nil {
		resp.Orientation = *output.OrientationCorrection
	}
	if len(output.FaceRecords) == 0 {
		return resp
	}

	list := make([]*FaceRecord, len(output.FaceRecords))
	for i, r := range output.FaceRecords {
		list[i] = &FaceRecord{
			Face:       NewFaceFromAWSFace(r.Face),
			FaceDetail: NewFaceDetailFromAWSFaceDetail(r.FaceDetail),
		}
	}
	resp.List = list
	return resp
}

// SearchFacesResponse is response data for `SearchFaces` operation.
type SearchFacesResponse struct {
	List []*FaceMatch

	// returned by `SearchFaces` operation
	FaceID string

	// returned by `SearchFacesByImage` operation
	HasFace        bool
	FaceConfidence float64
	BoundingHeight float64
	BoundingWidth  float64
	BoundingTop    float64
	BoundingLeft   float64
}

// FaceMatch has face data and similarity. (wrapper struct for *SDK.FaceMatch)
type FaceMatch struct {
	*Face
	Similarity float64
}

// NewSearchFacesResponseFromAWSOutput creates SearchFacesResponse from *SDK.SearchFacesOutput.
func NewSearchFacesResponseFromAWSOutput(output *SDK.SearchFacesOutput) *SearchFacesResponse {
	resp := &SearchFacesResponse{}
	if output.SearchedFaceId != nil {
		resp.FaceID = *output.SearchedFaceId
	}

	resp.List = newFaceMatchesFromAWSFaceMatches(output.FaceMatches)
	return resp
}

// NewSearchFacesResponseFromAWSOutputByImage creates SearchFacesResponse from *SDK.SearchFacesByImageOutput.
func NewSearchFacesResponseFromAWSOutputByImage(output *SDK.SearchFacesByImageOutput) *SearchFacesResponse {
	resp := &SearchFacesResponse{}
	if output.SearchedFaceConfidence != nil {
		resp.FaceConfidence = *output.SearchedFaceConfidence
	}
	if output.SearchedFaceBoundingBox != nil {
		resp.HasFace = true
		resp.BoundingHeight = *output.SearchedFaceBoundingBox.Height
		resp.BoundingWidth = *output.SearchedFaceBoundingBox.Width
		resp.BoundingTop = *output.SearchedFaceBoundingBox.Top
		resp.BoundingLeft = *output.SearchedFaceBoundingBox.Left
	}

	resp.List = newFaceMatchesFromAWSFaceMatches(output.FaceMatches)
	return resp
}

func newFaceMatchesFromAWSFaceMatches(list []*SDK.FaceMatch) []*FaceMatch {
	if len(list) == 0 {
		return nil
	}

	result := make([]*FaceMatch, len(list))
	for i, f := range list {
		match := &FaceMatch{
			Face: NewFaceFromAWSFace(f.Face),
		}
		if f.Similarity != nil {
			match.Similarity = *f.Similarity
		}
		result[i] = match
	}
	return result
}
