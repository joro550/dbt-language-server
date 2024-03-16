package main

import (
	"testing"

	p "github.com/tliron/glsp/protocol_3_16"
)

func TestThing(t *testing.T) {
	helper(p.Position{Line: 0, Character: 10}, "{{ ref('my_first_dbt_model') }}", "my_first_dbt_model", true, t)
	helper(p.Position{Line: 1, Character: 10}, "\n{{ ref('my_first_dbt_model') }}", "my_first_dbt_model", true, t)
	helper(p.Position{Line: 0, Character: 8}, "{{ ref('my_first_dbt_model') }}", "my_first_dbt_model", true, t)
	helper(p.Position{Line: 0, Character: 25}, "{{ ref('my_first_dbt_model') }}", "my_first_dbt_model", true, t)
	helper(p.Position{Line: 0, Character: 26}, "{{ ref('my_first_dbt_model') }}", "my_first_dbt_model", true, t)
	helper(p.Position{Line: 0, Character: 0}, "{{ ref('my_first_dbt_model') }}", "", false, t)
}

func helper(position p.Position, code, expectedName string, success bool, t *testing.T) {
	node := Node{RawCode: code}

	ok, value := node.DoThing(position)
	if ok != success || value != expectedName {
		t.Errorf("Something was not right %v", position)
	}
}
