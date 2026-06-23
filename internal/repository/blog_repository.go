package repository

import (
	"context"
	"errors"

	"hmdp-go/internal/model"

	"gorm.io/gorm"
)

type BlogRepository interface {
	// SaveBlog 保存一篇新博客。
	SaveBlog(ctx context.Context, blog *model.Blog) error
	// FindBlogByID 根据博客 id 查询一篇博客。
	FindBlogByID(ctx context.Context, id int64) (*model.Blog, error)
	// FindBlogsByHot 查询热门博客列表。
	FindBlogsByHot(ctx context.Context, current int) ([]model.Blog, error)
	// FindBlogsByUserID 查询某个用户发布的博客。
	FindBlogsByUserID(ctx context.Context, userID int64, current int) ([]model.Blog, error)
	// UpdateBlogLiked 修改博客点赞数，delta 可以是 +1 或 -1。
	UpdateBlogLiked(ctx context.Context, id int64, delta int) error
	//BlogRepository 批量查博客
	FindBlogsByIDs(ctx context.Context, ids []int64) ([]model.Blog, error)
}

type blogRepository struct {
	// db 是 GORM 数据库连接对象。
	db *gorm.DB
}

// NewBlogRepository 创建博客 Repository。
func NewBlogRepository(db *gorm.DB) BlogRepository {
	return &blogRepository{db: db}
}

// SaveBlog 负责向 tb_blog 插入一条博客记录。
func (r *blogRepository) SaveBlog(ctx context.Context, blog *model.Blog) error {
	// 使用 GORM 插入新记录，GORM 会自动将生成的主键自增 ID 写回 blog.ID
	return r.db.WithContext(ctx).Create(blog).Error
}

// FindBlogByID 负责按主键查询博客详情。
func (r *blogRepository) FindBlogByID(ctx context.Context, id int64) (*model.Blog, error) {
	var blog model.Blog
	err := r.db.WithContext(ctx).
		First(&blog, id).Error
	if err != nil {
		return nil, err
	}
	return &blog, nil
}

// FindBlogsByHot 查询热门博客。
//
// 当前实现先按 liked 点赞数倒序排序，每页 10 条。
// 对应 SQL 大概是:
// SELECT * FROM tb_blog ORDER BY liked DESC LIMIT 10 OFFSET ?;
func (r *blogRepository) FindBlogsByHot(ctx context.Context, current int) ([]model.Blog, error) {
	var blogs []model.Blog

	// current 是页码，前端没传时默认是 1。这里兜底避免非法页码。
	if current < 1 {
		current = 1
	}

	// pageSize 是每页数量；offset 是从第几条开始查。
	const pageSize = 10
	offset := (current - 1) * pageSize

	// Order("liked DESC") 表示点赞数从高到低；
	// Offset + Limit 是 MySQL 分页；
	// Find(&blogs) 查询多条博客记录。
	err := r.db.WithContext(ctx).Order("liked DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&blogs).Error

	if err != nil {
		return nil, err
	}
	return blogs, nil
}

// FindBlogsByUserID 后面用于个人主页，查询某个用户的博客列表。
func (r *blogRepository) FindBlogsByUserID(ctx context.Context, userID int64, current int) ([]model.Blog, error) {
	var blogs []model.Blog
	if current < 1 {
		current = 1
	}
	const pageSize = 10
	offset := (current - 1) * pageSize
	err := r.db.WithContext(ctx).Order("create_time DESC").
		Where("user_id = ?", userID).
		Offset(offset).
		Limit(pageSize).
		Find(&blogs).Error
	if err != nil {
		return nil, err
	}
	return blogs, nil
}

// UpdateBlogLiked 后面用于点赞/取消点赞时更新 tb_blog.liked。
// UpdateBlogLiked 更新博客点赞数。delta 传 1 代表增加，传 -1 代表减少。
func (r *blogRepository) UpdateBlogLiked(ctx context.Context, blogId int64, delta int) error {
	// 等价于 SQL: UPDATE tb_blog SET liked = liked + (delta) WHERE id = blogId;
	// 1. 初始化基础查询条件
	query := r.db.WithContext(ctx).
		Model(&model.Blog{}).
		Where("id = ?", blogId)

	// 2. 终极防线：如果是取消点赞 (delta < 0)，强制要求 liked 必须大于 0 才能扣减
	if delta < 0 {
		query = query.Where("liked > 0")
	}

	// 3. 执行原子更新，并将结果暂存到 res 变量中
	// 等价 SQL: UPDATE tb_blog SET liked = liked + (delta) WHERE id = ? [AND liked > 0];
	// res 是一个 *gorm.DB 对象，里面包含 .Error 和 .RowsAffected
	res := query.UpdateColumn("liked", gorm.Expr("liked + ?", delta))

	// 4. 先判断有没有系统级的 SQL 错误（比如断网、语法错）
	if res.Error != nil {
		return res.Error
	}

	// 5. 核心防线：判断到底有没有真正更新到数据！
	if res.RowsAffected == 0 {
		// 结合我们的条件，RowsAffected 为 0 只有两种可能：
		// a. 这篇 blogId 根本不存在！
		// b. 想要取消点赞，但数据库里 liked 已经是 0 了！
		// 无论哪种情况，都说明这条操作是无效的，人为抛出错误！
		return errors.New("操作失败: 博客不存在或状态已变更")
	}
	return nil
}

// FindBlogsByIDs批量查博客
func (r *blogRepository) FindBlogsByIDs(ctx context.Context, ids []int64) ([]model.Blog, error) {
	if len(ids) == 0 {
		return []model.Blog{}, nil
	}

	var blogs []model.Blog
	err := r.db.WithContext(ctx).
		Where("id IN ?", ids).
		Find(&blogs).Error

	if err != nil {
		return nil, err
	}
	return blogs, nil
}
