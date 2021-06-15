package bees

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBinaryConversions(t *testing.T) {
	testCases := []struct {
		dst    interface{}
		v      interface{}
		name   string
		assert func(t *testing.T, v interface{}, dst interface{})
	}{
		{
			name: "Test conversion from BinaryValue to []byte",
			dst:  &[]byte{},
			v:    &BinaryValue{Data: ioutil.NopCloser(bytes.NewBuffer([]byte{0x00, 0x01, 0xFF}))},
			assert: func(t *testing.T, v interface{}, dst interface{}) {
				b, ok := dst.(*[]byte)
				require.True(t, ok)
				assert.Len(t, *b, 3)
			},
		},

		{
			name: "Test conversion from []byte to *BinaryValue",
			dst:  &BinaryValue{},
			v:    []byte{0x00, 0x01, 0xFF},
			assert: func(t *testing.T, v interface{}, dst interface{}) {
				b, ok := dst.(*BinaryValue)
				require.True(t, ok)
				require.NotNil(t, b.Data)
				bb, _ := ioutil.ReadAll(b.Data)
				assert.Len(t, bb, 3)
			},
		},

		{
			name: "Test conversion from [][]byte to *[]*BinaryValue",
			dst:  &[]*BinaryValue{},
			v:    [][]byte{{0x00}, {0x01}, {0x02}},
			assert: func(t *testing.T, v interface{}, dst interface{}) {
				b, ok := dst.(*[]*BinaryValue)
				require.True(t, ok)

				require.Len(t, *b, 3)
				for _, bv := range *b {
					bb, _ := ioutil.ReadAll(bv.Data)
					assert.Len(t, bb, 1)
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					if err, ok := r.(error); ok {
						t.Fatalf("Failed to convert values: %s", err)
					} else {
						t.Fatalf("Failed to convert values: unknown reason")
					}
				}
			}()
			ConvertValue(tc.v, tc.dst)
			if tc.assert != nil {
				tc.assert(t, tc.v, tc.dst)
			}
		})
	}
}
