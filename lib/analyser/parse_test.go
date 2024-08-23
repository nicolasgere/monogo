package analyzer

import (
	"fmt"
	"testing"
)

func TestList(t *testing.T) {
	modules, err := ListModule("./__playground__/workspace/")
	if err != nil {
		t.Error(err)
	}
	if len(modules) != 2 {
		t.Error(":Module missing")
	}
	for _, m := range modules {
		var d []ModuleDetail
		d, err = GetDependencyForModule(m.Dir)
		if err != nil {
			fmt.Println(m.Dir, err)
		}
		fmt.Println(GetImportFromModules(d))
	}
}
