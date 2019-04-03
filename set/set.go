package set

type HashSet struct {
	items map[interface{}]status
}

type status bool

const statusExists status = true

func New() *HashSet {
	return &HashSet{items: make(map[interface{}]status)}
}


func (set *HashSet) Add(item interface{}) bool {
	if set.Contains(item) {
		return false
	}
	set.items[item] = statusExists
	return true
}

func (set *HashSet) Remove(item interface{}) bool {
	if !set.Contains(item) {
		return false
	}
	delete(set.items, item)
	return true
}

func (set *HashSet) Size() int {
	return len(set.items)
}

func (set *HashSet) Contains(item interface{}) bool {
	_, exists := set.items[item]
	return exists
}

func (set *HashSet) Values() []interface{} {
	values := make([]interface{}, len(set.items))
	count := 0
	for k, v := range set.items {
		if v {
			values[count] = k
			count++
		}

	}
	return values
}
