// Package game handles the game logic.
package game

import (
	"errors"
	"math/rand"
	"strings"
	"time"
)

// Team holds the team specific information for a given team
type Team struct {
	Clues         []Clue
	Code          string
	LastSeen      time.Time
	Solved        []int
	Unlocked      []int
	UnlockedCount int
	Events        []string
	Message       string
	Status        string
}

// Manager holds each of the teams
type Manager struct {
	Teams []Team
}

var clues []Clue

// Clue holds the next clue and any challenges associated with it
type Clue struct {
	Code      string
	Title     string
	Text      string
	Challenge bool
}

// GetClue returns whether or not the clue code is valid.
func (m Manager) GetClue(code string) (Clue, error) {
	for _, clue := range clues {
		if clue.Code == code {
			return clue, nil
		}
	}
	return Clue{}, errors.New("Clue does not exist")
}

// GetTeam returns whether or not the team code is valid.
func (m *Manager) GetTeam(teamCode string) (int, error) {
	teamCode = strings.ToUpper(teamCode)
	for index, team := range m.Teams {
		if team.Code == teamCode {
			return index, nil
		}
	}
	return -1, errors.New("Team could not be found")
}

// CheckIn updates the LastSeen time of a team
func (m *Manager) CheckIn(teamCode string) {
	teamCode = strings.ToUpper(teamCode)
	for _, team := range m.Teams {
		if team.Code == teamCode {
			team.LastSeen = time.Now()
		}
	}
}

// CheckIn updates the LastSeen time of a team
func (team *Team) CheckIn() {
	team.LastSeen = time.Now()
}

// Solve will check if the team can solve the game then update the game state.
func (team *Team) Solve(clueCode string) error {
	team.CheckIn()

	for i, pos := range team.Unlocked {
		if pos > len(team.Clues)-1 {
			continue
		}
		if team.Clues[pos].Code == clueCode {
			team.Solved = append(team.Solved, pos)
			team.UnlockedCount++
			team.Unlocked[i] = team.UnlockedCount - 1
			if team.Unlocked[i] >= len(team.Clues) {
				team.Unlocked = remove(team.Unlocked, i)
			}
			return nil
		}
	}
	return errors.New("this team has not unlocked this location")
}

// remove returns the []int with the item at index s removed.
func remove(slice []int, s int) []int {
	return append(slice[:s], slice[s+1:]...)
}

// FastForward will check if the team can solve the game then update the game state.
func (team *Team) FastForward() {
	team.CheckIn()

	for i, pos := range team.Unlocked {
		if team.UnlockedCount != len(team.Clues) {
			team.Solved = append(team.Solved, pos)
			team.UnlockedCount++
			team.Unlocked[i] = team.UnlockedCount - 1
		}
	}
}

// Hinder will check if the team can solve the game then update the game state.
func (team *Team) Hinder() {
	team.CheckIn()

	for i := 0; i < 3; i++ {
		team.Clues[team.Unlocked[i]], team.Clues[len(team.Clues)-1-i] = team.Clues[len(team.Clues)-1-i], team.Clues[team.Unlocked[i]]
	}
}

// The symbols team codes are created from.
// Confusing letters such as I and L, O and Q have one pair removed.
var symbols = []rune("ABCDEFGHJKLMNPRSTUVWXYZ")

func newCode() string {
	b := make([]rune, 4)
	for i := 0; i < 4; i++ {
		b[i] = symbols[rand.Intn(len(symbols))]
	}
	return string(b)
}

// CreateTeams will create however many teams are asked for.
// Num must be > 0
func (m *Manager) CreateTeams(num int) {
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < num; i++ {
		set := make([]Clue, len(clues))
		copy(set, clues)
		rand.Shuffle(len(set), func(i, j int) { set[i], set[j] = set[j], set[i] })
		m.Teams = append(m.Teams,
			Team{
				Clues:         set,
				Code:          newCode(),
				Solved:        []int{},
				Unlocked:      []int{0, 1, 2},
				UnlockedCount: 3,
			},
		)
	}
}

// find returns the index of a value in a slice of ints.
func find(slice []int, val int) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

func init() {
	rand.Seed(time.Now().UnixNano())
	// Clues written by the wonderful Tamika!
	// TODO: Load these from a file on run.
	// Clues must have the same ID and not change on reload. Posters must stay accurate.
	clues = append(clues, Clue{"QWOP4", "St Daves", "Named after a saint but if you have a class here you will know, saint-like is the least like what your grades will show", false})
	clues = append(clues, Clue{"DFG5J", "Student Health", "Pregnancy tests and STI checks please? Book an appointment here to put your mind at ease", false})
	clues = append(clues, Clue{"9LCZ4", "Campus Watch", "0800 479 5000. Save the contact.", false})
	clues = append(clues, Clue{"63FJC", "Mellor Labs", "A pretty new building with chemistry inside, find the molecules to know what it hides.", false})
	clues = append(clues, Clue{"9CV5A", "Radio One", "Their card gives you a bunch of sweet deals, tune in to 91fm to get in your feels", false})
	clues = append(clues, Clue{"WM6SQ", "Law Library", "A library for rules and cases of times passed, climb 6-9 flights of stairs if you can be assed.", false})
	clues = append(clues, Clue{"IOHRU", "Zoology Department", "Mammals, reptiles, amphibians you hear? Their kind is studied in this department I swear.", false})
	clues = append(clues, Clue{"GOG79", "Black Sale House", "Our Very own Sam Leaper belongs to a band, whose name is the same as this place I understand (note: check out his poster in locals HQ)", false})
	clues = append(clues, Clue{"NTODH", "Otago Business School", "Chads and Brads learn about stocks I wager, and something about hedges? (IDK im an english major)", false})
	clues = append(clues, Clue{"LDQJQ", "Touchstone by the Leith", "This place is likened to a moat, and reminds me of a movie quote: 'that is a nice boulder - Donkey'", true})
	clues = append(clues, Clue{"QT7KH", "Union Lawn", "A favourite place of our resident ducks, this lawn is where they give no fucks. Still stumped? Well worry no more, here you can usually find dumplings galore.", true})
	clues = append(clues, Clue{"3CZ3I", "Burns Lecture Theatres", "Named after our old pal Robbie in the octagon, here is where a love of words can be built on.", false})
	clues = append(clues, Clue{"B0Y2H", "Castle Lecture Theatres", "Though you may not find any knights here, this Castle is still home to things that should be feared…", false})
	clues = append(clues, Clue{"I0AHP", "Botany Department", "Rhymes aside, may I just ask, can these experts revive my houseplants?", false})
	clues = append(clues, Clue{"2YXNQ", "Central Library", "This library is large and perhaps the most popular, where basic bitches study, according to critic I gather.", false})
	clues = append(clues, Clue{"C8P8S", "Microbiology Building", "Study something small and living too? This building is the place for you", false})
	clues = append(clues, Clue{"E0MBX", "School of Physiotherapy", "Is it your neck, your back, or something else that cracks? Head to this school to get back on track.", false})
	clues = append(clues, Clue{"ZEEIB", "Locals HQ", "Third floor of union where you can go, for freshers who live in flats or at home.", true})
	clues = append(clues, Clue{"2SNU5", "Math / Stats Department", "I gave up on this subject long ago, but if you like numbers this place is your home.", false})
	clues = append(clues, Clue{"BS0MC", "Union Grill", "If a craving for burgers, fries, and iced coffee arrives, this grill on campus has what you need to survive.", true})
	clues = append(clues, Clue{"O0N66", "Science Library", "This library has diagrams and technical terms, so long to all the fiction from which we once learned. But as if that sting wasnt enough, the 7am-11pm hours will make you feel rough.", false})
	clues = append(clues, Clue{"FR4SQ", "OUSA Clubs and Socs", "Their dollar-deal lunches went up in price this year. Worth it? Ill let you judge here.", true})
	clues = append(clues, Clue{"H3903", "Critic Office", "Have an affinity for writing, drawing, or simply an opinion? Here, a love for journalism and magazines will make you fit in.", false})
	clues = append(clues, Clue{"4IYE5", "Owheo Building", "They put the it in IT and what else I dont really know, but the computer angels frequent here, to learn at Owheo.", false})
	clues = append(clues, Clue{"CFEK9", "Archway Lecture Theatres", "Consult your map for this one might be a little tricky, find the anagram of 'crawhAy' and you should be there in a jiffy.", false})
}
