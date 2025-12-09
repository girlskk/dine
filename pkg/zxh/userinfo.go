package zxh

import "context"

type UserInfo struct {
	ID          int64  `json:"uid"`
	Name        string `json:"name"`
	PhoneNumber string `json:"phoneNumber"`
}

func (c *Client) GetUserInfo(ctx context.Context, payCode string) (*UserInfo, error) {
	userInfo := new(UserInfo)
	req := newRequest(c, UserInfoPath, map[string]any{"payCode": payCode}, userInfo)
	if err := req.do(ctx); err != nil {
		return nil, err
	}

	return userInfo, nil
}
