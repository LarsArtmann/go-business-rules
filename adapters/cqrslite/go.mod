module github.com/LarsArtmann/go-business-rules/adapters/cqrslite

go 1.26.7

require (
	github.com/LarsArtmann/go-business-rules v0.0.0
	github.com/larsartmann/go-cqrs-lite/event/v4 v4.11.0
	github.com/larsartmann/go-cqrs-lite/event/v4/eventtest v0.4.0
	github.com/larsartmann/go-cqrs-lite/id/v4 v4.6.0
)

replace github.com/LarsArtmann/go-business-rules => ../..
