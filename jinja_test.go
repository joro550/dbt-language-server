package main

import "testing"

func TestGettingRefTags(t *testing.T) {
	testWrapper(`{{ ref('my_first_dbt_model')}}`, []string{"my_first_dbt_model"}, []int{0}, 1, t)
	testWrapper(`{{ ref('project', 'my_first_dbt_model')}}`, []string{"my_first_dbt_model"}, []int{0}, 1, t)
	testWrapper(`{{ ref('my_first_dbt_model')}} {{ ref('my_second_dbt_model')}}`, []string{"my_first_dbt_model", "my_second_dbt_model"}, []int{0, 31}, 2, t)
}

func TestMacroFunctionName(t *testing.T) {
	testMacroWrapper(`{{ hello() }}`, []string{"hello"}, []int{0}, 1, t)
}

func testWrapper(content string, expectedName []string, expectedPosition []int, expectedLength int, t *testing.T) {
	parser := NewJinjaParser()
	hasRefTag := content
	refs := parser.GetAllRefTags(hasRefTag)
	if len(refs) != expectedLength {
		t.Errorf("not enough refs")
	}

	for i, ref := range refs {
		if ref.ModelName != expectedName[i] {
			t.Errorf("got wrong name expected %v but got %v", expectedName[i], ref.ModelName)
		}
	}
	for i, ref := range refs {
		if ref.Range.Start != expectedPosition[i] {
			t.Errorf("got wrong position expected %v but got %v", expectedPosition[i], ref.Range.Start)
		}
	}
}

func testMacroWrapper(content string, expectedName []string, expectedPosition []int, expectedLength int, t *testing.T) {
	parser := NewJinjaParser()
	macros := parser.GetMacros(content)
	if len(macros) != expectedLength {
		t.Errorf("not enough refs")
	}

	for i, ref := range macros {
		if ref.ModelName != expectedName[i] {
			t.Errorf("got wrong name expected %v but got %v", expectedName[i], ref.ModelName)
		}
	}
	for i, ref := range macros {
		if ref.Range.Start != expectedPosition[i] {
			t.Errorf("got wrong position expected %v but got %v", expectedPosition[i], ref.Range.Start)
		}
	}
}
