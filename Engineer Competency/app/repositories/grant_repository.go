package repositories

import (
	"engineer-comp/app/global/pgsql"
	"engineer-comp/app/global/utils/helper"
	"engineer-comp/app/models"

	"fmt"
)

type GrantRepositoryInterface interface {
	GetPasswordByEmail(request string, resultChan chan *models.LoginResponseChan)
	InsertLoginSession(request *models.LoginToken, resultChan chan *models.LoginResponseChan)
	GetLoginSession(request string, resultChan chan *models.LoginResponseChan)
	GetRoleById(request string) ([]int64, int64, error)
}

type grantRepository struct {
	sqldb pgsql.SqlInterface
}

func InitGrantRepository(sqldb pgsql.SqlInterface) GrantRepositoryInterface {
	return &grantRepository{
		sqldb: sqldb,
	}
}

func (r *grantRepository) GetPasswordByEmail(request string, resultChan chan *models.LoginResponseChan) {
	response := &models.LoginResponseChan{}
	var total int

	sql := `select  email, user_name, password, user_status from users where email = '` + request + `' and user_status = 1`
	sqlCount := `select count(id) from users where user_status = 1 and email = '` + request + `'`

	err := r.sqldb.DB().QueryRow(sqlCount).Scan(&total)
	if err != nil {
		fmt.Println("error login email checking user email")
		fmt.Println(err.Error())
		response.Error = err
		resultChan <- response
		return
	}

	if total == 0 {
		fmt.Println("login, email not exist or user is deactive")
		err = helper.NewError("login email, user is not available")
		response.Error = err
		resultChan <- response
		return
	}

	usersQueue, err := r.sqldb.DB().Query(sql)

	if err != nil {
		fmt.Println(err)
		response.Error = err
		resultChan <- response
		return
	}

	var loginResponse models.LoginResponse
	for usersQueue.Next() {
		err = usersQueue.Scan(&loginResponse.Email, &loginResponse.UserName, &loginResponse.Password, &loginResponse.UserStatus)
		if err != nil {
			fmt.Println("login email, error scanning users")
			fmt.Println(err.Error())
			response.Error = err
			resultChan <- response
			return
		}
	}

	response.Login = &loginResponse
	response.Error = nil
	resultChan <- response
	return

}

func (r *grantRepository) InsertLoginSession(request *models.LoginToken, resultChan chan *models.LoginResponseChan) {
	response := &models.LoginResponseChan{}
	var total int
	sqlCount := `select count(id) from login_session where token = '` + request.Token + `'`
	err := r.sqldb.DB().QueryRow(sqlCount).Scan(&total)
	if err != nil {
		fmt.Println("error login email checking user email")
		fmt.Println(err.Error())
		response.Error = err
		resultChan <- response
		return
	}

	if total != 0 {
		fmt.Println("token already exist")
		err = helper.NewError("token already exist")
		response.Error = err
		resultChan <- response
		return
	}

	_, err = r.sqldb.DB().Query(`INSERT INTO public.login_session
													( "token", expire_time, expire_in, created_at) 
											VALUES ( $1, $2, $3, NOW())`, request.Token, request.ExpireTime, request.ExpireIn)

	if err != nil {
		fmt.Println(err)
		response.Error = err
		resultChan <- response
		return
	}

	response.Error = nil
	resultChan <- response
	return

}

func (r *grantRepository) GetLoginSession(request string, resultChan chan *models.LoginResponseChan) {
	response := &models.LoginResponseChan{}
	var total int
	sqlCount := `select count(id) from login_session where token = '` + request + `'`
	err := r.sqldb.DB().QueryRow(sqlCount).Scan(&total)
	if err != nil {
		fmt.Println("error login email checking user email")
		fmt.Println(err.Error())
		response.Error = err
		resultChan <- response
		return
	}

	if total == 0 {
		fmt.Println("token not found")
		err = helper.NewError("token not found")
		response.Error = err
		resultChan <- response
		return
	}

	loginSession, err := r.sqldb.DB().Query(`select token, expire_time, expire_in, created_at from login_session where token = $1`, request)
	if err != nil {
		fmt.Println(err)
		response.Error = err
		resultChan <- response
		return
	}

	var loginResponse models.LoginToken
	for loginSession.Next() {
		err = loginSession.Scan(&loginResponse.Token, &loginResponse.ExpireTime, &loginResponse.ExpireIn, &loginResponse.CreatedAt)
		if err != nil {
			fmt.Println("login email, error scanning users")
			fmt.Println(err.Error())
			response.Error = err
			resultChan <- response
			return
		}
	}

	requestLoginToken := &models.LoginToken{
		Token:      loginResponse.Token,
		ExpireTime: loginResponse.ExpireTime,
		ExpireIn:   loginResponse.ExpireIn,
		CreatedAt:  loginResponse.CreatedAt,
	}

	response.LoginToken = requestLoginToken
	response.Error = nil
	resultChan <- response
	return

}

func (r *grantRepository) GetRoleById(request string) ([]int64, int64, error) {
	var roleIDs []int64
	var total int
	var userId int64

	sql := `select ur.role_id, u.id from users u join user_roles ur on u.id= ur.user_id where u.email = '` + request + `' and u.user_status = 1`
	sqlCount := `select count(*) from users u join user_roles ur on u.id= ur.user_id where u.email = '` + request + `' and u.user_status = 1`

	err := r.sqldb.DB().QueryRow(sqlCount).Scan(&total)
	if err != nil {
		fmt.Println("error login email checking user email")
		fmt.Println(err.Error())
		return []int64{}, 0, err
	}

	if total == 0 {
		fmt.Println("login, email not exist or user is deactive")
		err = helper.NewError("login email, user is not available")
		return []int64{}, 0, err
	}

	usersQueue, err := r.sqldb.DB().Query(sql)

	if err != nil {
		fmt.Println(err)
		return []int64{}, 0, err
	}

	for usersQueue.Next() {
		var roleID int64
		err = usersQueue.Scan(&roleID, &userId)
		if err != nil {
			fmt.Println("error scanning users, email not found")
			fmt.Println(err.Error())
			return []int64{}, 0, err
		}
		roleIDs = append(roleIDs, roleID)
	}

	return roleIDs, userId, nil

}
