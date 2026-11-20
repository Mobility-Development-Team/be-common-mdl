package apiutil

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/Mobility-Development-Team/be-common-mdl/response"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// SummarizeValidationErrors converts a binding validation failure into a
// compact, user-readable summary such as "name is required; remark must be
// at most 200 characters". Field names prefer the JSON tags of obj so the
// message shows the names clients actually send. It returns an empty string
// when err is not a validation error (for example malformed JSON), so the
// caller can fall back to a generic message.
func SummarizeValidationErrors(obj interface{}, err error) string {
	var ve validator.ValidationErrors
	if !errors.As(err, &ve) {
		return ""
	}
	parts := make([]string, 0, len(ve))
	for _, fe := range ve {
		parts = append(parts, describeFieldError(obj, fe))
	}
	return strings.Join(parts, "; ")
}

// BindAndValidate binds the JSON request body into obj and validates it.
// On failure it writes the unified error response built from msg (the caller
// supplies the validation message code whose text carries one %s verb for the
// summary), aborts the context and returns false so the handler can stop
// without running any further logic.
func BindAndValidate(c *gin.Context, obj interface{}, msg response.Message) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		summary := SummarizeValidationErrors(obj, err)
		if summary == "" {
			summary = "invalid request body"
		}
		GenerateResponse(c, nil, msg, summary)
		return false
	}
	return true
}

// describeFieldError renders a single validation failure in a form suitable
// for the API message shown to end users.
func describeFieldError(obj interface{}, fe validator.FieldError) string {
	field := jsonFieldName(derefStructType(reflect.TypeOf(obj)), fe.Field())
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "max":
		return fmt.Sprintf("%s must be at most %s characters", field, fe.Param())
	default:
		return fmt.Sprintf("%s failed %s validation", field, fe.Tag())
	}
}

// jsonFieldName resolves the JSON field name for a Go field name. Struct
// field promotion resolves fields inherited through embedded structs.
func jsonFieldName(t reflect.Type, goField string) string {
	if t == nil {
		return goField
	}
	if sf, ok := t.FieldByName(goField); ok {
		if name, _, _ := strings.Cut(sf.Tag.Get("json"), ","); name != "" && name != "-" {
			return name
		}
	}
	return goField
}

// derefStructType dereferences pointers, slices and arrays down to the
// underlying struct type, returning nil for anything that is not a struct.
func derefStructType(t reflect.Type) reflect.Type {
	for t != nil && (t.Kind() == reflect.Ptr || t.Kind() == reflect.Slice || t.Kind() == reflect.Array) {
		t = t.Elem()
	}
	if t == nil || t.Kind() != reflect.Struct {
		return nil
	}
	return t
}
