// Package cqrslite bridges businessrules validation events onto a
// go-cqrs-lite event bus, turning validation runs into durable domain
// events that projections and read models can consume.
//
// # Wire contract
//
// Payloads are adapter-owned DTOs (RuleEvaluatedData, ValidationCompletedData)
// with plain fields, deliberately decoupled from the businessrules structs so
// the wire format stays stable while the library evolves. Events are encoded
// with the go-cqrs-lite default codec; decode with event.DecodePayloadAuto.
//
// Event types:
//   - "businessrules.rule.evaluated.v1"
//   - "businessrules.validation.completed.v1"
//
// # Usage
//
//	bus := watermill.NewEventBus()
//	listener := cqrslite.NewBusListener(bus, streamID, "Order",
//	    cqrslite.WithPublishErrorHandler(func(err error) { slog.Error(err.Error()) }))
//
//	businessrules.NewValidator().
//	    WithListener(listener).
//	    AddRules(rules...).
//	    Build()
//
// # Publishing constraints
//
// The businessrules Listener contract is synchronous and error-free, so
// publish failures cannot abort the validation run: they are reported through
// the WithPublishErrorHandler handler (default: silently ignored — always set
// one in production). Versions start at 1 and increase by one per published
// event, giving consumers a per-stream ordering guarantee.
//
// # Module note
//
// This is a nested module: it depends on go-cqrs-lite, which the root
// go-business-rules/v2 module deliberately does not. Within this repository
// the parent dependency is satisfied by a replace directive to ../..;
// publishing this module publicly first requires the parent module's v2.1.0
// tag to be pushed (see TODO_LIST.md).
package cqrslite
