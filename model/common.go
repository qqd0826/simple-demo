package model

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

type FavoriteData struct {
	UserId     int64
	VideoId    int64
	IsFavorite bool
	Time       int64
}

type Video struct {
	Id            int64  `json:"id,omitempty"`
	Author        User   `json:"author"`
	AuthorId      int64  //作者id
	PlayUrl       string `json:"play_url" json:"play_url,omitempty"`
	CoverUrl      string `json:"cover_url,omitempty"`
	FavoriteCount int64  `json:"favorite_count,omitempty"`
	CommentCount  int64  `json:"comment_count,omitempty"`
	IsFavorite    bool   `json:"is_favorite,omitempty"`
	UpLoadTime    int64  // 额外加的上传时间，用于排序
	Title         string `json:"title,omitempty"` // 视频标题
}

type Comment struct {
	Id         int64  `json:"id,omitempty"`
	VideoId    int64  `json:"-"`
	User       User   `json:"user"`
	UserId     int64  `json:"-"`
	Content    string `json:"content,omitempty"`
	CreateDate string `json:"create_date,omitempty"`
}

type User struct {
	Id            int64  `json:"id,omitempty"`       //唯一序号字段
	Username      string `json:"username,omitempty"` //用户名（唯一）
	Name          string `json:"name,omitempty"`     //用户昵称
	FollowCount   int64  `json:"follow_count,omitempty"`
	FollowerCount int64  `json:"follower_count,omitempty"`
	IsFollow      bool   `json:"is_follow,omitempty"`
	Password      string `json:"password,omitempty"` //密码

	WorkCount     int64 `json:"work_count,omitempty"`     // 作品数量
	FavoriteCount int64 `json:"favorite_count,omitempty"` // 点赞数量
}

type Message struct {
	Id         int64  `json:"id,omitempty"`
	Content    string `json:"content,omitempty"`
	CreateTime string `json:"create_time,omitempty"`
}

type MessageSendEvent struct {
	UserId     int64  `json:"user_id,omitempty"`
	ToUserId   int64  `json:"to_user_id,omitempty"`
	MsgContent string `json:"msg_content,omitempty"`
}

type MessagePushEvent struct {
	FromUserId int64  `json:"user_id,omitempty"`
	MsgContent string `json:"msg_content,omitempty"`
}
