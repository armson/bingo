package weixin

import(
	"github.com/armson/bingo/utils"
	"github.com/armson/bingo/attach"
	"errors"
	"fmt"
)

//判断素材类型
func CheckMaterialType(t string) bool {
	return utils.Slice.In(t, []string{"image", "video", "voice", "news"})
}

//获取素材总数
//例如：
//count , _ := weixin.New(c).MaterialCount()
func (bin *Interface) MaterialCount() (map[string]int , error) {
	accessToken, err := bin.Token()
	if err != nil {
		return nil, err
	}

	bin.Url = bin.cgiUrl("/cgi-bin/material/get_materialcount")
	bin.Method = "GET"
	bin.Queries["access_token"] =  accessToken

	js , err :=  bin.Send()
	if err != nil {
		return nil, requestWithoutResponse
	}

	errCode, err := js.Get("errcode").Int()
	if err == nil {
		errMsg, _ := js.Get("errmsg").String()
		return nil, &Error{errCode, errMsg}
	}

	m := map[string]int{}
	m["voice"] , _  = js.Get("voice_count").Int()
	m["video"] , _  = js.Get("video_count").Int()
	m["image"] , _  = js.Get("image_count").Int()
	m["news"] , _   = js.Get("news_count").Int()
	return m, nil
}

// 获取素材列表(未调试)
// t:素材的类型，图片（image）、视频（video）、语音 （voice）、图文（news）
// 例如：
// list, total, _ := weixin.New(c).MaterialList("news",0,10)
func (bin *Interface) MaterialList(t string, offset, limit int) ([]map[string]interface{}, int , error) {
	accessToken, err := bin.Token()
	if err != nil {
		return nil, 0, err
	}

	bin.Url = bin.cgiUrl("/cgi-bin/material/batchget_material", accessToken)
	bin.Method = "POST"
	body := map[string]interface{}{
		"type":t,
		"offset":offset,
		"count":limit,
	}
	bin.Raw(body)

	js , err :=  bin.Send()
	if err != nil {
		return nil, 0, requestWithoutResponse
	}

	if errCode, err := js.Get("errcode").Int(); err == nil {
		errMsg, _ := js.Get("errmsg").String()
		return nil, 0, &Error{errCode, errMsg}
	}


	total_count, err := js.Get("total_count").Int()
	if err != nil {
		return nil, 0 , &ParseError{"MaterialList","total_count"}
	}

	//解析
	m := []map[string]interface{}{}
	if total_count == 0 {
		return m , 0 , nil
	}

	itemsJs := js.Get("item")
	items ,err := itemsJs.Array()
	if err != nil || len(items) == 0 {
		return m , total_count , err
	}

	if t == "image" || t == "video" || t == "voice" {
		for k, _ := range items {
			item ,err := itemsJs.GetIndex(k).Map()
			if err != nil { continue }
			m = append(m, map[string]interface{}{
				"mediaId"     : item["media_id"],
				"name" 		  : item["name"],
				"updateTime"  : item["update_time"],
				"url" 		  : item["url"],
			})
		}

	}

	if t == "news" {
		for k, _ := range items {
			item, err := itemsJs.GetIndex(k).Get("content").Get("news_item").GetIndex(0).Map()
			if err != nil { continue }
			media_id, _ := itemsJs.GetIndex(k).Get("media_id").String()
			create_time, _ := itemsJs.GetIndex(k).Get("content").Get("create_time").Int()
			update_time, _ := itemsJs.GetIndex(k).Get("content").Get("update_time").Int()

			m = append(m, map[string]interface{}{
				"mediaId"     			: media_id,
				"title" 		 		: item["title"],
				"author"  				: item["author"],
				"digest" 		  		: item["digest"],
				"content" 		  		: item["content"],
				"contentSourceUrl" 		: item["content_source_url"],
				"thumbMediaId" 		  	: item["thumb_media_id"],
				"showCoverPic" 		  	: item["show_cover_pic"],
				"url" 		  			: item["url"],
				"thumbUrl" 		  		: item["thumb_url"],
				"needOpenComment" 		: item["need_open_comment"],
				"onlyFansCanComment"  	: item["only_fans_can_comment"],
				"createTime" 		  	: create_time,
				"updateTime" 		  	: update_time,
			})
		}
	}
	return m ,total_count, nil
}

//上传图文消息内的图片获取URL
//本接口所上传的图片不占用公众号的素材库中图片数量的5000个的限制。图片仅支持jpg/png格式，大小必须在1MB以下
//示例：
//file, _ := c.File("image")
//url , err := weixin.New(c).MaterialUploadImage(file)
func (bin *Interface) MaterialUploadImage(file *attach.Attachment) (string, error) {
	accessToken, err := bin.Token()
	if err != nil {
		return "", err
	}

	if file.Size > MaterialUploadImageMaxSize() {
		return "" ,errors.New("Image size exceed the maximum limit of 1 M")
	}

	if	!CheckMaterialImageMime(file) {
		return "", fmt.Errorf("This file's attribute is %s and %s Forbidden to upload the image type", file.Mime , file.Ext)
	}

	bin.Url = bin.cgiUrl("/cgi-bin/media/uploadimg?type=image", accessToken)
	bin.Method = "POST"
	bin.Reader("image", file.Name, file.Rc())
	js , err :=  bin.Send()
	if err != nil {
		return "", requestWithoutResponse
	}

	errCode, err := js.Get("errcode").Int()
	if err == nil {
		errMsg, _ := js.Get("errmsg").String()
		return "", &Error{errCode, errMsg}
	}
	url , err := js.Get("url").String()
	if err != nil {
		return "", &ParseError{"MaterialUploadImage","url"}
	}
	return url, nil
}

//上传永久图片
//file, err := c.File("image")
//item , err := weixin.New(c).MaterialUploadImageAlways(file)
func (bin *Interface) MaterialUploadImageAlways(file *attach.Attachment) (map[string]string, error) {
	accessToken, err := bin.Token()
	if err != nil {
		return nil, err
	}
	if file.Size > MaterialUploadImageAlwaysMaxSize() {
		return nil ,errors.New("Image size exceed the maximum limit of 2 M")
	}

	if	!CheckMaterialImageMime(file) {
		return nil, errors.New("Forbidden to upload the image type")
	}

	bin.Url = bin.cgiUrl("/cgi-bin/material/add_material?type=image", accessToken)
	bin.Method = "POST"
	bin.Reader("media", file.Name, file.Rc())
	js , err :=  bin.Send()
	if err != nil {
		return nil, requestWithoutResponse
	}

	errCode, err := js.Get("errcode").Int()
	if err == nil {
		errMsg, _ := js.Get("errmsg").String()
		return nil, &Error{errCode, errMsg}
	}

	item , err := js.Map()
	if err != nil {
		return nil, &ParseError{"MaterialUploadImageAlways","response"}
	}
	return map[string]string{
		"mediaId" 	: item["media_id"].(string),
		"url" 		: item["url"].(string),
	}, nil
}




//获取永久素材
// 仅能获取图文消息、视频，其他类型，如图片、音频暂不支持
// detail , err := weixin.New(c).MaterialDetail(c.GET("mediaId"))
func (bin *Interface) MaterialDetail(mediaId string) ([]map[string]interface{}, error) {
	accessToken, err := bin.Token()
	if err != nil {
		return nil, err
	}
	bin.Url = bin.cgiUrl("/cgi-bin/material/get_material", accessToken)
	bin.Method = "POST"
	bin.Raw(map[string]string{"media_id":mediaId})

	js , err :=  bin.Send()
	if err != nil {
		return nil, requestWithoutResponse
	}

	if errCode, err := js.Get("errcode").Int(); err == nil {
		errMsg, _ := js.Get("errmsg").String()
		return nil, &Error{errCode, errMsg}
	}

	items, err := js.Get("news_item").Array()
	if err != nil {
		return nil , &ParseError{"MaterialDetail","news_item"}
	}

	m := []map[string]interface{}{}
	for i := 0 ; i < len(items) ; i ++ {
		item, _ := js.Get("news_item").GetIndex(i).Map()
		if item != nil {
			create_time, _ := js.Get("create_time").Int()
			update_time, _ := js.Get("update_time").Int()
			m = append(m , map[string]interface{}{
				"mediaId"     			: mediaId,
				"title" 		 		: item["title"],
				"author"  				: item["author"],
				"digest" 		  		: item["digest"],
				"content" 		  		: item["content"],
				"contentSourceUrl" 		: item["content_source_url"],
				"thumbMediaId" 		  	: item["thumb_media_id"],
				"showCoverPic" 		  	: item["show_cover_pic"],
				"url" 		  			: item["url"],
				"thumbUrl" 		  		: item["thumb_url"],
				"needOpenComment" 		: item["need_open_comment"],
				"onlyFansCanComment"  	: item["only_fans_can_comment"],
				"createTime" 		  	: create_time,
				"updateTime" 		  	: update_time,
			})
		}
	}
	return m, nil
}

//新增永久图文素材
//例如：
//article := map[string]string{
//		"title" 				: c.POST("title"),
//		"thumb_media_id" 		: c.POST("thumbMediaId"),
//		"author" 				: c.POST("author"),
//		"digest" 				: c.POST("digest"),
//		"show_cover_pic" 		: c.POST("showCoverPic"),
//		"content" 				: c.POST("content"),
//		"content_source_url" 	: c.POST("contentSourceUrl"),
//}
//mediaId, err := weixin.New(c).MaterialArticleCreate(article)
func (bin *Interface) MaterialArticleCreate(article map[string]string) (string, error) {
	accessToken, err := bin.Token()
	if err != nil {
		return "", err
	}

	if article["thumb_media_id"] == "" {
		return "",errors.New("The lack of cover image")
	}

	if article["title"] == "" {
		return "",errors.New("The title can't be empty")
	}

	if article["content"] == "" {
		return "",errors.New("The content can't be empty")
	}

	bin.Url = bin.cgiUrl("/cgi-bin/material/add_news", accessToken)
	bin.Method = "POST"
	articles := []map[string]string{article}
	bin.Raw(map[string]interface{}{"articles":articles})

	js , err :=  bin.Send()
	if err != nil {
		return "", requestWithoutResponse
	}

	if errCode, err := js.Get("errcode").Int(); err == nil {
		errMsg, _ := js.Get("errmsg").String()
		return "", &Error{errCode, errMsg}
	}

	mediaId, err := js.Get("media_id").String()
	if err != nil {
		return "", &ParseError{"MaterialNewsCreate","media_id"}
	}
	return mediaId , nil
}

//更新图文素材
//例如：
//mediaId := c.POST("mediaId")
//		article := map[string]string{
//		"title" 				: c.POST("title"),
//		"thumb_media_id" 		: c.POST("thumbMediaId"),
//		"author" 				: c.POST("author"),
//		"digest" 				: c.POST("digest"),
//		"show_cover_pic" 		: c.POST("showCoverPic"),
//		"content" 				: c.POST("content"),
//		"content_source_url" 	: c.POST("contentSourceUrl"),
//}
//mediaId, err := weixin.New(c).MaterialArticleUpdate(mediaId, article)
func (bin *Interface) MaterialArticleUpdate(mediaId string, article map[string]string) (string, error) {
	accessToken, err := bin.Token()
	if err != nil {
		return "", err
	}

	if mediaId == "" {
		return "",errors.New("mediaId Can't be empty")
	}

	if article["thumb_media_id"] == "" {
		return "",errors.New("The lack of cover image")
	}

	if article["title"] == "" {
		return "",errors.New("The title can't be empty")
	}

	if article["content"] == "" {
		return "",errors.New("The content can't be empty")
	}

	bin.Url = bin.cgiUrl("/cgi-bin/material/update_news", accessToken)
	bin.Method = "POST"
	bin.Raw(map[string]interface{}{"media_id":mediaId,"index":0,"articles":article})

	js , err :=  bin.Send()
	if err != nil {
		return "", requestWithoutResponse
	}

	errCode, err := js.Get("errcode").Int();
	if  errCode != 0  {
		errMsg, _ := js.Get("errmsg").String()
		return "", &Error{errCode, errMsg}
	}
	return mediaId , nil
}


// 删除永久素材
// mediaId , err := weixin.New(c).MaterialDelete(c.POST("mediaId"))
func (bin *Interface) MaterialDelete(mediaId string) (string, error) {
	accessToken, err := bin.Token()
	if err != nil {
		return "", err
	}

	bin.Url = bin.cgiUrl("/cgi-bin/material/del_material", accessToken)
	bin.Method = "POST"
	bin.Raw(map[string]string{"media_id":mediaId})

	js , err :=  bin.Send()
	if err != nil {
		return "", requestWithoutResponse
	}

	errCode, err := js.Get("errcode").Int();
	if  errCode != 0  {
		errMsg, _ := js.Get("errmsg").String()
		return "", &Error{errCode, errMsg}
	}
	return mediaId , nil
}





