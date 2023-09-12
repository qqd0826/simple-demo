package util

import (
	"github.com/RaymondCode/simple-demo/dao"
	"github.com/RaymondCode/simple-demo/model"
	"github.com/dgrijalva/jwt-go"
	"log"
	"time"
)

type Claims struct {
	UserId             int64
	UserName           string
	jwt.StandardClaims // jwt中标准格式,主要是设置token的过期时间
}

// 密钥，不可泄露，后期改,并存其他地方
const secretkey = "1234567890abcdefghijk"

// 创建token
func CreateToken(id int64, username string) (string, error) {
	claims := &Claims{
		UserId:   id,
		UserName: username,
		StandardClaims: jwt.StandardClaims{
			//Audience:  "app项目",                              //颁发给谁，就是使用的一方
			ExpiresAt: time.Now().Add(time.Hour * 5).Unix(), //过期时间,暂设为1小时
			//Id:        "",//非必填
			IssuedAt: time.Now().Unix(), //颁发时间
			Issuer:   "dousheng",        //颁发者
			//NotBefore: 0,         //生效时间，就是在这个时间前不能使用，非必填
			//Subject: "111",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretkey))
	if err != nil {
		return "create token error:", err
	}
	//defer func() {
	//    e := recover()
	//    if e != nil {
	//        panic(e)
	//    }
	//}()

	return tokenString, nil
}

// 解析token
func ParseToken(tokenString string) (*Claims, error) {
	/*
	   jwt.ParseWithClaims有三个参数
	   第一个就是加密后的token字符串
	   第二个是Claims, &Claims 传指针的目的是。然让它在我这个里面写数据
	   第三个是一个自带的回调函数，将秘钥和错误return出来即可，要通过密钥解密
	*/
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		//if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		//    return nil, fmt.Errorf("算法类型: %v", token.Header["alg"])
		//}
		return []byte(secretkey), nil
	})
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		//a := token.Valid
		//b := claims.Valid()
		//fmt.Println(a, b)
		return claims, nil
	} else {
		return nil, err
	}
}
func GetUserByToken(tokenString string) model.User {
	if len(tokenString) == 0 {
		return model.User{}
	} else {
		claims, err := ParseToken(tokenString)
		if err != nil {
			log.Println(err)
			return model.User{}
		}
		return dao.GetUserById(claims.UserId)
	}
}
