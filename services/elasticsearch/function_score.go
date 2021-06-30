package elasticsearch

type FunctionScore struct {
	FieldValueFactor map[string]interface{}
	ScriptScore      map[string]interface{}
	BoostMode        *string
	MaxBoost         *float64
}

func NewFunctionScore() *FunctionScore {
	return new(FunctionScore)
}

func NewLog1pFunctionScore(field string) *FunctionScore {
	functionScore := new(FunctionScore)
	functionScore.SetFieldValueFactor(map[string]interface{}{
		"field":    field,
		"modifier": "log1p",
	})
	return functionScore
}

func (receiver *FunctionScore) SetFieldValueFactor(fieldValueFactor map[string]interface{}) *FunctionScore {
	receiver.FieldValueFactor = fieldValueFactor
	return receiver
}

func (receiver *FunctionScore) SetBoostMode(boostMode string) *FunctionScore {
	receiver.BoostMode = &boostMode
	return receiver
}

func (receiver *FunctionScore) SetMaxBoost(maxBoost float64) *FunctionScore {
	receiver.MaxBoost = &maxBoost
	return receiver
}

func (receiver *FunctionScore) SetScriptScore(scriptScore map[string]interface{}) *FunctionScore {
	receiver.ScriptScore = scriptScore
	return receiver
}

func (receiver *FunctionScore) Build() map[string]interface{} {
	data := make(map[string]interface{})
	if len(receiver.FieldValueFactor) > 0 {
		data["field_value_factor"] = receiver.FieldValueFactor
	}
	if len(receiver.ScriptScore) > 0 {
		data["script_score"] = receiver.ScriptScore
	}
	if receiver.MaxBoost != nil {
		data["max_boost"] = *receiver.MaxBoost
	}

	if receiver.BoostMode != nil {
		data["boost_mode"] = *receiver.BoostMode
	}
	return data
}
