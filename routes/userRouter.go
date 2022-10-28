package routes

import (
	"golang_cms/controller"
	// "golang_cms/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	// incomingRoutes.Use(middleware.Authentication)
	incomingRoutes.GET("/users", controller.GetUsers)
	incomingRoutes.GET("/users/:user_id", controller.GetUser)
	incomingRoutes.GET("/user/:Email", controller.GetUserEmail)
	incomingRoutes.PUT("/user/:update", controller.UpdateUser)

	//banner
	incomingRoutes.POST("/banner", controller.CreateBanner)             //memasukan data banner baru
	incomingRoutes.GET("/banner/:bannerId", controller.GetBanner)      //mengambil satu data menggunakan filter ID
	incomingRoutes.PUT("/banner/:bannerId", controller.EditBanner)     //mengedit satu data menggunaakn filter ID
	incomingRoutes.DELETE("/banner/:bannerId", controller.DeleteBanner) //menghapus satu data menggunaakn filter ID
	incomingRoutes.GET("/banners", controller.GetAllBanner)             // mengambil semuah data Banner
	//meta
	incomingRoutes.POST("/meta", controller.CreateMeta)           //memasukan data meta baru
	incomingRoutes.GET("/meta/:metaId", controller.GetMeta)      //mengambil satu data meta dengan filter ID
	incomingRoutes.PUT("/meta/:metaId", controller.EditMeta)     //mengedit satu data meta dengan filter ID
	incomingRoutes.DELETE("/meta/:metaId", controller.DeleteMeta) //menghapus satu data dengan filter ID
	incomingRoutes.GET("/metas", controller.GetAllMeta)           //mengambill semuah data meta
	//Kategori Produk main
	incomingRoutes.POST("kategori", controller.CreateKategori)               //memasukan data baru pada kategori_produk
	incomingRoutes.GET("kategori/:kategoriid", controller.GetKategori)      //memanggil satu data kategori denga filter ID
	incomingRoutes.PUT("kategori/:kategoriid", controller.EditKategori)      //mengedit satu data kategori dengan filter ID
	incomingRoutes.DELETE("kategori/:kategoriid", controller.DeleteKategori) //menghaspus satu data kategori dengan filter ID
	incomingRoutes.GET("kategori", controller.GetAllKategori)                //mengambil semuah data kategori produk
}