// Package procrastination implements a deterministic strategic-delay engine.
// It turns an obligation into a dependency graph of plausible but irrelevant
// side quests while keeping enough state to measure delay, theatrical output,
// achievements, and eventual returns to the original task.
package procrastination

import (
	"encoding/json"
	"errors"
	"fmt"
	"hash/fnv"
	"math"
	"sort"
	"strings"
	"sync"
	"time"
)

const (
	DefaultQuestCount = 5
	MaximumQuestCount = 12
	FocusDuration     = 5 * time.Minute
	ArchiveVersion    = 1
)

type Category string

const (
	CategoryEmail    Category = "email"
	CategoryStudy    Category = "study"
	CategoryCleaning Category = "cleaning"
	CategoryWriting  Category = "writing"
	CategoryWork     Category = "work"
	CategoryGeneral  Category = "general"
)

type Urgency uint8

const (
	UrgencyDecorative Urgency = iota
	UrgencyEventually
	UrgencySoon
	UrgencyAllegedlyNow
)

func (u Urgency) String() string {
	switch u {
	case UrgencyDecorative:
		return "decorative"
	case UrgencyEventually:
		return "eventually"
	case UrgencySoon:
		return "soon"
	case UrgencyAllegedlyNow:
		return "allegedly now"
	default:
		return "unfiled"
	}
}

type QuestKind string

const (
	QuestResearch       QuestKind = "research"
	QuestOrganization   QuestKind = "organization"
	QuestAesthetic      QuestKind = "aesthetic"
	QuestPreparation    QuestKind = "preparation"
	QuestMetaWork       QuestKind = "meta-work"
	QuestIntermission   QuestKind = "intermission"
	QuestDocumentation  QuestKind = "documentation"
	QuestClassification QuestKind = "classification"
)

type Obligation struct {
	ID          string    `json:"id"`
	Text        string    `json:"text"`
	Category    Category  `json:"category"`
	Urgency     Urgency   `json:"urgency"`
	SubmittedAt time.Time `json:"submitted_at"`
	CompletedAt time.Time `json:"completed_at,omitempty"`
	AbandonedAt time.Time `json:"abandoned_at,omitempty"`
}

func (o Obligation) Active() bool {
	return o.CompletedAt.IsZero() && o.AbandonedAt.IsZero()
}

func (o Obligation) Age(now time.Time) time.Duration {
	if o.SubmittedAt.IsZero() {
		return 0
	}
	if !o.CompletedAt.IsZero() {
		return o.CompletedAt.Sub(o.SubmittedAt)
	}
	if !o.AbandonedAt.IsZero() {
		return o.AbandonedAt.Sub(o.SubmittedAt)
	}
	return now.Sub(o.SubmittedAt)
}

type QuestTemplate struct {
	ID          string        `json:"id"`
	Category    Category      `json:"category"`
	Kind        QuestKind     `json:"kind"`
	Instruction string        `json:"instruction"`
	Rationale   string        `json:"rationale"`
	Minutes     int           `json:"minutes"`
	Theater     int           `json:"theater"`
	Difficulty  int           `json:"difficulty"`
	Tags        []string      `json:"tags"`
	Cooldown    time.Duration `json:"cooldown"`
}

type Quest struct {
	ID           string        `json:"id"`
	Sequence     int           `json:"sequence"`
	Instruction  string        `json:"instruction"`
	Rationale    string        `json:"rationale"`
	Kind         QuestKind     `json:"kind"`
	Estimated    time.Duration `json:"estimated"`
	TheaterValue int           `json:"theater_value"`
	Dependencies []string      `json:"dependencies"`
	UnlockedAt   time.Time     `json:"unlocked_at,omitempty"`
	StartedAt    time.Time     `json:"started_at,omitempty"`
	CompletedAt  time.Time     `json:"completed_at,omitempty"`
	SkippedAt    time.Time     `json:"skipped_at,omitempty"`
}

func (q Quest) Complete() bool { return !q.CompletedAt.IsZero() }
func (q Quest) Skipped() bool  { return !q.SkippedAt.IsZero() }
func (q Quest) Finished() bool { return q.Complete() || q.Skipped() }

func (q Quest) Duration(now time.Time) time.Duration {
	if q.StartedAt.IsZero() {
		return 0
	}
	if !q.CompletedAt.IsZero() {
		return q.CompletedAt.Sub(q.StartedAt)
	}
	if !q.SkippedAt.IsZero() {
		return q.SkippedAt.Sub(q.StartedAt)
	}
	return now.Sub(q.StartedAt)
}

type Plan struct {
	ID           string     `json:"id"`
	ObligationID string     `json:"obligation_id"`
	Generation   int        `json:"generation"`
	Quests       []Quest    `json:"quests"`
	CreatedAt    time.Time  `json:"created_at"`
	RetiredAt    *time.Time `json:"retired_at,omitempty"`
	Seed         uint64     `json:"seed"`
}

func (p Plan) CompletedCount() int {
	count := 0
	for _, quest := range p.Quests {
		if quest.Complete() {
			count++
		}
	}
	return count
}

func (p Plan) Finished() bool {
	if len(p.Quests) == 0 {
		return false
	}
	for _, quest := range p.Quests {
		if !quest.Finished() {
			return false
		}
	}
	return true
}

func (p Plan) EstimatedDelay() time.Duration {
	var delay time.Duration
	for _, quest := range p.Quests {
		delay += quest.Estimated
	}
	return delay
}

func (p Plan) ActualDelay(now time.Time) time.Duration {
	end := now
	if p.RetiredAt != nil {
		end = *p.RetiredAt
	}
	if p.CreatedAt.IsZero() || end.Before(p.CreatedAt) {
		return 0
	}
	return end.Sub(p.CreatedAt)
}

func (p Plan) QuestByID(id string) (Quest, bool) {
	for _, quest := range p.Quests {
		if quest.ID == id {
			return quest, true
		}
	}
	return Quest{}, false
}

func (p Plan) Available() []Quest {
	finished := make(map[string]bool, len(p.Quests))
	for _, quest := range p.Quests {
		finished[quest.ID] = quest.Finished()
	}
	available := make([]Quest, 0)
	for _, quest := range p.Quests {
		if quest.Finished() {
			continue
		}
		ready := true
		for _, dependency := range quest.Dependencies {
			if !finished[dependency] {
				ready = false
				break
			}
		}
		if ready {
			available = append(available, quest)
		}
	}
	return available
}

type AchievementCode string

const (
	AchievementMinorTaskMajorCeremony AchievementCode = "minor-task-major-ceremony"
	AchievementBusyBeyondRecognition  AchievementCode = "busy-beyond-recognition"
	AchievementDetourSpecialist       AchievementCode = "detour-specialist"
	AchievementZeroProgress           AchievementCode = "zero-progress-excellent-form"
	AchievementLaterAchieved          AchievementCode = "later-achieved"
	AchievementRecursiveDelay         AchievementCode = "recursive-delay"
	AchievementActuallyFinished       AchievementCode = "actually-finished"
	AchievementFiveMinuteMiracle      AchievementCode = "five-minute-miracle"
)

type Achievement struct {
	Code        AchievementCode `json:"code"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	AwardedAt   time.Time       `json:"awarded_at"`
}

type FocusSession struct {
	ID           string        `json:"id"`
	ObligationID string        `json:"obligation_id"`
	StartedAt    time.Time     `json:"started_at"`
	EndsAt       time.Time     `json:"ends_at"`
	StoppedAt    time.Time     `json:"stopped_at,omitempty"`
	Duration     time.Duration `json:"duration"`
	Completed    bool          `json:"completed"`
}

func (f FocusSession) Remaining(now time.Time) time.Duration {
	if f.Completed || !f.StoppedAt.IsZero() {
		return 0
	}
	remaining := f.EndsAt.Sub(now)
	if remaining < 0 {
		return 0
	}
	return remaining
}

type EventType string

const (
	EventObligationSubmitted EventType = "obligation_submitted"
	EventPlanGenerated       EventType = "plan_generated"
	EventPlanRetired         EventType = "plan_retired"
	EventQuestStarted        EventType = "quest_started"
	EventQuestCompleted      EventType = "quest_completed"
	EventQuestSkipped        EventType = "quest_skipped"
	EventAchievementAwarded  EventType = "achievement_awarded"
	EventFocusStarted        EventType = "focus_started"
	EventFocusStopped        EventType = "focus_stopped"
	EventObligationCompleted EventType = "obligation_completed"
	EventMachineReset        EventType = "machine_reset"
)

type Event struct {
	Type      EventType       `json:"type"`
	At        time.Time       `json:"at"`
	SubjectID string          `json:"subject_id"`
	Detail    string          `json:"detail,omitempty"`
	Metrics   map[string]any  `json:"metrics,omitempty"`
}

type Metrics struct {
	DetoursCompleted  int           `json:"detours_completed"`
	DetoursSkipped    int           `json:"detours_skipped"`
	TasksActuallyDone int           `json:"tasks_actually_done"`
	PlansGenerated    int           `json:"plans_generated"`
	TimeLost          time.Duration `json:"time_lost"`
	EstimatedLost     time.Duration `json:"estimated_lost"`
	Productivity      int           `json:"productivity_theater"`
	DelayIndex        int           `json:"delay_index"`
	FocusSessions     int           `json:"focus_sessions"`
	Achievements      int           `json:"achievements"`
}

type Snapshot struct {
	Version      int                    `json:"version"`
	Obligations  map[string]Obligation  `json:"obligations"`
	Plans        map[string]Plan        `json:"plans"`
	ActivePlanID string                 `json:"active_plan_id,omitempty"`
	Achievements map[AchievementCode]Achievement `json:"achievements"`
	Focus        *FocusSession          `json:"focus,omitempty"`
	Events       []Event                `json:"events"`
	Metrics      Metrics                `json:"metrics"`
}

type Engine struct {
	mu           sync.RWMutex
	clock        func() time.Time
	templates    []QuestTemplate
	obligations  map[string]Obligation
	plans        map[string]Plan
	activePlanID string
	achievements map[AchievementCode]Achievement
	focus        *FocusSession
	events       []Event
	nonce        uint64
}

func New() *Engine {
	return NewWithClock(time.Now)
}

func NewWithClock(clock func() time.Time) *Engine {
	if clock == nil {
		clock = time.Now
	}
	return &Engine{
		clock:        clock,
		templates:    DefaultTemplates(),
		obligations:  make(map[string]Obligation),
		plans:        make(map[string]Plan),
		achievements: make(map[AchievementCode]Achievement),
		events:       make([]Event, 0, 64),
	}
}

func (e *Engine) Submit(text string, urgency Urgency) (Obligation, error) {
	e.mu.Lock()
	defer e.mu.Unlock()
	text = strings.TrimSpace(text)
	if text == "" {
		return Obligation{}, errors.New("an obligation requires at least one visible character")
	}
	if len([]rune(text)) > 120 {
		return Obligation{}, errors.New("an obligation cannot exceed 120 characters")
	}
	now := e.clock()
	id := e.nextID("task", text)
	obligation := Obligation{
		ID:          id,
		Text:        text,
		Category:    Classify(text),
		Urgency:     urgency,
		SubmittedAt: now,
	}
	e.obligations[id] = obligation
	e.record(Event{Type: EventObligationSubmitted, At: now, SubjectID: id, Detail: text})
	return obligation, nil
}

func (e *Engine) Generate(obligationID string, count int) (Plan, error) {
	e.mu.Lock()
	defer e.mu.Unlock()
	obligation, ok := e.obligations[obligationID]
	if !ok {
		return Plan{}, fmt.Errorf("obligation %q does not exist", obligationID)
	}
	if !obligation.Active() {
		return Plan{}, errors.New("cannot generate detours for an inactive obligation")
	}
	if count <= 0 {
		count = DefaultQuestCount
	}
	if count > MaximumQuestCount {
		count = MaximumQuestCount
	}
	now := e.clock()
	generation := 1
	if e.activePlanID != "" {
		if previous, exists := e.plans[e.activePlanID]; exists {
			generation = previous.Generation + 1
			previous.RetiredAt = &now
			e.plans[previous.ID] = previous
			e.record(Event{Type: EventPlanRetired, At: now, SubjectID: previous.ID})
		}
	}
	seed := hashText(fmt.Sprintf("%s:%d", obligation.Text, generation))
	selected := e.selectTemplates(obligation.Category, count, seed)
	quests := make([]Quest, 0, len(selected))
	for index, template := range selected {
		id := e.nextID("detour", template.ID)
		quest := Quest{
			ID:           id,
			Sequence:     index + 1,
			Instruction:  template.Instruction,
			Rationale:    template.Rationale,
			Kind:         template.Kind,
			Estimated:    time.Duration(template.Minutes) * time.Minute,
			TheaterValue: template.Theater,
		}
		if index > 0 {
			quest.Dependencies = []string{quests[index-1].ID}
		} else {
			quest.UnlockedAt = now
		}
		quests = append(quests, quest)
	}
	plan := Plan{
		ID:           e.nextID("plan", obligationID),
		ObligationID: obligationID,
		Generation:   generation,
		Quests:       quests,
		CreatedAt:    now,
		Seed:         seed,
	}
	e.plans[plan.ID] = plan
	e.activePlanID = plan.ID
	e.record(Event{Type: EventPlanGenerated, At: now, SubjectID: plan.ID, Metrics: map[string]any{"quest_count": len(quests), "generation": generation}})
	if generation >= 3 {
		e.award(AchievementRecursiveDelay, now)
	}
	return plan, nil
}

func (e *Engine) StartQuest(questID string) (Quest, error) {
	e.mu.Lock()
	defer e.mu.Unlock()
	plan, index, err := e.findQuest(questID)
	if err != nil {
		return Quest{}, err
	}
	quest := plan.Quests[index]
	if quest.Finished() {
		return Quest{}, errors.New("finished detours cannot be restarted")
	}
	if !e.dependenciesComplete(plan, quest) {
		return Quest{}, errors.New("the preceding detour has not been ceremonially completed")
	}
	if quest.StartedAt.IsZero() {
		quest.StartedAt = e.clock()
		plan.Quests[index] = quest
		e.plans[plan.ID] = plan
		e.record(Event{Type: EventQuestStarted, At: quest.StartedAt, SubjectID: quest.ID})
	}
	return quest, nil
}

func (e *Engine) CompleteQuest(questID string) (Quest, []Achievement, error) {
	e.mu.Lock()
	defer e.mu.Unlock()
	plan, index, err := e.findQuest(questID)
	if err != nil {
		return Quest{}, nil, err
	}
	quest := plan.Quests[index]
	if quest.Complete() {
		return quest, nil, nil
	}
	if quest.Skipped() {
		return Quest{}, nil, errors.New("a skipped detour cannot later claim completion")
	}
	if !e.dependenciesComplete(plan, quest) {
		return Quest{}, nil, errors.New("detours must be performed in their unnecessary order")
	}
	now := e.clock()
	if quest.StartedAt.IsZero() {
		quest.StartedAt = now
	}
	quest.CompletedAt = now
	plan.Quests[index] = quest
	if index+1 < len(plan.Quests) {
		plan.Quests[index+1].UnlockedAt = now
	}
	e.plans[plan.ID] = plan
	e.record(Event{Type: EventQuestCompleted, At: now, SubjectID: quest.ID, Metrics: map[string]any{"theater": quest.TheaterValue}})
	awarded := e.evaluateAchievements(plan, now)
	return quest, awarded, nil
}

func (e *Engine) SkipQuest(questID string) (Quest, error) {
	e.mu.Lock()
	defer e.mu.Unlock()
	plan, index, err := e.findQuest(questID)
	if err != nil {
		return Quest{}, err
	}
	quest := plan.Quests[index]
	if quest.Finished() {
		return quest, nil
	}
	if !e.dependenciesComplete(plan, quest) {
		return Quest{}, errors.New("locked detours cannot be skipped from a safe distance")
	}
	now := e.clock()
	quest.SkippedAt = now
	plan.Quests[index] = quest
	if index+1 < len(plan.Quests) {
		plan.Quests[index+1].UnlockedAt = now
	}
	e.plans[plan.ID] = plan
	e.record(Event{Type: EventQuestSkipped, At: now, SubjectID: quest.ID})
	return quest, nil
}

func (e *Engine) StartFocus(obligationID string, duration time.Duration) (FocusSession, error) {
	e.mu.Lock()
	defer e.mu.Unlock()
	obligation, ok := e.obligations[obligationID]
	if !ok {
		return FocusSession{}, fmt.Errorf("obligation %q does not exist", obligationID)
	}
	if !obligation.Active() {
		return FocusSession{}, errors.New("focus sessions require an active obligation")
	}
	if duration <= 0 {
		duration = FocusDuration
	}
	if duration > 90*time.Minute {
		return FocusSession{}, errors.New("emergency focus is intentionally limited to ninety minutes")
	}
	now := e.clock()
	session := FocusSession{
		ID:           e.nextID("focus", obligationID),
		ObligationID: obligationID,
		StartedAt:    now,
		EndsAt:       now.Add(duration),
		Duration:     duration,
	}
	e.focus = &session
	e.record(Event{Type: EventFocusStarted, At: now, SubjectID: session.ID, Metrics: map[string]any{"seconds": int(duration.Seconds())}})
	return session, nil
}

func (e *Engine) TickFocus() (*FocusSession, []Achievement) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.focus == nil {
		return nil, nil
	}
	now := e.clock()
	if !e.focus.Completed && !now.Before(e.focus.EndsAt) {
		e.focus.Completed = true
		e.focus.StoppedAt = now
		e.record(Event{Type: EventFocusStopped, At: now, SubjectID: e.focus.ID, Detail: "completed"})
		award := e.award(AchievementFiveMinuteMiracle, now)
		if award.Code != "" {
			return cloneFocus(e.focus), []Achievement{award}
		}
	}
	return cloneFocus(e.focus), nil
}

func (e *Engine) StopFocus() (*FocusSession, error) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.focus == nil {
		return nil, errors.New("no focus session is active")
	}
	if e.focus.StoppedAt.IsZero() {
		e.focus.StoppedAt = e.clock()
		e.record(Event{Type: EventFocusStopped, At: e.focus.StoppedAt, SubjectID: e.focus.ID, Detail: "stopped"})
	}
	return cloneFocus(e.focus), nil
}

func (e *Engine) CompleteObligation(obligationID string) (Obligation, []Achievement, error) {
	e.mu.Lock()
	defer e.mu.Unlock()
	obligation, ok := e.obligations[obligationID]
	if !ok {
		return Obligation{}, nil, fmt.Errorf("obligation %q does not exist", obligationID)
	}
	if !obligation.CompletedAt.IsZero() {
		return obligation, nil, nil
	}
	now := e.clock()
	obligation.CompletedAt = now
	e.obligations[obligationID] = obligation
	e.record(Event{Type: EventObligationCompleted, At: now, SubjectID: obligationID})
	award := e.award(AchievementActuallyFinished, now)
	awards := make([]Achievement, 0, 1)
	if award.Code != "" {
		awards = append(awards, award)
	}
	return obligation, awards, nil
}

func (e *Engine) Metrics() Metrics {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.metricsLocked(e.clock())
}

func (e *Engine) Snapshot() Snapshot {
	e.mu.RLock()
	defer e.mu.RUnlock()
	obligations := make(map[string]Obligation, len(e.obligations))
	for id, obligation := range e.obligations {
		obligations[id] = obligation
	}
	plans := make(map[string]Plan, len(e.plans))
	for id, plan := range e.plans {
		plans[id] = clonePlan(plan)
	}
	achievements := make(map[AchievementCode]Achievement, len(e.achievements))
	for code, achievement := range e.achievements {
		achievements[code] = achievement
	}
	return Snapshot{
		Version:      ArchiveVersion,
		Obligations:  obligations,
		Plans:        plans,
		ActivePlanID: e.activePlanID,
		Achievements: achievements,
		Focus:        cloneFocus(e.focus),
		Events:       append([]Event(nil), e.events...),
		Metrics:      e.metricsLocked(e.clock()),
	}
}

func (e *Engine) Export() ([]byte, error) {
	return json.MarshalIndent(e.Snapshot(), "", "  ")
}

func Restore(data []byte, clock func() time.Time) (*Engine, error) {
	var snapshot Snapshot
	if err := json.Unmarshal(data, &snapshot); err != nil {
		return nil, fmt.Errorf("decode procrastination archive: %w", err)
	}
	if snapshot.Version != ArchiveVersion {
		return nil, fmt.Errorf("unsupported archive version %d", snapshot.Version)
	}
	engine := NewWithClock(clock)
	engine.obligations = snapshot.Obligations
	engine.plans = snapshot.Plans
	engine.activePlanID = snapshot.ActivePlanID
	engine.achievements = snapshot.Achievements
	engine.focus = snapshot.Focus
	engine.events = append([]Event(nil), snapshot.Events...)
	return engine, nil
}

func (e *Engine) Reset() {
	e.mu.Lock()
	defer e.mu.Unlock()
	now := e.clock()
	e.obligations = make(map[string]Obligation)
	e.plans = make(map[string]Plan)
	e.activePlanID = ""
	e.achievements = make(map[AchievementCode]Achievement)
	e.focus = nil
	e.events = []Event{{Type: EventMachineReset, At: now, SubjectID: "machine"}}
}

func (e *Engine) ActivePlan() (Plan, bool) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	plan, ok := e.plans[e.activePlanID]
	return clonePlan(plan), ok
}

func (e *Engine) Events(after int) []Event {
	e.mu.RLock()
	defer e.mu.RUnlock()
	if after < 0 {
		after = 0
	}
	if after >= len(e.events) {
		return nil
	}
	return append([]Event(nil), e.events[after:]...)
}

func Classify(text string) Category {
	normalized := strings.ToLower(text)
	tokens := tokenize(normalized)
	scores := map[Category]int{
		CategoryEmail:    scoreTokens(tokens, "email", "emails", "inbox", "reply", "respond", "message"),
		CategoryStudy:    scoreTokens(tokens, "study", "exam", "test", "homework", "class", "quiz", "read"),
		CategoryCleaning: scoreTokens(tokens, "clean", "room", "laundry", "closet", "dishes", "organize", "tidy"),
		CategoryWriting:  scoreTokens(tokens, "write", "report", "essay", "draft", "article", "paper", "chapter"),
		CategoryWork:     scoreTokens(tokens, "work", "project", "meeting", "presentation", "client", "deadline", "deck"),
	}
	best, bestScore := CategoryGeneral, 0
	order := []Category{CategoryEmail, CategoryStudy, CategoryCleaning, CategoryWriting, CategoryWork}
	for _, category := range order {
		if scores[category] > bestScore {
			best, bestScore = category, scores[category]
		}
	}
	return best
}

func tokenize(text string) map[string]int {
	fields := strings.FieldsFunc(text, func(r rune) bool {
		return !(r >= 'a' && r <= 'z') && !(r >= '0' && r <= '9')
	})
	tokens := make(map[string]int, len(fields))
	for _, field := range fields {
		tokens[field]++
	}
	return tokens
}

func scoreTokens(tokens map[string]int, words ...string) int {
	score := 0
	for _, word := range words {
		score += tokens[word]
	}
	return score
}

func DefaultTemplates() []QuestTemplate {
	return []QuestTemplate{
		{ID: "email-emotional-labels", Category: CategoryEmail, Kind: QuestOrganization, Instruction: "Reorganize your inbox labels by emotional tone", Rationale: "Future messages deserve a more nuanced filing climate.", Minutes: 11, Theater: 16, Difficulty: 2, Tags: []string{"inbox", "labels"}},
		{ID: "email-future-signature", Category: CategoryEmail, Kind: QuestAesthetic, Instruction: "Draft a signature for a future job title", Rationale: "The present reply will be stronger once your hypothetical promotion has branding.", Minutes: 9, Theater: 18, Difficulty: 1, Tags: []string{"signature", "identity"}},
		{ID: "email-best-research", Category: CategoryEmail, Kind: QuestResearch, Instruction: "Research whether 'Best' sounds passive-aggressive", Rationale: "Tone cannot be rushed, especially when the message itself can.", Minutes: 14, Theater: 21, Difficulty: 2, Tags: []string{"etiquette", "tone"}},
		{ID: "email-newsletter-archive", Category: CategoryEmail, Kind: QuestOrganization, Instruction: "Archive every newsletter from 2021", Rationale: "Historical inbox sediment threatens present-day correspondence.", Minutes: 17, Theater: 19, Difficulty: 3, Tags: []string{"archive", "cleanup"}},
		{ID: "email-unread-flag", Category: CategoryEmail, Kind: QuestAesthetic, Instruction: "Design a tiny flag for unread messages", Rationale: "A geopolitical symbol will clarify which emails remain unopened.", Minutes: 13, Theater: 23, Difficulty: 2, Tags: []string{"design", "unread"}},
		{ID: "study-font", Category: CategoryStudy, Kind: QuestAesthetic, Instruction: "Choose the perfect study font", Rationale: "Knowledge enters differently through a carefully selected lowercase a.", Minutes: 12, Theater: 20, Difficulty: 2, Tags: []string{"font", "notes"}},
		{ID: "study-highlighters", Category: CategoryStudy, Kind: QuestClassification, Instruction: "Rank every highlighter by seriousness", Rationale: "The exam material must understand the hierarchy of yellow.", Minutes: 10, Theater: 18, Difficulty: 1, Tags: []string{"stationery", "ranking"}},
		{ID: "study-playlist", Category: CategoryStudy, Kind: QuestPreparation, Instruction: "Make a playlist called Academic Weather", Rationale: "No learning can begin before the atmospheric conditions are correct.", Minutes: 19, Theater: 25, Difficulty: 2, Tags: []string{"music", "atmosphere"}},
		{ID: "study-yesterday", Category: CategoryStudy, Kind: QuestResearch, Instruction: "Calculate how long studying would take if started yesterday", Rationale: "Counterfactual scheduling may reveal the timeline you deserved.", Minutes: 8, Theater: 14, Difficulty: 2, Tags: []string{"math", "schedule"}},
		{ID: "study-title-page", Category: CategoryStudy, Kind: QuestAesthetic, Instruction: "Rewrite the title page with better margins", Rationale: "Content follows form, preferably after form has been adjusted several times.", Minutes: 15, Theater: 22, Difficulty: 2, Tags: []string{"formatting", "margins"}},
		{ID: "clean-archive-photo", Category: CategoryCleaning, Kind: QuestDocumentation, Instruction: "Photograph the mess for archival purposes", Rationale: "Removal without documentation would erase an important domestic era.", Minutes: 7, Theater: 17, Difficulty: 1, Tags: []string{"photo", "archive"}},
		{ID: "clean-personality", Category: CategoryCleaning, Kind: QuestClassification, Instruction: "Sort one drawer by object personality", Rationale: "Compatible objects may remain cleaner through improved social grouping.", Minutes: 18, Theater: 24, Difficulty: 3, Tags: []string{"drawer", "sorting"}},
		{ID: "clean-dust-history", Category: CategoryCleaning, Kind: QuestResearch, Instruction: "Research the history of dust", Rationale: "The material should be understood before its removal is authorized.", Minutes: 16, Theater: 26, Difficulty: 2, Tags: []string{"dust", "history"}},
		{ID: "clean-donation-labels", Category: CategoryCleaning, Kind: QuestOrganization, Instruction: "Create a donation box label system", Rationale: "No box can receive an object until its taxonomy is complete.", Minutes: 12, Theater: 21, Difficulty: 2, Tags: []string{"labels", "donation"}},
		{ID: "clean-supplies", Category: CategoryCleaning, Kind: QuestPreparation, Instruction: "Clean the cleaning supplies first", Rationale: "Unclean tools cannot ethically participate in cleanliness.", Minutes: 14, Theater: 23, Difficulty: 2, Tags: []string{"supplies", "recursive"}},
		{ID: "write-rename", Category: CategoryWriting, Kind: QuestOrganization, Instruction: "Rename the document three times", Rationale: "The draft needs an identity before it can develop a body.", Minutes: 8, Theater: 16, Difficulty: 1, Tags: []string{"filename", "identity"}},
		{ID: "write-unused-fact", Category: CategoryWriting, Kind: QuestResearch, Instruction: "Research a fact that will not appear in the draft", Rationale: "Invisible research provides structural moral support.", Minutes: 18, Theater: 24, Difficulty: 2, Tags: []string{"research", "fact"}},
		{ID: "write-cursor", Category: CategoryWriting, Kind: QuestAesthetic, Instruction: "Choose a more authoritative cursor color", Rationale: "The insertion point should command respect before producing prose.", Minutes: 7, Theater: 19, Difficulty: 1, Tags: []string{"cursor", "authority"}},
		{ID: "write-bibliography", Category: CategoryWriting, Kind: QuestPreparation, Instruction: "Format the bibliography before writing", Rationale: "Sources prefer to know where they will eventually be cited.", Minutes: 21, Theater: 27, Difficulty: 3, Tags: []string{"citations", "formatting"}},
		{ID: "write-acknowledgements", Category: CategoryWriting, Kind: QuestMetaWork, Instruction: "Draft an acknowledgements section", Rationale: "Gratitude can be completed even while the actual work cannot.", Minutes: 13, Theater: 23, Difficulty: 2, Tags: []string{"thanks", "preface"}},
		{ID: "work-meeting", Category: CategoryWork, Kind: QuestMetaWork, Instruction: "Create a meeting about avoiding meetings", Rationale: "The calendar should formally recognize its own overuse.", Minutes: 22, Theater: 30, Difficulty: 3, Tags: []string{"meeting", "calendar"}},
		{ID: "work-color-code", Category: CategoryWork, Kind: QuestClassification, Instruction: "Color-code the task by hypothetical urgency", Rationale: "Priority becomes credible only after receiving a hex value.", Minutes: 9, Theater: 18, Difficulty: 1, Tags: []string{"priority", "color"}},
		{ID: "work-dashboard", Category: CategoryWork, Kind: QuestMetaWork, Instruction: "Build a dashboard for one checkbox", Rationale: "A single binary state deserves operational visibility.", Minutes: 28, Theater: 35, Difficulty: 4, Tags: []string{"dashboard", "metrics"}},
		{ID: "work-easier-version", Category: CategoryWork, Kind: QuestResearch, Instruction: "Research an easier version of the task", Rationale: "Alternative scope must be thoroughly explored before current scope can be ignored.", Minutes: 17, Theater: 22, Difficulty: 2, Tags: []string{"scope", "research"}},
		{ID: "work-status", Category: CategoryWork, Kind: QuestDocumentation, Instruction: "Write a status update about beginning soon", Rationale: "Stakeholders require advance notice of future momentum.", Minutes: 11, Theater: 25, Difficulty: 1, Tags: []string{"status", "communication"}},
		{ID: "general-alphabetize", Category: CategoryGeneral, Kind: QuestOrganization, Instruction: "Alphabetize three nearby objects", Rationale: "Local order may radiate toward the obligation eventually.", Minutes: 8, Theater: 15, Difficulty: 1, Tags: []string{"objects", "alphabet"}},
		{ID: "general-flag", Category: CategoryGeneral, Kind: QuestAesthetic, Instruction: "Design a flag for the task", Rationale: "The obligation cannot proceed without a recognizable sovereign identity.", Minutes: 14, Theater: 22, Difficulty: 2, Tags: []string{"flag", "identity"}},
		{ID: "general-rectangle", Category: CategoryGeneral, Kind: QuestResearch, Instruction: "Find the quietest rectangle in the room", Rationale: "Geometric acoustic conditions may be responsible for current resistance.", Minutes: 12, Theater: 20, Difficulty: 2, Tags: []string{"geometry", "quiet"}},
		{ID: "general-playlist", Category: CategoryGeneral, Kind: QuestPreparation, Instruction: "Make an unnecessarily specific playlist", Rationale: "A precise soundtrack converts avoidance into curation.", Minutes: 19, Theater: 26, Difficulty: 2, Tags: []string{"music", "curation"}},
		{ID: "general-later-origin", Category: CategoryGeneral, Kind: QuestResearch, Instruction: "Research the origin of the word 'later'", Rationale: "Etymology may establish whether postponement is historically justified.", Minutes: 16, Theater: 27, Difficulty: 2, Tags: []string{"later", "etymology"}},
		{ID: "general-reorganize", Category: CategoryGeneral, Kind: QuestOrganization, Instruction: "Organize something that was already organized", Rationale: "Existing organization may conceal a superior organizational possibility.", Minutes: 13, Theater: 21, Difficulty: 2, Tags: []string{"organize", "repeat"}},
		{ID: "general-hydration", Category: CategoryGeneral, Kind: QuestIntermission, Instruction: "Take a strategic hydration intermission", Rationale: "Water is important and therefore technically adjacent to every task.", Minutes: 6, Theater: 12, Difficulty: 1, Tags: []string{"water", "break"}},
		{ID: "general-folders", Category: CategoryGeneral, Kind: QuestOrganization, Instruction: "Rename every folder involved", Rationale: "The file system deserves narrative consistency before work enters it.", Minutes: 15, Theater: 24, Difficulty: 2, Tags: []string{"folders", "names"}},
	}
}

func (e *Engine) selectTemplates(category Category, count int, seed uint64) []QuestTemplate {
	candidates := make([]QuestTemplate, 0, len(e.templates))
	for _, template := range e.templates {
		if template.Category == category || template.Category == CategoryGeneral {
			candidates = append(candidates, template)
		}
	}
	sort.SliceStable(candidates, func(i, j int) bool {
		left := mixedHash(seed, candidates[i].ID)
		right := mixedHash(seed, candidates[j].ID)
		return left < right
	})
	if count > len(candidates) {
		count = len(candidates)
	}
	selected := make([]QuestTemplate, 0, count)
	kinds := make(map[QuestKind]int)
	for _, template := range candidates {
		if len(selected) >= count {
			break
		}
		if kinds[template.Kind] >= 2 && len(candidates)-len(selected) > count {
			continue
		}
		selected = append(selected, template)
		kinds[template.Kind]++
	}
	return selected
}

func (e *Engine) findQuest(questID string) (Plan, int, error) {
	for _, plan := range e.plans {
		for index, quest := range plan.Quests {
			if quest.ID == questID {
				return plan, index, nil
			}
	}
	}
	return Plan{}, -1, fmt.Errorf("detour %q does not exist", questID)
}

func (e *Engine) dependenciesComplete(plan Plan, quest Quest) bool {
	for _, dependency := range quest.Dependencies {
		found, ok := plan.QuestByID(dependency)
		if !ok || !found.Finished() {
			return false
		}
	}
	return true
}

func (e *Engine) evaluateAchievements(plan Plan, now time.Time) []Achievement {
	awarded := make([]Achievement, 0, 3)
	completed := plan.CompletedCount()
	if completed == 1 {
		awarded = appendAward(awarded, e.award(AchievementMinorTaskMajorCeremony, now))
	}
	if completed >= 3 {
		awarded = appendAward(awarded, e.award(AchievementDetourSpecialist, now))
	}
	if completed == len(plan.Quests) {
		awarded = appendAward(awarded, e.award(AchievementBusyBeyondRecognition, now))
		awarded = appendAward(awarded, e.award(AchievementZeroProgress, now))
		awarded = appendAward(awarded, e.award(AchievementLaterAchieved, now))
	}
	return awarded
}

func appendAward(target []Achievement, achievement Achievement) []Achievement {
	if achievement.Code != "" {
		return append(target, achievement)
	}
	return target
}

func (e *Engine) award(code AchievementCode, now time.Time) Achievement {
	if _, exists := e.achievements[code]; exists {
		return Achievement{}
	}
	title, description := achievementCopy(code)
	achievement := Achievement{Code: code, Title: title, Description: description, AwardedAt: now}
	e.achievements[code] = achievement
	e.record(Event{Type: EventAchievementAwarded, At: now, SubjectID: string(code), Detail: title})
	return achievement
}

func achievementCopy(code AchievementCode) (string, string) {
	switch code {
	case AchievementMinorTaskMajorCeremony:
		return "Minor Task, Major Ceremony", "Completed one detour with disproportionate administrative dignity."
	case AchievementBusyBeyondRecognition:
		return "Busy Beyond Recognition", "Finished an entire plan without approaching the original obligation."
	case AchievementDetourSpecialist:
		return "Detour Specialist", "Completed three consecutively less relevant activities."
	case AchievementZeroProgress:
		return "Zero Progress, Excellent Form", "Produced visible effort while preserving the task exactly as submitted."
	case AchievementLaterAchieved:
		return "Later Achieved", "Successfully transformed now into an undefined future period."
	case AchievementRecursiveDelay:
		return "Recursive Delay", "Regenerated the avoidance plan three times."
	case AchievementActuallyFinished:
		return "Actually Finished", "Completed an obligation despite the machine's extensive support."
	case AchievementFiveMinuteMiracle:
		return "Five Minute Miracle", "Remained with the real task for one uninterrupted focus interval."
	default:
		return "Unfiled Achievement", "The office has not yet prepared the appropriate stamp."
	}
}

func (e *Engine) metricsLocked(now time.Time) Metrics {
	metrics := Metrics{Achievements: len(e.achievements)}
	theaterTotal := 0
	theaterPossible := 0
	for _, obligation := range e.obligations {
		if !obligation.CompletedAt.IsZero() {
			metrics.TasksActuallyDone++
		}
	}
	for _, plan := range e.plans {
		metrics.PlansGenerated++
		metrics.EstimatedLost += plan.EstimatedDelay()
		metrics.TimeLost += plan.ActualDelay(now)
		for _, quest := range plan.Quests {
			theaterPossible += quest.TheaterValue
			if quest.Complete() {
				metrics.DetoursCompleted++
				theaterTotal += quest.TheaterValue
			}
			if quest.Skipped() {
				metrics.DetoursSkipped++
			}
		}
	}
	for _, event := range e.events {
		if event.Type == EventFocusStarted {
			metrics.FocusSessions++
		}
	}
	if theaterPossible > 0 {
		metrics.Productivity = int(math.Round(float64(theaterTotal) / float64(theaterPossible) * 100))
	}
	minutes := int(metrics.TimeLost.Minutes())
	metrics.DelayIndex = clampInt(17+metrics.DetoursCompleted*14+minutes/2+metrics.PlansGenerated*3, 0, 99)
	return metrics
}

func (e *Engine) record(event Event) {
	e.events = append(e.events, event)
}

func (e *Engine) nextID(prefix, material string) string {
	e.nonce++
	return fmt.Sprintf("%s-%08x-%03d", prefix, hashText(material)&0xffffffff, e.nonce)
}

// MiniGameRule keeps the browser games deterministic when the Go engine is
// compiled to WebAssembly. The interface may reshuffle their order, but each
// game always retains a real objective and a measurable win condition.
type MiniGameRule struct {
	ID          string
	Target      int
	Instruction string
}

// MiniGameRules returns the complete activity cabinet. Keeping this catalog in
// the engine means analytics, saved archives, and the playful browser layer
// agree about what a completed detour actually represents.
func MiniGameRules() []MiniGameRule {
	return []MiniGameRule{
		{ID: "planes", Target: 6, Instruction: "catch six falling memos"},
		{ID: "clock", Target: 3, Instruction: "freeze three perfect moments"},
		{ID: "snail", Target: 1, Instruction: "reach the lamp through the maze"},
		{ID: "pairs", Target: 4, Instruction: "discover four object friendships"},
		{ID: "thread", Target: 6, Instruction: "connect six conspiracy pins"},
		{ID: "rhythm", Target: 6, Instruction: "repeat the six-beat desk rhythm"},
		{ID: "cabinet", Target: 6, Instruction: "file six wandering records"},
		{ID: "typewriter", Target: 3, Instruction: "repair three missing letters"},
		{ID: "spotlight", Target: 4, Instruction: "discover four midnight objects"},
		{ID: "balance", Target: 7, Instruction: "balance seven drifting pages"},
		{ID: "switchboard", Target: 9, Instruction: "align nine brass connections"},
	}
}

// MiniGameProgress normalizes a game's current state for the shared metrics
// display. Unknown games and negative values intentionally produce zero.
func MiniGameProgress(gameID string, current int) int {
	if current < 0 {
		return 0
	}
	for _, rule := range MiniGameRules() {
		if rule.ID != gameID {
			continue
		}
		if current >= rule.Target {
			return 100
		}
		return current * 100 / rule.Target
	}
	return 0
}

func hashText(text string) uint64 {
	hash := fnv.New64a()
	_, _ = hash.Write([]byte(text))
	return hash.Sum64()
}

func mixedHash(seed uint64, text string) uint64 {
	value := hashText(text) ^ seed
	value ^= value >> 30
	value *= 0xbf58476d1ce4e5b9
	value ^= value >> 27
	value *= 0x94d049bb133111eb
	return value ^ (value >> 31)
}

func clampInt(value, minimum, maximum int) int {
	if value < minimum {
		return minimum
	}
	if value > maximum {
		return maximum
	}
	return value
}

func clonePlan(plan Plan) Plan {
	plan.Quests = append([]Quest(nil), plan.Quests...)
	for index := range plan.Quests {
		plan.Quests[index].Dependencies = append([]string(nil), plan.Quests[index].Dependencies...)
	}
	return plan
}

func cloneFocus(focus *FocusSession) *FocusSession {
	if focus == nil {
		return nil
	}
	clone := *focus
	return &clone
}
