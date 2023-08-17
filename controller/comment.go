package controller

import (
	"github.com/RaymondCode/simple-demo/model"

	"github.com/RaymondCode/simple-demo/service"
	"github.com/RaymondCode/simple-demo/utils/sensetive"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"time"
)

type CommentListResponse struct {
	Response
	CommentList []service.Comment `json:"comment_list,omitempty"`
}

type CommentActionResponse struct {
	Response
	Comment service.Comment `json:"comment,omitempty"`
}

// CommentAction no practical effect, just check if token is valid
func CommentAction(c *gin.Context) {
	// 获取userId
	log.Println("Controller_Comment_Action: run") //函数已运行
	userId := c.GetInt64("userId")
	videoId, err := strconv.ParseInt(c.Query("video_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, CommentActionResponse{
			Response: Response{StatusCode: -1,
				StatusMsg: "comment videoId json invalid"},
		})
		return
	}
	actionType, err := strconv.ParseInt(c.Query("action_type"), 10, 64)
	if err != nil || actionType < 1 || actionType > 2 {
		c.JSON(http.StatusOK, CommentActionResponse{
			Response: Response{StatusCode: -1,
				StatusMsg: "comment actionType json invalid"},
		})
		return
	}
	commentService := service.GetCommentServiceInstance()
	//switch {
	// 评论
	if actionType == 1 {
		content := c.Query("comment_text")
		content = sensetive.Filter.Replace(content, '#')
		var comment model.Comment
		comment.UserID = userId
		comment.VideoID = videoId
		comment.Content = content
		timeNow := time.Now()
		comment.CreatedAt = timeNow
		comment.ActionType = strconv.Itoa(1)
		//调用service发表评论
		commentRes, err := commentService.CommentAction(comment)
		//发表品论失败
		if err != nil {
			c.JSON(http.StatusOK, CommentActionResponse{
				Response: Response{StatusCode: -1,
					StatusMsg: "comment failed"},
			})
			return
		}
		//发表评论成功
		c.JSON(http.StatusOK, CommentActionResponse{
			Response: Response{StatusCode: 0,
				StatusMsg: "comment success"},
			Comment: commentRes,
		})
		return

		// 取消评论
	} else {
		commentId, err := strconv.ParseInt(c.Query("comment_id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusOK, CommentActionResponse{
				Response: Response{StatusCode: -1,
					StatusMsg: "delete commentId invalid"},
			})
			return
		}
		//删除评论操作
		err = commentService.DeleteCommentAction(commentId)
		if err != nil {
			c.JSON(http.StatusOK, CommentActionResponse{
				Response: Response{StatusCode: -1,
					StatusMsg: err.Error()},
			})
			return
		}
		c.JSON(http.StatusOK, CommentActionResponse{
			Response: Response{StatusCode: 0,
				StatusMsg: "delete commentId success"},
		})
		log.Println("controller-comment_Action:  delete success")
		return
	}
}

// CommentList all videos have same demo comment list
func CommentList(c *gin.Context) {
	userId := c.GetInt64("userId")
	videoId, err := strconv.ParseInt(c.Query("video_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, CommentListResponse{
			Response: Response{StatusCode: -1,
				StatusMsg: "comment videoId json invalid"},
		})
		return
	}
	commentService := service.GetCommentServiceInstance()
	commentList, err := commentService.GetCommentList(videoId, userId)
	// 获取评论列表失败
	if err != nil {
		c.JSON(http.StatusOK, CommentListResponse{
			Response: Response{StatusCode: -1,
				StatusMsg: err.Error()},
		})
		return
	}
	// 获取评论列表成功
	c.JSON(http.StatusOK, CommentListResponse{
		Response:    Response{StatusCode: 0},
		CommentList: commentList,
	})
	return
}
