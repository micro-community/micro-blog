// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: tags.proto

package tags

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/golang/protobuf/ptypes"
)

// ensure the imports are used
var (
	_ = bytes.MinRead
	_ = errors.New("")
	_ = fmt.Print
	_ = utf8.UTFMax
	_ = (*regexp.Regexp)(nil)
	_ = (*strings.Reader)(nil)
	_ = net.IPv4len
	_ = time.Duration(0)
	_ = (*url.URL)(nil)
	_ = (*mail.Address)(nil)
	_ = ptypes.DynamicAny{}
)

// define the regex for a UUID once up-front
var _tags_uuidPattern = regexp.MustCompile("^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$")

// Validate checks the field values on Tag with the rules defined in the proto
// definition for this message. If any rules are violated, an error is returned.
func (m *Tag) Validate() error {
	if m == nil {
		return nil
	}

	// no validation rules for Type

	// no validation rules for Slug

	// no validation rules for Title

	// no validation rules for Description

	// no validation rules for Count

	return nil
}

// TagValidationError is the validation error returned by Tag.Validate if the
// designated constraints aren't met.
type TagValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e TagValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e TagValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e TagValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e TagValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e TagValidationError) ErrorName() string { return "TagValidationError" }

// Error satisfies the builtin error interface
func (e TagValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sTag.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = TagValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = TagValidationError{}

// Validate checks the field values on AddRequest with the rules defined in the
// proto definition for this message. If any rules are violated, an error is returned.
func (m *AddRequest) Validate() error {
	if m == nil {
		return nil
	}

	// no validation rules for ResourceID

	// no validation rules for Type

	// no validation rules for ResourceCreated

	return nil
}

// AddRequestValidationError is the validation error returned by
// AddRequest.Validate if the designated constraints aren't met.
type AddRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e AddRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e AddRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e AddRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e AddRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e AddRequestValidationError) ErrorName() string { return "AddRequestValidationError" }

// Error satisfies the builtin error interface
func (e AddRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sAddRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = AddRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = AddRequestValidationError{}

// Validate checks the field values on AddResponse with the rules defined in
// the proto definition for this message. If any rules are violated, an error
// is returned.
func (m *AddResponse) Validate() error {
	if m == nil {
		return nil
	}

	return nil
}

// AddResponseValidationError is the validation error returned by
// AddResponse.Validate if the designated constraints aren't met.
type AddResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e AddResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e AddResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e AddResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e AddResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e AddResponseValidationError) ErrorName() string { return "AddResponseValidationError" }

// Error satisfies the builtin error interface
func (e AddResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sAddResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = AddResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = AddResponseValidationError{}

// Validate checks the field values on RemoveRequest with the rules defined in
// the proto definition for this message. If any rules are violated, an error
// is returned.
func (m *RemoveRequest) Validate() error {
	if m == nil {
		return nil
	}

	// no validation rules for ResourceID

	// no validation rules for Type

	return nil
}

// RemoveRequestValidationError is the validation error returned by
// RemoveRequest.Validate if the designated constraints aren't met.
type RemoveRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e RemoveRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e RemoveRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e RemoveRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e RemoveRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e RemoveRequestValidationError) ErrorName() string { return "RemoveRequestValidationError" }

// Error satisfies the builtin error interface
func (e RemoveRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sRemoveRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = RemoveRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = RemoveRequestValidationError{}

// Validate checks the field values on RemoveResponse with the rules defined in
// the proto definition for this message. If any rules are violated, an error
// is returned.
func (m *RemoveResponse) Validate() error {
	if m == nil {
		return nil
	}

	return nil
}

// RemoveResponseValidationError is the validation error returned by
// RemoveResponse.Validate if the designated constraints aren't met.
type RemoveResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e RemoveResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e RemoveResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e RemoveResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e RemoveResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e RemoveResponseValidationError) ErrorName() string { return "RemoveResponseValidationError" }

// Error satisfies the builtin error interface
func (e RemoveResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sRemoveResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = RemoveResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = RemoveResponseValidationError{}

// Validate checks the field values on UpdateRequest with the rules defined in
// the proto definition for this message. If any rules are violated, an error
// is returned.
func (m *UpdateRequest) Validate() error {
	if m == nil {
		return nil
	}

	// no validation rules for Type

	// no validation rules for Description

	return nil
}

// UpdateRequestValidationError is the validation error returned by
// UpdateRequest.Validate if the designated constraints aren't met.
type UpdateRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e UpdateRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e UpdateRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e UpdateRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e UpdateRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e UpdateRequestValidationError) ErrorName() string { return "UpdateRequestValidationError" }

// Error satisfies the builtin error interface
func (e UpdateRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sUpdateRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = UpdateRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = UpdateRequestValidationError{}

// Validate checks the field values on UpdateResponse with the rules defined in
// the proto definition for this message. If any rules are violated, an error
// is returned.
func (m *UpdateResponse) Validate() error {
	if m == nil {
		return nil
	}

	return nil
}

// UpdateResponseValidationError is the validation error returned by
// UpdateResponse.Validate if the designated constraints aren't met.
type UpdateResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e UpdateResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e UpdateResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e UpdateResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e UpdateResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e UpdateResponseValidationError) ErrorName() string { return "UpdateResponseValidationError" }

// Error satisfies the builtin error interface
func (e UpdateResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sUpdateResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = UpdateResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = UpdateResponseValidationError{}

// Validate checks the field values on ListRequest with the rules defined in
// the proto definition for this message. If any rules are violated, an error
// is returned.
func (m *ListRequest) Validate() error {
	if m == nil {
		return nil
	}

	// no validation rules for ResourceID

	// no validation rules for Type

	// no validation rules for MinCount

	// no validation rules for MaxCount

	return nil
}

// ListRequestValidationError is the validation error returned by
// ListRequest.Validate if the designated constraints aren't met.
type ListRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ListRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ListRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ListRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ListRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ListRequestValidationError) ErrorName() string { return "ListRequestValidationError" }

// Error satisfies the builtin error interface
func (e ListRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sListRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ListRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ListRequestValidationError{}

// Validate checks the field values on ListResponse with the rules defined in
// the proto definition for this message. If any rules are violated, an error
// is returned.
func (m *ListResponse) Validate() error {
	if m == nil {
		return nil
	}

	for idx, item := range m.GetTags() {
		_, _ = idx, item

		if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return ListResponseValidationError{
					field:  fmt.Sprintf("Tags[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	return nil
}

// ListResponseValidationError is the validation error returned by
// ListResponse.Validate if the designated constraints aren't met.
type ListResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ListResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ListResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ListResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ListResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ListResponseValidationError) ErrorName() string { return "ListResponseValidationError" }

// Error satisfies the builtin error interface
func (e ListResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sListResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ListResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ListResponseValidationError{}
