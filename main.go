package main

import (
	"go-base/config"
	"go-base/database"
	_ "go-base/docs"
	"go-base/router"
	"go-base/utils"
	"os"
	"strconv"
)

// @title go-base API
// @version 1.0
// @description go-base project API documentation
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	if err := utils.InitLogger("gin.log"); err != nil {
		_, _ = os.Stderr.WriteString("failed to initialize logger: " + err.Error() + "\n")
		return
	}
	// 1. 初始化配置
	if err := config.InitConfig("config/config.yaml"); err != nil {
		utils.Errorf("❌ Failed to load config: %v", err)
		return
	}
	utils.Infof("✅ Config loaded successfully")

	// 2. 初始化数据库
	db, err := database.InitDatabase()
	if err != nil {
		utils.Errorf("❌ Failed to initialize database: %v", err)
		return
	}

	// 3. 设置路由
	r := router.SetupRouter(db)

	// 4. 启动服务器
	port := config.GlobalConfig.Server.Port
	utils.Infof("🚀 Server is running on http://localhost:%d", port)
	utils.Infof("Press Ctrl+C to stop the server")

	if err := r.Run(":" + strconv.Itoa(port)); err != nil {
		utils.Errorf("❌ Failed to start server: %v", err)
	}
}
