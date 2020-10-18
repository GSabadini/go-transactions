package domain

import "errors"

const (
	Debit  string = "DEBIT"
	Credit string = "CREDIT"

	CompraAVista    string = "fd426041-0648-40f6-9d04-5284295c5095"
	CompraParcelada string = "b03dcb59-006f-472f-a8f1-58651990dea6"
	Saque           string = "3f973e5b-cb9f-475c-b27d-8f855a0b90b0"
	Pagamento       string = "976f88ea-eb2f-4325-a106-26f9cb35810d"
)

var (
	ErrOperationInvalid = errors.New("operation type invalid")
)

type (
	// Operation defines the operation entity
	Operation struct {
		id          string
		description string
		opType      string
	}
)

// NewOperation checks if there is an operation and returns it
func NewOperation(id string) (Operation, error) {
	var operations = map[string]Operation{
		CompraAVista: {
			id:          CompraAVista,
			description: "COMPRA A VISTA",
			opType:      Debit,
		},
		CompraParcelada: {
			id:          CompraParcelada,
			description: "COMPRA PARCELADA",
			opType:      Debit,
		},
		Saque: {
			id:          Saque,
			description: "SAQUE",
			opType:      Debit,
		},
		Pagamento: {
			id:          Pagamento,
			description: "PAGAMENTO",
			opType:      Credit,
		},
	}

	operation, exists := operations[id]
	if exists {
		return operation, nil
	}

	return Operation{}, ErrOperationInvalid
}

// ID returns the id property
func (o Operation) ID() string {
	return o.id
}

// Description returns the description property
func (o Operation) Description() string {
	return o.description
}

// Type returns the opType property
func (o Operation) Type() string {
	return o.opType
}
