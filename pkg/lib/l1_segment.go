// Copyright 2020 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package lib

import (
	"reflect"
	"strings"
	"unicode/utf8"

	"github.com/moov-io/metro2/pkg/utils"
)

// L1Segment holds the l1 segment
type L1Segment struct {
	// Contains a constant of L1.
	SegmentIdentifier string `json:"segmentIdentifier"  validate:"required"`

	// Contains a code representing the change being reported.
	// Values available:
	//  1 = Consumer Account Number Change ONLY
	//  2 = Identification Number Change ONLY
	//  3 = Consumer Account Number AND Identification Number Change
	ChangeIndicator int `json:"changeIndicator"  validate:"required"`

	// Contains the new Account Number assigned to this account.
	// Do not include embedded blanks or special characters.
	// If field 2 = 2, this field should be blank filled.
	NewConsumerAccountNumber string `json:"newConsumerAccountNumber,omitempty"`

	// Contains the new Identification Number assigned to this account.
	// Do not include embedded blanks or special characters.
	// If field 2 = 1, this field should be blank filled.
	NewIdentificationNumber string `json:"balloonPaymentDueDate,omitempty"`

	converter
	validator
}

// Name returns name of L1 segment
func (s *L1Segment) Name() string {
	return L1SegmentName
}

// Parse takes the input record string and parses the l1 segment values
func (s *L1Segment) Parse(record string) (int, error) {
	if utf8.RuneCountInString(record) < L1SegmentLength {
		return 0, utils.ErrSegmentLength
	}

	fields := reflect.ValueOf(s).Elem()
	length, err := s.parseRecordValues(fields, l1SegmentFormat, record, &s.validator)
	if err != nil {
		return length, err
	}

	return L1SegmentLength, nil
}

// String writes the l1 segment struct to a 54 character string.
func (s *L1Segment) String() string {
	var buf strings.Builder
	specifications := s.toSpecifications(l1SegmentFormat)
	fields := reflect.ValueOf(s).Elem()
	buf.Grow(L1SegmentLength)
	for _, spec := range specifications {
		value := s.toString(spec.Field, fields.FieldByName(spec.Name))
		buf.WriteString(value)
	}

	return buf.String()
}

// Validate performs some checks on the record and returns an error if not Validated
func (s *L1Segment) Validate() error {
	return s.validateRecord(s, l1SegmentFormat)
}

// Length returns size of segment
func (s *L1Segment) Length() int {
	return L1SegmentLength
}

// validation of change indicator
func (s *L1Segment) ValidateChangeIndicator() error {
	switch s.ChangeIndicator {
	case ChangeIndicatorAccountNumber, ChangeIndicatorIdentificationNumber, ChangeIndicatorBothNumber:
		return nil
	}
	return utils.NewErrValidValue("change indicator")
}

// validation of new consumer account number
func (s *L1Segment) ValidateNewConsumerAccountNumber() error {
	if s.ChangeIndicator == ChangeIndicatorIdentificationNumber {
		if !validFilledString(s.NewConsumerAccountNumber) {
			return utils.NewErrValidValue("new consumer account number")
		}
	}
	return nil
}

// validation of new identification number
func (s *L1Segment) ValidateNewIdentificationNumber() error {
	if s.ChangeIndicator == ChangeIndicatorAccountNumber {
		if !validFilledString(s.NewIdentificationNumber) {
			return utils.NewErrValidValue("new identification number")
		}
	}
	return nil
}
