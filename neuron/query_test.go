package neuron

import (
	"testing"
)

func TestQuery(t *testing.T) {
	_, err := Query()
	if err != nil {
		t.Fatal(err)
	}
}
