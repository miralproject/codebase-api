package domain

type User struct {
	BaseDomain
	FirstName    string `gorm:"type:varchar(150);column:first_name;not null" json:"first_name"`
	LastName     string `gorm:"type:varchar(150);column:last_name;not null" json:"last_name"`
	Username     string `gorm:"type:varchar(150);column:username;not null" json:"username"`
	Email        string `gorm:"type:varchar(100);unique;not null" json:"email"`
	Password     string `gorm:"type:varchar(150);column:password;not null" json:"password"`
	Phone        string `gorm:"type:varchar(100);unique" json:"phone"`
	IsActive     bool   `gorm:"default:true;column:is_active" json:"is_active"`
	RefreshToken string `gorm:"type:text;column:refresh_token" json:"refresh_token"`
}
