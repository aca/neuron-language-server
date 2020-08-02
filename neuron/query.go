package neuron

import (
	"encoding/json"
	"os/exec"
)

func Query(arg ...string) (*QueryResult, error) {
	arg = append([]string{"query"}, arg...)
	cmd := exec.Command("neuron", arg...)
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	result := new(QueryResult)
	err = json.Unmarshal(output, result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
