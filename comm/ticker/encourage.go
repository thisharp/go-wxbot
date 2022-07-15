package ticker

import (
	"fmt"
	"os"
	"time"

	"github.com/eatmoreapple/openwechat"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"go-wxbot/openwechat/comm/global"
	"go-wxbot/openwechat/comm/image"
)

// 程序员鼓励师

func Encourage() {
	for {
		select {
		case t := <-time.After(1 * time.Minute):
			nowTime := t.Format("15:04")
			if nowTime != "11:55" {
				continue
			}

			var (
				err error
				message,
				imgURL, imgPath string
				groups openwechat.Groups
			)

			message = fmt.Sprintf("BUG 虽好，但不要贪多哦！程序员鼓励师提醒，该吃午饭了~")

			imgURL, err = image.GetImage()
			if err != nil {
				err = errors.Wrapf(err, "Encourage get image err")
				logrus.Error(err.Error())
				continue
			}

			imgPath, err = image.SaveEncourageImg(imgURL)
			if err != nil {
				err = errors.Wrapf(err, "Encourage save image err")
				logrus.Error(err.Error())
				continue
			}

			img, err := os.Open(imgPath)
			defer img.Close()
			if err != nil {
				err = errors.Wrapf(err, "reword open file err")
				logrus.Error(err.Error())
				continue
			}

			groups, err = global.WxSelf.Groups(true)
			if err != nil {
				err = errors.Wrapf(err, "SendMessageToFans get groups err")
				logrus.Error(err.Error())
				continue
			}

			// 后场村粉丝群
			groups.SearchByNickName(1, global.Conf.Keys.HouchangcunFans).SendText(message)
			groups.SearchByNickName(1, global.Conf.Keys.HouchangcunFans).SendImage(img)

			// 五壮士群
			for _, each := range groups {
				members, err := each.Members()
				if err != nil {
					err = errors.Wrapf(err, "SendMessageToFans get members err")
					logrus.Error(err.Error())
					continue
				}

				// 不能通过群备注来获取群，真是恶心
				var Is = false
				for _, member := range members {
					if member.NickName == "李欢庭" || member.NickName == "邢宇超" {
						Is = true
						break
					}
				}

				if Is {
					each.SendText(message)
					each.SendImage(img)
				}
			}

			os.Remove(imgPath)
		}
	}
}
