package brest

// Op operation filter type
type Op string

const (
	// And operation for group
	And Op = "and"
	// Or operation for group
	Or Op = "or"
	// Eq operation for attribute (? = ?)
	Eq Op = "eq"
	// Neq operation for attribute (? != ?)
	Neq Op = "neq"
	// In operation for attribute (? IN ?)
	In Op = "in"
	// Nin operation for attribute (? NOT IN ?)
	Nin Op = "nin"
	// Gt operation for attribute (? > ?)
	Gt Op = "gt"
	// Gte operation for attribute (? >= ?)
	Gte Op = "gte"
	// Lt operation for attribute (? < ?)
	Lt Op = "lt"
	// Lte operation for attribute (? <= ?)
	Lte Op = "lte"
	// Lk operation for attribute (? LIKE ?)
	Lk Op = "lk"
	// Nlk operation for attribute (? NOT LIKE ?)
	Nlk Op = "nlk"
	// Ilk operation for attribute (? LOWER LIKE ?)
	Llk Op = "ilk"
	// Nilk operation for attribute (? NOT LOWER LIKE ?)
	Nllk Op = "nilk"
	// Sim operation for attribute (? SIMILAR TO ?)
	Sim Op = "sim"
	// Nsim operation for attribute (? NOT SIMILAR TO ?)
	Nsim Op = "nsim"
	// Lulk operation for attribute (? LOWER UNACCENT LIKE ?)
	Lulk Op = "ilkua"
	// Nlulk operation for attribute (? NOT LOWER UNACCENT ILIKE ?)
	Nlulk Op = "nilkua"
	// Null operation for attribute (? IS NULL)
	Null Op = "null"
	// Nnull operation for attribute (? IS NOT NULL)
	Nnull Op = "nnull"
)

func (o Op) String() string {
	return string(o)
}
