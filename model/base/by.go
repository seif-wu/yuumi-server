package base

type ControlBy struct {
	CreatedBy int `json:"createdBy" gorm:"index;comment:创建者"`
	UpdatedBy int `json:"updatedBy" gorm:"index;comment:更新者"`
}

func (c *ControlBy) SetCreatedBy(createdBy int) {
	c.CreatedBy = createdBy
}

func (c *ControlBy) SetUpdatedBy(updatedBy int) {
	c.UpdatedBy = updatedBy
}
