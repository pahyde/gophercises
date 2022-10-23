package main

import (
    "fmt"
    "flag"
    "strings"
    "example.com/blackjack/deck"
)

type Game struct {
    Dealer      Dealer
    Players     []Player
    Deck        deck.Deck
    Round       int
    CurrPlayer  int
}

type Dealer struct {
    Cards  []deck.Card
    Chips  int
}

type Player struct {
    Type           PlayerType
    Hands          []Hand //multiple hands for splitting
    CurrHand       int
    Chips          int
}

type PlayerType uint8
const (
    User PlayerType = iota
    AI
)

type PlayerAction string
const (
    Hit PlayerAction = "HIT"
    Stand            = "STAND"
    DoubleDown       = "DOUBLEDOWN"
    Split            = "SPLIT"
)

type Hand struct {
    Cards    []deck.Card
    IsStand  bool
    Wager    int
}

func NewGame(players ...Player) *Game {
    return &Game{
        Dealer:   Dealer{},
        Players:  players,
        Round:    1,
    }
}

func (g *Game) InitRound() {
    dk, _ := deck.New(deck.WithRandomOrder())
    g.Deck = dk
    g.Dealer.Cards = []deck.Card{}
    for i := range g.Players {
        g.Players[i].InitHands()
    }
    g.Round++
    g.CurrPlayer = 0;
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
        g.Dealer.Cards = append(g.Dealer.Cards, c)
    }
}

func (g *Game) GetCurrPlayer() *Player {
    if g.CurrPlayer >= len(g.Players) {
        return nil
    }
    return &g.Players[g.CurrPlayer]
}

func (g *Game) NextPlayer() {
    g.CurrPlayer++
}

func (g *Game) Hit(p *Player) {
    hand := p.GetCurrHand()
    c, err := g.Deck.PullTop()
    if err != nil {
        panic(err)
    }
    hand.Hit(c)
}

func (g *Game) Stand(p *Player) {
    p.GetCurrHand().IsStand = true
}

func (g *Game) DoubleDown(p *Player) {
    g.Hit(p)
    hand := p.GetCurrHand()
    w := hand.Wager
    p.Chips    -= w
    hand.Wager += w
    g.Stand(p)
}

func (g *Game) Split(p *Player) {
    hands  := &p.Hands
    updated := make([]Hand, 0)
    for i, h := range *hands {
        if i == p.CurrHand {
            c0 := []deck.Card{h.Cards[0]}
            c1 := []deck.Card{h.Cards[1]}
            updated = append(updated, Hand{c0, false, h.Wager})
            updated = append(updated, Hand{c1, false, h.Wager})
        } else {
            updated = append(updated, h)
        }
    }
    *hands = updated
}

func (p *Player) InitHands() {
    p.CurrHand = 0
    p.Hands = []Hand{
        Hand{[]deck.Card{}, false, 0},
    }
}

func (p *Player) SetWager(w int) {
    p.Hands[p.CurrHand].Wager = w
}

func (p *Player) GetCurrHand() *Hand {
    return &p.Hands[p.CurrHand]
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
    game := NewGame(append(ai, user)...)
    
    // Game Loop
    for len(game.Players) > 0 {
        game.InitRound()

        // Wagers
        for i, p := range game.Players {
            var w int
            if p.Type == User {
                w = getUserWager(p.Chips)
            } else {
                w = 10
            }
            game.Players[i].SetWager(w)
        }

        game.Deal()

        // Game Round Loop
        for {
            p := game.GetCurrPlayer()
            if p == nil {
                break
            }

            // Deal new card if split hand
            if len(p.GetCurrHand().Cards) == 1 {
                game.Hit(p)
            }

            // Display game state
            displayGame(game)
            // Display player state
            displayPlayer(p, game.CurrPlayer)

            // Handle blackjacks
            if isBlackJack(p.GetCurrHand()) {
                handleBlackJack(p, game)
                continue
            }

            // Hit, Stand, Double Down, Split
            a := getPlayerAction(p, game)
            switch a {
            case Hit:
                handleHit(p, game)
            case Stand:
                handleStand(p, game)
            case DoubleDown:
                handleDoubleDown(p, game)
            case Split:
                handleSplit(p, game)
            }
        }
    }
}

func printC(s string, color string) {
    fmt.Printf(color, s)
}
const (
    green = "\033[38;2;20;190;50m%s\033[0m"
    white = "\033[38;2;255;255;255m%s\033[0m"
)
func displayGame(game *Game) {
    fmt.Println("Dealer:\n")
    fmt.Printf("_ %s\n\n", game.Dealer.Cards[1])
    fmt.Println("Players:\n")
    for i, p := range game.Players {
        var color string
        if i == game.CurrPlayer {
            color = green
        } else {
            color = white
        }
        if i == game.CurrPlayer {
            printC("> ", color)
        }
        for _, hand := range p.Hands {
            for _, card := range hand.Cards {
                printC(fmt.Sprintf("%s  ", card), color)
            }
            printC(fmt.Sprintf("  wager: %d", hand.Wager), color)
        }
        if i == game.CurrPlayer {
            printC(" <", color)
        }
        fmt.Println()
    }
    fmt.Println()
}

func displayPlayer(p *Player, pIdx int) {
    fmt.Printf("Player %d Turn\n\nYour hand:\n\n", pIdx+1)
    for _, h := range p.Hands {
        for _, c := range h.Cards {
            printC(fmt.Sprintf("%s ", c), green)
        }
        fmt.Println()
    }
}

func isBlackJack(h *Hand) bool {
    if len(h.Cards) != 2 {
        return false
    }
    c1, c2 := h.Cards[0], h.Cards[1]
    hasJack := c1.Rank() == deck.Jack || c2.Rank() == deck.Jack
    return h.Score() == 21 && hasJack
}

func handleBlackJack(p *Player, game *Game) {
    fmt.Println("\nBLACKJACK!\n")
    nextHand(p, game)
}

func handleHit(p *Player, game *Game) {
    game.Hit(p)
    if s := p.GetCurrHand().Score(); s > 21 {
        fmt.Printf("Hit with score: %d\n", s)
        fmt.Println("Bust!")
        nextHand(p, game)
    }
}

func handleStand(p *Player, game *Game) {
    game.Stand(p)
    nextHand(p, game)
}

func handleDoubleDown(p *Player, game *Game) {
    w := p.GetCurrHand().Wager
    fmt.Printf("Double Down wager: %d -> wager: %d", w, 2*w)
    game.DoubleDown(p)
    s := p.GetCurrHand().Score()
    fmt.Printf("Hit with score: %d\n", s)
    if s > 21 {
        fmt.Println("Bust!")
    }
    nextHand(p, game)
}

func handleSplit(p *Player, game *Game) {
    game.Split(p)
}

func nextHand(p *Player, game *Game) {
    p.CurrHand++
    if p.CurrHand == len(p.Hands) {
        game.NextPlayer()
    }
}
            /*

                Your Hand:

                > 2H,AS,6C     4D,JH
                  w: 20

                  2H,AS,6C   > 4D,JH
                               w: 20

----------------------------------------------------------------------------------

                Your Hand:

                > 2H,AS,6C     4D
                  w: 20        w: 20

                You have 964 chips available. Current wager: 20.

                Type one of the following: Hit, Stand, Double Down, Split (case insensitive)
                Enter your action: 

                Hit:
                    1) score <= 21 -> next action
                    2) score > 21  -> hit with score(2H, KS, 10D) = 22, -> bust
                Stand:
                    -> next hand
                Double Down:
                    -> Double down w: 20 -> w: 40 
                    -> hit with score(2H, ...) == 21
                    1) score <= 21 -> next hand
                    2) score > 21  -> bust
                Split:
                    -> next hand

                */
    

func getPlayerAction(p *Player, g *Game) PlayerAction {
    var a PlayerAction
    switch p.Type {
    case User:
        a = getUserAction(p, g)
    case AI:
        //TODO 
        //  getAIAction(p, d)
        a = getUserAction(p, g)
    }
    return a
}

func getUserAction(p *Player, game *Game) PlayerAction {
    fmt.Println()
    fmt.Printf("You have %d chips available.\n\n", p.Chips)
    fmt.Println("Type one of the following: 'Hit', 'Stand', 'Double Down', 'Split' (case insensitive)")
    var a PlayerAction
    for {
        fmt.Printf("Enter player action: ")
        if _, err := fmt.Scanf("%s", &a); err != nil {
            fmt.Println(err)
            continue
        }
        a = PlayerAction( strings.ToUpper(string(a)) )
        if a == Hit || a == Stand || a == DoubleDown || a == Split {
            break
        }
    }
    fmt.Println()
    return a
}

func getUserWager(chips int) int {
    fmt.Printf("You have %d chips available.\n\n", chips)
    fmt.Printf("Enter bet amount: ")
    var w int
    fmt.Scanf("%d", &w)
    return w
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
