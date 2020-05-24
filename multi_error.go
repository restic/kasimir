package main

// MultiError bundles several errors.
type MultiError struct {
	list []string
}

func (m *MultiError) Insert(err error) {
	if err == nil {
		return
	}

	m.list = append(m.list, err.Error())
}

func (m *MultiError) Error() (s string) {
	for _, text := range m.list {
		s += text + "\n"
	}

	return s
}

func (m *MultiError) Length() int {
	return len(m.list)
}
