package rekognition

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFaceDetail_FilterFaceByConfidence(t *testing.T) {
	a := assert.New(t)

	tests := []struct {
		Threshold    float64
		ExpectedSize int
	}{
		{0.0, 11},
		{0.1, 10},
		{0.2, 9},
		{0.3, 8},
		{0.4, 7},
		{0.5, 6},
		{0.6, 5},
		{0.7, 4},
		{0.8, 3},
		{0.9, 2},
		{1.0, 1},
	}

	// nolint:gofmt
	list := []*FaceDetail{
		&FaceDetail{FaceConfidence: 0.0},
		&FaceDetail{FaceConfidence: 0.11},
		&FaceDetail{FaceConfidence: 0.22},
		&FaceDetail{FaceConfidence: 0.33},
		&FaceDetail{FaceConfidence: 0.44},
		&FaceDetail{FaceConfidence: 0.55},
		&FaceDetail{FaceConfidence: 0.66},
		&FaceDetail{FaceConfidence: 0.77},
		&FaceDetail{FaceConfidence: 0.88},
		&FaceDetail{FaceConfidence: 0.99},
		&FaceDetail{FaceConfidence: 1.0},
	}
	resp := FaceDetailResponse{
		List: list,
	}

	for i, tt := range tests {
		target := fmt.Sprintf("[#%d] %+v", i, tt)

		result := resp.FilterFaceByConfidence(tt.Threshold)
		a.Equal(tt.ExpectedSize, len(result), target)
	}
}

func TestFaceDetail_FilterFaceBySize(t *testing.T) {
	a := assert.New(t)

	tests := []struct {
		Threshold    float64
		ExpectedSize int
	}{
		{0.0, 11},
		{0.1, 10},
		{0.2, 9},
		{0.3, 8},
		{0.4, 7},
		{0.5, 6},
		{0.6, 5},
		{0.7, 4},
		{0.8, 3},
		{0.9, 2},
		{1.0, 1},
	}

	// nolint:gofmt
	list := []*FaceDetail{
		&FaceDetail{BoundingHeight: 0.0, BoundingWidth: 0.01},
		&FaceDetail{BoundingHeight: 0.11, BoundingWidth: 0.1},
		&FaceDetail{BoundingHeight: 0.22, BoundingWidth: 0.2},
		&FaceDetail{BoundingHeight: 0.33, BoundingWidth: 0.3},
		&FaceDetail{BoundingHeight: 0.44, BoundingWidth: 0.4},
		&FaceDetail{BoundingHeight: 0.55, BoundingWidth: 0.5},
		&FaceDetail{BoundingHeight: 0.66, BoundingWidth: 0.6},
		&FaceDetail{BoundingHeight: 0.77, BoundingWidth: 0.7},
		&FaceDetail{BoundingHeight: 0.88, BoundingWidth: 0.8},
		&FaceDetail{BoundingHeight: 0.99, BoundingWidth: 0.9},
		&FaceDetail{BoundingHeight: 1.0, BoundingWidth: 1.0},
	}
	resp := FaceDetailResponse{
		List: list,
	}

	for i, tt := range tests {
		target := fmt.Sprintf("[#%d] %+v", i, tt)

		result := resp.FilterFaceBySize(tt.Threshold)
		a.Equal(tt.ExpectedSize, len(result), target)
	}
}

func TestFaceDetail_FilterFaceByConfidenceAndSize(t *testing.T) {
	a := assert.New(t)

	tests := []struct {
		Confidence   float64
		Size         float64
		ExpectedSize int
	}{
		{0.0, 0.0, 11},
		{0.1, 0.1, 10},
		{0.2, 0.2, 9},
		{0.3, 0.3, 8},
		{0.4, 0.4, 7},
		{0.5, 0.5, 3},
		{0.6, 0.6, 0},
		{0.7, 0.7, 0},
		{0.8, 0.8, 0},
		{0.9, 0.9, 0},
		{1.0, 1.0, 0},
		{0.0, 0.5, 6},
		{0.1, 0.5, 6},
		{0.2, 0.5, 6},
		{0.3, 0.5, 6},
		{0.4, 0.5, 6},
		{0.6, 0.5, 0},
		{0.7, 0.5, 0},
		{0.8, 0.5, 0},
		{0.9, 0.5, 0},
		{1.0, 0.5, 0},
		{0.0, 0.9, 2},
		{0.1, 0.9, 2},
		{0.2, 0.9, 2},
		{0.3, 0.9, 2},
		{0.4, 0.9, 2},
		{0.5, 0.9, 1},
		{0.6, 0.9, 0},
		{0.7, 0.9, 0},
		{0.8, 0.9, 0},
		{1.0, 0.9, 0},
	}

	// nolint:gofmt
	list := []*FaceDetail{
		&FaceDetail{FaceConfidence: 0.49, BoundingHeight: 0.0, BoundingWidth: 0.01},
		&FaceDetail{FaceConfidence: 0.51, BoundingHeight: 0.11, BoundingWidth: 0.1},
		&FaceDetail{FaceConfidence: 0.49, BoundingHeight: 0.22, BoundingWidth: 0.2},
		&FaceDetail{FaceConfidence: 0.51, BoundingHeight: 0.33, BoundingWidth: 0.3},
		&FaceDetail{FaceConfidence: 0.49, BoundingHeight: 0.44, BoundingWidth: 0.4},
		&FaceDetail{FaceConfidence: 0.51, BoundingHeight: 0.55, BoundingWidth: 0.5},
		&FaceDetail{FaceConfidence: 0.49, BoundingHeight: 0.66, BoundingWidth: 0.6},
		&FaceDetail{FaceConfidence: 0.51, BoundingHeight: 0.77, BoundingWidth: 0.7},
		&FaceDetail{FaceConfidence: 0.49, BoundingHeight: 0.88, BoundingWidth: 0.8},
		&FaceDetail{FaceConfidence: 0.51, BoundingHeight: 0.99, BoundingWidth: 0.9},
		&FaceDetail{FaceConfidence: 0.49, BoundingHeight: 1.0, BoundingWidth: 1.0},
	}
	resp := FaceDetailResponse{
		List: list,
	}

	for i, tt := range tests {
		target := fmt.Sprintf("[#%d] %+v", i, tt)

		result := resp.FilterFaceByConfidenceAndSize(tt.Confidence, tt.Size)
		a.Equal(tt.ExpectedSize, len(result), target)
	}
}

func TestFaceDetail_FilterSmileByConfidence(t *testing.T) {
	a := assert.New(t)

	tests := []struct {
		Threshold    float64
		ExpectedSize int
	}{
		{0.0, 11},
		{0.1, 10},
		{0.2, 9},
		{0.3, 8},
		{0.4, 7},
		{0.5, 6},
		{0.6, 5},
		{0.7, 4},
		{0.8, 3},
		{0.9, 2},
		{1.0, 1},
	}

	// nolint:gofmt
	list := []*FaceDetail{
		&FaceDetail{HasSmile: true, SmileConfidence: 0.0},
		&FaceDetail{HasSmile: true, SmileConfidence: 0.11},
		&FaceDetail{HasSmile: true, SmileConfidence: 0.22},
		&FaceDetail{HasSmile: true, SmileConfidence: 0.33},
		&FaceDetail{HasSmile: true, SmileConfidence: 0.44},
		&FaceDetail{HasSmile: true, SmileConfidence: 0.55},
		&FaceDetail{HasSmile: true, SmileConfidence: 0.66},
		&FaceDetail{HasSmile: true, SmileConfidence: 0.77},
		&FaceDetail{HasSmile: true, SmileConfidence: 0.88},
		&FaceDetail{HasSmile: true, SmileConfidence: 0.99},
		&FaceDetail{HasSmile: true, SmileConfidence: 1.0},
		&FaceDetail{HasSmile: false, SmileConfidence: 0.0},
		&FaceDetail{HasSmile: false, SmileConfidence: 0.11},
		&FaceDetail{HasSmile: false, SmileConfidence: 0.22},
		&FaceDetail{HasSmile: false, SmileConfidence: 0.33},
		&FaceDetail{HasSmile: false, SmileConfidence: 0.44},
		&FaceDetail{HasSmile: false, SmileConfidence: 0.55},
		&FaceDetail{HasSmile: false, SmileConfidence: 0.66},
		&FaceDetail{HasSmile: false, SmileConfidence: 0.77},
		&FaceDetail{HasSmile: false, SmileConfidence: 0.88},
		&FaceDetail{HasSmile: false, SmileConfidence: 0.99},
		&FaceDetail{HasSmile: false, SmileConfidence: 1.0},
	}
	resp := FaceDetailResponse{
		List: list,
	}

	for i, tt := range tests {
		target := fmt.Sprintf("[#%d] %+v", i, tt)

		result := resp.FilterSmileByConfidence(tt.Threshold)
		a.Equal(tt.ExpectedSize, len(result), target)
	}
}

func TestFaceDetail_IsFaceConfidenceGTE(t *testing.T) {
	a := assert.New(t)

	tests := []struct {
		Expected       bool
		FaceConfidence float64
		Threshold      float64
	}{
		{true, 0, 0},
		{true, 1, 0},
		{true, 0.5, 0.4},
		{true, 0.5, 0.4999999},
		{true, 0.49999999, 0.4999999},
		{true, 0.499999999, 0.49999998},
		{true, 0.4999999901, 0.4999999},
		{false, 0, 0.00001},
		{false, 0, 1},
		{false, 0.4999999, 0.5},
		{false, 0.4999998, 0.4999999},
	}

	for i, tt := range tests {
		target := fmt.Sprintf("[#%d] %+v", i, tt)

		f := FaceDetail{
			FaceConfidence: tt.FaceConfidence,
		}
		a.Equal(tt.Expected, f.IsFaceConfidenceGTE(tt.Threshold), target)
	}
}

func TestFaceDetail_IsBoundingGTE(t *testing.T) {
	a := assert.New(t)

	tests := []struct {
		Expected   bool
		FaceHeight float64
		FaceWidth  float64
		Size       float64
	}{
		{true, 0, 0, 0},
		{true, 0.5, 0.5, 0},
		{true, 0.5, 0.5, 0.4},
		{true, 0.5, 0.5, 0.49},
		{true, 0.5, 0.5, 0.499999},
		{true, 0.49999999, 0.49999999, 0.4999999},
		{true, 0.499999999, 0.5, 0.49999998},
		{true, 0.4999999901, 0.5, 0.4999999},
		{true, 0.5, 0.499999999, 0.49999998},
		{true, 0.5, 0.4999999901, 0.4999999},
		{false, 0, 0, 0.00001},
		{false, 0, 0, 1},
		{false, 0.4999999, 0.4999999, 0.5},
		{false, 0.5, 0.4999999, 0.5},
		{false, 0.4999999, 0.5, 0.5},
		{false, 0.4999998, 0.4999998, 0.4999999},
		{false, 0.4999998, 0.4999999, 0.4999999},
		{false, 0.4999999, 0.4999998, 0.4999999},
	}

	for i, tt := range tests {
		target := fmt.Sprintf("[#%d] %+v", i, tt)

		f := FaceDetail{
			BoundingHeight: tt.FaceHeight,
			BoundingWidth:  tt.FaceWidth,
		}
		a.Equal(tt.Expected, f.IsBoundingGTE(tt.Size), target)
	}
}

func TestFaceDetail_IsSmileConfidenceGTE(t *testing.T) {
	a := assert.New(t)

	hasSmile := true
	tests := []struct {
		Expected        bool
		HasSmile        bool
		SmileConfidence float64
		Threshold       float64
	}{
		{true, hasSmile, 0, 0},
		{true, hasSmile, 1, 0},
		{true, hasSmile, 0.5, 0.4},
		{true, hasSmile, 0.5, 0.4999999},
		{true, hasSmile, 0.49999999, 0.4999999},
		{true, hasSmile, 0.499999999, 0.49999998},
		{true, hasSmile, 0.4999999901, 0.4999999},
		{false, hasSmile, 0, 0.00001},
		{false, hasSmile, 0, 1},
		{false, hasSmile, 0.4999999, 0.5},
		{false, hasSmile, 0.4999998, 0.4999999},
		{false, !hasSmile, 0, 0},
		{false, !hasSmile, 1, 0},
		{false, !hasSmile, 0.5, 0.4},
		{false, !hasSmile, 0.5, 0.4999999},
		{false, !hasSmile, 0.49999999, 0.4999999},
		{false, !hasSmile, 0.499999999, 0.49999998},
		{false, !hasSmile, 0.4999999901, 0.4999999},
	}

	for i, tt := range tests {
		target := fmt.Sprintf("[#%d] %+v", i, tt)

		f := FaceDetail{
			HasSmile:        tt.HasSmile,
			SmileConfidence: tt.SmileConfidence,
		}
		a.Equal(tt.Expected, f.IsSmileConfidenceGTE(tt.Threshold), target)
	}
}
