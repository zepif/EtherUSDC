package handlers

import (
	"context"
	"net/http"

	"gitlab.com/distributed_lab/logan/v3"
    "github.com/zepif/EtherUSDC/internal/data"
    "github.com/zepif/EtherUSDC/internal/config"
)

type ctxKey int

const (
	logCtxKey ctxKey = iota
    dbCtxKey
    ethConfigCtxKey
)

func CtxLog(entry *logan.Entry) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, logCtxKey, entry)
	}
}

func Log(r *http.Request) *logan.Entry {
	return r.Context().Value(logCtxKey).(*logan.Entry)
}

func CtxDB(entry data.MasterQ) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, dbCtxKey, entry)
	}
}

func DB(r *http.Request) data.MasterQ {
	return r.Context().Value(dbCtxKey).(data.MasterQ).New()
}

func CtxEthConfig(entry *config.EthConfig) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, ethConfigCtxKey, entry)
	}
}

func EthConfig(r *http.Request) *config.EthConfig {
	return r.Context().Value(ethConfigCtxKey).(*config.EthConfig)
}

