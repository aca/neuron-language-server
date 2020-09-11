package neuron

import (
	"encoding/json"
	"os/exec"
)

func Query(arg ...string) (*QueryResult, error) {
	// TODO: to avoid creating .neuron directory
	// This should be optional with flags
	// path, err := os.Getwd()
	// if err != nil {
	// 	return nil, err
	// }
	//
	// _, err = os.Stat(filepath.Join(path, ".neuron"))
	// if os.IsNotExist(err) {
	// 	return nil, nil
	// } else if err != nil {
	// 	return nil, err
	// }

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
