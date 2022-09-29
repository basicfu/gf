package builder

// Operator is interface that should implement by each struct
// that want to be an operator.
type Operator interface {
	GetKey() string
	GetVal() interface{}
}

// BaseOperator is simple base operator that implemented Operator.
type BaseOperator struct {
	key string
	val interface{}
}

func (operator *BaseOperator) GetKey() string {
	return operator.key
}

func (operator *BaseOperator) GetVal() interface{} {
	return operator.val
}

var _ Operator = &BaseOperator{}
