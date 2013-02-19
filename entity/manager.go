package entity

import "container/list"

type Manager struct {
	entities list.List
	nextId   int
	ids      map[int]*Entity
	tags     map[string]*Entity
	groups   map[string]map[*Entity]bool
}

// NewManager returns an initialised manager.
func NewManager() *Manager {
	m := new(Manager)
	m.ids = make(map[int]*Entity)
	m.tags = make(map[string]*Entity)
	m.groups = make(map[string]map[*Entity]bool)
	return m
}

// All returns a slice of all active entities.
func (m *Manager) All() []*Entity {
	entities := make([]*Entity, 0, m.entities.Len())
	for e := m.entities.Front(); e != nil; e = e.Next() {
		if ent, ok := e.Value.(*Entity); ok {
			entities = append(entities, ent)
		}
	}
	return entities
}

// Get returns the unique entity with the given id.
func (m *Manager) Get(id int) *Entity {
	return m.ids[id]
}

// Tag returns the unique entity with the given tag.
func (m *Manager) Tag(id string) *Entity {
	return m.tags[id]
}

// Group returns a slice of all entities in the named group.
func (m *Manager) Group(id string) []*Entity {
	g := make([]*Entity, 0, 10)
	group := m.groups[id]
	for e := range group {
		g = append(g, e)
	}
	return g
}

// InGroup returns true if the given entity is in the named group and false otherwise.
func (m *Manager) InGroup(entity *Entity, id string) bool {
	g, ok := m.groups[id]
	in := false
	if ok {
		_, in = g[entity]
	}
	return in
}

// New creates a new entity, assigning it to the next available id and
// returning the created entity.
func (m *Manager) New() *Entity {
	e := new(Entity)
	e.id = m.nextId
	e.Components = make(map[string]interface{})
	m.nextId++
	m.entities.PushBack(e)
	m.ids[e.id] = e
	return e
}

// Delete removes the given entity from the manager.
func (m *Manager) Delete(entity *Entity) {
	e := m.entities.Front()
	for ; e != nil; e = e.Next() {
		if e.Value == entity {
			break
		}
	}
	if e != nil {
		m.entities.Remove(e)
	}
	delete(m.ids, entity.id)
	for tag, e := range m.tags {
		if e == entity {
			delete(m.tags, tag)
		}
	}
	for _, g := range m.groups {
		delete(g, entity)
	}
}

// SetTag sets the given tag to refer to the given entity.
func (m *Manager) SetTag(entity *Entity, tag string) {
	m.tags[tag] = entity
}

// SetGroup adds the given entity to the given group(s).
func (m *Manager) SetGroup(entity *Entity, groups ...string) {
	for _, g := range groups {
		group, ok := m.groups[g]
		if !ok {
			group = make(map[*Entity]bool)
			m.groups[g] = group
		}
		group[entity] = true
	}
}

// ClearTag removes the given tag.
func (m *Manager) ClearTag(tag string) {
	delete(m.tags, tag)
}

// ClearGroup removes the given entity from the given group(s).
func (m *Manager) ClearGroup(entity *Entity, groups ...string) {
	for _, g := range groups {
		group, ok := m.groups[g]
		if ok {
			delete(group, entity)
		}
	}
}
