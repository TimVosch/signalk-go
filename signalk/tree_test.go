package signalk_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"signalk/signalk"
)

func TestTreeSerializes(t *testing.T) {
	tree := signalk.CreateTree()
	err := tree.Set("a.b.c", "1")
	require.NoError(t, err)
	err = tree.Set("a.b.d", "2")
	require.NoError(t, err)
	err = tree.Set("a.b.d", "3")
	require.NoError(t, err)

	jsonBytes, err := json.Marshal(tree)
	require.NoError(t, err)
	expectedJSON := `{"a":{"b":{"c":"1","d":"3"}}}`
	assert.Equal(t, expectedJSON, string(jsonBytes))

	fmt.Printf("\n\n\n")
	cNode, err := tree.Get(signalk.CreatePath("a.b.c"))
	require.NoError(t, err)
	assert.Equal(t, "1", cNode.GetValue())

	assert.Equal(t, "a.b.c", cNode.Path().String())

	cNode.SetValue("5")
	jsonBytes, err = json.Marshal(tree)
	require.NoError(t, err)
	expectedJSON = `{"a":{"b":{"c":"5","d":"3"}}}`
	assert.Equal(t, expectedJSON, string(jsonBytes))
}
