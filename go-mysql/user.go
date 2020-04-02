package model

type User struct {
	ID            int
	UserId        int
	Name          string
	Phone         int
	City          int
	Address       string
	Status        int
	Mtime         string `gorm:"default:'galeone'"`
	Ctime         string `gorm:"default:'galeone'"`
}

func (User) TableName() string {
	return "user"
}

func AddUser(sql string, values ...interface{}) error {
	if err := DB.Exec(sql, values).Error; err != nil {
		return err
	}
	return nil
}

func SetUserPhone(userId int, attr map[string]interface{}) error {
	if err := DB.Model(User{}).Where("user_id = ?", userId).Update(attr).Error; err != nil {
		return err
	}
	return nil
}