package signalk_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"signalk/signalk"
)

func TestVesselDataTreeShouldAddLeafOnCorrectPath(t *testing.T) {
	tree := signalk.NewVesselDataTree()

	path1 := signalk.CreatePath("a.b.c")
	value1 := &signalk.VesselDataEntry{
		SourceRef: "test",
		Value:     signalk.DataValueFromNumerical(15),
		Timestamp: time.Now(),
	}
	path2 := signalk.CreatePath("a.b.d")
	value2 := &signalk.VesselDataEntry{
		SourceRef: "test",
		Value:     signalk.DataValueFromNumerical(55),
		Timestamp: time.Now(),
	}
	tree.Add(path1, value1)
	tree.Add(path2, value2)

	entry, retTree, err := tree.Get(path1)
	require.NoError(t, err)
	assert.Nil(t, retTree)
	assert.Equal(t, value1, entry)

	entry, retTree, err = tree.Get(signalk.CreatePath("a.b"))
	require.NoError(t, err)
	assert.Nil(t, entry)

	entry1, treeShouldBeNil, err := retTree.Get(signalk.CreatePath("c"))
	require.NoError(t, err)
	assert.Nil(t, treeShouldBeNil)
	entry2, treeShouldBeNil, err := retTree.Get(signalk.CreatePath("d"))
	require.NoError(t, err)
	assert.Nil(t, treeShouldBeNil)
	assert.Equal(t, value1, entry1)
	assert.Equal(t, value2, entry2)
}
