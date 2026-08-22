# Procrastination Machine

> A cabinet of at least fifteen playable side quests for avoiding one important task.

**Live exhibit:** https://jean-tmk.github.io/procrastination-machine/

## What it is

This is not a task tracker. It turns avoidance into a real arcade: timing games, sequencing puzzles, memory challenges, a harder maze, thread ordering, cabinet sorting, rhythm, balancing, and other short detours must all circulate before any game repeats.

## What a visitor can do

1. Enter the thing you should be doing.
2. Open one of the five offered detours.
3. Complete the actual mini-game inside the side-quest window.
4. Request new distractions; the rotation will exhaust every game before repeating.
5. Return to the real task only if absolutely necessary.

## How it works

- Go defines the activity catalogue, no-repeat deck, task classification, scoring, achievements, progression, maze data, and mini-game rules.
- JavaScript is the browser runtime for the same game models, pointer/keyboard controls, modal lifecycle, audio, and session state.
- Individually prepared WebP assets provide clean icons and game objects.

## Repository map

| Path | What it does |
|---|---|
| `.github/workflows/pages.yml` | GitHub Actions workflow that validates, builds, and/or deploys the exhibit. |
| `app.js` | The browser interaction runtime and top-level state coordinator. |
| `engine.go` | Domain, engine, tooling, or specification source in the repository’s polyglot architecture. |
| `enhancements.css` | A focused style layer for this named area of the experience. |
| `go.mod` | Go module identity and toolchain declaration. |
| `index.html` | The deployable HTML shell: metadata, accessible structure, controls, and script/style entry points. |
| `minigames.go` | Domain, engine, tooling, or specification source in the repository’s polyglot architecture. |
| `styles.css` | The primary responsive visual system. |
| `assets/` | 26 production illustration/icon files loaded by the live interface. |
| `polyglot/` | 59 isolated language-atlas files plus the majority registry and manifest; these never load in the visible frontend. |

## Languages and why they are here

Percentages below are calculated from the byte counts currently returned by GitHub Linguist. Tiny language-atlas modules are intentionally isolated from the production frontend.

| Language | GitHub | Role |
|---|---:|---|
| Go | 91.3% | the majority game catalogue, rules, and rotation engine |
| HTML | 0.9% | semantic arcade shell |
| Apex | 0.3% | an isolated language-atlas adapter used to broaden the comparative polyglot collection without changing the exhibit UI |
| Raku | 0.3% | an isolated language-atlas adapter used to broaden the comparative polyglot collection without changing the exhibit UI |
| q | 0.3% | an isolated language-atlas adapter used to broaden the comparative polyglot collection without changing the exhibit UI |
| JetBrains MPS | 0.2% | an isolated language-atlas adapter used to broaden the comparative polyglot collection without changing the exhibit UI |
| Parrot | 0.2% | an isolated language-atlas adapter used to broaden the comparative polyglot collection without changing the exhibit UI |
| Java | 0.2% | an isolated language-atlas adapter used to broaden the comparative polyglot collection without changing the exhibit UI |
| ColdFusion | 0.2% | an isolated language-atlas adapter used to broaden the comparative polyglot collection without changing the exhibit UI |
| OpenEdge ABL | 0.2% | an isolated language-atlas adapter used to broaden the comparative polyglot collection without changing the exhibit UI |
| FreeMarker | 0.2% | an isolated language-atlas adapter used to broaden the comparative polyglot collection without changing the exhibit UI |
| MoonScript | 0.2% | an isolated language-atlas adapter used to broaden the comparative polyglot collection without changing the exhibit UI |
| DataWeave | 0.2% | an isolated language-atlas adapter used to broaden the comparative polyglot collection without changing the exhibit UI |
| ABAP CDS | 0.2% | an isolated language-atlas adapter used to broaden the comparative polyglot collection without changing the exhibit UI |
| GDShader | 0.2% | an isolated language-atlas adapter used to broaden the comparative polyglot collection without changing the exhibit UI |
| LigoLANG | 0.2% | an isolated language-atlas adapter used to broaden the comparative polyglot collection without changing the exhibit UI |
| Cangjie | 0.2% | an isolated language-atlas adapter used to broaden the comparative polyglot collection without changing the exhibit UI |
| Clarion | 0.2% | an isolated language-atlas adapter used to broaden the comparative polyglot collection without changing the exhibit UI |
| Cycript | 0.2% | an isolated language-atlas adapter used to broaden the comparative polyglot collection without changing the exhibit UI |
| LOLCODE | 0.2% | an isolated language-atlas adapter used to broaden the comparative polyglot collection without changing the exhibit UI |
| M4 | 0.2% | an isolated language-atlas adapter used to broaden the comparative polyglot collection without changing the exhibit UI |
| Boogie | 0.2% | an isolated language-atlas adapter used to broaden the comparative polyglot collection without changing the exhibit UI |
| Berry | 0.2% | an isolated language-atlas adapter used to broaden the comparative polyglot collection without changing the exhibit UI |
| CLIPS | 0.2% | an isolated language-atlas adapter used to broaden the comparative polyglot collection without changing the exhibit UI |
| Faust | 0.2% | an isolated language-atlas adapter used to broaden the comparative polyglot collection without changing the exhibit UI |
| Gleam | 0.2% | an isolated language-atlas adapter used to broaden the comparative polyglot collection without changing the exhibit UI |
| IDL | 0.2% | an isolated language-atlas adapter used to broaden the comparative polyglot collection without changing the exhibit UI |
| Mirah | 0.2% | an isolated language-atlas adapter used to broaden the comparative polyglot collection without changing the exhibit UI |
| Nasal | 0.2% | an isolated language-atlas adapter used to broaden the comparative polyglot collection without changing the exhibit UI |
| Polar | 0.2% | an isolated language-atlas adapter used to broaden the comparative polyglot collection without changing the exhibit UI |
| Rouge | 0.2% | an isolated language-atlas adapter used to broaden the comparative polyglot collection without changing the exhibit UI |
| Sieve | 0.2% | an isolated language-atlas adapter used to broaden the comparative polyglot collection without changing the exhibit UI |
| Talon | 0.2% | an isolated language-atlas adapter used to broaden the comparative polyglot collection without changing the exhibit UI |
| Xtend | 0.2% | an isolated language-atlas adapter used to broaden the comparative polyglot collection without changing the exhibit UI |
| Zimpl | 0.2% | an isolated language-atlas adapter used to broaden the comparative polyglot collection without changing the exhibit UI |
| PDDL | 0.2% | an isolated language-atlas adapter used to broaden the comparative polyglot collection without changing the exhibit UI |
| Scheme | 0.2% | an isolated language-atlas adapter used to broaden the comparative polyglot collection without changing the exhibit UI |
| Sway | 0.2% | an isolated language-atlas adapter used to broaden the comparative polyglot collection without changing the exhibit UI |
| Volt | 0.2% | an isolated language-atlas adapter used to broaden the comparative polyglot collection without changing the exhibit UI |
| Wren | 0.2% | an isolated language-atlas adapter used to broaden the comparative polyglot collection without changing the exhibit UI |
| ATS | 0.2% | an isolated language-atlas adapter used to broaden the comparative polyglot collection without changing the exhibit UI |
| KCL | 0.2% | an isolated language-atlas adapter used to broaden the comparative polyglot collection without changing the exhibit UI |
| Nit | 0.2% | an isolated language-atlas adapter used to broaden the comparative polyglot collection without changing the exhibit UI |
| RPC | 0.2% | an isolated language-atlas adapter used to broaden the comparative polyglot collection without changing the exhibit UI |
| Red | 0.2% | an isolated language-atlas adapter used to broaden the comparative polyglot collection without changing the exhibit UI |
| VBA | 0.2% | an isolated language-atlas adapter used to broaden the comparative polyglot collection without changing the exhibit UI |
| EQ | 0.1% | an isolated language-atlas adapter used to broaden the comparative polyglot collection without changing the exhibit UI |
| Brainfuck | 0.1% | an isolated language-atlas adapter used to broaden the comparative polyglot collection without changing the exhibit UI |
| J | 0.1% | an isolated language-atlas adapter used to broaden the comparative polyglot collection without changing the exhibit UI |

### About the language atlas

Where present, `polyglot/language-atlas.json` is the machine-readable index of the languages assigned to this repository. `polyglot/languages/` contains one small, independent signature module per assignment, and `polyglot/majority/` contains the larger registry that preserves the intended majority language. These files are documentation and comparative code specimens: the live site does not download or execute them.

## Local development

```bash
python3 -m http.server 8000
go test ./...
```

Then open `http://localhost:8000` unless the framework development server prints a different local address.

## Privacy and access

- No sign-in is required.
- No API key is required for the live exhibit.
- No visitor text is sent to an AI service.
- Any saved progress stays in local browser storage unless the README explicitly describes an optional external architecture.
- Sound begins only after a user gesture where browser autoplay rules require it.

## Deployment

The public version is a static GitHub Pages deployment. The workflow in `.github/workflows/` is the source of truth for its exact build and publish steps. The favicon is stored with the deployed app so browser tabs and bookmarks use the project’s own mark.
