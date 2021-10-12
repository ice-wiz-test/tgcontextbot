package chatStorer

var Chats = ChatStorer{}

func BuildStorer() {
	Chats.presentIn = make(map[int64] bool)
}

func IsInChatDir(id int64) bool{
	return IsPresentIn(id, &Chats)
}


func AddToChatDir(id int64) {
	AddToChat(id, &Chats)
}
type ChatStorer struct {
	presentIn map[int64] bool
}

func(s ChatStorer)  isIn(id int64) bool {
	return s.presentIn[id]
}

func IsPresentIn(id int64, s *ChatStorer) bool {
	return s.isIn(id)
}

func AddToChat(id int64, s *ChatStorer) {
	s.presentIn[id] = true
}
