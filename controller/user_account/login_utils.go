package user_account

import (
	"errors"

	"github.com/TianYi2000/simple-tiktok/models"
)

// AccountHandler 注册用户并得到token和id
func AccountHandler(username, password string) (*AccountResponse, error) {
	return NewAccountHandlerFlow(username, password).Do()
}

func NewAccountHandlerFlow(username, password string) *AccountHandlerFlow {
	return &AccountHandlerFlow{username: username, password: password}
}

type AccountHandlerFlow struct {
	username string
	password string

	data   *AccountResponse
	userid int64
	token  string
}

func (q *AccountHandlerFlow) Do() (*AccountResponse, error) {
	//对参数进行合法性验证
	if err := q.checkNum(); err != nil {
		return nil, err
	}

	//更新数据到数据库
	if err := q.updateData(); err != nil {
		return nil, err
	}

	//打包response
	if err := q.packResponse(); err != nil {
		return nil, err
	}
	return q.data, nil
}

func (q *AccountHandlerFlow) checkNum() error {
	if q.username == "" {
		return errors.New("用户名为空")
	}
	if len(q.username) > MaxUsernameLength {
		return errors.New("用户名长度超出限制")
	}
	if q.password == "" {
		return errors.New("密码为空")
	}
	return nil
}

func (q *AccountHandlerFlow) updateData() error {

	//准备好userInfo,默认name为username
	userLogin := models.UserLogin{Username: q.username, Password: q.password}
	userinfo := models.UserInfo{User: &userLogin, Name: q.username}

	//判断用户名是否已经存在
	userLoginDAO := models.NewUserLoginDao()
	if userLoginDAO.IsUserExistByUsername(q.username) {
		return errors.New("用户名已存在")
	}

	//更新操作，由于userLogin属于userInfo，故更新userInfo即可，且由于传入的是指针，所以插入的数据内容也是清楚的
	userInfoDAO := models.NewUserInfoDAO()
	err := userInfoDAO.AddUserInfo(&userinfo)
	if err != nil {
		return err
	}

	//颁发token
	// token, err := middleware.ReleaseToken(userLogin)
	token := q.username + q.password
	if err != nil {
		return err
	}
	q.token = token
	q.userid = userinfo.Id
	return nil
}

func (q *AccountHandlerFlow) packResponse() error {
	q.data = &AccountResponse{
		UserId: q.userid,
		Token:  q.token,
	}
	return nil
}

const (
	MaxUsernameLength = 32
	MaxPasswordLength = 32
	MinPasswordLength = 6
)

type AccountResponse struct {
	UserId int64  `json:"user_id"`
	Token  string `json:"token"`
}

// QueryUserLogin 查询用户是否存在，并返回token和id
func QueryUserLogin(username, password string) (*AccountResponse, error) {
	return NewQueryUserLoginFlow(username, password).Do()
}

func NewQueryUserLoginFlow(username, password string) *QueryUserLoginFlow {
	return &QueryUserLoginFlow{username: username, password: password}
}

type QueryUserLoginFlow struct {
	username string
	password string

	data   *AccountResponse
	userid int64
	token  string
}

func (q *QueryUserLoginFlow) Do() (*AccountResponse, error) {
	//对参数进行合法性验证
	if err := q.checkNum(); err != nil {
		return nil, err
	}
	//准备好数据
	if err := q.prepareData(); err != nil {
		return nil, err
	}
	//打包最终数据
	if err := q.packData(); err != nil {
		return nil, err
	}
	return q.data, nil
}

func (q *QueryUserLoginFlow) checkNum() error {
	if q.username == "" {
		return errors.New("用户名为空")
	}
	if len(q.username) > MaxUsernameLength {
		return errors.New("用户名长度超出限制")
	}
	if q.password == "" {
		return errors.New("密码为空")
	}
	if len(q.password) > MaxPasswordLength {
		return errors.New("密码长度超出限制")
	}
	if len(q.password) < MinPasswordLength {
		return errors.New("密码需要大于5位")
	}
	return nil
}

func (q *QueryUserLoginFlow) prepareData() error {
	userLoginDAO := models.NewUserLoginDao()
	var login models.UserLogin
	//准备好userid
	err := userLoginDAO.QueryUserLogin(q.username, q.password, &login)
	if err != nil {
		return err
	}
	q.userid = login.UserInfoId

	//准备颁发token
	// token, err := middleware.ReleaseToken(login)
	token := q.username + q.password
	if err != nil {
		return err
	}
	q.token = token
	return nil
}

func (q *QueryUserLoginFlow) packData() error {
	q.data = &AccountResponse{
		UserId: q.userid,
		Token:  q.token,
	}
	return nil
}
