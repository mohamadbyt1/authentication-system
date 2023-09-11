package api
import (
	"context"
	"encoding/json"
	"errors"
	"github.com/mohamadbyt1/authentication-system/models"
	"net/http"
	"strconv"
	"time"
	"github.com/golang-jwt/jwt/v5"
)
func (s *ApiServer) Signup(w http.ResponseWriter, r *http.Request)error{
	if r.Method != "POST" {
		return ResponseJson(w,http.StatusBadRequest,map[string]string{"msg":"bad request"})
	}
	ctx , cancel:= context.WithTimeout(context.Background(),time.Millisecond *100)
	defer cancel()
	User := new(models.UserSignup)
	if err := json.NewDecoder(r.Body).Decode(User) ; err != nil {
		return err
	}
	if err := User.ValidateSignup(); err != nil {
		return err
	}
	ex,_ := s.Storage.UserCheck(ctx,User.Username)
	if ex {
		return ResponseJson(w,http.StatusBadRequest,map[string]string{"msg":"user exists!"})
	}
	hashePass , err := hashPassword(User.Password)
	if err != nil {
		return err
	}
	NewUser := models.NewUser(User.Username,User.FirstName,User.LastName,string(hashePass))
	if err := s.Storage.CreaateUser(ctx,NewUser); err != nil {
		return err
	}
	return ResponseJson(w,http.StatusOK,map[string]string{"msg":"some good msg"})
}
func (s *ApiServer) Login(w http.ResponseWriter, r *http.Request)error{
	if r.Method != "POST" {
		return ResponseJson(w,http.StatusBadRequest,map[string]string{"msg":"bad request"})
	}
	ctx , cancel:= context.WithTimeout(context.Background(),time.Millisecond *100)
	defer cancel()
var requestData map[string]interface{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&requestData); err != nil {
		return err
	}
	if len(requestData) != 2 || requestData["username"] == nil || requestData["password"] == nil {
		return errors.New("not valid")
	}
	userBody := new(models.UserLogin)
	userBody.Password = requestData["password"].(string)
	userBody.Username = requestData["username"].(string)
	user ,err := s.Storage.GetUser(ctx ,userBody.Username)
	if err != nil{
		return err
	}
	if user == nil {
        return ResponseJson(w, http.StatusUnauthorized, map[string]string{"error": "User not found"})
    }
	if err := compare(userBody.Password,user.Password) ; err != nil{
		return err
	}
	jwtToken,err :=generateJWT(user.Id,user.Username)
	if err != nil{
		return err
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "authToken",
		Value:    jwtToken,
		MaxAge:   int(time.Hour.Seconds() * 24),
		Path:     "/",
		Secure:   false,
		HttpOnly: true,
	})
	UserResponse , err:= s.Storage.GetUser(ctx ,user.Username)
	if err != nil{
		return err
	}
	return ResponseJson(w,http.StatusOK,UserResponse)
}
func (s *ApiServer) DeleteUser(w http.ResponseWriter , r *http.Request)error{
	if r.Method != http.MethodDelete {
		return 	errors.New("bad request")
	}
	ctx , cancel:= context.WithTimeout(context.Background(),time.Millisecond *100)
	defer cancel()

	claims, ok := r.Context().Value("claims").(jwt.MapClaims)
	if !ok {
		return errors.New("unauthorized")
	}
	id := r.URL.Query().Get("id")
	if id == "" {
		return errors.New("missing id parameter")
	}
	idClaim := claims["id"].(float64)
	idClaimStr := strconv.FormatFloat(idClaim, 'f', -1, 64)
	if id != idClaimStr {
		return errors.New("not a valid id")
	}
	exist , _:= s.Storage.IdCheck(ctx ,id)
	if !exist {
		return errors.New("id does not exist")
	}
	err := s.Storage.DeleteUser(ctx ,id)
	if err != nil {
		return err
	}
	return ResponseJson(w,http.StatusOK,map[string]string{"msg":"user deleated"})
}
func (s *ApiServer) UpdateUser(w http.ResponseWriter, r *http.Request)error{
	if r.Method != http.MethodPatch{
		return ResponseJson(w,http.StatusBadRequest,map[string]string{"msg":"bad request"})
	}
	ctx , cancel:= context.WithTimeout(context.Background(),time.Millisecond *100)
	defer cancel()
	userBody := new(models.UserSignup)
	err := json.NewDecoder(r.Body).Decode(userBody)
	if err != nil {
		return err
	}
	if err :=  userBody.ValidateSignup(); err != nil{
		return err
	}
	claims, ok := r.Context().Value("claims").(jwt.MapClaims)
	if !ok {
		return errors.New("unauthorized")
	}
	if claims["username"] != userBody.Username {
		return errors.New("unauthorized")
	}
	exist , _ := s.Storage.UserCheck(ctx ,userBody.Username)
	if !exist{
		return errors.New("failed to update user information")
	}
	hashPass , err := hashPassword(userBody.Password)
	if err != nil {
		return err
	}
	userBody.Password = string(hashPass)
	err =  s.Storage.UpdateUser(ctx ,userBody)
	if err != nil {
		return err
	}
	return	ResponseJson(w,http.StatusOK,map[string]string{"msg":"user information updated"})
}