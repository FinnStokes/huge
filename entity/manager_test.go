package entity

import (
	"sort"
	"testing"
)

func TestNew(t *testing.T) {
	m := NewManager()
	entities := make([]*Entity, 100)
	for i := range entities {
		entities[i] = m.New()
		if entities[i] == nil {
			t.Errorf("Entity %v is nil", i)
		}
	}
	for i, e := range m.All() {
		if e != entities[i] {
			t.Errorf("Entity %v does not match", i)
		}
	}
}

func TestGet(t *testing.T) {
	m := NewManager()
	entities := make([]*Entity, 100)
	for i := range entities {
		entities[i] = m.New()
	}
	for i, e := range entities {
		if m.Get(e.Id()) != e {
			t.Errorf("Entity %v does not match (%v vs %v)", i, e.Id(), m.Get(e.Id()).Id())
		}
	}
}

func TestTag(t *testing.T) {
	m := NewManager()
	entities := make([]*Entity, 100)
	for i := range entities {
		entities[i] = m.New()
	}
	tags := map[string]int{
		"five":       5,
		"fifty":      50,
		"49":         49,
		"ninetynine": 99,
		"99":         99,
	}
	for tag, id := range tags {
		m.SetTag(entities[id], tag)
	}
	for tag, id := range tags {
		if m.Tag(tag) != entities[id] {
			t.Errorf("Tag %v does not match (%v vs %v)", tag, entities[id].Id(), m.Tag(tag).Id())
		}
	}
	for tag := range tags {
		m.SetTag(entities[0], tag)
		if m.Tag(tag) != entities[0] {
			t.Errorf("Overridden tag %v does not match (%v vs %v)", tag, entities[0].Id(), m.Tag(tag).Id())
		}
	}
	for tag := range tags {
		m.ClearTag(tag)
		if m.Tag(tag) != nil {
			t.Errorf("Tag %v not cleared", tag)
		}
	}
	for tag, id := range tags {
		m.SetTag(entities[id], tag)
		if m.Tag(tag) != entities[id] {
			t.Errorf("Redefined tag %v does not match (%v vs %v)", tag, entities[id].Id(), m.Tag(tag).Id())
		}
	}
	for tag := range tags {
		m.ClearTag(tag)
	}
	for tag := range tags {
		if m.Tag(tag) != nil {
			t.Errorf("Tag %v not cleared", tag)
		}
	}
}

func sliceEqual(slice1 []int, slice2 []int) bool {
	if len(slice1) != len(slice2) {
		return false
	}
	for i := range slice1 {
		if slice1[i] != slice2[i] {
			return false
		}
	}
	return true
}

func inSlice(elt int, slice []int) bool {
	for _, v := range slice {
		if v == elt {
			return true
		}
	}
	return false
}

func TestGroup(t *testing.T) {
	m := NewManager()
	entities := make([]*Entity, 100)
	for i := range entities {
		entities[i] = m.New()
	}
	groups := map[string][]int{
		"a":    []int{5, 6, 7, 8, 9},
		"b":    []int{50, 30, 80},
		"9":    []int{49},
		"asdf": []int{},
	}
	for _, v := range groups {
		sort.Ints(v)
	}
	for group, ids := range groups {
		for _, id := range ids {
			m.SetGroup(m.Get(id), group)
		}
	}
	for group, ents := range groups {
		g := make([]int, len(m.Group(group)))
		for i, v := range m.Group(group) {
			g[i] = v.Id()
		}
		sort.Ints(g)
		if !sliceEqual(g, ents) {
			t.Errorf("Group %v does not match (%v vs %v)", group, g, ents)
		}
		for _, e := range ents {
			if !m.InGroup(entities[e], group) {
				t.Errorf("Entity %v not in group %v", e, group)
			}
		}
	}
	groupNames := make([]string, 0, len(groups))
	for group, _ := range groups {
		groupNames = append(groupNames, group)
	}
	m.SetGroup(entities[99], groupNames...)
	m.SetGroup(entities[98], groupNames...)
	for group, ents := range groups {
		g := make([]int, len(m.Group(group)))
		for i, v := range m.Group(group) {
			g[i] = v.Id()
		}
		sort.Ints(g)
		ents = append(ents, 99, 98)
		sort.Ints(ents)
		groups[group] = ents
		if !sliceEqual(g, ents) {
			t.Errorf("Group %v does not match (%v vs %v)", group, g, ents)
		}
		for _, e := range ents {
			if !m.InGroup(entities[e], group) {
				t.Errorf("Entity %v not in group %v", e, group)
			}
		}
	}
	m.ClearGroup(entities[98], groupNames...)
	for group, ents := range groups {
		m.ClearGroup(entities[ents[len(ents)-1]], group)
		g := make([]int, len(m.Group(group)))
		for i, v := range m.Group(group) {
			g[i] = v.Id()
		}
		sort.Ints(g)
		ents = ents[:len(ents)-2]
		groups[group] = ents
		if !sliceEqual(g, ents) {
			t.Errorf("Group %v does not match (%v vs %v)", group, g, ents)
		}
		for _, e := range ents {
			if !m.InGroup(entities[e], group) {
				t.Errorf("Entity %v not in group %v", e, group)
			}
		}
	}
	m.SetGroup(entities[98], groupNames[0])
	m.ClearGroup(entities[98], groupNames...)
	for group, ents := range groups {
		g := make([]int, len(m.Group(group)))
		for i, v := range m.Group(group) {
			g[i] = v.Id()
		}
		sort.Ints(g)
		if !sliceEqual(g, ents) {
			t.Errorf("Group %v does not match (%v vs %v)", group, g, ents)
		}
		for _, e := range ents {
			if !m.InGroup(entities[e], group) {
				t.Errorf("Entity %v not in group %v", e, group)
			}
		}
	}
	for group, ents := range groups {
		for len(ents) > 0 {
			m.ClearGroup(entities[ents[0]], group)
			g := make([]int, len(m.Group(group)))
			for i, v := range m.Group(group) {
				g[i] = v.Id()
			}
			sort.Ints(g)
			ents = ents[1:]
			groups[group] = ents
			if !sliceEqual(g, ents) {
				t.Errorf("Group %v does not match (%v vs %v)", group, g, ents)
			}
			for _, e := range ents {
				if !m.InGroup(entities[e], group) {
					t.Errorf("Entity %v not in group %v", e, group)
				}
			}
		}
	}
}

func TestDelete(t *testing.T) {
	m := NewManager()
	entities := make([]*Entity, 100, 110)
	for i := range entities {
		entities[i] = m.New()
	}
	tags := map[string]int{
		"five":       5,
		"fifty":      50,
		"49":         49,
		"ninetynine": 99,
		"99":         99,
	}
	groups := map[string][]int{
		"a":    []int{5, 6, 7, 8, 9},
		"b":    []int{50, 30, 80},
		"9":    []int{49},
		"asdf": []int{},
	}
	del := []int{2, 4, 5, 6, 8, 10, 30, 35, 46, 79, 80, 99}
	for tag, id := range tags {
		m.SetTag(entities[id], tag)
	}
	for _, v := range groups {
		sort.Ints(v)
	}
	for group, ids := range groups {
		for _, id := range ids {
			m.SetGroup(m.Get(id), group)
		}
	}
	for _, id := range del {
		m.Delete(entities[id])
	}
	ents := m.All()
	for i, j := 0, 0; i < len(entities) && j < len(ents); i++ {
		if entities[i] != ents[j] {
			if !inSlice(i, del) {
				t.Errorf("Entity %v does not match", i)
			}
		} else {
			j++
			if inSlice(i, del) {
				t.Errorf("Entity %v was not deleted", i)
			}
		}
	}
	for i := 0; i < 10; i++ {
		entities = append(entities, m.New())
	}
	ents = m.All()
	for i, j := 0, 0; i < len(entities) && j < len(ents); i++ {
		if entities[i] != ents[j] {
			if !inSlice(i, del) {
				t.Errorf("Entity %v does not match", i)
			}
		} else {
			j++
			if inSlice(i, del) {
				t.Errorf("Entity %v was not deleted", i)
			}
		}
	}
	for _, id := range del {
		if m.Get(id) != nil {
			t.Errorf("Entity %v is still gettable", id)
		}
	}
	for tag, id := range tags {
		if inSlice(id, del) {
			if m.Tag(tag) != nil {
				t.Errorf("Tag %v is still gettable", tag)
			}
		} else {
			if m.Tag(tag) != entities[id] {
				t.Errorf("Tag %v does not match (%v vs %v)", tag, m.Tag(tag).Id(), id)
			}
		}
	}
	for group := range groups {
		g := m.Group(group)
		for _, e := range g {
			if inSlice(e.Id(), del) {
				t.Errorf("Entity %v is still in group %v", e.Id(), group)
			}
		}
	}
}
