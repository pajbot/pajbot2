package pbtwitter

import "time"

// Follow follows given user on twitter with the given account
func (c *Client) Follow(targetUser string) error {
	for _, u := range c.followedUsers {
		if u == targetUser {
			// return if already following given user
			return nil
		}
	}
	for !c.doneLoading {
		time.Sleep(5 * time.Second)
		log.Debug("waiting to load all users before following new ones")
	}
	user, err := c.Rest.FollowUser(targetUser)
	if err != nil {
		return err
	}
	log.Debug("followed ", user.Name)
	c.followedUsers = append(c.followedUsers, targetUser)
	c.redis.SaveTwitterFollows(c.followedUsers)
	return nil
}
