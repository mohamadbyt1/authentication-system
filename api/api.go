package api

import (
	"encoding/json"
	"fmt"
	"log"
	"github.com/mohamadbyt1/authentication-system/storage"
	"net/http"
	"time"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)
type ApiServer struct {
	ListenAddr string
	Storage storage.Storage
}
func NewApiServer(addr string , store storage.Storage) *ApiServer {

	return &ApiServer{
		ListenAddr: addr,
		Storage: store,
	}
}
func (s ApiServer)Start()error{
	r := http.DefaultServeMux
	r.HandleFunc("/delete",AuthMiddleware(handleWithErr(s.DeleteUser)))
	r.HandleFunc("/update",AuthMiddleware(handleWithErr(s.UpdateUser)))
	r.HandleFunc("/signup",handleWithErr(s.Signup))
	r.HandleFunc("/login",handleWithErr(s.Login))
	log.Println("running on :",s.ListenAddr)
	err:= http.ListenAndServe(s.ListenAddr,nil)
	if err != nil {
	return err	
	}
	return nil
}
func ResponseJson(w http.ResponseWriter, status int, data interface{})error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}
func generateJWT(id int , username string)(string , error){
	//generate JWT token and returning it
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = id
	claims["username"] = username
	claims["exp"] = time.Now().Add(60 * time.Second *24 *30).Unix()
	claims["authorized"] = true
	var sampleSecretKey = []byte("secretkey")
	tokenString, err := token.SignedString(sampleSecretKey)
	if err != nil{
		return "",err
	}
	return tokenString, nil
}
func hashPassword(pass string)([]byte,error){
	//hash password and return byte
	hashPassword , err := bcrypt.GenerateFromPassword([]byte(pass),bcrypt.DefaultCost)
	 if err != nil {
		fmt.Println(err)
		return nil,err
	}
	return hashPassword,nil
}
func compare(pass string,hash string)error{

	if err := bcrypt.CompareHashAndPassword([]byte(hash),[]byte(pass)); err != nil {
		fmt.Println("err is:",err)

		return err
	}
	return nil
}
type ApiError struct {
	Error string `json:"error"`
}

func  handleWithErr(Handler func(http.ResponseWriter , *http.Request)error) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
		if err :=  Handler(w, r); err != nil {

			ResponseJson(w,http.StatusBadRequest,ApiError{Error: err.Error()})
		}
	}
}  