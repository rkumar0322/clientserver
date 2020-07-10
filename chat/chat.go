package chat

import (
	"errors"
	"golang.org/x/net/context"
	"log"
	"strings"
	"time"
)

type Server struct {
	users    []*user
	channels []*channel
}

type user struct {
	username    string
	createdtime time.Time
}

type channel struct {
	name   string
	public bool
	admins []*user
	banned []*user
	users  []*user
	cm     []*channelMessage
}
type channelMessage struct {
	username *user
	msg      string
}

/*
AddUser: "AddUser,username"

AddingChannel: "AddChannel,username,channelname"

AddUsersToChannel "AddUsertoChanel,username,othername,channelname"

JoinChannel "JoinChannel,username,channelname"

SendMessage "sendmessage,username,channelname, message"

showchannel

showusers
*/

func (s *Server) SayHello(ctx context.Context, in *Message) (*Message, error) {
	log.Printf("Receive message body from client: %s", in.Body)
	msg := in.Body
	msgArr := strings.Split(msg, ",")
	if msgArr[0] == "adduser" {
		return s.addUser(msgArr[1])
	} else if msgArr[0] == "addchannel" {
		return s.addChannel(msgArr[2], msgArr[1], msgArr[3])
	} else if msgArr[0] == "removechannel" {
		return s.removeChannel(msgArr[2], msgArr[1])
	} else if msgArr[0] == "addusertochannel" {
		return s.AddUsersToChannelServ(msgArr[1], msgArr[2], msgArr[3])
	} else if msgArr[0] == "removeuserfromchannel" {
		return s.removeUsersFromChannelServ(msgArr[1], msgArr[2], msgArr[3])
	} else if msgArr[0] == "banuserfromchannel" {
		return s.AddUsersToBannedServ(msgArr[1], msgArr[2], msgArr[3])
	} else if msgArr[0] == "removebanuser" {
		return s.removeUsersFromBannedServ(msgArr[1], msgArr[2], msgArr[3])
	} else if msgArr[0] == "joinchannel" {
		return s.joinChannelServ(msgArr[1], msgArr[2])
	} else if msgArr[0] == "leavechannel" {
		return s.leaveChannelServ(msgArr[1], msgArr[2])
	} else if msgArr[0] == "sendMessage" {
		return s.addMessageServ(msgArr[1], msgArr[2], msgArr[3])
	} else if msgArr[0] == "showWorkspace" {
		return &Message{Body: s.toString(msgArr[1])}, nil
	} else if msgArr[0] == "showChannel" {
		ch, _, _ := s.grabChannel(msgArr[2])
		return &Message{Body: ch.toString(msgArr[1])}, nil
	} else if msgArr[0] == "userExists" {
		if s.userExists(msgArr[1]) {
			return &Message{Body: "true"}, nil
		} else {
			return &Message{Body: "false"}, nil
		}
	} else {
		return &Message{Body: "User Does Not Exist"}, nil
	}
}

func (s *Server) toString(user1 string) string {
	str := "\nUsers:\n"
	for i := 0; i < len(s.users); i++ {
		str += s.users[i].toString() + "\n"
	}
	str += "\nExisting Channels:\n" + s.showChannels(user1)
	str += "\nChannels To Join:\n" + s.showChannelsToJoin(user1)
	return str
}

func (s *Server) joinChannelServ(user1 string, channel2 string) (*Message, error) {
	channel1, _, err := s.grabChannel(channel2)
	if err != nil {
		return &Message{Body: "Channel IS NOT VALID"}, nil
	}
	user2, _, err2 := s.grabUser(user1)
	if err2 != nil {
		return &Message{Body: "User IS NOT VALID"}, nil
	}
	return channel1.joinChannel(user2)
}

func (s *Server) leaveChannelServ(user1 string, channel2 string) (*Message, error) {
	channel1, _, err := s.grabChannel(channel2)
	if err != nil {
		return &Message{Body: "Channel IS NOT VALID"}, nil
	}
	user2, _, err2 := s.grabUser(user1)
	if err2 != nil {
		return &Message{Body: "User IS NOT VALID"}, nil
	}
	return channel1.leaveChannel(user2)
}

func (s *Server) AddUsersToChannelServ(userAdding1 string, userToAdd1 string, channel2 string) (*Message, error) {
	channel1, _, err := s.grabChannel(channel2)
	if err != nil {
		return &Message{Body: "Channel IS NOT VALID"}, nil
	}
	userToAdd, _, err2 := s.grabUser(userToAdd1)
	if err2 != nil {
		return &Message{Body: "User IS NOT VALID"}, nil
	}

	userAdding, _, err3 := s.grabUser(userAdding1)
	if err3 != nil {
		return &Message{Body: "User IS NOT VALID"}, nil
	}
	return channel1.addUsersToChannel(userAdding, userToAdd)
}

func (s *Server) removeUsersFromChannelServ(userAdding1 string, userToAdd1 string, channel2 string) (*Message, error) {
	channel1, _, err := s.grabChannel(channel2)
	if err != nil {
		return &Message{Body: "channel IS NOT VALID"}, nil
	}
	userToAdd, _, err2 := s.grabUser(userToAdd1)
	if err2 != nil {
		return &Message{Body: "User IS NOT VALID"}, nil
	}

	userAdding, _, err3 := s.grabUser(userAdding1)
	if err3 != nil {
		return &Message{Body: "User IS NOT VALID"}, nil
	}
	return channel1.removeUsersFromChannel(userAdding, userToAdd)
}

func (s *Server) AddUsersToBannedServ(userAdding1 string, userToAdd1 string, channel2 string) (*Message, error) {
	channel1, _, err := s.grabChannel(channel2)
	if err != nil {
		return &Message{Body: "Channel IS NOT VALID"}, nil
	}
	userToAdd, _, err2 := s.grabUser(userToAdd1)
	if err2 != nil {
		return &Message{Body: "User IS NOT VALID"}, nil
	}

	userAdding, _, err3 := s.grabUser(userAdding1)
	if err3 != nil {
		return &Message{Body: "User IS NOT VALID"}, nil
	}
	return channel1.addUsersToBanned(userAdding, userToAdd)
}

func (s *Server) removeUsersFromBannedServ(userAdding1 string, userToAdd1 string, channel2 string) (*Message, error) {
	channel1, _, err := s.grabChannel(channel2)
	if err != nil {
		return &Message{Body: "User IS NOT VALID"}, nil
	}
	userToAdd, _, err2 := s.grabUser(userToAdd1)
	if err2 != nil {
		return &Message{Body: "User IS NOT VALID"}, nil
	}

	userAdding, _, err3 := s.grabUser(userAdding1)
	if err3 != nil {
		return &Message{Body: "User IS NOT VALID"}, nil
	}
	return channel1.removeUsersFromBanned(userAdding, userToAdd)
}

func (s *Server) addMessageServ(user1 string, channel2 string, m string) (*Message, error) {
	channel1, _, err := s.grabChannel(channel2)
	if err != nil {
		return &Message{Body: "User IS NOT VALID"}, nil
	}
	user2, _, err2 := s.grabUser(user1)
	if err2 != nil {
		return &Message{Body: "Channel IS NOT VALID"}, nil
	}
	return channel1.addMessageToChannel(user2, m)
}

func (s *Server) addChannel(chan1 string, username string, public string) (*Message, error) {
	if s.channelExists(chan1) {
		return &Message{Body: "channel Exists"}, nil
	}
	user1, _, err1 := s.grabUser(username)
	if err1 != nil {
		return &Message{Body: "User Does Not Exist"}, nil
	}
	var public1 bool
	if public == "yes" {
		public1 = true
	} else if public == "no" {
		public1 = false
	} else {
		return &Message{Body: "Invalid Input"}, nil
	}
	var users []*user
	users = append(users, user1)
	var admins []*user
	admins = append(admins, user1)
	var banned []*user
	var messages []*channelMessage
	s.channels = append(s.channels, &channel{chan1, public1, admins, banned, users, messages})
	return &Message{Body: "SUCCESS: Creaeted Channel " + chan1}, nil
}

func (s *Server) removeChannel(chan1 string, username string) (*Message, error) {
	user1, _, err1 := s.grabUser(username)
	if err1 != nil {
		return &Message{Body: "User Does Not Exist"}, nil
	}
	ch1, pos, err2 := s.grabChannel(chan1)
	if err2 != nil {
		return &Message{Body: "Channel Does Not exist"}, nil
	}
	if ch1.isUserAdmin(user1) == false {
		return &Message{Body: "Admins Can only remove the channel"}, nil
	}
	s.channels, _ = deleteChannel(s.channels, pos)
	return &Message{Body: "SUCCESS: Deleted Channel " + chan1}, nil
}

func (s *Server) addUser(user1 string) (*Message, error) {
	if s.userExists(user1) {
		return &Message{Body: "USER IS NOT VALID"}, nil
	}
	s.users = append(s.users, &user{user1, time.Now()})
	return &Message{Body: "SUCCESS: User is " + user1}, nil
}

/*
Helper Functions
*/

func (s *Server) showChannels(user1 string) string {
	user2, _, _ := s.grabUser(user1)
	str := "\n"
	for i := 0; i < len(s.channels); i++ {
		if s.channels[i].isUserInChannel(user2) {
			str += s.channels[i].name + "\n"
		}
	}
	return str
}

func (s *Server) showAdminChannels(user1 string) string {
	user2, _, _ := s.grabUser(user1)
	str := "\n"
	for i := 0; i < len(s.channels); i++ {
		if s.channels[i].isUserAdmin(user2) {
			str += s.channels[i].name + "\n"
		}
	}
	return str
}

func (s *Server) showChannelsToJoin(user1 string) string {
	user2, _, _ := s.grabUser(user1)
	str := "\n"
	for i := 0; i < len(s.channels); i++ {
		if s.channels[i].isUserInChannel(user2) == false && s.channels[i].public && s.channels[i].isUserBanned(user2) == false {
			str += s.channels[i].name + "\n"
		}
	}
	return str
}

func (s *Server) channelExists(channel1 string) bool {
	for i := 0; i < len(s.channels); i++ {
		if s.channels[i].name == channel1 {
			return true
		}
	}
	return false
}

func (s *Server) grabUser(user1 string) (*user, int, error) {
	for i := 0; i < len(s.users); i++ {
		if s.users[i].username == user1 {
			return s.users[i], i, nil
		}
	}
	return nil, -1, errors.New("")
}



func (s *Server) grabChannel(channel1 string) (*channel, int, error) {
	for i := 0; i < len(s.channels); i++ {
		if s.channels[i].name == channel1 {
			return s.channels[i], i, nil
		}
	}
	return nil, -1, errors.New("")
}

func (s *Server) userExists(user1 string) bool {
	for i := 0; i < len(s.users); i++ {
		if s.users[i].username == user1 {
			return true
		}
	}
	return false
}

func (s *channel) grabUser(user1 string) (*user, int, error) {
	for i := 0; i < len(s.users); i++ {
		if s.users[i].username == user1 {
			return s.users[i], i, nil
		}
	}
	return nil, -1, errors.New("")
}

func (s *channel) grabBannedUser(user1 string) (*user, int, error) {
	for i := 0; i < len(s.banned); i++ {
		if s.banned[i].username == user1 {
			return s.banned[i], i, nil
		}
	}
	return nil, -1, errors.New("")
}

func (s *channel) isUserInChannel(user1 *user) bool {
	for i := 0; i < len(s.users); i++ {
		if s.users[i].username == user1.username {
			return true
		}
	}
	return false
}

func (s *channel) isUserAdmin(user1 *user) bool {
	for i := 0; i < len(s.admins); i++ {
		if s.admins[i].username == user1.username {
			return true
		}
	}
	return false
}

func (s *channel) isUserBanned(user1 *user) bool {
	for i := 0; i < len(s.banned); i++ {
		if s.banned[i].username == user1.username {
			return true
		}
	}
	return false
}

func (u *user) toString() string {
	return u.username + ": " + u.createdtime.Format(time.RFC1123)
}

func (u *channelMessage) toString() string {
	return u.username.username + ": " + u.msg
}

func deleteUser(a []*user, i int) ([]*user, error) {
	if i < len(a) {
		copy(a[i:], a[i+1:]) // Shift a[i+1:] left one index.
		a[len(a)-1] = nil    // Erase last element (write zero value).
		a = a[:len(a)-1]     // Truncate slice.
		return a, nil
	} else {
		return a, errors.New("Record Not Found")
	}
}

func deleteChannel(a []*channel, i int) ([]*channel, error) {
	if i < len(a) {
		copy(a[i:], a[i+1:]) // Shift a[i+1:] left one index.
		a[len(a)-1] = nil    // Erase last element (write zero value).
		a = a[:len(a)-1]     // Truncate slice.
		return a, nil
	} else {
		return a, errors.New("Record Not Found")
	}
}

func (s *channel) addUsersToChannel(userAdding *user, userToAdd *user) (*Message, error) {
	//checks if the user thats adding someone else is in the channel
	if s.isUserInChannel(userAdding) == false {
		return &Message{Body: "USER IS NOT VALID"}, nil
	}

	//checks if the user thats adding someone else is in the channel
	if s.isUserAdmin(userAdding) == false {
		return &Message{Body: "USER IS NOT VALID"}, nil
	}
	//checks if the user that will be added to the channel isn't already in the channel
	if s.isUserInChannel(userToAdd) == true {
		return &Message{Body: "USER IS NOT VALID"}, nil
	}

	if s.isUserBanned(userToAdd) == true {
		_, pos, _ := s.grabBannedUser(userToAdd.username)
		s.banned, _ = deleteUser(s.banned, pos)
	}
	s.users = append(s.users, userToAdd)
	return &Message{Body: "SUCCESS: Added UserToChannel " + userToAdd.username}, nil
}

func (s *channel) removeUsersFromChannel(userAdding *user, userToAdd *user) (*Message, error) {
	//checks if the user thats adding someone else is in the channel
	if s.isUserInChannel(userAdding) == false {
		return &Message{Body: "USER IS NOT VALID"}, nil
	}

	//checks if the user thats adding someone else is in the channel
	if s.isUserAdmin(userAdding) == false {
		return &Message{Body: "USER IS NOT VALID"}, nil
	}
	//checks if the user that will be added to the channel isn't already in the channel
	if s.isUserInChannel(userToAdd) != true {
		return &Message{Body: "USER IS NOT VALID"}, nil
	}
	_,pos,_ := s.grabUser(userToAdd.username)
	s.users, _ = deleteUser(s.users, pos)
	return &Message{Body: "SUCCESS: Added UserToChannel " + userToAdd.username}, nil
}

func (s *channel) addUsersToBanned(userAdding *user, userToAdd *user) (*Message, error) {
	//checks if the user thats adding someone else is in the channel
	if s.isUserInChannel(userAdding) == false {
		return &Message{Body: "USER IS NOT VALID"}, nil
	}

	//checks if the user thats adding someone else is in the channel
	if s.isUserAdmin(userAdding) == false {
		return &Message{Body: "USER IS NOT VALID"}, nil
	}
	//checks if the user that will be added to the channel isn't already in the channel
	if s.isUserInChannel(userToAdd) == true {
		_,pos,_ := s.grabUser(userToAdd.username)
		s.users, _ = deleteUser(s.users, pos)
	}

	if s.isUserBanned(userToAdd) == true {
		return &Message{Body: "USER IS NOT VALID"}, nil
	}
	s.banned = append(s.banned, userToAdd)
	return &Message{Body: "SUCCESS: Banned " + userToAdd.username}, nil
}

func (s *channel) removeUsersFromBanned(userAdding *user, userToAdd *user) (*Message, error) {
	//checks if the user thats adding someone else is in the channel
	if s.isUserInChannel(userAdding) == false {
		return &Message{Body: "USER IS NOT VALID"}, nil
	}

	//checks if the user thats adding someone else is in the channel
	if s.isUserAdmin(userAdding) == false {
		return &Message{Body: "USER IS NOT VALID"}, nil
	}

	if s.isUserBanned(userToAdd) == true {
		_,pos,_:=s.grabBannedUser(userToAdd.username)
		s.banned, _ = deleteUser(s.banned, pos)
		return &Message{Body: "SUCCESS: UnBanned " + userToAdd.username}, nil
	}
	return &Message{Body: "USER IS NOT VALID"}, nil
}

func (s *channel) joinChannel(userAdding *user) (*Message, error) {
	//checks if the user thats adding someone else is in the channel
	if s.isUserInChannel(userAdding) == true {
		return &Message{Body: "USER IS NOT VALID"}, nil
	}
	if s.isUserBanned(userAdding) == true {
		return &Message{Body: "USER IS NOT VALID"}, nil
	}
	s.users = append(s.users, userAdding)
	return &Message{Body: "SUCCESS: Added UserToChannel " + userAdding.username}, nil
}

func (s *channel) leaveChannel(userAdding *user) (*Message, error) {
	//checks if the user thats adding someone else is in the channel
	if s.isUserInChannel(userAdding) == false {
		return &Message{Body: "USER IS NOT VALID"}, nil
	}
	_,pos,_ := s.grabUser(userAdding.username)
	s.users, _ = deleteUser(s.users, pos)
	return &Message{Body: "SUCCESS: Added UserToChannel " + userAdding.username}, nil
}

func (s *channel) addMessageToChannel(user1 *user, m string) (*Message, error) {
	if s.isUserInChannel(user1) == false {
		return &Message{Body: "USER IS NOT VALID"}, nil
	}
	s.cm = append(s.cm, &channelMessage{username: user1, msg: m})
	return &Message{Body: s.toString(user1.username)}, nil
}

func (s *channel) toString(user string) string {
	user1, _, _ := s.grabUser(user)
	if user1 == nil {
		return "You Aren't in this channel"
	}
	if s != nil {
		str := "Users: \n"
		for i := 0; i < len(s.users); i++ {
			str += s.users[i].toString() + "\n"
		}

		str += "Admins: \n"
		for i := 0; i < len(s.admins); i++ {
			str += s.admins[i].toString() + "\n"
		}
		if s.isUserAdmin(user1) {
			str += "\nBanned: \n"
			for i := 0; i < len(s.banned); i++ {
				str += s.banned[i].toString() + "\n"
			}
		}

		str += "\nChat: \n"
		for i := 0; i < len(s.cm); i++ {
			str += s.cm[i].toString() + "\n"
		}
		return str
	}
	return "INVALID CHANNEL"
}


