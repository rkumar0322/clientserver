package chat
import (
	"errors"
	"log"
	"strings"
	"time"
	"golang.org/x/net/context"
)
type Server struct {
	users    []*user
	channels []*channel
}

func (s *Server) toString(user1 string) string{
	str := "\nUsers:\n"
	for i := 0; i < len(s.users);i++ {
		str += s.users[i].toString() + "\n"
	}
	str += "\nExisting Channels:\n" + s.showChannels(user1)
	str += "\nChannels To Join:\n" + s.showChannelsToJoin(user1)
	return str
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
	msgArr := strings.Split(msg,",")
	if msgArr[0] == "adduser" {
		return s.addUser(msgArr[1])
	} else if msgArr[0] == "addchannel" {
		return s.addChannel(msgArr[2],msgArr[1],msgArr[3])
	} else if msgArr[0] == "addusertochannel" {
		return s.AddUsersToChannelServ(msgArr[1],msgArr[2],msgArr[3])
	} else if msgArr[0] == "joinchannel" {
		return s.joinChannelServ(msgArr[1],msgArr[2])
	} else if msgArr[0] == "sendMessage" {
		return s.addMessageServ(msgArr[1],msgArr[2],msgArr[3])
	} else if msgArr[0] == "showWorkspace" {
		return &Message{Body: s.toString(msgArr[1])},nil
	} else if msgArr[0] == "userExists" {
		if s.userExists(msgArr[1]) {
			return &Message{Body: "true"},nil
		} else {
			return &Message{Body: "false"},nil
		}
	} else {
		return &Message{Body: "User Does Not Exist"},nil
	}
}
type user struct {
	username    string
	createdtime time.Time
}

func (u *user) toString() string {
	return u.username + ": " + u.createdtime.Format(time.RFC1123)
}

type channel struct {
	name string
	public bool
	users []*user
	cm    []*channelMessage
}
type channelMessage struct {
	username *user
	msg      string
}

func (s *Server) showChannels(user1 string) string{
	user2,_ := s.grabUser(user1)
	str := "\n"
	for i:= 0; i < len(s.channels);i++ {
		if s.channels[i].isUserInChannel(user2) {
			str += s.channels[i].name + "\n"
		}
	}
	return str
}

func (s *Server) showChannelsToJoin(user1 string) string{
	user2,_ := s.grabUser(user1)
	str := "\n"
	for i:= 0; i < len(s.channels);i++ {
		if s.channels[i].isUserInChannel(user2) == false && s.channels[i].public {
			str += s.channels[i].name + "\n"
		}
	}
	return str
}

func (u *channelMessage) toString() string {
	return u.username.username + ": " + u.msg
}

func (s *Server) userExists(user1 string) bool{
	for i := 0; i < len(s.users);i++ {
		if s.users[i].username == user1 {
			return true
		}
	}
	return false
}

func (s *Server) channelExists(user1 string) bool{
	for i := 0; i < len(s.channels);i++ {
		if s.channels[i].name == user1 {
			return true
		}
	}
	return false
}

func (s *Server) grabUser(user1 string) (*user,error){
	for i := 0; i < len(s.users);i++ {
		if s.users[i].username == user1 {
			return s.users[i],nil
		}
	}
	return nil,errors.New("")
}

func (s *Server) grabChannel(user1 string) (*channel,error){
	for i := 0; i < len(s.channels);i++ {
		if s.channels[i].name == user1 {
			return s.channels[i],nil
		}
	}
	return nil,errors.New("")
}
func (s *Server) addUser(user1 string) (*Message, error) {
	if s.userExists(user1) {
		return &Message{Body: "USER IS NOT VALID"},nil
	}
	s.users = append(s.users,&user{user1,time.Now()})
	return &Message{Body: "SUCCESS: User is " + user1},nil
}

func (s *Server) addChannel(chan1 string, username string,public string) (*Message, error) {
	if s.channelExists(chan1) {
		return &Message{Body: "channel Exists"},nil
	}
	user1,err1 := s.grabUser(username)
	if err1 != nil {
		return &Message{Body: "User Does Not Exist"},nil
	}
	var public1 bool
	if public == "yes" {
		public1 = true
	} else if public == "no" {
		public1 = false
	} else {
		return &Message{Body: "Invalid Input"},nil
	}
	var users []*user
	users = append(users,user1)
	var messages []*channelMessage
	s.channels = append(s.channels,&channel{chan1,public1,users,messages})
	return &Message{Body: "SUCCESS: Creaeted Channel " + chan1},nil
}

func (s *Server) joinChannelServ(user1 string, channel2 string) (*Message, error) {
	channel1,err := s.grabChannel(channel2)
	if err != nil{
		return &Message{Body: "User IS NOT VALID"},nil
	}
	user2,err2 := s.grabUser(user1)
	if err2 != nil{
		return &Message{Body: "Channel IS NOT VALID"},nil
	}
	return channel1.joinChannel(user2)
}

func (s *Server) AddUsersToChannelServ(userAdding1 string,userToAdd1 string, channel2 string) (*Message, error) {
	channel1,err := s.grabChannel(channel2)
	if err != nil{
		return &Message{Body: "User IS NOT VALID"},nil
	}
	userToAdd,err2 := s.grabUser(userToAdd1)
	if err2 != nil{
		return &Message{Body: "User IS NOT VALID"},nil
	}

	userAdding,err3 := s.grabUser(userAdding1)
	if err3 != nil{
		return &Message{Body: "User IS NOT VALID"},nil
	}
	return channel1.addUsersToChannel(userAdding,userToAdd)
}

func (s *Server) addMessageServ(user1 string, channel2 string,m string) (*Message, error) {
	channel1,err := s.grabChannel(channel2)
	if err != nil{
		return &Message{Body: "User IS NOT VALID"},nil
	}
	user2,err2 := s.grabUser(user1)
	if err2 != nil{
		return &Message{Body: "Channel IS NOT VALID"},nil
	}
	return channel1.addMessageToChannel(user2,m)
}

func (s *channel) isUserInChannel(user1 *user) bool{
	exists1 := false
	for i := 0; i < len(s.users);i++ {
		if s.users[i].username == user1.username {
			return true
		}
	}
	return exists1
}

func (s *channel) addUsersToChannel(userAdding *user,userToAdd *user) (*Message, error){
	//checks if the user thats adding someone else is in the channel
	if s.isUserInChannel(userAdding) == false {
		return &Message{Body: "USER IS NOT VALID"},nil
	}
	//checks if the user that will be added to the channel isn't already in the channel
	if s.isUserInChannel(userToAdd) == true {
		return &Message{Body: "USER IS NOT VALID"},nil
	}
	s.users = append(s.users,userToAdd)
	return &Message{Body: "SUCCESS: Added UserToChannel " + userToAdd.username},nil
}

func (s *channel) joinChannel(userAdding *user) (*Message, error){
	//checks if the user thats adding someone else is in the channel
	if s.isUserInChannel(userAdding) == true {
		return &Message{Body: "USER IS NOT VALID"},nil
	}
	s.users = append(s.users,userAdding)
	return &Message{Body: "SUCCESS: Added UserToChannel " + userAdding.username},nil
}

func (s *channel) addMessageToChannel(user1 *user,m string) (*Message, error){
	if s.isUserInChannel(user1) == false {
		return &Message{Body: "USER IS NOT VALID"},nil
	}
	s.cm = append(s.cm,&channelMessage{username: user1,msg:m})
	return &Message{Body: s.toString()},nil
}

func (s *channel) toString() string {
	str := "Users: \n"
	for i := 0; i < len(s.users);i++ {
		str += s.users[i].toString() + "\n"
	}
	str += "\nChat: \n"
	for i := 0; i < len(s.cm);i++ {
		str += s.cm[i].toString() + "\n"
	}
	return str
}