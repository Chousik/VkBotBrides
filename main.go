package main

import (
	"context"
	"fmt"
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/events"
	"github.com/SevereCloud/vksdk/v2/longpoll-bot"
	"log"
	"math/rand"
	"strconv"
	"strings"
)

type Player struct {
	id     int
	number int
	ans    string
}
type Game struct {
	Round      int
	Game_Stage int
	KIK        int
	Propusk    int
	WantGroom  []int
	Play       []int
	Brides     []Player
	Groom      int
	Question   string
}

func NewGame() *Game {
	g := &Game{}
	g.KIK = 1
	g.Propusk = 1
	return g
}

const (
	ChatID  = 2000000005
	ModerId = 411634368
)

func (g *Game) AddBrides(id int, number int) {
	g.Brides = append(g.Brides, Player{id, number, ""})
}
func (g *Game) DeleteBrides(id int, number int) bool {
	if id == 0 {
		for i, p := range g.Brides {
			if p.number == number {
				g.Brides = append(g.Brides[:i], g.Brides[i+1:]...)
				return true
			}
		}
	} else {
		for i, p := range g.Brides {
			if p.id == id {
				g.Brides = append(g.Brides[:i], g.Brides[i+1:]...)
				return true
			}
		}
	}
	return false
}
func (g *Game) GetIdNumberBrides(id int, number int) int {
	if id == 0 {
		for _, p := range g.Brides {
			fmt.Println(number)
			fmt.Println(p.number)
			if p.number == number {
				fmt.Println(p.id)
				return p.id
			}
		}
	} else {
		for _, p := range g.Brides {
			if p.id == id {
				return p.number
			}
		}
	}
	return 0
}
func (g *Game) AddPlayer(id int) bool {
	for _, p := range g.Play {
		if p == id {
			return false
		}
	}
	g.Play = append(g.Play, id)
	return true
}
func (g *Game) DeletePlayer(id int) bool {
	for i, p := range g.Play {
		if p == id {
			g.Play = append(g.Play[:i], g.Play[i+1:]...)
			return true
		}
	}
	return false
}
func (g *Game) AddWantGroom(id int) bool {
	for _, p := range g.WantGroom {
		if p == id {
			return false
		}
	}
	for _, p := range g.Play {
		if p == id {
			g.WantGroom = append(g.WantGroom, id)
			return true
		}
	}
	return false
}
func (g *Game) DeleteWantGroom(id int) bool {
	for i, p := range g.WantGroom {
		if p == id {
			g.WantGroom = append(g.WantGroom[:i], g.WantGroom[i+1:]...)
			return true
		}
	}
	return false
}
func (g *Game) StartGame() bool {
	if len(g.WantGroom) == 0 || len(g.Play) == 0 {
		return false
	}
	o := rand.Intn(len(g.WantGroom))
	g.Groom = g.WantGroom[o]
	g.DeletePlayer(g.Groom)
	g.DeletePlayer(g.Groom)
	s := len(g.Play)
	for i := 1; i <= s; i++ {
		l := rand.Intn(len(g.Play))
		BridesID := g.Play[l]
		g.Play = append(g.Play[:l], g.Play[l+1:]...)
		g.AddBrides(BridesID, i)
	}
	g.Round = 1
	g.Game_Stage = 1
	return true
}
func (g *Game) HePlay(id int) bool {
	for _, p := range g.Brides {
		if p.id == id {
			return true
		}
	}
	fmt.Println("One")
	return false
}
func (g *Game) AddAnswer(id int, answer string) bool {
	fmt.Println(answer)
	for i, p := range g.Brides {
		fmt.Println(p)
		if p.id == id {
			g.Brides[i].ans = answer
			fmt.Println(p.ans)
			return true
		}
	}
	return false
}
func (g *Game) ReturnAnwer() string {
	answers := "Ответы невест!:\n"
	for _, p := range g.Brides {
		answers += strconv.Itoa(p.id) + ":" + p.ans
	}
	return answers
}
func (g *Game) AddQuest(q []string) {
	g.Question = strings.Join(q, " ")
}
func (g *Game) EndQuest(vk *api.VK, yes bool) {
	qust := ""
	if yes == false {
		qust += "Время закончилось, ответы прекрасных невест:\n"
	} else {
		qust += "Все невесты ответили!\n"
	}
	for i, p := range g.Brides {
		if p.ans != "" {
			qust += strconv.Itoa(p.number) + ": " + p.ans + "\n"
			g.Brides[i].ans = ""
		} else {
			qust += "Невеста" + strconv.Itoa(p.number) + "проспала! Ею была @id" + strconv.Itoa(p.id) + "\n"
			g.Brides = append(g.Brides[:i], g.Brides[i+1:]...)
		}
		g.Game_Stage = 4
	}
	SendMessege(vk, qust, ChatID)
	SendMessege(vk, qust, g.Groom)
	SendMessege(vk, "Для того что бы кикнуть нвесту используйте !кик <номер>", g.Groom)
}
func (g *Game) KikOne(vk *api.VK, mess []string) bool {
	numb := strings.Join(mess[:1], "")
	mess = mess[1:]
	number, err := strconv.Atoi(numb)
	if err != nil {
		return false
	}
	id := g.GetIdNumberBrides(0, number)
	i := g.DeleteBrides(0, number)
	if i == false {
		return false
	}
	SendMessege(vk, "К сожаление нас покидает невеста номер "+numb+". Ею была прекрасная @id"+strconv.Itoa(id)+"\n", ChatID)
	SendMessege(vk, "Вы кикнули невесту номер "+numb+". Ею была прекрасная @id"+strconv.Itoa(id)+"\n", g.Groom)
	if len(mess) >= 1 {
		komment := strings.Join(mess, " ")
		SendMessege(vk, "Комментарий жениха: "+komment, ChatID)
		SendMessege(vk, "Комментарий жениха: "+komment, g.Groom)
	}
	return true

}
func (g *Game) ConstPlayer() int {
	ans := 0
	for i, _ := range g.Brides {
		ans += 1
		if i > 0 {
		}
	}
	return ans
}
func (g *Game) ConstPlayer2() int {
	ans := 0
	for i, _ := range g.Play {
		ans += 1
		if i > 0 {
		}
	}
	return ans
}
func (g *Game) AllEND() bool {
	for _, p := range g.Brides {
		fmt.Println(p.ans)
		if p.ans == "" {
			return false
		}
	}
	return true
}
func (g *Game) NewRound(vk *api.VK) {
	g.Round += 1
	g.Game_Stage = 2
	SendMessege(vk, "Раунд "+strconv.Itoa(g.Round)+" начинается!", ChatID)
	SendMessege(vk, "Раунд "+strconv.Itoa(g.Round)+" начинается!\n Для вопроса используйте !вопрос", g.Groom)
	for _, p := range g.Brides {
		SendMessege(vk, "Поздравляем, вы прошли в раунд "+strconv.Itoa(g.Round)+". Ожидайте вопроса жениха!", p.id)
	}

}
func (g *Game) GetQuest(vk *api.VK, obj events.MessageNewObject) {
	mes := strings.Split(obj.Message.Text, " ")
	if mes[0] == "!вопрос" {
		ans := mes[1:]
		if len(ans) > 0 {
			g.AddQuest(ans)
			SendMessege(vk, "Вопрос задание! Ожидайте ответа невест!", g.Groom)
			SendMessege(vk, "Вопрос жениха №"+strconv.Itoa(g.Round)+"\n"+g.Question, ChatID)
			for _, p := range g.Brides {
				SendMessege(vk, "Вопрос жениха №"+strconv.Itoa(g.Round)+"\n"+g.Question, p.id)
				SendMessege(vk, "Для ответа на вопрос используйте !ответ <ответ> ", p.id)

			}
			g.Game_Stage = 3
		} else {
			SendMessege(vk, "Что бы задать вопрос испольуйте \n !вопрос <вопрос>", g.Groom)
		}
	} else {
		SendMessege(vk, "Что бы задать вопрос испольуйте \n !вопрос <вопрос>", g.Groom)
	}

}
func (g *Game) ENDGAME(vk *api.VK) *Game {
	bride := g.Brides[0].id
	groom := g.Groom
	SendMessege(vk, "Поздравляем наших молодоженов! Прекрасную @id"+strconv.Itoa(bride)+" и прекрасного @id"+strconv.Itoa(groom), ChatID)
	SendMessege(vk, "Игра окончена! Что бы начать новую напишите !невесты", ChatID)
	g = NewGame()
	return g
}
func main() {
	g := NewGame()
	token := "vk1.a.bCp6Il1J3O9wePPDU9ElvErT85_z7SiQW-OPF45Ui5zdtsJq2r8HsRvywZ03F4x1RBef7yGeovc34H6iKzWv5ium29LVbJkefFtR7em7Qt0VrUbjum6PAIVoXTp4KNTf6jO-IKTGaZkVduFzcl11SZNFjiwylQrJhyOYX4aqhLf-bAswUkkXn5LBAe02kd8R" // use os.Getenv("TOKEN")
	vk := api.NewVK(token)
	GameTrue := false
	// get information about the group
	group, err := vk.GroupsGetByID(nil)
	if err != nil {
		log.Fatal(err)
	}
	// Initializing Long Poll
	lp, err := longpoll.NewLongPoll(vk, group[0].ID)
	if err != nil {
		log.Fatal(err)
	}

	// New message event
	lp.MessageNew(func(_ context.Context, obj events.MessageNewObject) {
		log.Printf("%d: %s", obj.Message.PeerID, obj.Message.Text)
		if GameTrue {
			g = NewGame()
		} else {
			g = Games(g, vk, obj)
		}
	})

	// Run Bots Long Poll
	log.Println("Start Long Poll")
	if err := lp.Run(); err != nil {
		log.Fatal(err)
	}
}
func SendMessege(vk *api.VK, text string, id int) {
	b := params.NewMessagesSendBuilder()
	b.Message(text)
	b.RandomID(0)
	b.PeerID(id)
	_, err := vk.MessagesSend(b.Params)
	if err != nil {
		log.Fatal(err)
	}
}
func Games(g *Game, vk *api.VK, obj events.MessageNewObject) *Game {
	switch g.Game_Stage {
	case 0:
		if obj.Message.PeerID == ChatID {
			{
				if obj.Message.Text == "!невесты" {
					g.Game_Stage = 1
					SendMessege(vk, "Набор на игру начат!", ChatID)
				}
			}
		}
	case 1:
		{
			switch obj.Message.PeerID {
			case ChatID:
				if obj.Message.Text == "!старт" {
					i := g.StartGame()
					if i == false {
						SendMessege(vk, "Никто не хочет быть женихом или недостаточно невест!", ChatID)
						return g
					}
					SendMessege(vk, "Игра начинается. Ожидайте вопроса жениха!", ChatID)
					for _, p := range g.Brides {
						fmt.Println(p.id)
						SendMessege(vk, "Поздравляем, вы невеста "+strconv.Itoa(p.number)+". Ожидайте вопроса жениха", p.id)
					}
					fmt.Println(g.Groom)
					SendMessege(vk, "Поздравляем, вы жених! Что бы задать вопрос используйте \n!вопрос <вопрос>", g.Groom)
					g.Game_Stage = 2
				}
			default:
				switch obj.Message.Text {
				case "!+":
					{
						i := g.AddPlayer(obj.Message.FromID)
						if i == false {
							SendMessege(vk, "Вы уже приняли игру!", obj.Message.FromID)
							return g
						}
						SendMessege(vk, "Кол-во игроков: "+strconv.Itoa(g.ConstPlayer2()), ChatID)
						SendMessege(vk, "Вы успешно приняты на игру! Ожидайте!", obj.Message.FromID)
					}
				case "!-":
					{
						i := g.DeletePlayer(obj.Message.FromID)
						if i == false {
							SendMessege(vk, "Вы и так не играете!", obj.Message.FromID)
							return g
						}
						SendMessege(vk, "Кол-во игроков:"+strconv.Itoa(g.ConstPlayer2()), ChatID)
						SendMessege(vk, "Вы успешно вышли с набора!", obj.Message.FromID)
					}
				case "!+жених":
					{
						i := g.AddWantGroom(obj.Message.FromID)
						if i == false {
							SendMessege(vk, "Вы либо не приняли игру, либо уже можете стать женихом!", obj.Message.FromID)
							return g
						}
						SendMessege(vk, "Теперь вы можете стать женихом!", obj.Message.FromID)

					}
				case "!-жених":
					{
						i := g.DeleteWantGroom(obj.Message.FromID)
						if i == false {
							SendMessege(vk, "Вы и так не можете стать женихом!", obj.Message.FromID)
							return g
						}
						SendMessege(vk, "Теперь вы не можете стать женихом!", obj.Message.FromID)
					}
				}
			}

		}
	case 2:
		{
			switch obj.Message.PeerID {
			case g.Groom:
				{
					g.GetQuest(vk, obj)
				}
			case ChatID:
				{
					return g
				}
			default:
				if g.HePlay(obj.Message.FromID) {
					SendMessege(vk, "Ожидайте вопроса жениха!", obj.Message.FromID)
				} else {
					SendMessege(vk, "Набор на игру уже закончен! Ожидайте следующей", obj.Message.FromID)
				}
			}

		}
	case 3:
		{
			switch obj.Message.PeerID {
			case ChatID:
				{
					if obj.Message.FromID == ModerId && obj.Message.Text == "!скип" {
						g.EndQuest(vk, false)
					}
				}
			case g.Groom:
				{
					SendMessege(vk, "Ожидайте ответа невест!", g.Groom)
				}

			default:
				{
					fmt.Println()
					if g.HePlay(obj.Message.PeerID) {
						mes := strings.Split(obj.Message.Text, " ")
						if mes[0] == "!ответ" {
							if len(mes[1:]) > 0 {
								answ := strings.Join(mes[1:], " ")
								i := g.AddAnswer(obj.Message.FromID, answ)
								if i == true {
									SendMessege(vk, "Ваш ответ успешно принят!", obj.Message.FromID)
									println(g.Brides)
									if g.AllEND() {
										g.EndQuest(vk, true)
									}
								} else {
									SendMessege(vk, "Для ответа используйте \n !ответ <ответ>", obj.Message.FromID)
								}
							} else {
								SendMessege(vk, "Для ответа используйте \n !ответ <ответ>", obj.Message.FromID)
							}
						} else {
							SendMessege(vk, "Для ответа используйте \n !ответ <ответ>", obj.Message.FromID)
						}
					} else {
						SendMessege(vk, "Набор на игру уже закончен! Ожидайте следующей", obj.Message.FromID)
					}
				}

			}
		}
	case 4:
		switch obj.Message.PeerID {
		case ChatID:
			{
				return g

			}
		case g.Groom:
			{
				mes := strings.Split(obj.Message.Text, " ")
				if mes[0] == "!кик" && len(mes[1:]) > 1 {
					i := g.KikOne(vk, mes[1:])
					if i == false {
						SendMessege(vk, "Неправильный номер", g.Groom)
					}
					if g.ConstPlayer() >= 2 {
						g.NewRound(vk)
					} else {
						g = g.ENDGAME(vk)
					}

				} else {
					SendMessege(vk, "Для кика используйте !кик <номер> <коммент>", g.Groom)
				}
			}
		default:
			{
				if g.HePlay(obj.Message.FromID) {
					SendMessege(vk, "Ожидайте когда жених кикнет невесту!", obj.Message.FromID)
				}
			}

		}

	}
	return g
}
