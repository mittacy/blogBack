
package data

import (
	"encoding/json"
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
	"strings"
)

// 实现service层中的data接口

type User struct {
	db 	   *gorm.DB
	cache  cache.CustomRedis
	logger *logger.CustomLogger
}

func NewUser(db *gorm.DB, cacheConn *redis.Pool, logger *logger.CustomLogger) service.IUserData {
	r := cache.ConnRedisByPool(cacheConn, "user")

	return &User{
		db:    	db,
		cache: 	r,
		logger: logger,
	}
}

// Create 创建用户
// @param user 用户信息
// @return error
func (ctl *User) Create(user *model.User) error {
	if err := ctl.db.Create(user).Error; err != nil {
		if strings.Contains(err.Error(), "Duplicate") {
			if strings.Contains(err.Error(), model.UserIdxName) {
				return apierr.ErrUserNameExist
			} else if strings.Contains(err.Error(), model.UserIdxEmail) {
				return apierr.ErrUserEmailExist
			}
		}

		return errors.WithStack(err)
	}

	return nil
}

// UpdatesById 更新用户信息
// @param user 用户信息
// @param updateFields 更新字段
// @param isCleanCache 是否清除缓存
// @return error
func (ctl *User) UpdatesById(user model.User, updateFields []string, isCleanCache bool) error {
	if err := ctl.db.Select(updateFields).Updates(&user).Error; err != nil {
		return errors.WithStack(err)
	}

	if isCleanCache {
		if err := ctl.cache.Del(ctl.cacheUserKey(user.Id)); err != nil {
			ctl.logger.CacheErrLog(err)
		}
	}

	return nil
}

// Get 查询用户记录
// @param id 用户id
// @return *model.User 用户信息
// @return error
func (ctl *User) Get(id int64) (*model.User, error) {
	/*
	 * 1. 从 redis 读取
	 * 2. 不存在 -> 数据库查询, 存入缓存
	 * 3. 返回
	 */
	user, err := ctl.GetFromCache(id)
	redisNormal := true

	if err != nil && !errors.Is(err, apierr.ErrUserNoExist){
		ctl.logger.CacheErrLog(err)
		redisNormal = false
	}

	// 从缓存库查询失败，从数据库查询并存入redis
	if errors.Is(err, apierr.ErrUserNoExist) {
		user, err = ctl.GetFromDB(id)
		if err != nil {
			return nil, err
		}

		// 缓存数据，如果前面的redis查询错误，就不要再尝试缓存
		if redisNormal {
			if err = ctl.CacheById(user); err != nil {
				// 缓存不成功记录日志，但可以返回成功
				ctl.logger.CacheErrLog(err)
				return user, nil
			}
		}
	}
	return user, nil
}

// GetByName 使用name查询用户记录
// @param name 用户name
// @return *model.User 用户信息
// @return error
func (ctl *User) GetByName(name string) (*model.User, error) {
	// 1. 查询name用户的id
	userId, err := ctl.GetIdByName(name)

	if err != nil {
		return nil, err
	}

	// 2. 使用id查询用户
	return ctl.Get(userId)
}

// GetByEmail 使用email查询用户记录
// @param email 用户email
// @return *model.User 用户信息
// @return error
func (ctl *User) GetByEmail(email string) (*model.User, error) {
	// 1. 查询email用户的id
	userId, err := ctl.GetIdByEmail(email)

	if err != nil {
		return nil, err
	}

	// 2. 使用id查询用户
	return ctl.Get(userId)
}

// GetIdByName 使用name查询用户id
// @param name 用户name
// @return int64 用户id
// @return error
func (ctl *User) GetIdByName(name string) (int64, error) {
	cacheKey := ctl.cacheIdByNameKey(name)

	userId, err := redis.Int64(ctl.cache.Do("get", cacheKey))
	redisNormal := true

	if err != nil && !errors.Is(err, redis.ErrNil) {
		ctl.logger.CacheErrLog(err)
		redisNormal = false
	}

	// 从缓存库查询失败，从数据库查询并存入 redis
	if errors.Is(err, redis.ErrNil) {
		userId, err = ctl.GetIdFromDBByName(name)
		if err != nil {
			return 0, err
		}

		if redisNormal {
			if err = ctl.cache.CacheString(cacheKey, strconv.FormatInt(userId, 10)); err != nil {
				ctl.logger.CacheErrLog(err)
				return userId, nil
			}
		}
	}
	return userId, nil
}

// GetIdByEmail 使用email查询用户id
// @param email 用户email
// @return int64 用户id
// @return error
func (ctl *User) GetIdByEmail(email string) (int64, error) {
	cacheKey := ctl.cacheIdByEmailKey(email)

	userId, err := redis.Int64(ctl.cache.Do("get", cacheKey))
	redisNormal := true

	if err != nil && !errors.Is(err, redis.ErrNil) {
		ctl.logger.CacheErrLog(err)
		redisNormal = false
	}

	// 从缓存库查询失败，从数据库查询并存入 redis
	if errors.Is(err, redis.ErrNil) {
		userId, err = ctl.GetIdFromDBByEmail(email)
		if err != nil {
			return 0, err
		}

		if redisNormal {
			if err = ctl.cache.CacheString(cacheKey, strconv.FormatInt(userId, 10)); err != nil {
				ctl.logger.CacheErrLog(err)
				return userId, nil
			}
		}
	}
	return userId, nil
}

// GetFromDB 通过id从数据库查询用户，不涉及缓存
// @param id 用户id
// @return *model.User
// @return error
func (ctl *User) GetFromDB(id int64) (*model.User, error) {
	user := model.User{Id: id}

	if err := ctl.db.First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apierr.ErrUserNoExist
		}
		return nil, errors.WithStack(err)
	}

	return &user, nil
}

// GetIdFromDBByName 通过name从数据库查询用户id，不涉及缓存
// @param name 用户name
// @return int64 用户id
// @return error
func (ctl *User) GetIdFromDBByName(name string) (int64, error) {
	user := model.User{}

	if err := ctl.db.Select("id").Where("name = ?", name).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, apierr.ErrUserNoExist
		}
		return 0, errors.WithStack(err)
	}

	return user.Id, nil
}

// GetIdFromDBByEmail 通过email从数据库查询用户id，不涉及缓存
// @param email 用户email
// @return int64 用户id
// @return error
func (ctl *User) GetIdFromDBByEmail(email string) (int64, error) {
	user := model.User{}

	if err := ctl.db.Select("id").Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, apierr.ErrUserNoExist
		}
		return 0, errors.WithStack(err)
	}

	return user.Id, nil
}

// GetFromCache 只从缓存查询
// @param id
// @return *model.User
// @return error
func (ctl *User) GetFromCache(id int64) (*model.User, error) {
	u, err := redis.Bytes(ctl.cache.Do("get", ctl.cacheUserKey(id)))
	if err != nil {
		if err == redis.ErrNil {
			return nil, apierr.ErrUserNoExist
		}
		return nil, errors.WithStack(err)
	}

	user := model.User{}
	if err := json.Unmarshal(u, &user); err != nil {
		return nil, errors.WithStack(err)
	}

	return &user, nil
}

// CacheById 缓存用户
// @param user
// @return error
func (ctl *User) CacheById(user *model.User) error {
	data, err := json.Marshal(user)
	if err != nil {
		return errors.WithStack(err)
	}

	if err = ctl.cache.CacheString(ctl.cacheUserKey(user.Id), string(data)); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// cacheUserKey 缓存用户，区分为用户id
// @param id 用户id
// @return string 缓存完整键
func (ctl *User) cacheUserKey(id int64) string {
	return fmt.Sprintf("%s:id#%s", ctl.cache.CachePrefixKey(), strconv.FormatInt(id, 10))
}

// cacheIdByNameKey 缓存用户name和id的映射关系
// @param name 用户name
// @return string 缓存完整键
func (ctl *User) cacheIdByNameKey(name string) string {
	return fmt.Sprintf("%s:name#%s", ctl.cache.CachePrefixKey(), name)
}

// cacheIdByEmailKey 缓存用户email和id的映射关系
// @param email 用户email
// @return string 缓存完整键
func (ctl *User) cacheIdByEmailKey(email string) string {
	return fmt.Sprintf("%s:email#%s", ctl.cache.CachePrefixKey(), email)
}

