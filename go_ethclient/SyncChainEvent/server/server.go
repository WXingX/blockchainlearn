package service

import (
	"SyncChainEvent/config"
	"SyncChainEvent/logger/xzap"
	"SyncChainEvent/model"
	"SyncChainEvent/server/eventsync"
	"SyncChainEvent/server/pointcount"
	"context"
	"errors"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Service struct {
	ctx              context.Context
	config           *config.Config
	db               *gorm.DB
	syncServer       *eventsync.SyncServer
	pointCountServer *pointcount.CountServer
}

func New(ctx context.Context, cfg *config.Config) (*Service, error) {
	if cfg.ChainCfg == nil {
		return nil, errors.New("chain cfg is empty")
	}

	db := model.NewDB(cfg.DBCfg)
	if db == nil {
		return nil, errors.New("db is nil")
	}

	syncServer, err := eventsync.New(ctx, cfg.ChainCfg.Name, cfg.ChainCfg.ID, cfg.ChainCfg.TokenAddress, cfg.ChainCfg.DeployerAddress, cfg.ChainCfg.RPCUrl, cfg.ChainCfg.DelayBlockNum, db)
	if err != nil {
		xzap.WithContext(ctx).Error("Failed to create eventsync server.", zap.Error(err))
		return nil, err
	}

	pointCountServer, err := pointcount.New(ctx, cfg.ChainCfg.Name, cfg.ChainCfg.ID, cfg.ChainCfg.RPCUrl, db)
	if err != nil {
		xzap.WithContext(ctx).Error("Failed to create pointcount.", zap.Error(err))
		return nil, err
	}

	return &Service{ctx: ctx, config: cfg, db: db, syncServer: syncServer, pointCountServer: pointCountServer}, nil
}

func (s *Service) Start() error {
	if err := s.syncServer.Start(); err != nil {
		xzap.WithContext(s.ctx).Error("eventsync server run failed.", zap.Error(err))
		return err
	}
	s.pointCountServer.Start()
	return nil
}
