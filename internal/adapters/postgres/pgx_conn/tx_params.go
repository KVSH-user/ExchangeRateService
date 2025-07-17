package pgx_conn

import "github.com/jackc/pgx/v5"

// TxParams interface defines a method for applying transaction options.
type TxParams interface {
	Apply(opts pgx.TxOptions) pgx.TxOptions
}

// TxParamsAccessMode holds the access mode for the transaction.
type TxParamsAccessMode struct {
	TxAccessMode pgx.TxAccessMode
}

// Apply method sets the access mode in the provided transaction options.
func (p *TxParamsAccessMode) Apply(opts pgx.TxOptions) pgx.TxOptions {
	opts.AccessMode = p.TxAccessMode
	return opts
}

// TxParamsIsoLevel holds the isolation level for the transaction.
type TxParamsIsoLevel struct {
	TxIsoLevel pgx.TxIsoLevel
}

// Apply method sets the isolation level in the provided transaction options.
func (p *TxParamsIsoLevel) Apply(opts pgx.TxOptions) pgx.TxOptions {
	opts.IsoLevel = p.TxIsoLevel
	return opts
}

// TxParamsDeferrableMode holds the deferrable mode for the transaction.
type TxParamsDeferrableMode struct {
	TxDeferrableMode pgx.TxDeferrableMode
}

// Apply method sets the deferrable mode in the provided transaction options.
func (p *TxParamsDeferrableMode) Apply(opts pgx.TxOptions) pgx.TxOptions {
	opts.DeferrableMode = p.TxDeferrableMode
	return opts
}

// TxParamsBeginQuery holds the query to be executed at the beginning of the transaction.
type TxParamsBeginQuery struct {
	BeginQuery string
}

// Apply method sets the begin query in the provided transaction options.
func (p *TxParamsBeginQuery) Apply(opts pgx.TxOptions) pgx.TxOptions {
	opts.BeginQuery = p.BeginQuery
	return opts
}
