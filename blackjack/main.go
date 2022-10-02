package main

import (
    "fmt"
    "flag"
    "example.com/blackjack/deck"
)

type Game struct {
    Dealer   []deck.Card
    Players  []Player
    Deck     deck.Deck
    Round    int
}

func NewGame(players ...Player) Game {
    return Game{
        Dealer:  []deck.Card{},
        Players: players,
    }
}

func (g *Game) InitRound() {
    dk, _ := deck.New(deck.WithRandomOrder())
    g.Deck = dk
    g.Dealer = []deck.Card{}
    for i := range g.Players {
        g.Players[i].InitHands()
    }
    g.Round++
}

func (g *Game) Deal() {
    for i := 0; i < 2; i++ {
        for i := range g.Players {
            c, err := g.Deck.PullTop()
            if err != nil {
                panic("this shouldn't happen. deck shouldn't run out of cards")
            }
            g.Players[i].Hands[0].Hit(c)
        }
        c, err := g.Deck.PullTop()
        if err != nil {
            panic("this shouldn't happen. deck shouldn't run out of cards")
        }
        g.Dealer = append(g.Dealer, c)
    }
}

type PlayerAction string
const (
    Hit PlayerAction = "Hit"
    Stand            = "Stand"
    DoubleDown       = "DoubleDown"
    Split            = "Split"
)

func (g *Game) ProcessAction(a PlayerAction, playerid int) string {
    switch a {
    case Hit:
        return g.Hit(playerid)
    case Stand:
        return g.Stand(playerid)
    case DoubleDown:
        return g.DoubleDown(playerid)
    case Split:
        return g.Split(playerid)
    }
    return ""
}


func (g *Game) Hit(playerid int) string {
    player := &g.Players[playerid]
    g.hit(player)
    if hand.Score() > 21 {
        player.Currhand++
        return "BUST!"
    }
    return ""
}

func (g *Game) Stand(playerid int) string {
    player := &g.Players[playerid]
    res := fmt.Sprintf("Stand w/ score: %d", player[player.CurrHand].Score())
    player.CurrHand++
}

func (g *Game) DoubleDown(playerid int) string {
    player := &g.Players[playerid]
    g.hit(player)
    res := fmt.Sprintf("Double Down w/ score: %d", player[player.CurrHand].Score())
    player.CurrHand++
}

func (g *Game) Split(playerid int) string {
    player := &g.Players[playerid]
    hands  := player.Hands

    updated := make([]Hand, 0)
    for i, h := range *hands {
        if i == player.CurrHand {
            c0 := []deck.Card{h.Cards[0]}
            c1 := []deck.Card{h.Cards[1]}
            updated = append(updated, Hand{c0, 0})
            updated = append(updated, Hand{c1, 0})
        } else {
            updated = append(updated, h)
        }
    }
    *hands = updated
}

func (g *Game) hit(p Player) error {
    hand := &p.Hands[p.Currhand]
    c, err := g.Deck.PullTop()
    if err != nil {
        return err
    }
    hand.Hit(c)
    return nil
}

type PlayerType uint8
const (
    User PlayerType = iota
    AI
)

type Player struct {
    Type     PlayerType
    Hands    []Hand //multiple hands for splitting
    CurrHand int
    Chips    int
}

func (p *Player) TurnFinished() bool {
    return p.CurrHand == len(p.Hands)
}

func (p *Player) InitHands() {
    p.Hands = []Hand{
        Hand{[]deck.Card{}, 0},
    }
}

func (p *Player) SetWager(w int) {
    p.Hands[p.CurrHand].Wager = w
}

func (p *Player) GetCurrHand() Hand {
    return p.Hands[p.CurrHand]
}

type Hand struct {
    Cards  []deck.Card
    Wager  int
}

func (h *Hand) Hit(c deck.Card) {
    h.Cards = append(h.Cards, c)
}

func (h Hand) Score() int {
    total := 0
    aces  := false
    for _, c := range h.Cards {
        if c.Rank() == deck.Ace {
            aces = true
        }
        total += scoreCard(c)
    }
    if total <= 11 && aces {
        total += 10
    }
    return total
}

func scoreCard(c deck.Card) int {
    if c.Rank() == deck.Ace {
        return 1
    }
    if c.Rank() > deck.Nine {
        return 10
    }
    return int(c.Rank())
}


func main() {
    chips := flag.Int(
        "chips", 
        1000, 
        "The number of chips to start the game with. Enter a non-negative integer",
    )
    n := flag.Int(
        "players", 
        3, 
        "The number of players at the table.",
    )
    flag.Parse()

    user := NewPlayer(User, *chips)
    ai   := NewAi(*n-1)
    g    := NewGame(append(ai, user)...)
    
    for len(g.Players) > 0 {
        g.InitRound()
        for i, p := range g.Players {
            var w int
            if p.Type == User {
                w = getUserWager()
            } else {
                w = 10
            }
            g.Player[i].SetWager(w)
        }
        g.Deal()
        fmt.Println(g)

        for i := range g.Players {
            player := &g.Players[i]

            for !player.TurnFinished() {
                if player.TurnFinished() {
                    break
                }

                if len(player.GetCurrHand().Cards) == 1 {
                    // previously split hand
                    g.hit(player)
                    if err != nil {
                        panic(err)
                    }
                }
                // Display game state
                result := ""
                if isBlackJack(player.GetCurrHand().Cards) {
                    player.CurrHand++
                    result = "BLACKJACK!"
                } else {
                    // Hit, Stand, Double Down, Split
                    action := getPlayerAction(player, g.Dealer)
                    result = g.ProcessAction(action,i)
                }
                fmt.Println(result)
            }
        }
    }
}

func getPlayerAction(p Player, d []deck.Card) {
    var a PlayerAction
    switch player.Type {
    case User:
        a = getUserAction(p, d)
    case AI:
        //TODO 
        //  getAIAction(p, d)
        a = getUserAction(p, d)
    }
    return a
}

func getUserAction(p Player, d []deck.Card) {
}

func NewAi(n int) []Player {
    ai := make([]Player, n)
    for i := 0; i < n; i++ {
        ai[i] = NewPlayer(AI, 1000)
    }
    return ai
}

func NewPlayer(t PlayerType, chips int) Player {
    return Player{
        Chips: chips,
        Type:  t,
    }
}

/*

Enter starting amount: non-negative integer

chips := getStartingChips()

dealer := Hand{}
u := NewPlayer(User, chips)

players := make([]Player, 3)
for i := 1; i < 3; i++ {
    players[i] = NewPlayer(User, chips)
}


1) New Round

w := getWager(chips)

dealer := []Card
player := NewPlayer(chips)
player := Player{
    []Hand{make(Hand, 0)}, 
    chips,
}

deal(player, dealer)

1-a) Player's Hand:

    for {

    }

1-b) Dealer's Hand:
    
    - reveal hole card

    for score(dealer) < 17 {

    }





*/
