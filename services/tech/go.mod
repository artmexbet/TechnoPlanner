module tech

go 1.25.4

require (
	github.com/artmexbet/TechnoPlanner/libs/config v0.0.0
	github.com/artmexbet/TechnoPlanner/libs/dto v0.0.0
	github.com/go-playground/validator/v10 v10.30.1
	github.com/google/uuid v1.6.0
	github.com/jackc/pgx/v5 v5.7.6
	github.com/mailru/easyjson v0.9.1
	github.com/nats-io/nats.go v1.48.0
)

replace (
	github.com/artmexbet/TechnoPlanner/libs/broker => ../../libs/broker
	github.com/artmexbet/TechnoPlanner/libs/config => ../../libs/config
	github.com/artmexbet/TechnoPlanner/libs/dto => ../../libs/dto
)

require (
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/gabriel-vasile/mimetype v1.4.13 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/klauspost/compress v1.18.4 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/nats-io/nkeys v0.4.15 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/stretchr/testify v1.11.1 // indirect
	golang.org/x/crypto v0.48.0 // indirect
	golang.org/x/sync v0.19.0 // indirect
	golang.org/x/sys v0.41.0 // indirect
	golang.org/x/text v0.34.0 // indirect
)
