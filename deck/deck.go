package deck

import (
    "fmt"
    "math/rand"
    "sort"
)

type Suit int
type Rank int

const (
    Spades Suit = iota 
    Diamonds
    Clubs
    Hearts
    JokerSuit
)

const (
    Two Rank = iota + 2
    Three
    Four
    Five
    Six
    Seven
    Eight
    Nine
    Ten
    Jack
    Queen
    King
    Ace
    JokerRank
)

type Card struct {
    r Rank
    s Suit
}

func NewCard(r Rank, s Suit) (Card, error) {
    if r < Two || r > JokerRank {
        return Card{}, fmt.Errorf("undefined Rank %s", r)
    }
    if s < Spades || s > JokerSuit {
        return Card{}, fmt.Errorf("undefined Suit %s", s)
    }
    if (r == JokerRank) != (s == JokerSuit) {
        return Card{}, fmt.Errorf("can't create mixed joker with rank: %s and suit: %s", r, s)
    }
    return Card{r, s}, nil
}

func NewJoker() Card {
    c, _ := NewCard(JokerRank, JokerSuit)
    return c
}

func (c Card) Rank() Rank {
    return c.r
}

func (c Card) Suit() Suit {
    return c.s
}

func (c Card) String() string {
    if c.r == JokerRank {
        return "Joker"
    }
    return fmt.Sprintf("%s of %s", c.r, c.s)
}

type Deck []Card

func (d *Deck) Shuffle() {
    rand.Shuffle(len(*d), func(i, j int) {
        (*d)[i], (*d)[j] = (*d)[j], (*d)[i]
    })
}

func (d *Deck) Sort(less func(x, y Card) bool) {
    sort.Slice(*d, func(i, j int) bool {
        return less((*d)[i], (*d)[j])
    })
}

func (d *Deck) PullTop() (Card, error) {
    card, err := d.getCardByIdx(0)
    if err != nil {
        return Card{}, fmt.Errorf("Can't Pull Top Card. %v", err)
    }
    *d = (*d)[1:]
    return card, nil
}

func (d *Deck) PullCard(i int) (Card, error) {
    card, err := d.getCardByIdx(i)
    if err != nil {
        return Card{}, fmt.Errorf("Can't Pull Card at index %d. %v", i, err)
    }
    *d = append((*d)[:i], (*d)[i+1:]...)
    return card, nil
}

func (d *Deck) Peek(i int) (Card, error) {
    card, err := d.getCardByIdx(i)
    if err != nil {
        return Card{}, fmt.Errorf("Can't Peek Card at index %d. %v", i, err)
    }
    return card, nil
}

func (d *Deck) getCardByIdx(i int) (Card, error) {
    if len(*d) == 0 {
        return Card{}, fmt.Errorf("The deck is empty.")
    } 
    if i < 0 || i >= len(*d) {
        return Card{}, fmt.Errorf("Index %d out of range for deck with %d cards", i, len(*d))
    }
    return (*d)[i], nil
}

func (d *Deck) PutTop(c Card) {
    d.Put(c, 0)
}

func (d *Deck) PutBottom(c Card) {
    d.Put(c, len(*d))
}

func (d *Deck) Put(c Card, i int) error {
    if i < 0 || i > len(*d) {
        return fmt.Errorf(
            "Can't Put Card at index %d. Out of range for deck with %d cards", 
            i, len(*d))
    } 
    newLength := len(*d) + 1
    result := make(Deck, newLength)
    for j := 0; j < newLength; j++ {
        switch {
        case j < i:
            result[j] = (*d)[j]
        case j == i:
            result[j] = c
        case j > i:
            result[j] = (*d)[j-1]
        }
    }
    *d = result
    return nil
}

func Join(decks ...Deck) Deck {
    d := Deck{}
    for _, deck := range decks {
        for _, card := range deck {
            d = append(d, card)
        }
    }
    return d
}

type DeckOption func(d *Deck) error

func WithCardsRemoved(toRemove func(c Card) bool) DeckOption {
    return func(d *Deck) error {
        keep := Deck{}
        for _, card := range *d {
            if !toRemove(card) {
                keep.PutBottom(card)
            }
        }
        *d = keep
        return nil
    }
}

func WithJokers(n int) DeckOption {
    return func(d *Deck) error {
        if n < 0 {
            return fmt.Errorf("Bad input n: %d. Can't add %d Jokers.", n, n)
        }
        for i := 0; i < n; i++ {
            d.PutBottom(Card{JokerRank, JokerSuit})
        }
        return nil
    }
}

func WithSortedOrder(less func(x, y Card) bool) DeckOption {
    return func(d *Deck) error {
        d.Sort(less)
        return nil
    }
}

func WithRandomOrder() DeckOption {
    return func(d *Deck) error {
        d.Shuffle()
        return nil
    }
}

func New(opts ...DeckOption) (Deck, error) {
    d := make(Deck, 0, 52)
    for s := Spades; s <= Hearts; s++ {
        for r := Two; r <= Ace; r++ {
            d = append(d, Card{r, s})
        }
    }
    for _, opt := range opts {
        if err := opt(&d); err != nil {
            return Deck{}, err
        }
    }
    return d, nil
}

