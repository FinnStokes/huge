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
func (m *EntityManager) All() []*Entity {
	entities := make([]*Entity, 0, m.entities.Len())
	for e := m.entities.Front(); e != nil; e = e.Next() {
		append(entities, e.Value)
	}
	return entities
}

// Get returns the unique entity with the given id.
func (m *EntityManager) Get(id int) *Entity {
	return m.ids[id]
}

// Tag returns the unique entity with the given tag.
func (m *EntityManager) Tag(id string) *Entity {
	return m.tags[id]
}

// Group returns a slice of all entities in the named group.
func (m *EntityManager) Group(id string) []*Entity {
	g := make([]*Entity, 0, 10)
	group := m.groups[id]
	for e := range group {
		append(g, e)
	}
	return g
}

// New creates a new entity, assigning it to the next available id and
// returning the created entity.
func (m *EntityManager) New() *Entity {
	e := new(Entity)
	e.id = m.nextId
	m.nextId++
	m.entities.PushBack(e)
	m.ids[e.id] = e
	return e
}

// Delete removes the given entity from the manager.
func (m *EntityManager) Delete(entity *Entity) {
	for e := m.entities.Front(); e != nil; e = e.Next() {
		if e.value == entity {
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
func (m *EntityManager) SetTag(entity *Entity, tag string) {
	m.tags[tag] = entity
}

// SetGroup adds the given entity to the given group(s).
func (m *EntityManager) SetGroup(entity *Entity, groups ...string) {
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
func (m *EntityManager) ClearTag(tag string) {
	delete(m.tags, tag)
}

// ClearGroup removes the given entity from the given group(s).
func (m *EntityManager) ClearGroup(entity *Entity, groups ...string) {
	for _, g := range groups {
		group, ok := m.groups[g]
		if ok {
			delete(group, entity)
		}
	}
}
