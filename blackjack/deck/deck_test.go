package deck

import (
    "testing"
    "fmt"
)

func ExampleCards() {
    ex := []struct{
        Rank
        Suit
    }{
        {Ace, Spades},
        {Three, Hearts},
        {Jack, Clubs},
        {16, Clubs},
        {Ten, Diamonds},
        {JokerRank, JokerSuit},
        {Two, JokerSuit},
        {King, Hearts},
    }
    for _, c := range ex {
        card, err := NewCard(c.Rank, c.Suit)
        if err != nil {
            fmt.Printf("error creating card with rank: %s, suit: %s\n", c.Rank, c.Suit)
            continue
        }
        fmt.Println(card)
    }
    fmt.Println(NewJoker())
    
    //Output:
    //Ace of Spades
    //Three of Hearts
    //Jack of Clubs
    //error creating card with rank: Rank(16), suit: Clubs
    //Ten of Diamonds
    //Joker
    //error creating card with rank: Two, suit: JokerSuit
    //King of Hearts
    //Joker
}

func TestDeckNoOptions(t *testing.T) {
    // create deck
    // it should return nil err 
    // it should have 52 cards
    // all combinations of suits and ranks should be included
    d, err := New()
    if err != nil {
        t.Errorf("non nil error creating new deck")
    }
    t.Run("deck length", func(t *testing.T) {
        ans := len(d)
        want := 52
        if ans != want {
            t.Errorf("got %d, want %d", ans, want)
        }
    })
    // for each suit + rank cobination
    for s := Spades; s <= Hearts; s++ {
        for r := Two; r <= Ace; r++ {
            // check if deck contains card
            t.Run(fmt.Sprintf("has %s of %s", r, s), func(t *testing.T) {
                for _, c := range d {
                    if c.Suit() == s && c.Rank() == r {
                        return
                    }
                }
                t.Errorf("%s of %s not found", r, s)

            })
        }
    }
}


func TestSort(t *testing.T) {
    d, _ := New()
    type lessfn func(x, y Card) bool
    tests := map[string]lessfn{
        "sort by rank":          func(x, y Card) bool {return x.Rank() < y.Rank()},
        "sort by rank (desc)":   func(x, y Card) bool {return x.Rank() > y.Rank()},
        "sort by suit":          func(x, y Card) bool {return x.Suit() < y.Suit()},
        "sort by suit (desc)":   func(x, y Card) bool {return x.Suit() > y.Suit()},
        "sort by rank and suit": func(x, y Card) bool {
            if x.Suit() != y.Suit() {
                return x.Suit() < y.Suit() 
            }
            return x.Rank() < y.Rank()
        },
    }

    for name, lf := range tests {
        t.Run(name, func(t *testing.T) {
            d.Sort(lf)
            for i := 0; i < len(d)-1; i++ {
                if lf(d[i],d[i+1]) == lf(d[i+1],d[i]) {
                    // d[i] and d[i+1] have equal order
                    continue
                }
                if !lf(d[i], d[i+1]) {
                    t.Fatalf("less function %v failed to sort deck", name)
                    return
                }
            }
            return
        })
    }
}

func TestPullTop(t *testing.T) {
    // it should return the first card in the deck
    // it should remove the card from the deck
    t.Run("Pull top card", func(t *testing.T) {
        d, _ := New()
        l := len(d)
        want := d[0]
        got, err := d.PullTop()
        if err != nil {
            t.Errorf("pulling top card from full deck returned non nil error")
        }
        if want != got {
            t.Errorf("Didn't pull top card. want %v, got %v", want, got)
        }
        if len(d) != l-1 {
            t.Errorf("got %d deck length after pulling top card. want %d", len(d), l-1)
        }
    })
    // it should return error if deck is empty
    t.Run("Pull top card empty deck", func(t *testing.T) {
        d := Deck{}
        _, err := d.PullTop()
        if err != nil {
            return
        }
        t.Errorf("Pulling from empty deck returned nil error")
    })
}

/*
TODO
func TestShuffle(t *testing.T) {
}

func TestPullCard(t *testing.T) {

}

func TestPeek(t *testing.T) {

}

func TestPutTop(t *testing.T) {

}

func TestPutBottom(t *testing.T) {

}

func TestPut(t *testing.T) {

}

func TestJoin(t *testing.T) {

}
*/
