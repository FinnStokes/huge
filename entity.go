package gauge

type Entity struct {
	id         int
	components map[string]Component
}

type EntityManager struct {
	entities []*Entity
	ids      map[int]*Entity
	tags     map[string]*Entity
	groups   map[string][]*Entity
	delete   []int
}

// All returns a list of all active Entities
func (m *EntityManager) All() []*Entity {
	entities := make([]*Entity, len(m.entities))
	copy(m.entities, entities)
	return entities
}

// Get returns the unique Entity with the given id
func (m *EntityManager) Get(id int) *Entity {
	return m.ids[id]
}

// Tag returns the unique Entity with the given tag
func (m *EntityManager) Tag(id string) *Entity {
	return m.tags[id]
}

// Group returns a slice of all entities in the named group
func (m *EntityManager) Group(id string) []*Entity {
	group := m.groups[id]
	entities := make([]*Entity, len(group))
	copy(group, entities)
	return entities
}

// Delete marks the given entity id for deletion
