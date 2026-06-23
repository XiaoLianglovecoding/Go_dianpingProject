package repository

import (
	"context"
	"log"

	"hmdp-go/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type FollowRepository interface {
	// FindFollow 查询当前用户是否关注了目标用户。
	FindFollow(ctx context.Context, userID int64, followUserID int64) (*model.Follow, error)
	// SaveFollow 新增关注关系。
	SaveFollow(ctx context.Context, follow *model.Follow) error
	// DeleteFollow 删除关注关系。
	DeleteFollow(ctx context.Context, userID int64, followUserID int64) error
	// FindCommonFollows 查询共同关注。
	FindCommonFollows(ctx context.Context, userID int64, otherUserID int64) ([]model.User, error)
	//FollowRepository 增加查粉丝方法
	FindFollowerIDs(ctx context.Context, followUserID int64) ([]int64, error)
}

type followRepository struct {
	// db 是 GORM 数据库连接对象。
	db *gorm.DB
}

// NewFollowRepository 创建关注 Repository。
func NewFollowRepository(db *gorm.DB) FollowRepository {
	return &followRepository{db: db}
}

// FindFollow 后面会查询 tb_follow 是否存在一条关注记录。
func (r *followRepository) FindFollow(ctx context.Context, userID int64, followUserID int64) (*model.Follow, error) {
	var follow model.Follow
	//SELECT * FROM tb_follow WHERE user_id = ? AND follow_user_id = ? LIMIT 1;
	res := r.db.WithContext(ctx).
		Where("user_id = ? AND follow_user_id = ?", userID, followUserID).
		Limit(1).
		Find(&follow)

	// 2. 拦截真正的数据库系统级错误（如断网、表不存在）
	if res.Error != nil {
		return nil, res.Error
	}
	// 3. 核心逻辑：判断有没有查到数据
	// 如果 RowsAffected 为 0，说明没查到记录，即【未关注】
	if res.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound // 手动抛出标准错误，让 Service 层去处理
	}
	// 4. 如果走到了这里，说明查到了数据，即【已关注】
	return &follow, nil
}

// SaveFollow 保存关注关系（新增一条关注记录）。
// 相当于 SQL: INSERT INTO tb_follow (user_id, follow_user_id, create_time) VALUES (?, ?, ?);
func (r *followRepository) SaveFollow(ctx context.Context, follow *model.Follow) error {
	// GORM 的精髓：直接传入结构体指针，它会自动解析字段并执行 INSERT
	//  健壮性防范：防并发双击 (Insert Ignore)
	// 使用 Clauses(clause.OnConflict{DoNothing: true})
	// 相当于 SQL: INSERT IGNORE INTO tb_follow ...
	// 如果用户狂点关注导致触发了联合唯一索引冲突，MySQL 会静默忽略，绝不抛出错误！
	err := r.db.WithContext(ctx).
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(follow).Error
	return err
}

// DeleteFollow 取消关注（硬删除这条关注记录）。
// 相当于 SQL: DELETE FROM tb_follow WHERE user_id = ? AND follow_user_id = ?;
func (r *followRepository) DeleteFollow(ctx context.Context, userID int64, followUserID int64) error {
	// 注意：Delete() 里面必须传一个空的 &model.Follow{}，这样 GORM 才能知道去哪张表删数据
	res := r.db.WithContext(ctx).
		Where("user_id = ? AND follow_user_id = ?", userID, followUserID).
		Delete(&model.Follow{})

	if res.Error != nil {
		return res.Error
	}

	// 健壮性防范：天然幂等校验
	// 如果 RowsAffected 为 0，意味着在删除那一刻记录已经不存在了（比如前端连发了两次取关请求）。
	// 在分布式系统中，这也算作一种“成功”，我们不返回 error，直接放行即可。
	if res.RowsAffected == 0 {
		// 这里可以留个口子，如果未来业务要求严格，可以打条 Warn 日志，但绝对不应该阻断流程。
		log.Printf("记录已经不存在")
	}
	return nil
}

// FindCommonFollows 查询两个用户都关注的人。
func (r *followRepository) FindCommonFollows(ctx context.Context, userID int64, otherUserID int64) ([]model.User, error) {
	var users []model.User

	// 相当于执行以下 SQL 语句（求两个用户共同关注的交集）：
	// SELECT DISTINCT u.* // FROM tb_user AS u
	// JOIN tb_follow AS f1 ON f1.follow_user_id = u.id AND f1.user_id = ?
	// JOIN tb_follow AS f2 ON f2.follow_user_id = u.id AND f2.user_id = ?
	err := r.db.WithContext(ctx).
		Table("tb_user AS u").
		Select("DISTINCT u.*").
		Joins("JOIN tb_follow AS f1 ON f1.follow_user_id = u.id AND f1.user_id = ?", userID).
		Joins("JOIN tb_follow AS f2 ON f2.follow_user_id = u.id AND f2.user_id = ?", otherUserID).
		Scan(&users).Error

	if err != nil {
		return nil, err // 真正的数据库错误
	}
	return users, nil
}

// FollowRepository 查询粉丝IDS
func (r *followRepository) FindFollowerIDs(ctx context.Context, followUserID int64) ([]int64, error) {
	var ids []int64

	err := r.db.WithContext(ctx).
		Model(&model.Follow{}).
		Where("follow_user_id = ?", followUserID).
		Pluck("user_id", &ids).Error

	if err != nil {
		return nil, err
	}
	return ids, nil
}
