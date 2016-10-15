package service

type repository interface {
	addToken(token Token) (error)
	addUser(user User) (error)
	getUserByUsername(username string) (User, error)
    getTokenByKey(key string) (Token, error)
}

type repoHandler struct {}

func (r *repoHandler) addToken(token Token) (error) {
    return DB.Create(&token).Error
}

func (r *repoHandler) addUser(user User) (error) {
    return DB.Create(&user).Error
}

func (r *repoHandler) getUserByUsername(username string) (User, error) {
    var user User
    err := DB.Where("username=?", username).First(&user).Error
    return user, err
}

func (r *repoHandler) getTokenByKey(key string) (Token, error) {
    var token Token
    err := DB.Where("key=?", key).First(&token).Error
    return token, err
}
