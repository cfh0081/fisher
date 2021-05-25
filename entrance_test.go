package fisher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckArgments(t *testing.T) {

	arg0 := []string{"a", "b", "c"}
	rtn0 := checkArgments(arg0)
	assert.NotNil(t, rtn0)

	arg1 := []string{"=a"}
	rtn1 := checkArgments(arg1)
	assert.NotNil(t, rtn1)

	arg2 := []string{"cx=bc", "a-)x="}
	rtn2 := checkArgments(arg2)
	assert.NotNil(t, rtn2)

	arg3 := []string{"a=b", "b==", "c=.", "a-)x=="}
	rtn3 := checkArgments(arg3)
	assert.Nil(t, rtn3)
}

func TestGetCustomArgs(t *testing.T) {
	expected := map[string]string{"a": "b", "b": "=", "c": ".", "a-)x": "="}
	arg1 := []string{"a=b", "b==", "c=.", "a-)x=="}
	rtn1 := getCustomArgs(arg1)
	assert.Equal(t, expected, rtn1)
}
