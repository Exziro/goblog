package article

import (
	"goblog/app/models"
	"goblog/pkg/logger"
	"goblog/pkg/pagination"
	"goblog/pkg/route"
	"net/http"

	"goblog/pkg/types"
)

//GET通过ID获取文章
func Get(idstr string) (Article, error) {
	var article Article
	id := types.StringToUint64(idstr)
	if err := models.DB.Preload("User").First(&article, id).Error; err != nil {
		return article, err
	}
	return article, nil
}

//获取全部文章
func GetAll(r *http.Request, Perpage int) ([]Article, pagination.ViewData, error) {
	//初始化分页
	db := models.DB.Model(Article{}).Order("created_at desc")
	_pager := pagination.New(r, db, route.Name2URL("articles.index"), Perpage)
	//获取视图数据
	viewData := _pager.Paging()
	//获取数据
	var articles []Article
	_pager.Results(&articles)

	return articles, viewData, nil
}
func (articles *Article) Create() (err error) {
	result := models.DB.Create(&articles)
	if err := result.Error; err != nil {
		logger.LogError(err)
		return err
	}
	return nil
}
func (articles *Article) Update() (rowsAffected int64, err error) {
	result := models.DB.Save(&articles)
	if err = result.Error; err != nil {
		return 0, err
	}
	return result.RowsAffected, nil
}
func (articles *Article) Delete() (rowsAffected int64, err error) {
	result := models.DB.Delete(&articles)
	if err = result.Error; err != nil {
		return 0, err
	}
	return result.RowsAffected, nil
}

func GetByUserID(uid string) ([]Article, error) {
	var articles []Article
	if err := models.DB.Where("user_id = ?", uid).Preload("User").Find(&articles).Error; err != nil {
		return articles, err
	}
	return articles, nil
}

//GetByCategoryID 根据分类查找文章
func GetByCategoryID(cid string, r *http.Request, perPage int) ([]Article, pagination.ViewData, error) {
	// 1. 初始化分页实例
	db := models.DB.Model(Article{}).Where("category_id = ?", cid).Order("created_at desc")
	_pager := pagination.New(r, db, route.Name2URL("categories.show", "id", cid), perPage)

	// 2. 获取视图数据
	viewData := _pager.Paging()

	// 3. 获取数据
	var articles []Article
	_pager.Results(&articles)

	return articles, viewData, nil

}
