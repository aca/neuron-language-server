package neuron

// import (
// 	"sort"
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// )

// func TestGetHeader(t *testing.T) {
// 	input := `---
// title: Installing NixOS on OVH dedicated servers
// date: "2020-05-26"
// tags:
//     - timeline 
//     - nixos
// ---

// I recently setup a [AMD Ryzen 7 3700 PRO](https://www.ovh.com/ca/en/dedicated-servers/infra/infra-limited-edition-2/) 
// `

// 	header, err := GetHeader([]byte(input))
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	sort.Strings(header.Tags)

// 	assert.Equal(t, header.Title, "Installing NixOS on OVH dedicated servers")
// 	assert.Equal(t, header.Date, "2020-05-26")
// 	assert.NotEqual(t, sort.SearchStrings(header.Tags, "timeline"), len(header.Tags))
// 	assert.NotEqual(t, sort.SearchStrings(header.Tags, "nixos"), 2)
// }

// func TestGetHeader2(t *testing.T) {
// 	input := `I recently setup a [AMD Ryzen 7 3700 PRO](https://www.ovh.com/ca/en/dedicated-servers/infra/infra-limited-edition-2/) 
// ---
// title: Installing NixOS on OVH dedicated servers
// date: "2020-05-26"
// tags:
//     - timeline 
//     - nixos
// ---
// `

// 	_, err := GetHeader([]byte(input))
// 	if err == nil {
// 		t.Error("test should fail")
// 	}
// }
