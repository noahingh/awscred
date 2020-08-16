package ini

import (
	"fmt"
	"reflect"

	"gopkg.in/ini.v1"
)

// DeepEqual check the both ini file is equal or not.
func DeepEqual(f *ini.File, cmp *ini.File) error {
	if err := equalStrings(f.SectionStrings(), cmp.SectionStrings()); err != nil {
		return fmt.Errorf("ini.DeepEqual sections is not equal: %s", err)
	}

	for _, sec := range f.Sections() {
		cs, _ := cmp.GetSection(sec.Name())
		if err := equalStrings(sec.KeyStrings(), cs.KeyStrings()); err != nil {
			return fmt.Errorf("ini.DeepEqual the keys of \"%s\" is not equal: %s", sec.Name(), err)
		}

		for _, key := range sec.Keys() {
			ck, _ := cs.GetKey(key.Name())

			if !reflect.DeepEqual(key.Value(), ck.Value()) {
				return fmt.Errorf("ini.DeepEqual the value of %s is not equal: %s != %s", key.Name(), key.Value(), ck.Value())
			}
		}
	}
	return nil
}

func equalStrings(m, n []string) error {
	has := make(map[string]bool)

	if len(m) != len(n) {
		return fmt.Errorf("The len of slice is different: %s != %s", m, n)
	}

	for _, nk := range n {
		has[nk] = true
	}

	for _, mk := range m {
		if _, ok := has[mk]; !ok {
			return fmt.Errorf("The %s doesn't exist in %s", mk, n)
		}
	}
	return nil
}
