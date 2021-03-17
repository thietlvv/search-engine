package log

import (
	"auth/models"
	dt "auth/utils/datastore"
	"context"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/labstack/gommon/log"
)

func LogSMS(client *datastore.Client, brandname string, mtid string, smsType string, phone string) {
	smsLog := new(models.SMSLog)
	smsLog.Brandname = brandname
	smsLog.MTID = mtid
	smsLog.Type = smsType
	smsLog.Phone = phone
	smsLog.SentAt = time.Now()
	ctx := context.Background()
	_, err := dt.Put(client, ctx, "sms_logs", smsLog)
	if err != nil {
		log.Error(err)
	}
}
