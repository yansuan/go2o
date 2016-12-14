/**
 * Copyright 2015 @ z3q.net.
 * name : snapshot
 * author : jarryliu
 * date : 2016-06-28 21:41
 * description :
 * history :
 */
package item

type (
    // 快照服务
    ISnapshotManager interface {
        // 生成商品快照
        GenerateSnapshot() (int32, error)

        // 获取最新的快照
        GetLatestSnapshot() *Snapshot

        // 获取最新的商品销售快照,如果商品有更新,则更新销售快照
        GetLatestSaleSnapshot() *SalesSnapshot

        // 根据KEY获取已销售商品的快照
        GetSaleSnapshotByKey(key string) *SalesSnapshot

        // 根据ID获取已销售商品的快照
        GetSaleSnapshot(id int32) *SalesSnapshot
    }

    // 商品快照
    Snapshot struct {
        //SKU编号
        SkuId        int32 `db:"sku_id" auto:"no" pk:"no"`
        //快照编号: 商户编号+g商品编号+快照时间戳
        Key          string `db:"snapshot_key"`
        //供应商编号
        VendorId     int32 `db:"vendor_id"`
        //商品编号
        //GoodsId int `db:"goods_id"`
        //商品标题
        GoodsTitle   string `db:"goods_title"`
        //小标题
        ShortTitle   string `db:"small_title"`
        //货号
        GoodsNo      string `db:"goods_no"`
        //货品编号
        ItemId       int32 `db:"item_id"`
        //分类编号
        CategoryId   int32 `db:"cat_id"`
        //SKU  todo:????
        Sku          string `db:"-"`
        //运费模板编号
        ExpressTplId int32 `db:"express_tid"`
        //图片
        Image        string `db:"img"`
        // 供货价
        Cost         float32 `db:"cost"`
        //定价
        RetailPrice  float32 `db:"retail_price"`
        //销售价
        SalePrice    float32 `db:"sale_price"`
        // 单件重量,单位:千克(kg)
        Weight       int32 `db:"weight"`
        // 体积:毫升(ml)
        Bulk         int32 `db:"bulk"`
        //是否有会员价
        LevelSales   int `db:"level_sales"`
        //销售数量
        SaleNum      int32 `db:"sale_num"`
        //库存
        StockNum     int32 `db:"stock_num"`
        // 是否上架
        ShelveState  int32 `db:"shelve_state"`
        //快照时间
        UpdateTime   int64 `db:"update_time"`
    }

    // 已销售商品快照
    SalesSnapshot struct {
        //快照编号
        Id          int32 `db:"id" auto:"yes" pk:"yes"`
        //商品SKU编号
        SkuId       int32 `db:"sku_id"`
        //快照编码: 商户编号+g商品编号+快照时间戳
        SnapshotKey string `db:"snap_key"`
        // 卖家编号
        SellerId    int32 `db:"seller_id"`
        // 卖家名称
        //SellerName  string `db:"seller_name"`
        //商品标题
        GoodsTitle  string `db:"goods_title"`
        //小标题
        //SmallTitle  string `db:"-"`
        //货号
        GoodsNo     string `db:"goods_no"`
        //货品编号
        ItemId      int32 `db:"item_id"`
        //分类编号
        CategoryId  int32 `db:"cat_id"`
        //SKU  todo:????
        Sku         string `db:"sku"`
        //图片
        Image       string `db:"img"`
        // 供货价
        Cost        float32 `db:"cost"`
        //销售价
        Price       float32 `db:"price"`
        // 快照时间
        CreateTime  int64 `db:"create_time"`
    }
)