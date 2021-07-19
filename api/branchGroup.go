package api

import (
	"git.npcompete.com/OSPSM_Servers/src/api/handlers"
	"github.com/labstack/echo"
)

func BranchGroup(branchGroup *echo.Group) {
	branchGroup.POST("/setpassword", handlers.RegisterBranch)
	branchGroup.POST("/updatebranch", handlers.UpdateBranch)
	branchGroup.GET("/branchdetail", handlers.GetBranchDetail)
	branchGroup.DELETE("/deletebranch", handlers.DeleteBranch)
	branchGroup.POST("/addproducts", handlers.AddProduct)
	branchGroup.GET("/listproducts", handlers.GetProductList)
	branchGroup.DELETE("/deleteproduct", handlers.DeleteProduct)
	branchGroup.PUT("/updateproduct", handlers.UpdateProduct)
	branchGroup.POST("/makeinvoice", handlers.MakeInvoice)
	// adminGroup.GET("/employeelist", handlers.GetEmployeeList)

}
