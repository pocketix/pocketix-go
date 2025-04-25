package models

type ProcedureStore struct {
	Procedures map[string]Procedure
}

func NewProcedureStore() *ProcedureStore {
	return &ProcedureStore{
		Procedures: make(map[string]Procedure),
	}
}

func (ps *ProcedureStore) AddProcedure(procedure Procedure) {
	ps.Procedures[procedure.Name] = procedure
}

func (ps *ProcedureStore) Get(name string) Procedure {
	procedure := ps.Procedures[name]
	return procedure
}

func (ps *ProcedureStore) Has(name string) bool {
	_, exists := ps.Procedures[name]
	return exists
}
