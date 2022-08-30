package core

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNullID_MarshalJSON_Valid(t *testing.T) {
	nullID := NewNullID("592b56f9-1c75-444f-8eb6-88d40d8314fc")

	bytes, err := json.Marshal(&nullID)

	require.NoError(t, err)
	assert.Equal(t, `"592b56f9-1c75-444f-8eb6-88d40d8314fc"`, string(bytes))
}

func TestNullID_MarshalJSON_NotValid(t *testing.T) {
	nullID := NewNullID("")

	bytes, err := json.Marshal(&nullID)

	require.NoError(t, err)
	assert.Equal(t, `null`, string(bytes))
}

func TestNullID_UnmarshalJSON_Null(t *testing.T) {
	var nullID NullID

	err := json.Unmarshal([]byte("null"), &nullID)

	require.NoError(t, err)
	assert.False(t, nullID.Valid)
	assert.Equal(t, "", nullID.ID.String())
}

func TestNullID_UnmarshalJSON_ID(t *testing.T) {
	var nullID NullID

	err := json.Unmarshal([]byte(`"592b56f9-1c75-444f-8eb6-88d40d8314fc"`), &nullID)

	require.NoError(t, err)
	assert.True(t, nullID.Valid)
	assert.Equal(t, "592b56f9-1c75-444f-8eb6-88d40d8314fc", nullID.ID.String())
}

func TestNullID_UnmarshalJSON_NoValue(t *testing.T) {
	var s struct {
		NullID NullID
	}

	err := json.Unmarshal([]byte(`{}`), &s)

	require.NoError(t, err)
	assert.False(t, s.NullID.Valid)
	assert.Equal(t, "", s.NullID.ID.String())
}
