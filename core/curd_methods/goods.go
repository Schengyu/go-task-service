package curd_methods

import (
	"context"
	"fmt"
	"github.com/sbigtree/go-db-model/models"
	"go-task-service/cmd/global"
	"time"

	"gorm.io/gorm"
)

// 定义一个通用的事务执行函数，接受一个处理函数作为参数
func ExecuteTransaction(ctx context.Context, f func(tx *gorm.DB) error) error {
	db := global.DB
	tx := db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback()
			panic(r)
		}
	}()
	if err := f(tx); err != nil {
		if rbErr := tx.Rollback().Error; rbErr != nil {
			return fmt.Errorf("failed to rollback transaction: %v, original error: %v", rbErr, err)
		}
		return err
	}

	return tx.Commit().Error
}

// 添加商品的sku
func AddSellerGoodSku(SellerGood models.YYMSellerGoodsSKU, tx *gorm.DB) (models.YYMSellerGoodsSKU, error) {
	err := tx.Create(&SellerGood).Error
	return SellerGood, err
}

// 删除卖家商品的某个sku
func DeleteSellerGoodSku(skuId int) error {
	err := global.DB.Where("id = ?", skuId).Delete(&models.YYMSellerGoodsSKU{}).Error
	return err
}

// 修改卖家商品sku的价格
func UpdateSellerGoodSkuPrice(skuId int, price float64) error {
	err := global.DB.Model(&models.YYMSellerGoodsSKU{}).Where("id = ?", skuId).Update("price", price).Error
	return err
}

// 添加卖家商品表
func AddSellerGood(sellerGood models.YYMSellerGoods, tx *gorm.DB) (models.YYMSellerGoods, error) {
	err := tx.Create(&sellerGood).Error
	return sellerGood, err
}

// 修改卖家商品表当中的上下架状态
func UpdateSellerGoodStatus(sellerGood models.YYMSellerGoods) error {
	err := global.DB.Model(&models.YYMSellerGoods{}).Where("id = ?", sellerGood.ID).Update("is_listing", sellerGood.IsListing).Error
	return err
}

// 修改卖家商品的名称
func UpdateSellerGoodName(sellerGood models.YYMSellerGoods) error {
	err := global.DB.Model(&models.YYMSellerGoods{}).Where("id = ?", sellerGood.ID).Update("goods_name", sellerGood.GoodName).Error
	return err
}

// 查询所有上架的卖家商品

// 添加一条库存记录
func AddInventoryLog(inventoryLog models.YYMCSGOBoxInventory) (models.YYMCSGOBoxInventory, error) {
	err := global.DB.Create(&inventoryLog).Error
	return inventoryLog, err
}

// 根据skuCode去查询masterGoods
func GetMasterGoodsSkuByCode(skuCode string) (models.YYMGoodsSKU, error) {
	var masterGood models.YYMGoodsSKU
	err := global.DB.Where("sku_code = ?", skuCode).Find(&masterGood).Limit(1).Error
	return masterGood, err
}

// 添加一个masterGoodSSku
func AddMasterGoodSSku(masterGoodSSku models.YYMGoodsSKU, tx *gorm.DB) (models.YYMGoodsSKU, error) {
	err := tx.Create(&masterGoodSSku).Error
	return masterGoodSSku, err
}

// 根据属性id 查询出该记录  并推算出他的父级id
func GetAttributeById(attributeId int) (models.YYMAttribute, error) {
	var attribute models.YYMAttribute
	err := global.DB.Where("id = ?", attributeId).Find(&attribute).Error
	return attribute, err
}

// 添加卖家商品活动表
func AddSellerGoodActivity(sellerGoodActivity models.YYMSellerGoodsActivity, tx *gorm.DB) (models.YYMSellerGoodsActivity, error) {
	err := tx.Create(&sellerGoodActivity).Error
	return sellerGoodActivity, err
}

// 添加卖家商品保障表
func AddSellerGoodsGuarantee(sellerGoodsGuarantee models.YYMSellerGoodsGuarantee, tx *gorm.DB) (models.YYMSellerGoodsGuarantee, error) {
	err := tx.Create(&sellerGoodsGuarantee).Error
	return sellerGoodsGuarantee, err
}

// 根据类别查询商品
func GetGoodsByCategory(categoryId int) ([]models.YYMGoodsMaster, error) {
	var goods []models.YYMGoodsMaster
	err := global.DB.Where("category_id = ?", categoryId).Find(&goods).Error
	return goods, err
}

// 根据pid查询出保障的属性
func GetGuaranteeAttrByParentId(parentId int) ([]models.YYMGuaranteeAttribute, error) {
	var attributes []models.YYMGuaranteeAttribute
	err := global.DB.Where("parent_id =?", parentId).Find(&attributes).Error
	return attributes, err
}

// 根据pid查询出属性值
func GetAttrValueByParentId(parentId int) ([]models.YYMAttribute, error) {
	var values []models.YYMAttribute
	err := global.DB.Where("parent_id =?", parentId).Find(&values).Error
	return values, err
}

// 根据id查询出卖家商品
func GetSellerGoodsInfoById(id int) (models.YYMSellerGoods, error) {
	var sellerGoods models.YYMSellerGoods
	err := global.DB.Where("id = ?", id).Find(&sellerGoods).Error
	return sellerGoods, err
}

// 根据卖家商品id查询出该商品所有的sku
func GetGoodsSkuById(id int) ([]models.YYMSellerGoodsSKU, error) {
	var sellerGoodsSku []models.YYMSellerGoodsSKU
	err := global.DB.Where("seller_goods_id =?", id).Find(&sellerGoodsSku).Error
	return sellerGoodsSku, err
}

// 根据商品ID查询出活动商品
func GetActivityGoodsByGoodsId(goodsId int) (models.YYMSellerGoodsActivity, error) {
	var activityGoods models.YYMSellerGoodsActivity
	err := global.DB.Where("seller_good_sku_id =?", goodsId).Find(&activityGoods).Error
	return activityGoods, err
}

// 根据ID查询出该商品的所有保障信息
func GetGuaranteeGoodsByGoodsId(goodsId int) ([]models.YYMSellerGoodsGuarantee, error) {
	var guaranteeGoods []models.YYMSellerGoodsGuarantee
	err := global.DB.Where("seller_goods_id =?", goodsId).Find(&guaranteeGoods).Error
	return guaranteeGoods, err
}

// 根据保障ID查询出该保障信息
func GetGuaranteeInfoById(id int) (models.YYMGuaranteeAttribute, error) {
	var guaranteeInfo models.YYMGuaranteeAttribute
	err := global.DB.Where("id =?", id).Find(&guaranteeInfo).Error
	return guaranteeInfo, err
}

// 添加商品版本记录
func AddGoodsVersion(goodsVersion models.YYMSellerGoodsVersionRecord) (models.YYMSellerGoodsVersionRecord, error) {
	err := global.DB.Create(&goodsVersion).Error
	return goodsVersion, err
}

// 修改卖家商品的版本信息
func UpdateSellerGoodsVersion(id int, version string) error {
	err := global.DB.Model(&models.YYMSellerGoods{}).Where("id =?", id).Update("version_hash", version).Error
	return err
}

// 根据版本hash去查询表当中是否存在该版本
func GetGoodsVersionByHash(version string) (models.YYMSellerGoodsVersionRecord, error) {
	var goodsVersion models.YYMSellerGoodsVersionRecord
	err := global.DB.Where("version_hash =?", version).Find(&goodsVersion).Error
	return goodsVersion, err
}

// 根据商品ID查询出商品详情
func QueryMasterGoodInfo(goodID int) (models.YYMGoodsMaster, error) {
	var good models.YYMGoodsMaster
	err := global.DB.Where("id = ?", goodID).Find(&good).Error
	return good, err
}

// 添加库存表
func AddBoxInventory(inventory models.YYMBoxInventory, tx *gorm.DB) (models.YYMBoxInventory, error) {
	err := tx.Create(&inventory).Error
	return inventory, err
}

// 添加订单
func AddOrder(order models.YYMOrderMaster, tx *gorm.DB) (models.YYMOrderMaster, error) {
	err := tx.Create(&order).Error
	return order, err
}

// 添加箱子订单
func AddBoxOrder(boxOrder models.YYMBoxOrder, tx *gorm.DB) (models.YYMBoxOrder, error) {
	err := tx.Create(&boxOrder).Error
	return boxOrder, err
}

// 根据sku_code和user_id查询出seller_goods_sku表当中的price
func QuerySellerGoodsSkuBySkuCodeAndUserID(skuCode string, userId int) (models.YYMSellerGoodsSKU, error) {
	var sku models.YYMSellerGoodsSKU
	err := global.DB.Where("sku_code = ? and user_id = ?", skuCode, userId).Find(&sku).Error
	return sku, err
}

// 根据userID查询出该用户的信息
func QueryUserInfoById(userId int) (models.User, error) {
	var user models.User
	err := global.DB.Where("id = ?", userId).Find(&user).Error
	return user, err
}

// // 根据订单编号去修改order_master的支付订单
func UpdateOrderMasterPayInfo(orderNo string, payNo string, payStatus int, tx *gorm.DB) error {
	err := tx.Model(&models.YYMOrderMaster{}).
		Where("order_no = ?", orderNo).
		Updates(map[string]interface{}{
			"pay_order_no": payNo,
			"pay_status":   payStatus,
			"payment_at":   time.Now(),
			"order_status": 1,
		}).Error
	return err
}

// UpdateBoxOrderPayInfo 根据订单编号更新支付订单编号和支付状态
func UpdateBoxOrderPayInfo(orderNo string, payNo string, payStatus int, tx *gorm.DB) error {
	err := tx.Model(&models.YYMBoxOrder{}).
		Where("order_no = ?", orderNo).
		Updates(map[string]interface{}{
			"pay_order_no": payNo,
			"pay_status":   payStatus,
		}).Error
	return err
}

// 修改库存状态为
func UpdateInventoryStatus(inventoryId int, status int, tx *gorm.DB) error {
	err := tx.Model(&models.YYMBoxInventory{}).Where("id = ?", inventoryId).Update("sell_status", status).Error
	return err
}

// 根据订单编号查询订单信息
func QueryOrderInfoByOrderNo(orderNo string) (models.YYMOrderMaster, error) {
	var order models.YYMOrderMaster
	err := global.DB.Where("order_no = ?", orderNo).Find(&order).Error
	return order, err
}

// 查询库存信息
func QueryInventoryInfoById(userId int) ([]models.YYMBoxInventory, error) {
	var inventory []models.YYMBoxInventory
	err := global.DB.Where("user_id = ?", userId).Preload("Steam").Find(&inventory).Error
	return inventory, err
}

// 根据steam账号查询出该账号的信息
func QuerySteamAccountInfoById(steamAccountId int) (models.SteamAccount, error) {
	var steamAccount models.SteamAccount
	err := global.DB.Where("id = ?", steamAccountId).Find(&steamAccount).Error
	return steamAccount, err
}

// 根据user_id查询出该用户的钱包信息
func QueryUserWalletById(userId int) (models.YYMWallet, error) {
	var userWallet models.YYMWallet
	err := global.DB.Where("user_id = ?", userId).Find(&userWallet).Error
	return userWallet, err
}

// 根据sku_code和goods_id查询出该商品的库存
func QueryInventoryBySkuCodeAndGoodsId(skuCode string, goodsId int) (models.YYMBoxInventory, error) {
	var inventory models.YYMBoxInventory
	err := global.DB.Where("sku_code = ? and goods_id = ? and sell_status = 0", skuCode, goodsId).Find(&inventory).Error
	return inventory, err
}

// 根据订单编号查询订单信息
func QueryBoxOrderInfoByOrderNo(orderNo string) (models.YYMBoxOrder, error) {
	var order models.YYMBoxOrder
	err := global.DB.Where("order_no = ?", orderNo).Find(&order).Error
	return order, err
}

// 写入发货记录表
func AddDeliveryRecord(deliveryRecord models.YYMShippingRecord, tx *gorm.DB) (models.YYMShippingRecord, error) {
	err := tx.Create(&deliveryRecord).Error
	return deliveryRecord, err
}
func QueryOutdatedInventory() ([]models.YYMBoxInventory, error) {
	var result []models.YYMBoxInventory
	sixHoursAgo := time.Now().Add(-6 * time.Hour)
	err := global.DB.Where("updated_at < ?", sixHoursAgo).Find(&result).Error
	return result, err
}
