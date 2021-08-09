
package data

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/mittacy/blogBack/apierr"
	"github.com/mittacy/blogBack/app/model"
	"github.com/mittacy/blogBack/app/service"
	"github.com/mittacy/blogBack/pkg/logger"
	"github.com/mittacy/blogBack/pkg/store/cache"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"strconv"
)

// 实现service层中的data接口

type Article struct {
	db 	   *gorm.DB
	cache  cache.CustomRedis
	logger *logger.CustomLogger
}

func NewArticle(db *gorm.DB, cacheConn *redis.Pool, logger *logger.CustomLogger) service.IArticleData {
	r := cache.ConnRedisByPool(cacheConn, "article")

	return &Article{
		db:    	db,
		cache: 	r,
		logger: logger,
	}
}

func (ctl *Article) Insert(article *model.Article) error {
	tx := ctl.db.Begin()

	// 创建文章
	if err := tx.Create(&article).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 分类文章+1
	category := model.Category{Id: article.CategoryId}

	result := tx.Model(&category).Update("article_count", gorm.Expr("article_count + ?", 1))

	if result.Error != nil || result.RowsAffected == 0 {
		tx.Rollback()
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return apierr.ErrCategoryNoExist
		}
	}

	tx.Commit()

	if err := ctl.cache.Del(ctl.cacheSumKey(), ctl.cacheSumByCategoryKey(article.CategoryId)); err != nil {
		ctl.logger.CacheErrLog(err)
	}

	return nil
}

func (ctl *Article) Delete(id int64) error {
	article, err := ctl.Get(id)
	if err != nil {
		return err
	}

	// 删除文章
	article.Deleted = model.ArticleDeletedYes

	if err := ctl.UpdateById(article, []string{"deleted"}); err != nil {
		return err
	}

	// 分类减1
	category := model.Category{Id: article.CategoryId}
	if err := ctl.db.Model(&category).Update("article_count", gorm.Expr("article_count - ?", 1)); err != nil {
		ctl.logger.Sugar().Errorf("update category articleCount err: %s", err.Error)
	}

	if err = ctl.cache.Del(ctl.cacheSumKey()); err != nil {
		ctl.logger.CacheErrLog(err)
	}

	return nil
}

func (ctl *Article) UpdateById(article *model.Article, updateFields []string) error {
	if err := ctl.db.Select(updateFields).Updates(article).Error; err != nil {
		return errors.WithStack(err)
	}

	if err := ctl.cache.Del(ctl.cacheByIdKey(article.Id)); err != nil {
		ctl.logger.CacheErrLog(err)
	}

	return nil
}

func (ctl *Article) Get(id int64) (*model.Article, error) {
	// todo 解决查询和更新view+1清空缓存的冲突
	///*
	// * 1. 从 redis 读取
	// * 2. 不存在 -> 数据库查询, 存入缓存
	// * 3. 返回
	// */
	//cacheKey := ctl.cacheByIdKey(id)
	//article := &model.Article{}
	//
	//// 从缓存库查询
	//cacheData, err := redis.Bytes(ctl.cache.Do("get", cacheKey))
	//
	//if err != nil && !errors.Is(err, redis.ErrNil) {
	//	return nil, errors.WithStack(err)
	//}
	//
	//// 反序列化失败，重新从DB查询并缓存
	//if err = json.Unmarshal(cacheData, article); err != nil {
	//	err = redis.ErrNil
	//}
	//
	//
	//// 不存在/反序列化 失败，从数据库查询并存入redis
	//if errors.Is(err, redis.ErrNil) {
	//	article, err = ctl.GetFromDB(id)
	//	if err != nil {
	//		return nil, errors.WithStack(err)
	//	}
	//
	//	// json序列化失败，记录日志，但可以返回成功
	//	cacheData, err = json.Marshal(article)
	//	if err != nil {
	//		ctl.logger.JsonMarshalErrLog(err)
	//		return article, nil
	//	}
	//
	//	// 缓存不成功记录日志，但可以返回成功
	//	if err = ctl.cache.CacheString(ctl.cacheByIdKey(id), string(cacheData)); err != nil {
	//		ctl.logger.CacheErrLog(err)
	//		return article, nil
	//	}
	//}
	//
	//// 返回结果
	//return article, nil
	return ctl.GetFromDB(id)
}

func (ctl *Article) GetFromDB(id int64) (*model.Article, error) {
	article := model.Article{Id: id}

	if err := ctl.db.Where("deleted = ?", model.ArticleDeletedNo).First(&article).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apierr.ErrArticleNoExist
		}
		return nil, errors.WithStack(err)
	}

	return &article, nil
}

func (ctl *Article) GetSum() (int64, error) {
	/*
	 * 1. 从 redis 读取
	 * 2. 不存在 => 数据库查询，存入 redis
	 * 3. 返回
	 */
	// 从缓存库查询
	count, err := redis.Int64(ctl.cache.Do("get", ctl.cacheSumKey()))

	// 缓存查询出错
	if err != nil && !errors.Is(err, redis.ErrNil) {
		return 0, errors.WithStack(err)
	}

	// 不存在, 从数据库查询并存入redis
	if errors.Is(err, redis.ErrNil) {
		count, err = ctl.GetSumFromDB()
		if err != nil {
			return 0, err
		}

		// 缓存不成功只记录错误日志，但可以返回成功
		if err = ctl.cache.CacheString(ctl.cacheSumKey(), strconv.FormatInt(count, 10)); err != nil {
			ctl.logger.CacheErrLog(err)
		}
	}

	// 返回
	return count, nil
}

func (ctl *Article) GetSumByCategory(categoryId int64) (int64, error) {
	/*
	 * 1. 从 redis 读取
	 * 2. 不存在 => 数据库查询，存入 redis
	 * 3. 返回
	 */
	// 从缓存库查询
	count, err := redis.Int64(ctl.cache.Do("get", ctl.cacheSumByCategoryKey(categoryId)))

	// 缓存查询出错
	if err != nil && !errors.Is(err, redis.ErrNil) {
		return 0, errors.WithStack(err)
	}

	// 不存在, 从数据库查询并存入redis
	if errors.Is(err, redis.ErrNil) {
		count, err = ctl.GetSumByCategoryFromDB(categoryId)
		if err != nil {
			return 0, err
		}

		// 缓存不成功只记录错误日志，但可以返回成功
		if err = ctl.cache.CacheString(ctl.cacheSumByCategoryKey(categoryId), strconv.FormatInt(count, 10)); err != nil {
			ctl.logger.CacheErrLog(err)
		}
	}

	// 返回
	return count, nil
}

func (ctl *Article) GetSumFromDB() (int64, error) {
	article := model.Article{}
	var count int64

	err := ctl.db.Model(&article).Select("count(*)").Where("deleted = ?", model.ArticleDeletedNo).Find(&count).Error

	if err != nil {
		return 0, err
	}

	return count, nil
}

func (ctl *Article) GetSumByCategoryFromDB(categoryId int64) (int64, error) {
	article := model.Article{}
	var count int64

	err := ctl.db.Model(&article).Select("count(*)").
		Where("category_id = ? and deleted = ?", categoryId, model.ArticleDeletedNo).Find(&count).Error

	if err != nil {
		return 0, err
	}

	return count, nil
}

func (ctl *Article) List(selectFields []string, page, pageSize int) ([]model.Article, error) {
	startIndex := (page - 1) * pageSize
	var articles []model.Article

	err := ctl.db.Select(selectFields).Where("deleted != ?", 1).
		Offset(startIndex).Limit(pageSize).Order("created_at desc").Find(&articles).Error
	if err != nil {
		return nil, err
	}

	return articles, nil
}

func (ctl *Article) ListByCategory(selectFields []string, categoryId int64, page, pageSize int) ([]model.Article, error) {
	startIndex := (page - 1) * pageSize
	var articles []model.Article

	err := ctl.db.Select(selectFields).Where("category_id = ? and deleted = ?", categoryId, model.ArticleDeletedNo).
		Offset(startIndex).Limit(pageSize).Order("created_at desc").Find(&articles).Error
	if err != nil {
		return nil, err
	}

	return articles, nil
}

func (ctl *Article) ListByWeight(selectFields []string, count int) ([]model.Article, error) {
	var articles []model.Article
	err := ctl.db.Select(selectFields).Where("weight > 0").Order("weight desc, created_at desc").Limit(count).Find(&articles).Error

	if err != nil {
		return nil, err
	}

	return articles, nil
}

func (ctl *Article) IncrView(id int64) error {
	article := model.Article{Id: id}

	if err := ctl.db.Model(&article).Update("views", gorm.Expr("views + ?", 1)).Error; err != nil {
		return err
	}

	// todo 解决查询和更新view+1清空缓存的冲突
	//ctl.cache.Del(ctl.cacheByIdKey(id))

	return nil
}

func (ctl *Article) cacheByIdKey(id int64) string {
	return fmt.Sprintf("%s:id#%d", ctl.cache.CachePrefixKey(), id)
}

func (ctl *Article) cacheSumKey() string {
	return fmt.Sprintf("%s:sum", ctl.cache.CachePrefixKey())
}

func (ctl *Article) cacheSumByCategoryKey(categoryId int64) string {
	return fmt.Sprintf("%s:sum:categoryId#%d", ctl.cache.CachePrefixKey(), categoryId)
}

