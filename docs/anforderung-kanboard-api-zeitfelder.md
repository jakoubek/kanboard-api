# Anforderung: Zeitfelder im `Task`-Typ der `kanboard-api`-Bibliothek

> **Umsetzung im externen Projekt** `~/Dev/kanboard-api` (Modul
> `code.beautifulmachines.dev/jakoubek/kanboard-api`), **nicht** in hqcli. Dieses Dokument ist die
> Vorlage, die in einer eigenen Session im dortigen Projekt abgearbeitet wird (Plan-Mode, manueller
> Test, dann Commit/Release gemäß dessen `AGENTS.md`).

## Ziel

Der `Task`-Typ soll die Kanboard-Zeitfelder **`time_estimated`** (geschätzte Zeit) und
**`time_spent`** (verbrauchte Zeit) tragen, damit hqcli seine Schätz-Reports (`estsum`, `estprj`,
`esttick`) und die Ticket-Detailanzeige/den Export **ohne direkten MySQL-Zugriff** über die API
umsetzen kann.

## Kontext / Warum

hqcli löst seinen SSH-getunnelten MySQL-Zugriff ab (Phase 1 der hqcli-Konsolidierung, siehe
`docs/hqcli-konsolidierung.md`). Alle betroffenen Reports lassen sich auf `GetAllProjects` /
`GetAllTasks` / `GetColumns` / `GetTask` portieren — **einziger Blocker**: Die Bibliothek mappt
`time_estimated`/`time_spent` aktuell nicht. Konkret genutzt werden sie in hqcli so:

- `SUM(time_estimated)` → Gesamt-Schätzung in Stunden; `SUM(time_estimated)/8.0` → Personentage
- `COUNT(... time_estimated > 0 / = 0 ...)` → Anzahl Tickets mit/ohne Schätzung
- `time_spent` in Ticket-Detail und Export

Die Aggregation (Summe, Division durch 8, Zählung) bleibt **Client-Logik in hqcli** — die
Bibliothek muss die Werte nur bereitstellen.

## Ist-Zustand (Bibliothek)

- `types.go` — `Task`-Struct (aktuell ohne Zeitfelder) und die Custom-Unmarshaler `StringInt`,
  `StringInt64`, `StringBool`. **Kein `StringFloat` vorhanden.**
- `tasks.go` — `GetTask`, `GetAllTasks`, `GetTaskByReference`, `SearchTasks` deserialisieren per
  `c.call(...)` direkt in `Task`. Ein Ergänzen der Struct-Felder wirkt daher automatisch für alle
  Lese-Methoden; **kein zusätzliches Mapping nötig.**

## Anforderungen

### R1 — Neuer Typ `StringFloat`
Ein `float64`, der aus JSON als **String *oder* Zahl** deserialisiert werden kann, analog zu
`StringInt` (`types.go`). Kanboard liefert numerische DB-Felder oft als String (z. B. `"8"`,
`"2.5"`, `"0.00"`), teils als Zahl. Leerer String → `0`.

Referenz-Implementierung (am Muster von `StringInt` orientiert):

```go
// StringFloat is a float64 that can be unmarshaled from a JSON string or number.
type StringFloat float64

// UnmarshalJSON implements json.Unmarshaler.
func (f *StringFloat) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		// Try as raw number
		var num float64
		if err := json.Unmarshal(data, &num); err != nil {
			return err
		}
		*f = StringFloat(num)
		return nil
	}
	if s == "" {
		*f = 0
		return nil
	}
	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return err
	}
	*f = StringFloat(val)
	return nil
}
```

`MarshalJSON` ist **nicht** erforderlich (nur Lesen; `StringInt` hat ebenfalls kein Marshal). Falls
später Schreibsupport gewünscht wird, kann es analog zu `StringBool.MarshalJSON` nachgezogen werden.

### R2 — Felder am `Task`-Typ ergänzen
In `types.go` im `Task`-Struct:

```go
TimeEstimated StringFloat `json:"time_estimated"`
TimeSpent     StringFloat `json:"time_spent"`
```

Einheit: **Stunden** (float). Keine Umrechnung in der Bibliothek — Werte 1:1 durchreichen.

### R3 — Tests
Konsistent zum bestehenden Teststil (`types_test.go`, `tasks_test.go`):

- **`StringFloat`-Unmarshal:** String `"8"`, `"2.5"`, `"0.00"`, `""` (→ 0), Zahl `8`, Zahl `2.5`,
  sowie ein ungültiger Wert (`"abc"` → Fehler).
- **`Task`-Decoding:** In einem `getTask`/`getAllTasks`-Fixture die Felder `time_estimated` und
  `time_spent` mitgeben und prüfen, dass `TimeEstimated`/`TimeSpent` korrekt befüllt sind — inkl.
  eines Falls mit fehlenden Feldern (→ 0) für Abwärtskompatibilität.

## Nicht im Scope

- **Schreibsupport:** `CreateTaskRequest`/`UpdateTaskRequest` bleiben unverändert. (Kanboard erlaubt
  `time_estimated`/`time_spent` beim Schreiben; hqcli braucht das aber nicht. → offene Frage 1.)
- Die Report-Aggregation selbst — die bleibt in hqcli.

## Abnahmekriterien

1. `mage Test` und `mage Lint` grün.
2. Neuer `StringFloat`-Typ mit Tests wie in R3.
3. `Task.TimeEstimated`/`Task.TimeSpent` werden aus echten/gefixten `getTask`-Responses korrekt
   befüllt; fehlende Felder ergeben `0` (keine Fehler, keine Breaking Changes an bestehenden Feldern).

## Zielprojekt-Workflow (aus dessen `AGENTS.md`)

- Vorgehen: Plan-Mode → implementieren → **Operator testet manuell** → **erst dann** committen.
- Conventional Commits, English: `feat: add time_estimated/time_spent to Task` (+ ggf. separater
  `feat: add StringFloat type` — „one commit per concern"). Kein Co-Author.
- **CHANGELOG:** Eintrag unter **Added** (Datei `CHANGELOG.md` existiert noch nicht → beim ersten
  user-facing Change anlegen, Format Keep a Changelog). Pflege via `/update-changelog`.
- **Release:** `feat:` ⇒ MINOR-Bump ⇒ neues Tag **`v1.6.0`** via `/tag-version`.

## Nach dem Release → zurück in hqcli

- In hqcli `go.mod`/`go.sum` auf `kanboard-api v1.6.0` aktualisieren (`go get -u code.beautifulmachines.dev/jakoubek/kanboard-api@v1.6.0`, dann `mage tidy`).
- Danach die Phase-1-Tasks „Schätz-Reports auf API portieren" und „`task show`/`export-tickets` auf
  reine API" umsetzen (siehe `docs/hqcli-konsolidierung.md`).

## Offene Fragen

1. Sollen `time_estimated`/`time_spent` auch **schreibbar** werden (in `CreateTaskRequest`/
   `UpdateTaskRequest`)? Vorschlag: nein, außerhalb des aktuellen Bedarfs — bei Bedarf separate
   Anforderung.
2. `StringFloat` **nur** `UnmarshalJSON` (Lesen) — ausreichend? Vorschlag: ja.
