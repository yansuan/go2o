/**
 * Copyright 2015 @ z3q.net.
 * name : category_manager.go
 * author : jarryliu
 * date : 2016-06-04 13:40
 * description :
 * history :
 */
package product

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/jsix/gof/algorithm/iterator"
	"go2o/core/domain/interface/product"
	"go2o/core/domain/interface/valueobject"
	"go2o/core/infrastructure/domain"
	"sort"
	"strconv"
	"strings"
	"time"
)

var _ product.ICategory = new(categoryImpl)

// 分类实现
type categoryImpl struct {
	value           *product.Category
	rep             product.ICategoryRepo
	parentIdChanged bool
	childIdArr      []int32
	opt             domain.IOptionStore
}

func newCategory(rep product.ICategoryRepo,
	v *product.Category) product.ICategory {
	return &categoryImpl{
		value: v,
		rep:   rep,
	}
}

func (c *categoryImpl) GetDomainId() int32 {
	return c.value.Id
}

func (c *categoryImpl) GetValue() *product.Category {
	return c.value
}

func (c *categoryImpl) GetOption() domain.IOptionStore {
	if c.opt == nil {
		opt := newCategoryOption(c)
		if err := opt.Stat(); err != nil {
			opt.Set(product.C_OptionViewName, &domain.Option{
				Key:   product.C_OptionViewName,
				Type:  domain.OptionTypeString,
				Must:  false,
				Title: "显示页面",
				Value: "goods_list.html",
			})
			opt.Set(product.C_OptionDescribe, &domain.Option{
				Key:   product.C_OptionDescribe,
				Type:  domain.OptionTypeString,
				Must:  false,
				Title: "描述",
				Value: "",
			})
			opt.Flush()
		}
		c.opt = opt
	}
	return c.opt
}

// 检查上级分类是否正确
func (c *categoryImpl) checkParent(parentId int32) error {
	if id := c.GetDomainId(); id > 0 && parentId > 0 {
		//检查上级栏目是否存在
		p := c.rep.GlobCatService().GetCategory(parentId)
		if p == nil {
			return product.ErrNoSuchCategory
		}
		// 检查上级分类
		if p.GetValue().ParentId == id {
			return product.ErrCategoryCycleReference
		}
	}
	return nil
}

// 设置值
func (c *categoryImpl) SetValue(v *product.Category) error {
	val := c.value
	if val.Id == v.Id {
		val.Enabled = v.Enabled
		val.Name = v.Name
		val.SortNum = v.SortNum
		val.Icon = v.Icon
		val.ProModel = v.ProModel
		if val.ParentId != v.ParentId {
			c.parentIdChanged = true
		} else {
			c.parentIdChanged = false
		}

		if c.parentIdChanged {
			err := c.checkParent(v.ParentId)
			if err != nil {
				return err
			}
			val.ParentId = v.ParentId
		}
	}
	return nil
}

// 获取子栏目的编号
func (c *categoryImpl) GetChildes() []int32 {
	if c.childIdArr == nil {
		childCats := c.getChildCategories(c.GetDomainId())
		c.childIdArr = make([]int32, len(childCats))
		for i, v := range childCats {
			c.childIdArr[i] = v.Id
		}
	}
	return c.childIdArr
}
func (c *categoryImpl) setCategoryLevel() {
	var mchId int32 = 0
	list := c.rep.GetCategories(mchId)
	c.parentWalk(list, mchId, &c.value.Level)
	//todo: 未实现
}

func (c *categoryImpl) parentWalk(list []*product.Category,
	parentId int32, level *int) {
	*level += 1
	if parentId <= 0 {
		return
	}
	for _, v := range list {
		if v.Id == v.ParentId {
			panic(errors.New(fmt.Sprintf(
				"Bad category , id is same of parent id , id:%s",
				v.Id)))
		} else if v.Id == parentId {
			c.parentWalk(list, v.ParentId, level)
			break
		}
	}
}

func (c *categoryImpl) Save() (int32, error) {
	//if c._manager.ReadOnly() {
	//    return c.GetDomainId(), product.ErrReadonlyCategory
	//}
	c.setCategoryLevel()
	id, err := c.rep.SaveCategory(c.value)
	if err == nil {
		c.value.Id = id
		if len(c.value.Url) == 0 || (c.parentIdChanged &&
			strings.HasPrefix(c.value.Url, "/c-")) {
			c.value.Url = c.getAutomaticUrl(id)
			c.parentIdChanged = false
			return c.Save()
		}
	}
	return id, err
}

// 获取子栏目
func (c *categoryImpl) getChildCategories(catId int32) []*product.Category {
	var all []*product.Category = c.rep.GetCategories(0)
	var newArr []*product.Category = []*product.Category{}

	var cdt iterator.Condition = func(v, v1 interface{}) bool {
		return v1.(*product.Category).ParentId == v.(*product.Category).Id
	}
	var start iterator.WalkFunc = func(v interface{}, level int) {
		c := v.(*product.Category)
		if c.Id != catId {
			newArr = append(newArr, c)
		}
	}

	var arr []interface{} = make([]interface{}, len(all))
	for i := range arr {
		arr[i] = all[i]
	}

	iterator.Walk(arr, &product.Category{Id: catId}, cdt, start, nil, 1)

	return newArr
}

// 获取与栏目相关的栏目
func (c *categoryImpl) getRelationCategories(catId int32) []*product.Category {
	var all []*product.Category = c.rep.GetCategories(0)
	var newArr []*product.Category = []*product.Category{}
	var isMatch bool
	var pid int32
	var l int = len(all)

	for i := 0; i < l; i++ {
		if !isMatch && all[i].Id == catId {
			isMatch = true
			pid = all[i].ParentId
			newArr = append(newArr, all[i])
			i = -1
		} else {
			if all[i].Id == pid {
				newArr = append(newArr, all[i])
				pid = all[i].ParentId
				i = -1
				if pid == 0 {
					break
				}
			}
		}
	}
	return newArr
}

func (c *categoryImpl) getAutomaticUrl(id int32) string {
	relCats := c.getRelationCategories(id)
	var buf *bytes.Buffer = bytes.NewBufferString("/c")
	var l int = len(relCats)
	for i := l; i > 0; i-- {
		buf.WriteString("-" + strconv.Itoa(int(relCats[i-1].Id)))
	}
	buf.WriteString(".htm")
	return buf.String()
}

var _ domain.IOptionStore = new(categoryOption)

// 分类数据选项
type categoryOption struct {
	domain.IOptionStore
	_c *categoryImpl
}

func newCategoryOption(c *categoryImpl) domain.IOptionStore {
	file := fmt.Sprintf("conf/core/sale/cate_opt_%d", c.GetDomainId())
	return &categoryOption{
		_c:           c,
		IOptionStore: domain.NewOptionStoreWrapper(file),
	}
}

var _ product.IGlobCatService = new(categoryManagerImpl)

//当商户共享系统的分类时,没有修改的权限,既只读!
type categoryManagerImpl struct {
	_readonly      bool
	_rep           product.ICategoryRepo
	_valRepo       valueobject.IValueRepo
	_mchId         int32
	lastUpdateTime int64
	_categories    []product.ICategory
}

func NewCategoryManager(mchId int32, rep product.ICategoryRepo,
	valRepo valueobject.IValueRepo) product.IGlobCatService {
	c := &categoryManagerImpl{
		_rep:     rep,
		_mchId:   mchId,
		_valRepo: valRepo,
	}
	return c.init()
}

func (c *categoryManagerImpl) init() product.IGlobCatService {
	mchConf := c._valRepo.GetPlatformConf()
	if !mchConf.MchGoodsCategory && c._mchId > 0 {
		c._readonly = true
		c._mchId = 0
	}
	return c
}

// 获取栏目关联的编号,系统用0表示
func (c *categoryManagerImpl) getRelationId() int32 {
	return c._mchId
}

// 清理缓存
func (c *categoryManagerImpl) clean() {
	c._categories = nil
}

// 是否只读,当商户共享系统的分类时,
// 没有修改的权限,即只读!
func (c *categoryManagerImpl) ReadOnly() bool {
	return c._readonly
}

// 创建分类
func (c *categoryManagerImpl) CreateCategory(v *product.Category) product.ICategory {
	if v.CreateTime == 0 {
		v.CreateTime = time.Now().Unix()
	}
	return newCategory(c._rep, v)
}

// 获取分类
func (c *categoryManagerImpl) GetCategory(id int32) product.ICategory {
	v := c._rep.GetCategory(c.getRelationId(), id)
	if v != nil {
		return c.CreateCategory(v)
	}
	return nil
}

// 获取所有分类
func (c *categoryManagerImpl) GetCategories() []product.ICategory {
	var list product.CategoryList = c._rep.GetCategories(c.getRelationId())
	sort.Sort(list)
	slice := make([]product.ICategory, len(list))
	for i, v := range list {
		slice[i] = c.CreateCategory(v)
	}
	return slice
}

// 删除分类
func (c *categoryManagerImpl) DeleteCategory(id int32) error {
	cat := c.GetCategory(id)
	if cat == nil {
		return product.ErrNoSuchCategory
	}
	if len(cat.GetChildes()) > 0 {
		return product.ErrHasChildCategories
	}
	if c._rep.CheckGoodsContain(c.getRelationId(), id) {
		return product.ErrCategoryContainGoods
	}
	err := c._rep.DeleteCategory(c.getRelationId(), id)
	if err == nil {
		err = cat.GetOption().Destroy()
		cat = nil
	}
	return err
}
