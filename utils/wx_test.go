package utils

import (
	"encoding/base64"
	"testing"
)

func TestWxEncryptAndWxDecrypt(t *testing.T) {
	data := "<xml><ToUserName><![CDATA[oia2TjjewbmiOUlr6X-1crbLOvLw]]></ToUserName><FromUserName><![CDATA[gh_7f083739789a]]></FromUserName><CreateTime>1407743423</CreateTime><MsgType>  <![CDATA[video]]></MsgType><Video><MediaId><![CDATA[eYJ1MbwPRJtOvIEabaxHs7TX2D-HV71s79GUxqdUkjm6Gs2Ed1KF3ulAOA9H1xG0]]></MediaId><Title><![CDATA[testCallBackReplyVideo]]></Title><Descript  ion><![CDATA[testCallBackReplyVideo]]></Description></Video></xml>"
	key := "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
	appid := "wxxxxxxxxxxxxxxxxx"

	encrypt, err := WxEncrypt([]byte(data), key, appid)
	if err != nil {
		t.Error(err)
	}
	t.Log(base64.StdEncoding.EncodeToString(encrypt))

	decrypt, err := WxDecrypt(encrypt, key)
	if err != nil {
		t.Error(err)
	}

	if string(decrypt) != data {
		t.Fail()
	}

	log.Info().Msg(string(decrypt))
}

func TestWxDecryptUserInfo(t *testing.T) {
	data := "CiyLU1Aw2KjvrjMdj8YKliAjtP4gsMZM" +
		"QmRzooG2xrDcvSnxIMXFufNstNGTyaGS" +
		"9uT5geRa0W4oTOb1WT7fJlAC+oNPdbB+" +
		"3hVbJSRgv+4lGOETKUQz6OYStslQ142d" +
		"NCuabNPGBzlooOmB231qMM85d2/fV6Ch" +
		"evvXvQP8Hkue1poOFtnEtpyxVLW1zAo6" +
		"/1Xx1COxFvrc2d7UL/lmHInNlxuacJXw" +
		"u0fjpXfz/YqYzBIBzD6WUfTIF9GRHpOn" +
		"/Hz7saL8xz+W//FRAUid1OksQaQx4CMs" +
		"8LOddcQhULW4ucetDf96JcR3g0gfRK4P" +
		"C7E/r7Z6xNrXd2UIeorGj5Ef7b1pJAYB" +
		"6Y5anaHqZ9J6nKEBvB4DnNLIVWSgARns" +
		"/8wR2SiRS7MNACwTyrGvt9ts8p12PKFd" +
		"lqYTopNHR1Vf7XjfhQlVsAJdNiKdYmYV" +
		"oKlaRv85IfVunYzO0IKXsyl7JCUjCpoG" +
		"20f0a04COwfneQAGGwd5oa+T8yO5hzuy" +
		"Db/XcxxmK01EpqOyuxINew=="
	key := "tiihtNczf5v6AKRyjwEUhQ=="
	iv := "r7BXXKkLb8qrSNn05n0qiA=="

	decrypt, err := WxDecryptUserInfo(data, key, iv)
	if err != nil {
		t.Error(err)
	}

	expect := `{"openId":"oGZUI0egBJY1zhBYw2KhdUfwVJJE","nickName":"Band","gender":1,"language":"zh_CN","city":"Guangzhou","province":"Guangdong","country":"CN","avatarUrl":"http://wx.qlogo.cn/mmopen/vi_32/aSKcBBPpibyKNicHNTMM0qJVh8Kjgiak2AHWr8MHM4WgMEm7GFhsf8OYrySdbvAMvTsw3mo8ibKicsnfN5pRjl1p8HQ/0","unionId":"ocMvos6NjeKLIBqg5Mr9QjxrP1FA","watermark":{"timestamp":1477314187,"appid":"wx4f4bc4dec97d474b"}}`
	if string(decrypt) != expect {
		t.Fail()
	}

	log.Info().Msg(string(decrypt))
}
