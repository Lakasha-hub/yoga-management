package class

import (
	"gorm.io/gorm"
)

type ClassRepository interface {
	GetAll() ([]Class, error)
	GetOne(id int) (*Class, error)
	Insert(in *Class) (*Class, error)
	Update(id int, up *Class) (*Class, error)
	Delete(id int) error
}

type ClassMysqlRepository struct {
	db gorm.DB
}

func NewClassMysqlRepository(db gorm.DB) ClassRepository {
	return &ClassMysqlRepository{db: db}
}

func (r *ClassMysqlRepository) GetAll() ([]Class, error) {
	var classes []Class
	result := r.db.Find(&classes)
	if result.Error != nil {
		return nil, result.Error
	}
	return classes, nil
}

func (r *ClassMysqlRepository) GetOne(id int) (*Class, error) {
	class := Class{}
	result := r.db.First(&class, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, result.Error
	}
	return &class, nil
}

func (r *ClassMysqlRepository) Insert(in *Class) (*Class, error) {

	var class_exists Class
	if err := r.db.Where("name_class = ? AND professor = ? AND date_class = ?", in.NameClass, in.Professor, in.DateClass).First(&class_exists).Error; err == nil {
		return nil, gorm.ErrDuplicatedKey
	}

	if err := r.db.Save(&in).Error; err != nil {
		return nil, err
	}
	return in, nil
}

func (r *ClassMysqlRepository) Update(id int, up *Class) (*Class, error) {
	var class Class
	if err := r.db.First(&class, id).Error; err != nil {
		return nil, gorm.ErrRecordNotFound
	}

	if err := r.db.Model(&class).Updates(up).Error; err != nil {
		return nil, err
	}
	return &class, nil
}

func (r *ClassMysqlRepository) Delete(id int) error {
	var class Class
	if err := r.db.First(&class, id).Error; err != nil {
		return gorm.ErrRecordNotFound
	}

	if err := r.db.Delete(&class).Error; err != nil {
		return err
	}
	return nil
}
