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

	value1 := &signalk.VesselDataEntry{
		Path:      signalk.CreatePath("a.b.c"),
		SourceRef: "test",
		Value:     signalk.DataValueFromNumerical(15),
		Timestamp: time.Now(),
	}
	value2 := &signalk.VesselDataEntry{
		Path:      signalk.CreatePath("a.b.d"),
		SourceRef: "test",
		Value:     signalk.DataValueFromNumerical(55),
		Timestamp: time.Now(),
	}
	tree.Add(value1.Path, value1)
	tree.Add(value2.Path, value2)

	entry, retTree, err := tree.Get(value1.Path)
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
