package edge

import (
	"sync"
)

var stringJoinRelKey = "Rel"

type graphManager struct {
	locker   *sync.Mutex
	group    *sync.WaitGroup
	mapEdges map[string]*edgeImpl
}

func (m *graphManager) addEdge(edgeName Type, edgeTypeSchema Type, edgeClass *edgeImpl) {
	m.locker.Lock()
	defer m.locker.Unlock()
	encodeEdgeName := string(edgeName) + stringJoinRelKey + string(edgeTypeSchema)
	if _, ok := m.mapEdges[encodeEdgeName]; !ok {
		m.mapEdges[encodeEdgeName] = edgeClass
	}
}

func (m *graphManager) getEdge(edgeName Type, edgeTypeSchema Type) *edgeImpl {
	encodeEdgeName := string(edgeName) + stringJoinRelKey + string(edgeTypeSchema)
	if _, ok := m.mapEdges[encodeEdgeName]; ok {
		return m.mapEdges[encodeEdgeName]
	}
	return nil
}

func (m *graphManager) GetWaitGroup() *sync.WaitGroup {
	return m.group
}

func (m *graphManager) GetLocker() *sync.Mutex {
	return m.locker
}

func (m *graphManager) LoopEdges(do_fn func(edge *edgeImpl)) {
	for _, edge := range m.mapEdges {
		m.locker.Lock()
		do_fn(edge)
		m.locker.Unlock()
	}
}

// skeleton class object
var BaseGraphManager *graphManager = &graphManager{
	locker:   &sync.Mutex{},
	group:    &sync.WaitGroup{},
	mapEdges: make(map[string]*edgeImpl),
}
