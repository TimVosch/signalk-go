package tree_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"signalk/tree"
)

func TestTreeSerializes(t *testing.T) {
	root := tree.Create()
	err := root.Set(tree.CreatePath("a.b.c"), "1")
	require.NoError(t, err)
	err = root.Set(tree.CreatePath("a.b.d"), "2")
	require.NoError(t, err)
	err = root.Set(tree.CreatePath("a.b.d"), "3")
	require.NoError(t, err)

	jsonBytes, err := json.Marshal(root)
	require.NoError(t, err)
	expectedJSON := `{"a":{"b":{"c":"1","d":"3"}}}`
	assert.Equal(t, expectedJSON, string(jsonBytes))

	fmt.Printf("\n\n\n")
	cNode, err := root.Get(tree.CreatePath("a.b.c"))
	require.NoError(t, err)
	assert.Equal(t, "1", cNode.GetValue())

	assert.Equal(t, "a.b.c", cNode.Path().String())

	cNode.SetValue("5")
	jsonBytes, err = json.Marshal(root)
	require.NoError(t, err)
	expectedJSON = `{"a":{"b":{"c":"5","d":"3"}}}`
	assert.Equal(t, expectedJSON, string(jsonBytes))
}
