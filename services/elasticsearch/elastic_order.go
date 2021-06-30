package elasticsearch

const (
	DESC = "desc"
	ASC  = "asc"
)

type ElasticOrder struct {
	Order Meta
}

func NewElasticOrder() *ElasticOrder {
	return &ElasticOrder{Order: make(Meta, 0)}
}

func (o *ElasticOrder) Add(field string, order ...string) *ElasticOrder {
	finalOrder := DESC
	if len(order) > 0 {
		finalOrder = order[0]
	}
	o.Order = append(o.Order, map[string]interface{}{field: map[string]interface{}{"order": finalOrder}})
	return o
}

func (o *ElasticOrder) AddScore(order ...string) *ElasticOrder {
	o.Add("_score", order...)
	return o
}

func (o *ElasticOrder) Clear() *ElasticOrder {
	o.Order = make(Meta, 0)
	return o
}
