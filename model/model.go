package model

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID       uint   `gorm:"primaryKey;autoIncrement;comment:用户ID" json:"id"`
	Username string `gorm:"type:varchar(50);uniqueIndex;not null;comment:用户名" json:"username"`
	Password string `gorm:"type:varchar(255);not null;comment:密码" json:"-"`
	Email    string `gorm:"type:varchar(100);comment:邮箱" json:"email"`
}

// TableName 指定表名
func (User) TableName() string {
	return "user" // 用户表
}

// Post 博客文章模型
type Post struct {
	ID        uint      `gorm:"primaryKey;autoIncrement;comment:文章ID" json:"id"`
	Title     string    `gorm:"type:varchar(200);not null;comment:文章标题" json:"title"`
	Content   string    `gorm:"type:text;comment:文章内容" json:"content"`
	UserID    uint      `gorm:"not null;index;comment:用户ID" json:"user_id"`
	CreatedAt time.Time `gorm:"autoCreateTime;comment:创建时间" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;comment:更新时间" json:"updated_at"`

	// 关联关系
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName 指定表名
func (Post) TableName() string {
	return "post" // 博客文章表
}

// Comment 评论模型
type Comment struct {
	ID        uint      `gorm:"primaryKey;autoIncrement;comment:评论ID" json:"id"`
	Content   string    `gorm:"type:text;not null;comment:评论内容" json:"content"`
	UserID    uint      `gorm:"not null;index;comment:用户ID" json:"user_id"`
	PostID    uint      `gorm:"not null;index;comment:文章ID" json:"post_id"`
	CreatedAt time.Time `gorm:"autoCreateTime;comment:创建时间" json:"created_at"`

	// 关联关系
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Post Post `gorm:"foreignKey:PostID" json:"post,omitempty"`
}

// TableName 指定表名
func (Comment) TableName() string {
	return "comment" // 评论表
}

// InitTableComments 初始化表注释（MySQL）
func InitTableComments(db *gorm.DB) error {
	// 为表添加注释
	if err := db.Exec("ALTER TABLE user COMMENT '用户表'").Error; err != nil {
		return err
	}
	if err := db.Exec("ALTER TABLE post COMMENT '博客文章表'").Error; err != nil {
		return err
	}
	if err := db.Exec("ALTER TABLE comment COMMENT '评论表'").Error; err != nil {
		return err
	}
	return nil
}
