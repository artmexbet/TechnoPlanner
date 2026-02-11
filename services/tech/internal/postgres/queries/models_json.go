package queries

//go:generate easyjson -all models_json.go

type AdditionalCharacteristics struct {
	Chars map[string]string
}
