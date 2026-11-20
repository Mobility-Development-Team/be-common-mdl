package apiutil

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Mobility-Development-Team/be-common-mdl/response"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type summarizeSampleEmbedded struct {
	Inner string `json:"inner" binding:"required"`
}

type summarizeSampleReq struct {
	summarizeSampleEmbedded
	Name   string `json:"name" binding:"required"`
	RefKey string `json:"refKey" binding:"required"`
	Plain  string `binding:"required"`
	Remark string `json:"remark" binding:"max=5"`
}

func TestSummarizeValidationErrors(t *testing.T) {
	// Fails on embedded required, top-level required, required without a
	// JSON tag and max length, all in one struct.
	req := summarizeSampleReq{
		RefKey: "present",
		Remark: "too long remark",
	}
	v := validator.New()
	v.SetTagName("binding") // gin registers the "binding" tag name.
	err := v.Struct(req)

	want := "inner is required; name is required; Plain is required; remark must be at most 5 characters"
	if got := SummarizeValidationErrors(req, err); got != want {
		t.Errorf("SummarizeValidationErrors() = %q, want %q", got, want)
	}

	// Valid struct produces no validation error.
	if err := v.Struct(summarizeSampleReq{
		summarizeSampleEmbedded: summarizeSampleEmbedded{Inner: "a"},
		Name: "n", RefKey: "r", Plain: "p", Remark: "ok",
	}); err != nil {
		t.Errorf("unexpected validation error for valid struct: %v", err)
	}

	// Non validation errors fall back to an empty summary.
	if got := SummarizeValidationErrors(req, errors.New("some other error")); got != "" {
		t.Errorf("SummarizeValidationErrors() = %q, want empty for non validation error", got)
	}
}

type bindSampleEmbedded struct {
	Inner string `json:"inner" binding:"required"`
}

type bindSampleReq struct {
	bindSampleEmbedded
	Name   string `json:"name" binding:"required"`
	Remark string `json:"remark" binding:"max=5"`
}

var bindTestMsg = response.NewMessage(http.StatusBadRequest, "40001", "Invalid parameters: %s")

func newBindTestContext(t *testing.T, body string) (*gin.Context, *httptest.ResponseRecorder) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")
	return c, w
}

func TestBindAndValidate(t *testing.T) {
	tests := []struct {
		name             string
		body             string
		wantOK           bool
		wantAborted      bool
		wantBodySnippets []string
	}{
		{
			name:   "valid body passes",
			body:   `{"inner":"a","name":"n","remark":"ok"}`,
			wantOK: true,
		},
		{
			name:        "required and max failures are reported",
			body:        `{"inner":"","name":"","remark":"too long"}`,
			wantOK:      false,
			wantAborted: true,
			wantBodySnippets: []string{
				`"msgCode":"40001"`,
				"Invalid parameters: inner is required; name is required; remark must be at most 5 characters",
			},
		},
		{
			name:        "malformed json falls back to generic summary",
			body:        `{"inner":`,
			wantOK:      false,
			wantAborted: true,
			wantBodySnippets: []string{
				`"msgCode":"40001"`,
				"Invalid parameters: invalid request body",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, w := newBindTestContext(t, tt.body)
			var req bindSampleReq
			got := BindAndValidate(c, &req, bindTestMsg)
			if got != tt.wantOK {
				t.Errorf("BindAndValidate() = %v, want %v", got, tt.wantOK)
				return
			}
			if c.IsAborted() != tt.wantAborted {
				t.Errorf("context aborted = %v, want %v", c.IsAborted(), tt.wantAborted)
			}
			for _, snippet := range tt.wantBodySnippets {
				if !strings.Contains(w.Body.String(), snippet) {
					t.Errorf("response body = %q, want it to contain %q", w.Body.String(), snippet)
				}
			}
		})
	}
}
