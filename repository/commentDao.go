package repository

import (
	"bytedancedemo/dao"
	"bytedancedemo/model"
	"github.com/gookit/slog"
)

type Comment struct {
	//Id         int64     //评论id
	//UserId     int64     //评论用户id
	//VideoId    int64     //视频id
	//Content    string    //评论内容
	//ActionType int64     //发布评论为1，取消评论为2
	//CreatedAt  time.Time //评论发布的日期mm-dd
	//UpdatedAt  time.Time
	model.Comment
}

// TableName 修改表名映射
func (Comment) TableName() string {
	return "comment"
}

func InsertComment(comment model.Comment) (model.Comment, error) {
	c := dao.Comment
	err := c.Create(&comment)
	if err != nil {
		slog.Error(err)
		return model.Comment{}, nil
	}
	return comment, nil
}

func DeleteComment(commentId int64) error {
	c := dao.Comment
	_, err := c.Where(c.ID.Eq(commentId)).First()
	if err != nil {
		return err
	}
	_, err = c.Where(c.ID.Eq(commentId)).Update(c.ActionType, 2)
	if err != nil {
		return err
	}
	return nil
	//var comment Comment
	//// 先查询是否有此评论
	//result := Db.Where("id = ?", commentId).
	//	First(&comment)
	//if result.Error != nil {
	//	return errors.New("del comment is not exist")
	//}
	//// 删除评论，将action_type置为2
	//result = Db.Model(Comment{}).
	//	Where("id=?", commentId).
	//	Update("action_type", 2)
	//if result.Error != nil {
	//	log.Println("Dao-DeleteComment: return del comment failed")
	//	return result.Error
	//}
	//return nil
}

func GetCommentList(videoId int64) ([]*model.Comment, error) {
	c := dao.Comment
	resList, err := c.Where(c.VideoID.Eq(videoId), c.ActionType.Eq("1")).Order(c.CreatedAt).Find()
	if err != nil {
		return nil, err
	}
	return resList, nil
	//var commentList []Comment
	//result := Db.Model(Comment{}).Where(map[string]interface{}{"video_id": videoId, "action_type": 1}).
	//	Order("created_at desc").
	//	Find(&commentList)
	//if result.Error != nil {
	//	log.Println(result.Error)
	//	return commentList, errors.New("get comment list failed")
	//}
	//return commentList, nil
}

func GetCommentCnt(videoId int64) (int64, error) {
	c := dao.Comment
	cnt, err := c.Where(c.VideoID.Eq(videoId), c.ActionType.Eq("1")).Count()
	if err != nil {
		return 0, err
	}
	return cnt, nil
	//var count int64
	//result := Db.Model(Comment{}).Where(map[string]interface{}{"video_id": videoId, "action_type": 1}).
	//	Count(&count)
	//if result.Error != nil {
	//	return 0, errors.New("find comments count failed")
	//}
	//return count, nil
}
