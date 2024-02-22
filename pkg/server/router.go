package server

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"os-virt/pkg/server/apis"
	"os-virt/pkg/server/handler"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func InitContext(ctx context.Context, client client.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("client", client)
		c.Set("ctx", ctx)
	}
}
func NewRouter(ctx context.Context, client client.Client) *gin.Engine {
	router := gin.New()

	router.Use(handler.GlobalExceptionCatch)
	router.Use(InitContext(ctx, client))

	router.GET("/healthz", func(context *gin.Context) {
		context.String(http.StatusOK, "status ok")
	})

	router.GET("/readyz", func(context *gin.Context) {
		context.String(http.StatusOK, "status ready")
	})
	//操作日志处理
	//opLogHandler := &oplog.OpLogHandler{}
	//认证处理
	//auth := &auth.UserAuth{}
	// 虚拟机
	vmGroup := router.Group("/vms")
	{
		vmApi := apis.NewVmApi()
		vmGroup.GET("/:vmName", vmApi.GetVmDetail)
		vmGroup.GET("", vmApi.GetVms)
		vmGroup.DELETE("/:vmName", vmApi.DeleteVm)
		vmGroup.POST("", vmApi.CreateVm)
	}

	// 镜像
	imageGroup := router.Group("/images")
	{
		imageApi := apis.NewImageApi()
		imageGroup.GET("/detail", imageApi.GetImageDetail)
	}

	// 模板
	templateGroup := router.Group("/templates")
	{
		templateApi := apis.NewTemplateApi()
		templateGroup.POST("/detail", templateApi.GetTemplateDetail)
	}

	instanceGroup := router.Group("/instances")
	{
		instanceApi := apis.NewInstanceApi()
		instanceGroup.POST("/detail", instanceApi.GetInstanceDetail)
	}

	// 网络
	networkGroup := router.Group("/networks")
	{
		networkApi := apis.NewNetworkApi()
		networkGroup.POST("/detail", networkApi.GetNetworkDetail)
	}

	volumeGroup := router.Group("/volumes")
	{
		volumeApi := apis.NewVolumeApi()
		volumeGroup.POST("/detail", volumeApi.GetVolumeDetail)
	}

	return router
}
