package edge

import (
	"runtime"
)

type Type string // implement string as edge type name
type IEdge interface {
	RefOn(edgeName Type, edgeTypeSchema Type) IEdge
	Unique() IEdge
}

type edgeImpl struct {
	unique     bool
	keyColHere string
	edgeName   Type
	toTable    Type
	refEdge    *edgeImpl
}

func (e *edgeImpl) RefOn(refEdgeName Type, refEdgeSchemaType Type) IEdge {
	defer func() {
		go func() {
			defer BaseGraphManager.GetWaitGroup().Done()
			for BaseGraphManager.getEdge(refEdgeName, refEdgeSchemaType) == nil {
				runtime.Gosched() // switch other context goroutine implement, save cpu
			}
			refEdge := BaseGraphManager.getEdge(refEdgeName, refEdgeSchemaType)
			e.refEdge, refEdge.refEdge = refEdge, e
			e.keyColHere = string(e.edgeName) + stringJoinRelKey + "Id"
			refEdge.keyColHere = "Id"
			return
		}()
	}()
	return e
}

func (e *edgeImpl) Unique() IEdge {
	//TODO implement me
	panic("implement me")
}

// constructor edgeImpl
func PointTo(edgeName Type, schemaEdgeType Type) IEdge {
	e := &edgeImpl{
		edgeName:   edgeName,
		unique:     false,
		keyColHere: "",
		toTable:    schemaEdgeType,
	}
	BaseGraphManager.addEdge(edgeName, schemaEdgeType, e)
	return e
}

func PointBack(edgeName Type, schemaEdgeType Type) IEdge {
	BaseGraphManager.group.Add(1)
	e := &edgeImpl{
		unique:     false,
		keyColHere: "",
		edgeName:   edgeName,
		toTable:    schemaEdgeType,
	}
	BaseGraphManager.addEdge(edgeName, schemaEdgeType, e)
	return e
}
