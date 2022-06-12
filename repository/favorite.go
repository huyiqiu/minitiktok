package repository

import (
	"strconv"
	"sync"

	"github.com/go-redis/redis"
)

//"gorm.io/gorm"

type Favorite struct {
	UserId  int `gorm:"column:user_id"`
	VideoId int `gorm:"column:video_id"`
}

type FavDao struct{
}

var favDao *FavDao
var favOnce sync.Once

func NewFavDaoInstance() *FavDao{
	favOnce.Do(func() {
		db.AutoMigrate(&Favorite{})
		favDao = &FavDao{}
	})
	return favDao
}

func (*FavDao) CreateLike(userId int, videoId int) error{
	// 记录点赞数 和 点赞关系
	relationship := strconv.Itoa(userId) + ":" + strconv.Itoa(videoId)
	go func ()  { // not sure this is a right opration
		// 持久化到数据库
		favorite := &Favorite{
			UserId: userId,
			VideoId: videoId,
		}
		db.Create(&favorite)
		video := &Video{
			Id: videoId,
			UserId: userId,
		}
		db.Model(&video).Update("favorite_count", "favorite_count + 1")
	}()
	client.Set(relationship, 1, 0)
	favCnt := "LikeCntOfVlog:" + strconv.Itoa(videoId)
	cnt, err := client.Get(favCnt).Int()
	if err == redis.Nil {
		client.Set(favCnt, 1, 0)
		return nil
	}
	cnt ++
	client.Set(favCnt, cnt, 0)
	return nil
}

func (*FavDao) CancelLike(userId int, videoId int) error {
	relationship := strconv.Itoa(userId) + ":" + strconv.Itoa(videoId)
	client.Del(relationship)
	favCnt := "LikeCntOfVlog:" + strconv.Itoa(videoId)
	cnt, err := client.Get(favCnt).Int()
	if err != nil {
		println(err)
		return nil
	}
	cnt --
	client.Set(favCnt, cnt, 0)
	go func(){
		favorite := &Favorite{
			UserId: userId,
			VideoId: videoId,
		}
		db.Delete(&favorite)
		video := &Video{
			Id: videoId,
			UserId: userId,
		}
		db.Model(&video).Update("favorite_count", "favorite_count - 1")
	}()
	return nil
}



func IsFavorite(userId int, videoId int) bool {
	// var favs []*Favorite
	// db.Where("user_id = ? and video_id = ?", userId, videoId).Find(&favs)
	// return len(favs) != 0
	relationship := strconv.Itoa(userId) + ":" + strconv.Itoa(videoId)
	isFav, err := client.Exists(relationship).Result()
	if err != nil {
		println(err)
	}
	return isFav == 1
}


func (*FavDao) QueryFavList(userId int) ([]*Video, error){
	var videos []*Video
	db.Raw("select * from videos v left join favorites f where f.user_id = ? on f.video_id = v.id ORDER BY v.created_time desc", userId).Scan(&videos)
	return videos, nil
}