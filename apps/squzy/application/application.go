package application

import (
	"fmt"
	serverPb "github.com/squzy/squzy_generated/generated/server/proto/v1"
	"google.golang.org/grpc"
	"net"
	"squzy/apps/internal/httpTools"
	scheduler_storage "squzy/apps/internal/scheduler-storage"
	"squzy/apps/internal/semaphore"
	sitemap_storage "squzy/apps/internal/sitemap-storage"
	"squzy/apps/internal/storage"
	"squzy/apps/squzy/server"
)

type app struct {
	schedulerStorage scheduler_storage.SchedulerStorage
	externalStorage  storage.Storage
	siteMapStorage   sitemap_storage.SiteMapStorage
	tool             httpTools.HttpTool
	semaphoreFactory semaphore.SemaphoreFactory
}

func New(
	schedulerStorage scheduler_storage.SchedulerStorage,
	externalStorage storage.Storage,
	siteMapStorage sitemap_storage.SiteMapStorage,
	tool httpTools.HttpTool,
	semaphoreFactory semaphore.SemaphoreFactory,
) *app {
	return &app{
		schedulerStorage,
		externalStorage,
		siteMapStorage,
		tool,
		semaphoreFactory,
	}
}

func (s *app) Run(port int32) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	grpcServer := grpc.NewServer()
	serverPb.RegisterServerServer(
		grpcServer,
		server.New(
			s.schedulerStorage,
			s.externalStorage,
			s.siteMapStorage,
			s.tool,
			s.semaphoreFactory,
		),
	)
	return grpcServer.Serve(lis)
}
