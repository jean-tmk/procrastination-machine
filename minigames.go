package procrastination

import (
	"errors"
	"math"
	"sort"
	"strings"
	"time"
)

// GameKind identifies one of the independent interactive detours presented by
// the browser. It is deliberately separate from Quest so the Go engine can
// validate real play rather than treating every click as completion.
type GameKind string

const (
	GamePlanes      GameKind = "planes"
	GameClock       GameKind = "clock"
	GameSnail       GameKind = "snail"
	GamePairs       GameKind = "pairs"
	GameThread      GameKind = "thread"
	GameRhythm      GameKind = "rhythm"
	GameCabinet     GameKind = "cabinet"
	GameTypewriter  GameKind = "typewriter"
	GameSpotlight   GameKind = "spotlight"
	GameBalance     GameKind = "balance"
	GameSwitchboard GameKind = "switchboard"
)

// GameDefinition is the engine-owned contract used by the interface, archive,
// and metrics layers. Target is the number of successful actions required.
type GameDefinition struct {
	Kind        GameKind
	Title       string
	Instruction string
	Target      int
	Difficulty  int
}

// GameCatalog returns a stable copy so callers cannot mutate global state.
func GameCatalog() []GameDefinition {
	return []GameDefinition{
		{GamePlanes, "Catch the runaway memo flock", "Catch six falling memos in the moving inbox.", 6, 2},
		{GameClock, "Stop the clock at precisely later", "Freeze three moments inside the timing window.", 3, 3},
		{GameSnail, "Guide a snail through the inbox maze", "Explore a shifting fifteen by fifteen maze and reach the red thread.", 1, 4},
		{GamePairs, "Match the secret office friendships", "Discover all four illustrated pairs.", 4, 2},
		{GameThread, "Untangle the red-thread conspiracy", "Follow the object clue instead of spatial order.", 6, 3},
		{GameRhythm, "Play the desk-lamp rhythm", "Repeat the complete six-beat office pattern.", 6, 3},
		{GameCabinet, "Sort the impossible filing cabinet", "File six records using their strange clues.", 6, 2},
		{GameTypewriter, "Repair the forgetful typewriter", "Restore the missing letter in three words.", 3, 2},
		{GameSpotlight, "Search the desk after midnight", "Reveal four objects with the wandering lamp.", 4, 2},
		{GameBalance, "Balance unfinished business", "Build a stable seven-page stack.", 7, 4},
		{GameSwitchboard, "Reroute the office daydream", "Align all nine brass connections.", 9, 3},
	}
}

// GameAction is a transport-neutral command suitable for JSON or WASM calls.
type GameAction struct {
	Name      string
	Value     int
	Secondary int
	Text      string
	At        time.Time
}

// GameResult describes the state after one action.
type GameResult struct {
	Accepted bool
	Complete bool
	Progress int
	Attempts int
	Message  string
}

// GameSession owns deterministic state for one miniature game.
type GameSession struct {
	Kind      GameKind
	StartedAt time.Time
	Progress  int
	Attempts  int
	Complete  bool
	Values    []int
	Flags     map[int]bool
	Text      []string
}

// NewGameSession validates a kind before allocating game-specific state.
func NewGameSession(kind GameKind, now time.Time) (*GameSession, error) {
	if now.IsZero() {
		now = time.Now()
	}
	var definition *GameDefinition
	for _, candidate := range GameCatalog() {
		if candidate.Kind == kind {
			copy := candidate
			definition = &copy
			break
		}
	}
	if definition == nil {
		return nil, errors.New("unknown miniature game")
	}
	session := &GameSession{Kind: kind, StartedAt: now, Flags: make(map[int]bool)}
	switch kind {
	case GameSnail:
		session.Values = []int{10}
	case GameThread:
		session.Values = []int{1, 4, 0, 2, 5, 3}
	case GameRhythm:
		session.Values = []int{0, 2, 3, 1, 2, 0}
	case GamePairs:
		session.Values = []int{0, 1, 2, 3, 0, 1, 2, 3}
	case GameTypewriter:
		session.Text = []string{"LATER", "MAYBE", "PAUSE"}
	case GameSwitchboard:
		session.Values = []int{1, 2, 3, 1, 2, 3, 1, 2, 3}
	}
	return session, nil
}

// Apply routes an action to the corresponding ruleset.
func (s *GameSession) Apply(action GameAction) GameResult {
	if s == nil {
		return GameResult{Message: "missing game session"}
	}
	if s.Complete {
		return s.result(false, "this detour is already complete")
	}
	s.Attempts++
	var accepted bool
	var message string
	switch s.Kind {
	case GamePlanes:
		accepted, message = s.applyPlane(action)
	case GameClock:
		accepted, message = s.applyClock(action)
	case GameSnail:
		accepted, message = s.applySnail(action)
	case GamePairs:
		accepted, message = s.applyPair(action)
	case GameThread:
		accepted, message = s.applyThread(action)
	case GameRhythm:
		accepted, message = s.applyRhythm(action)
	case GameCabinet:
		accepted, message = s.applyCabinet(action)
	case GameTypewriter:
		accepted, message = s.applyTypewriter(action)
	case GameSpotlight:
		accepted, message = s.applySpotlight(action)
	case GameBalance:
		accepted, message = s.applyBalance(action)
	case GameSwitchboard:
		accepted, message = s.applySwitchboard(action)
	default:
		message = "unknown ruleset"
	}
	target := targetFor(s.Kind)
	s.Complete = target > 0 && s.Progress >= target
	return s.result(accepted, message)
}

func (s *GameSession) result(accepted bool, message string) GameResult {
	return GameResult{
		Accepted: accepted,
		Complete: s.Complete,
		Progress: s.Progress,
		Attempts: s.Attempts,
		Message:  message,
	}
}

func (s *GameSession) applyPlane(action GameAction) (bool, string) {
	if action.Name != "catch" {
		return false, "move the inbox beneath a falling memo"
	}
	if action.Value < 0 || action.Value > 100 || action.Secondary < 0 || action.Secondary > 100 {
		return false, "the memo missed the inbox"
	}
	if absInt(action.Value-action.Secondary) > 13 {
		return false, "the memo drifted past"
	}
	s.Progress++
	return true, "memo safely delayed"
}

func (s *GameSession) applyClock(action GameAction) (bool, string) {
	if action.Name != "freeze" {
		return false, "the clock is still moving"
	}
	angle := action.Value % 360
	if angle < 0 {
		angle += 360
	}
	distance := int(math.Min(float64(angle), float64(360-angle)))
	if distance > 38 {
		return false, "that moment was much too useful"
	}
	s.Progress++
	return true, "perfectly later"
}

func (s *GameSession) applySnail(action GameAction) (bool, string) {
	if action.Name != "move" || len(s.Values) == 0 {
		return false, "choose a controller direction"
	}
	delta := action.Value
	if delta != -9 && delta != 9 && delta != -1 && delta != 1 {
		return false, "invalid controller direction"
	}
	position := s.Values[0]
	next := position + delta
	if next < 0 || next >= 81 {
		return false, "the snail found the edge of the desk"
	}
	if (delta == -1 || delta == 1) && position/9 != next/9 {
		return false, "the snail cannot wrap around the maze"
	}
	if snailWall(next) {
		return false, "a filing-cabinet wall blocks the route"
	}
	s.Values[0] = next
	if next == 70 {
		s.Progress = 1
		return true, "the snail reached the thread"
	}
	return true, "the snail continues"
}

func snailWall(position int) bool {
	layout := []string{
		"#########",
		"#S#.....#",
		"#.#.###.#",
		"#...#...#",
		"###.#.###",
		"#...#...#",
		"#.#####.#",
		"#......G#",
		"#########",
	}
	row, column := position/9, position%9
	return row < 0 || row >= len(layout) || layout[row][column] == '#'
}

func (s *GameSession) applyPair(action GameAction) (bool, string) {
	if action.Name != "reveal" || action.Value < 0 || action.Value >= 8 {
		return false, "choose a face-down card"
	}
	if s.Flags[action.Value] {
		return false, "that friendship is already known"
	}
	if len(s.Text) == 0 {
		s.Text = []string{string(rune('0' + action.Value))}
		return true, "remember this object"
	}
	first := int(s.Text[0][0] - '0')
	s.Text = nil
	if first == action.Value {
		return false, "a card cannot befriend itself"
	}
	if s.Values[first] != s.Values[action.Value] {
		return false, "those objects remain polite acquaintances"
	}
	s.Flags[first] = true
	s.Flags[action.Value] = true
	s.Progress++
	return true, "a secret friendship discovered"
}

func (s *GameSession) applyThread(action GameAction) (bool, string) {
	if action.Name != "pin" || s.Progress >= len(s.Values) {
		return false, "choose the next object from the clue"
	}
	if action.Value != s.Values[s.Progress] {
		return false, "the conspiracy refuses that route"
	}
	s.Progress++
	return true, "the thread found its next pin"
}

func (s *GameSession) applyRhythm(action GameAction) (bool, string) {
	if action.Name != "beat" || s.Progress >= len(s.Values) {
		return false, "listen to the desk orchestra"
	}
	if action.Value != s.Values[s.Progress] {
		s.Progress = 0
		return false, "the office rhythm wandered away"
	}
	s.Progress++
	return true, "beat remembered"
}

func (s *GameSession) applyCabinet(action GameAction) (bool, string) {
	if action.Name != "file" {
		return false, "select a record and a drawer"
	}
	expected := action.Text
	if expected == "" {
		expected = []string{"A", "B", "C"}[action.Value%3]
	}
	if string(rune('A'+action.Secondary)) != expected {
		return false, "the drawer politely rejects that record"
	}
	s.Progress++
	return true, "record filed somewhere improbable"
}

func (s *GameSession) applyTypewriter(action GameAction) (bool, string) {
	if action.Name != "letter" || s.Progress >= len(s.Text) {
		return false, "press one of the typewriter vowels"
	}
	expected := []string{"A", "A", "U"}[s.Progress]
	if strings.ToUpper(strings.TrimSpace(action.Text)) != expected {
		return false, "that letter belongs to another excuse"
	}
	s.Progress++
	return true, "the word can breathe again"
}

func (s *GameSession) applySpotlight(action GameAction) (bool, string) {
	if action.Name != "discover" || action.Value < 0 || action.Value >= 4 {
		return false, "sweep the lamp across the dark desk"
	}
	if s.Flags[action.Value] {
		return false, "that object is already in the light"
	}
	s.Flags[action.Value] = true
	s.Progress++
	return true, "another midnight object revealed"
}

func (s *GameSession) applyBalance(action GameAction) (bool, string) {
	if action.Name != "drop" || action.Value < 0 || action.Value > 100 {
		return false, "wait for the drifting page"
	}
	last := 50
	if len(s.Values) > 0 {
		last = s.Values[len(s.Values)-1]
	}
	if absInt(action.Value-last) > 24 {
		s.Values = nil
		s.Progress = 0
		return false, "the unfinished business toppled"
	}
	s.Values = append(s.Values, action.Value)
	s.Progress++
	return true, "paper balanced"
}

func (s *GameSession) applySwitchboard(action GameAction) (bool, string) {
	if action.Name != "rotate" || action.Value < 0 || action.Value >= len(s.Values) {
		return false, "choose a brass switch"
	}
	s.Values[action.Value] = (s.Values[action.Value] + 1) % 4
	s.Progress = 0
	for _, rotation := range s.Values {
		if rotation == 0 {
			s.Progress++
		}
	}
	if s.Progress == len(s.Values) {
		return true, "the daydream has a complete circuit"
	}
	return true, "one connection rerouted"
}

// OrderedGameKinds exposes a stable menu ordering for archives and tests.
func OrderedGameKinds() []GameKind {
	definitions := GameCatalog()
	sort.SliceStable(definitions, func(i, j int) bool {
		if definitions[i].Difficulty == definitions[j].Difficulty {
			return definitions[i].Title < definitions[j].Title
		}
		return definitions[i].Difficulty < definitions[j].Difficulty
	})
	kinds := make([]GameKind, 0, len(definitions))
	for _, definition := range definitions {
		kinds = append(kinds, definition.Kind)
	}
	return kinds
}

func targetFor(kind GameKind) int {
	for _, definition := range GameCatalog() {
		if definition.Kind == kind {
			return definition.Target
		}
	}
	return 0
}

func absInt(value int) int {
	if value < 0 {
		return -value
	}
	return value
}
