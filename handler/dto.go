package handler

import "test-git/model"

// CreateBookRequest 创建书籍的请求体
type CreateBookRequest struct {
	Title       string  `json:"title" binding:"required"`  // 书名（必填）
	Author      string  `json:"author" binding:"required"` // 作者（必填）
	Price       float64 `json:"price"`                     // 价格
	Description string  `json:"description"`               // 描述
}

// UpdateBookRequest 更新书籍的请求体
type UpdateBookRequest struct {
	Title       string  `json:"title"`       // 书名（可选，不填则不更新）
	Author      string  `json:"author"`      // 作者（可选）
	Price       float64 `json:"price"`       // 价格（可选）
	Description string  `json:"description"` // 描述（可选）
}

// BookResponse 书籍的响应体（返回给前端的数据）
type BookResponse struct {
	ID          uint    `json:"id"`          // 主键ID
	Title       string  `json:"title"`       // 书名
	Author      string  `json:"author"`      // 作者
	Price       float64 `json:"price"`       // 价格
	Description string  `json:"description"` // 描述
	CreatedAt   string  `json:"created_at"`  // 创建时间
	UpdatedAt   string  `json:"updated_at"`  // 更新时间
}

type BookListResponse struct {
	Total int            `json:"total"` // 总条数
	List  []BookResponse `json:"list"`  // 分页数据列表
}

type RolePreviewResponse struct {
	Name        string `json:"name"`        // 角色名称
	Description string `json:"description"` // 角色描述
	AvatarURL   string `json:"avatar_url"`  // 头像URL
}

func toBookResponse(book model.Book) BookResponse {
	return BookResponse{
		ID:          book.ID,
		Title:       book.Title,
		Author:      book.Author,
		Price:       book.Price,
		Description: book.Description,
		CreatedAt:   book.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   book.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

// COC角色卡主结构体
type COCRoleCard struct {
	BasicInfo      BasicInfo      `json:"basic_info"`      // 基本信息
	Attributes     Attributes     `json:"attributes"`      // 属性值
	Skills         Skills         `json:"skills"`          // 技能
	Inventory      Inventory      `json:"inventory"`       // 物品与财富
	PersonalTraits PersonalTraits `json:"personal_traits"` // 个人特征
	Status         Status         `json:"status"`          // 当前状态
}

// 基本信息
type BasicInfo struct {
	AvatarURL  string `json:"avatar_url,omitempty"` // 头像URL
	RoleName   string `json:"name"`                 // 角色名
	Gender     string `json:"gender"`               // 性别
	Age        int    `json:"age"`                  // 年龄
	Occupation string `json:"occupation"`           // 职业
	Alignment  string `json:"alignment"`            // 阵营
	Race       string `json:"race"`                 // 种族
	Appearance string `json:"appearance"`           // 外貌描述
	Backstory  string `json:"backstory"`            // 背景故事
}

// 属性值（含派生属性）
type Attributes struct {
	Strength     int               `json:"(STR)"`   // 力量
	Constitution int               `json:"(CON)"`   // 体质
	Size         int               `json:"(SIZ)"`   // 体型
	Dexterity    int               `json:"(DEX)"`   // 敏捷
	Appearance   int               `json:"(APP)"`   // 外貌
	Intelligence int               `json:"(INT)"`   // 智力
	Willpower    int               `json:"(POW)"`   // 意志
	Education    int               `json:"(EDU)"`   // 教育
	Luck         int               `json:"(LUK)"`   // 幸运
	Derived      DerivedAttributes `json:"derived"` // 派生属性
}

// 派生属性
type DerivedAttributes struct {
	HP        int `json:"(HP)"`          // 生命值
	SAN       int `json:"(SAN)"`         // 理智值
	MP        int `json:"(MP)"`          // 魔法值
	MOV       int `json:"(MOV)"`         // 移动力
	Actions   int `json:"actions"`       // 行动数
	LoadLimit int `json:"loadlimit(kg)"` // 负重上限(kg)
}

// 技能（职业技能/通用技能/魔法技能）
type Skills struct {
	Occupational []Skill `json:"occupational"` // 职业技能
	General      []Skill `json:"general"`      // 通用技能
	Magic        []Skill `json:"magic"`        // 魔法技能
}

// 单个技能
type Skill struct {
	Name   string `json:"name"`             // 技能名称
	Value  int    `json:"value"`            // 技能数值
	Remark string `json:"remark,omitempty"` // 备注（可选）
}

// 物品与财富
type Inventory struct {
	Equipments []Equipment `json:"equipments"` // 装备列表
	Wealth     Wealth      `json:"wealth"`     // 财富信息
}

// 单个装备
type Equipment struct {
	Name     string `json:"name"`             // 装备名称
	Quantity int    `json:"quantity"`         // 数量
	Ammo     int    `json:"ammo,omitempty"`   // 弹药（仅武器有，可选）
	Remark   string `json:"remark,omitempty"` // 备注（可选）
}

// 财富信息
type Wealth struct {
	Cash        int    `json:"cash"`        // 现金
	Assets      string `json:"assets"`      // 资产
	CreditScore int    `json:"creditScore"` // 信用评级
}

// 个人特征
type PersonalTraits struct {
	Personality     string `json:"personality"`     // 个性特点
	ImportantPerson string `json:"importantPerson"` // 重要之人
	ImportantItem   string `json:"importantItem"`   // 重要物品
	SpecialAbility  string `json:"specialAbility"`  // 特殊能力
}

// 当前状态
type Status struct {
	CurrentSAN int    `json:"currentSAN"` // 当前理智值
	CurrentHP  int    `json:"currentHP"`  // 当前生命值
	IsInjured  bool   `json:"isInjured"`  // 是否受伤
	IsInsane   bool   `json:"isInsane"`   // 是否疯狂
	Remark     string `json:"remark"`     // 备注
}
