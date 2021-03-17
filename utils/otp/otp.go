package otp

import (
	"auth/config"
	dt "auth/utils/datastore"
	ulg "auth/utils/log"
	"auth/utils/sms"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/go-redis/redis"
	"github.com/labstack/gommon/log"
	uuid "github.com/satori/go.uuid"
)

func WhiteListsPhoneNumber(client *datastore.Client, ctx context.Context, phone_number string) bool {
	type resq_struct struct {
		WhiteLists []string `datastore:"white_lists"`
	}
	filter := []dt.DatastoreFilter{}
	docs := new([]resq_struct)
	err := dt.Get(client, ctx, "phones_config", docs, filter, 1)
	if err != nil {
		log.Error(err)
		return false
	}
	if len(*docs) > 0 {
		fmt.Println("White lists: ", (*docs)[0].WhiteLists)
		for _, s := range (*docs)[0].WhiteLists {
			if s == phone_number {
				return true
			}
		}
	}
	return false
}

func GenerateMega1OTP(client *datastore.Client, ctx context.Context, phone_number string) string {
	if WhiteListsPhoneNumber(client, ctx, phone_number) == true {
		return "0000"
	}
	rand.Seed(time.Now().UnixNano())
	min := 1000
	max := 9999
	OTP := rand.Intn(max-min+1) + min
	mtid := uuid.NewV4()
	brandname := "MEGA1MALL"

	body := fmt.Sprintf(`{
    "dest": "%s",
    "brandname": "%s",
    "msgbody": "OTP la %d de xac thuc TK MEGA1MALL, ma dung 1 lan & hieu luc trong 3 phut. KHONG chia se ma cho nguoi khac. Chi tiet hotline 1900633996",
    "content_type": "text",
    "serviceid": "Gapit",
    "mtid": "%s",
    "cpid": %d
	}`, phone_number, brandname, OTP, mtid, config.C.SMSServices.GAPIT.Cpid)
	content := sms.SendOTP(body)

	var resq_struct struct {
		MTid   string `json:"mtid"`
		Status int    `json:"status"`
	}
	err := json.Unmarshal(content, &resq_struct)
	if err != nil {
		log.Error(err)
		return ""
	}
	if resq_struct.Status == 200 {
		go ulg.LogSMS(client, brandname, resq_struct.MTid, "OTP", phone_number)
		return fmt.Sprintf(`%d`, OTP)
	}
	return ""
}

func CheckLockPhoneNumber(redisClient *redis.Client, phone_number string, increase int) error {
	key_lock := fmt.Sprintf(`BLOCK_PHONE_NUMBER_%s`, phone_number)
	key_lock_val, err := redisClient.Get(key_lock).Result()
	if err == nil {
		ttl, _ := redisClient.TTL(key_lock).Result()
		i, _ := strconv.Atoi(key_lock_val)
		if i == 3 {
			redisClient.Del(fmt.Sprintf(`OTP_%s`, phone_number))
			return errors.New("Số điện thoại này đã bị khoá.")
		} else {
			redisClient.Set(key_lock, i+increase, ttl).Err()
		}
	} else if increase > 0 {
		redisClient.Set(key_lock, increase, 24*time.Hour).Err()
	}
	return nil
}

func VerifyGenOTP(redisClient *redis.Client, phone string, otp string) bool {
	key := fmt.Sprintf(`OTP_GEN_%s`, phone)
	val_otp, err := redisClient.Get(key).Result()
	if err != nil {
		return VerifyOTP(redisClient, phone, otp)
	}
	if val_otp != otp {
		CheckLockPhoneNumber(redisClient, phone, 1)
		return VerifyOTP(redisClient, phone, otp)
	}
	return true
}

func VerifyOTP(redisClient *redis.Client, phone string, otp string) bool {
	key := fmt.Sprintf(`OTP_%s`, phone)
	val_otp, err := redisClient.Get(key).Result()
	if err != nil {
		return false
	}
	if val_otp != otp {
		CheckLockPhoneNumber(redisClient, phone, 1)
		return false
	}
	return true
}

func DelOTP(redisClient *redis.Client, phone string) {
	redisClient.Del(fmt.Sprintf(`OTP_GEN_%s`, phone)).Err()
	redisClient.Del(fmt.Sprintf(`OTP_%s`, phone)).Err()
}
