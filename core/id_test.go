package core

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/pborman/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var validInternalIDs = []ID{
	ID("00000000-0000-0000-0000-00005ac3fbfb"),
	ID("00000000-0000-0000-0000-00015ab3f6fb"),
	ID("00000000-0000-0000-0000-00025a43f8fb"),
	ID("00000000-0000-0000-0000-00035ac3fbab"),
}

var notValidInternalIDs = []ID{
	ID("10000000-0000-0000-0000-00035ac3fbab"),
	ID("01000000-0000-0000-0000-00035ac3fbab"),
	ID("00100000-0000-0000-0000-00035ac3fbab"),
	ID("00010000-0000-0000-0000-00035ac3fbab"),
	ID("00001000-0000-0000-0000-00035ac3fbab"),
	ID("00000100-0000-0000-0000-00035ac3fbab"),
	ID("00000010-0000-0000-0000-00035ac3fbab"),
	ID("00000001-0000-0000-0000-00035ac3fbab"),
	ID("00000000-1000-0000-0000-00035ac3fbab"),
	ID("00000000-0100-0000-0000-00035ac3fbab"),
	ID("00000000-0010-0000-0000-00035ac3fbab"),
	ID("00000000-0001-0000-0000-00035ac3fbab"),
	ID("00000000-0000-1000-0000-00035ac3fbab"),
	ID("00000000-0000-0100-0000-00035ac3fbab"),
	ID("00000000-0000-0010-0000-00035ac3fbab"),
	ID("00000000-0000-0001-0000-00035ac3fbab"),
	ID("00000000-0000-0000-1000-00035ac3fbab"),
	ID("00000000-0000-0000-0100-00035ac3fbab"),
	ID("00000000-0000-0000-0010-00035ac3fbab"),
	ID("00000000-0000-0000-0001-00035ac3fbab"),
	ID("00000000-0000-0000-0000-10035ac3fbab"),
	ID("00000000-0000-0000-0000-01035ac3fbab"),
	ID("00000000-0000-0000-0000-00135ac3fbab"),
}

func TestNewIDLength(t *testing.T) {
	require.Equal(t, len(uuid.New()), len(NewID()))
	require.Equal(t, len(uuid.New()), len(NewID()))
	require.Equal(t, len(uuid.New()), len(NewID()))
	require.Equal(t, len(uuid.New()), len(NewID()))
	require.Equal(t, len(uuid.New()), len(NewID()))
	require.Equal(t, len(uuid.New()), len(NewID()))
	require.Equal(t, len(uuid.New()), len(NewID()))
	require.Equal(t, len(uuid.New()), len(NewID()))
}

func TestID_WithTime(t *testing.T) {

	t.Parallel()

	t.Run("OK", func(t *testing.T) {
		timee := time.Now()
		id := NewID()
		idWithTime, err := id.WithTime(timee)
		assert.NoError(t, err)

		lastPartUint, err := idWithTime.extractLastPartAsUint()
		assert.NoError(t, err)
		lastPartUint -= uint64(timee.Unix())

		oldLastPartUint, err := id.extractLastPartAsUint()
		assert.NoError(t, err)
		assert.Equal(t, oldLastPartUint, lastPartUint)
	})

	t.Run("invalid ID should fail", func(t *testing.T) {
		id := ID("ffffffff-ffff-ffff-ffff-xxxxxxxxxxxx")
		newID, err := id.WithTime(time.Now())
		assert.Equal(t, id.String(), newID.String())
		assert.Error(t, err)
	})

	t.Run("malformed ID should not panic", func(t *testing.T) {
		id := ID("fffffff-fff-fff-fff-")
		newID, err := id.WithTime(time.Now())
		assert.Equal(t, id.String(), newID.String())
		assert.Error(t, err)
	})

}

func TestID_AbTest(t *testing.T) {
	cases := []struct {
		expected bool
		lastNum  int64
		percent  int
	}{
		{
			expected: true,
			lastNum:  0,
			percent:  100,
		},
		{
			expected: false,
			lastNum:  99,
			percent:  0,
		},
		{
			expected: false,
			lastNum:  99,
			percent:  99,
		},
		{
			expected: false,
			lastNum:  01,
			percent:  0,
		},
		{
			expected: true,
			lastNum:  99,
			percent:  100,
		},
		{
			expected: false,
			lastNum:  1,
			percent:  1,
		},
		{
			expected: false,
			lastNum:  2,
			percent:  1,
		},
		{
			expected: true,
			lastNum:  0,
			percent:  1,
		},
	}

	for _, c := range cases {
		id, err := idWithLastNum(NewID(), c.lastNum)
		if err != nil {
			t.Errorf("Unexpected ID format, %v", err)
		}
		if id.AbTest(c.percent) != c.expected {
			t.Errorf("Unexpected result, %v, %v", c, id)
		}
	}
}

func TestABTestRange(t *testing.T) {
	type testCase struct {
		comment  string
		id       ID
		from, to int
		expected bool
	}

	testCases := []testCase{
		{
			comment:  "value strictly inside range",
			id:       ID("db2e2e9e-7b84-46f0-93e7-000000000007"),
			from:     5,
			to:       10,
			expected: true,
		},
		{
			comment:  "value on the left bound of range",
			id:       ID("db2e2e9e-7b84-46f0-93e7-000000000005"),
			from:     5,
			to:       10,
			expected: true,
		},
		{
			comment:  "value on the right bound of range",
			id:       ID("db2e2e9e-7b84-46f0-93e7-00000000000a"),
			from:     5,
			to:       10,
			expected: false,
		},
		{
			comment:  "value strictly ouside of range on the right",
			id:       ID("db2e2e9e-7b84-46f0-93e7-00000000000f"),
			from:     5,
			to:       10,
			expected: false,
		},
		{
			comment:  "value strictly ouside of range on the left",
			id:       ID("db2e2e9e-7b84-46f0-93e7-000000000002"),
			from:     5,
			to:       10,
			expected: false,
		},
		{
			comment:  "id is not a valid uuid",
			id:       ID("hello world"),
			from:     0,
			to:       10,
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.comment, func(t *testing.T) {
			actual := tc.id.InABTestRange(tc.from, tc.to)
			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestID_ABTestGroup(t *testing.T) {
	tests := []struct {
		lastNum       int64
		expectedGroup int
	}{
		{0, 0},
		{10, 10},
		{43, 43},
		{51, 51},
		{99, 99},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test case #%d", i), func(t *testing.T) {
			id, err := idWithLastNum(NewID(), tt.lastNum)
			if err != nil {
				t.Errorf("Unexpected ID format, %v", err)
			}
			if id.ABTestGroup() != tt.expectedGroup {
				t.Fatalf("Unexpected result for ID %q, %v", id, tt)
			}
		})
	}
}

func TestNewInternalID(t *testing.T) {
	ids := map[ID]struct{}{}
	for i := 0; i < 16; i++ {
		internalID := NewInternalID(i)

		if uuid.Parse(internalID.String()) == nil {
			t.Errorf("Not valid uuid: %v", internalID)
		}
		if !internalID.IsInternal() {
			t.Errorf("The new internal id is not internal: %v", internalID)
		}
		// Some crazy test on uniqueness
		if _, ok := ids[internalID]; ok {
			t.Errorf("Interanl id is duplicate: %v", internalID)
		}

		ids[internalID] = struct{}{}
	}
}

func TestIDIsInternalTrue(t *testing.T) {
	for _, id := range validInternalIDs {
		if uuid.Parse(id.String()) == nil {
			t.Errorf("Not valid uuid: %v", id)
		}
		if id.IsInternal() == false {
			t.Errorf("Expected true but got false for id: %v", id)
		}
	}
}

func TestIDIsInternalFalse(t *testing.T) {
	for _, id := range notValidInternalIDs {
		if uuid.Parse(id.String()) == nil {
			t.Errorf("Not valid uuid: %v", id)
		}
		if id.IsInternal() == true {
			t.Errorf("Expected false but got true for id: %v", id)
		}
	}
}

func TestIDScan(t *testing.T) {
	table := map[ID]interface{}{
		ID("00000000-0000-0000-0000-00035ac3fbab"): "00000000-0000-0000-0000-00035ac3fbab",
		ID("00000000-0000-0000-0000-00035ac3fbac"): []byte("00000000-0000-0000-0000-00035ac3fbac"),
	}

	var i ID
	for expected, v := range table {
		err := (&i).Scan(v)
		if err != nil {
			t.Error(err)
		}

		if i != expected {
			t.Error("Expected:", expected, "got: ", i)
		}
	}

	errorTable := []interface{}{
		0.35,
		"bad id",
		nil,
	}

	for _, errorVal := range errorTable {
		if err := i.Scan(errorVal); err == nil {
			t.Errorf("expected an error for input %q", errorVal)
		}
	}
}

func TestResetChunkCounter(t *testing.T) {
	id := ID("10000000-8000-0000-5000-000000002015")
	id.ResetChunkCounter()

	expected := ID("10000000-8000-0000-5000-000000000000")
	if id != expected {
		t.Errorf("Expected %q, got %q", expected, id)
	}
}

func TestParseID(t *testing.T) {
	idStr := "10000000-8000-0000-5000-000000002015"
	id, ok := ParseID(idStr)
	require.True(t, ok)
	require.Equal(t, idStr, id.String())

	_, ok = ParseID("10000000-800-0000-5000-000000002015")
	require.False(t, ok)
}

func TestIncChunkCounter(t *testing.T) {
	id := ID("10000000-8000-0000-5000-000000002015")
	id.ResetChunkCounter()

	for i := 0; i < 255; i++ {
		id.Inc()
	}

	expected := ID("10000000-8000-0000-5000-0000000000ff")
	if id != expected {
		t.Errorf("Expected %q, got %q", expected, id)
	}
}

func TestIsEqualUpperCaseAndLowerCase(t *testing.T) {
	upperCaseID := ID("8D29309C-307B-4DA3-AEE1-3B01251EFE66")
	lowerCaseID := ID("8d29309c-307b-4da3-aee1-3b01251efe66")

	require.True(t, upperCaseID.Eq(lowerCaseID), "Expected true for equal upper case and lower case ID")
}

func TestToNullID(t *testing.T) {
	t.Run("nil ID", func(t *testing.T) {
		var id *ID
		nullID := id.ToNullID()
		assert.False(t, nullID.Valid)
		assert.Error(t, nullID.ID.validate())
	})

	t.Run("new ID", func(t *testing.T) {
		id := NewID()
		nullID := id.ToNullID()
		assert.True(t, nullID.Valid)
		assert.NoError(t, nullID.ID.validate())
	})
}

func BenchmarkIDMarshalJSON(b *testing.B) {
	id := NewID()

	for n := 0; n < b.N; n++ {
		_, err := id.MarshalJSON()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func idWithLastNum(id ID, lastNum int64) (ID, error) {
	s := strings.Split(id.String(), "-")
	if len(s) < 5 {
		return ID(""), fmt.Errorf("Unexpected id format, %v", id)
	}
	value, err := strconv.ParseInt(s[4], 16, 64)
	if err != nil {
		return ID(""), fmt.Errorf("Unexpected id section format, %v", s[4])
	}
	s[4] = strconv.FormatInt((value/100)*100+lastNum, 16)
	return ID(strings.Join(s, "-")), nil
}
