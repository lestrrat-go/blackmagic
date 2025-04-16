package blackmagic_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/lestrrat-go/blackmagic"
	"github.com/stretchr/testify/require"
)

func TestAssignment(t *testing.T) {
	testcases := []struct {
		Name        string
		Error       bool
		ErrorCheck  func(error) error
		Value       interface{}
		Destination func() interface{}
	}{
		{
			Name:  `empty struct`,
			Error: false,
			Value: struct{}{},
			Destination: func() interface{} {
				var v interface{}
				return &v
			},
		},
		{
			Name:  `non pointer destination`,
			Error: true,
			Value: &struct{}{},
		},
		{
			Name:  `assign empty struct to int`,
			Error: true,
			Value: &struct{}{},
			Destination: func() interface{} {
				var v int
				return &v
			},
		},
		{
			Name:  `source is nil`,
			Error: true,
			Value: nil,
			Destination: func() interface{} {
				var v interface{}
				return &v
			},
			ErrorCheck: func(err error) error {
				if !errors.Is(err, blackmagic.InvalidValueError()) {
					return fmt.Errorf(`error should be InvalidValueError, but got %v`, err)
				}
				return nil
			},
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			var dst interface{}
			if dstFunc := tc.Destination; dstFunc != nil {
				dst = dstFunc()
			}
			err := blackmagic.AssignIfCompatible(dst, tc.Value)
			if tc.Error {
				require.Error(t, err, `blackmagic.AssignIfCompatible should fail`)
				if check := tc.ErrorCheck; check != nil {
					if checkErr := check(err); checkErr != nil {
						require.NoError(t, checkErr, `check function should succeed`)
					}
				}

			} else {
				require.NoError(t, err, `blackmagic.AssignIfCompatible should succeed`)
			}
		})
	}
}

func TestAssignOptionalField(t *testing.T) {
	var f struct {
		Foo *string
		Bar *int
	}

	require.NoError(t, blackmagic.AssignOptionalField(&f.Foo, "Hello"), `blackmagic.AssignOptionalField should succeed`)
	require.Equal(t, *(f.Foo), "Hello")
	require.NoError(t, blackmagic.AssignOptionalField(&f.Bar, 1), `blackmagic.AssignOptionalField should succeed`)
	require.Equal(t, *(f.Bar), 1)
}
