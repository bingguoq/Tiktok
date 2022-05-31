package dao

import (
	"Project/models"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
	"strconv"
	"time"
)

// UserLogin 用户登录，查询数据库中的信息是否匹配
// 参数 :
//      userId : 返回的结果
//		username : 请求的用户名
// 		userPassword : 请求的用户密码
// 返回值：
//		如果请求的用户名和密码正确，返回 nil ，否则返回 错误

func UserLogin(userId *int64, username string, userPassword string) error {
	var user models.User
	// 指定只查找 user_id 和 password，其他数据不需要
	DB.Select([]string{"user_id", "password"}).
		// 筛选 name = username 的数据
		Where(&models.User{Name: username}).
		First(&user)
	// 判断请求的密码是否和数据库中的密码相匹配
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userPassword))
	if err != nil {
		// 如果密码不匹配，返回 错误
		log.Println(err)
		return err
	}
	// 设置返回的 user_id
	*userId = user.ID
	return nil

}

// UserRegister 用户注册，查询用户名是否存在，不存在注册新用户
// 参数 :
//		username : 请求的用户名
// 		userPassword : 请求的用户密码
// 返回值：
//		如果注册成功，返回 userid ，否则返回 -1，标明用户名已存在，注册失败

func UserRegister(username string, password string) int64 {
	var user models.User
	// 查询数据库username是否存在
	err := DB.Select("user_id").Where("name=?", username).Take(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) { // 用户名不存在,进行注册
		// 使用bcrypt对密码加密
		pd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		// 密码加密失败
		if err != nil {
			log.Println(err)
			return -1 //	注册失败返回-1
		}
		// 在数据库中创建用户数据
		newUser := models.User{
			ID:            -1,
			Name:          username,
			Password:      string(pd),
			FollowCount:   0,
			FollowerCount: 0,
			CreateTime:    time.Now(),
			UpdateTime:    time.Now(),
		}
		if newUser.ID, err = RegisterIdWorker.NextId(); err != nil {
			log.Println(err)
			return -1 //	ID 生成失败返回-1
		}
		// 插入新用户进 user 表
		DB.Create(&newUser)
		return newUser.ID
	} else { // 用户名已存在或者查询失败
		if err != nil { // 查询失败，记录err
			log.Println(err)
		}
		return -1 //	注册失败返回-1
	}
}

// GetUserInfo 查询用户信息，返回用户的所有公开信息
// 参数 :
//      queryId : 被查询者的 id
//		userId : 用户的 id
// 返回值：
//		查询出来用户，如果查询出错或无该用户，返回 nil

func GetUserInfo(queryId int64, userId int64) models.User {
	var user models.User // 结果
	// 查看是否关注该用户
	queryFollow := DB.Select("follow.user_id, 1 as is_follow").Where("follower_id = ?", userId).Table("follow")
	// 查询该用户信息
	DB.Select("user.*, is_follow").
		// 条件筛选，按 user_id 查找
		Where("user.user_id = ?", queryId).
		// 联结是否关注
		Joins("LEFT JOIN (?) AS fo ON fo.user_id = user.user_id", queryFollow).
		First(&user)
	return user
}

// GetFollowerUserList 查询粉丝信息，返回粉丝列表
// 参数 :
//		userId : 用户的 id
// 返回值：
//		返回根据用户id查询出的粉丝列表
func GetFollowerUserList(id string) []models.User {
	var user []models.User // 结果
	// 查询用户的粉丝id
	var userId int64
	userId, _ = strconv.ParseInt(id, 10, 64)
	// 查询 follow
	DB.Table("user").Select("user.user_id,user.name,user.`follow_count`,user.`follower_count`,1 AS is_follow").
		Joins("LEFT JOIN `follow` AS f ON user.`user_id`=f.`follower_id` AND user.`user_id`= ? ", userId).Find(&user)
	return user
}
