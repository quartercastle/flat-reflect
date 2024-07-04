package flat

type Map map[string]*Token

func (m Map) Filter(fn func(string, *Token) bool) Map {
	result := Map{}
	for key := range m {
		if fn(key, m[key]) {
			result[key] = m[key]
		}
	}
	return result
}

func Merge(a, b Map) Map {
	for k, v := range b {
		a[k] = v
	}
	return a
}
