# Anforderung: Link-Typen (`getAllLinks`) in der `kanboard-api`-Bibliothek

> **Umsetzung im externen Projekt** `~/Dev/kanboard-api` (Modul
> `code.beautifulmachines.dev/jakoubek/kanboard-api`), **nicht** in hqcli. Dieses Dokument ist die
> Vorlage, die in einer eigenen Session im dortigen Projekt abgearbeitet wird (Plan-Mode, manueller
> Test, dann Commit/Release gemäß dessen `AGENTS.md`).

## Ziel

Die Bibliothek soll die in einer Kanboard-Instanz **konfigurierten Link-Typen** auslesen können,
damit hqcli beim Anlegen einer Ticket-Verknüpfung die passende `link_id` zur Laufzeit auflösen kann
statt sie hart zu kodieren.

## Kontext / Warum

`kb task link add --relation <name>` in hqcli (`internal/kanboard/links.go`) übersetzt Relationsnamen
über die Liste `linkRelations` in numerische IDs:

```
relates 1 · blocks 2 · is_blocked_by 3 · duplicates 4 · is_duplicated_by 5 ·
is_child_of 6 · is_parent_of 7 · targets_milestone 8 · is_milestone_of 9 ·
fixes 10 · is_fixed_by 11
```

Das sind Kanboards **Default-Seed-Werte** der Tabelle `links`. Sie stimmen nur, solange in der
Instanz keine Link-Typen gelöscht, ergänzt oder umsortiert wurden. Bei abweichender Konfiguration
legt hqcli stillschweigend den falschen Beziehungstyp an — ein Fehler, der beim Anlegen nicht
auffällt, weil die API den Aufruf klaglos akzeptiert.

Beim **Lesen** existiert das Problem nicht: `getAllTaskLinks` liefert das Klartext-`Label` mit,
`TaskLink` trägt keine `link_id` (siehe `docs/anforderung-tasklink-mapping.md`).

## Ist-Zustand (Bibliothek)

- `links.go` — nur `GetAllTaskLinks`, `CreateTaskLink(taskID, oppositeTaskID, linkID)`,
  `RemoveTaskLink`. Keine Methode für die Link-**Typen**.
- `jsonrpc.go:49` — `call` ist unexported, ebenso `callObjectOrFalse`. Es gibt **keinen** exportierten
  Escape-Hatch, mit dem sich `getAllLinks` von außen aufrufen ließe. Die Ergänzung muss daher in der
  Bibliothek erfolgen.

## Verifizierte Response

Live per curl gegen die produktive Instanz geprüft (`{"jsonrpc":"2.0","method":"getAllLinks","id":1}`,
keine Parameter). Vollständiges `result` — die Instanz entspricht exakt dem Kanboard-Default-Seed:

```json
[
  {"id": 1,  "label": "relates to",        "opposite_id": 0},
  {"id": 2,  "label": "blocks",            "opposite_id": 3},
  {"id": 3,  "label": "is blocked by",     "opposite_id": 2},
  {"id": 4,  "label": "duplicates",        "opposite_id": 5},
  {"id": 5,  "label": "is duplicated by",  "opposite_id": 4},
  {"id": 6,  "label": "is a child of",     "opposite_id": 7},
  {"id": 7,  "label": "is a parent of",    "opposite_id": 6},
  {"id": 8,  "label": "targets milestone", "opposite_id": 9},
  {"id": 9,  "label": "is a milestone of", "opposite_id": 8},
  {"id": 10, "label": "fixes",             "opposite_id": 11},
  {"id": 11, "label": "is fixed by",       "opposite_id": 10}
]
```

Beobachtungen:

- Genau **drei** Felder, keine weiteren.
- `id` und `opposite_id` kommen als **echte JSON-Zahlen**, nicht als Strings (anders als bei vielen
  anderen Kanboard-Endpunkten).
- Ein Link-Typ ohne Gegenrichtung liefert `opposite_id: 0`, **nicht** `null` — hier „relates to".
- Die Labels sind unabhängig von der UI-Sprache englisch.

## Anforderungen

### R1 — Typ `Link`

```go
// Link represents a task link type as configured in Kanboard's "links" table.
// A link type without an opposite direction (e.g. "relates to") has OppositeID 0.
type Link struct {
	ID         StringInt `json:"id"`
	Label      string    `json:"label"`
	OppositeID StringInt `json:"opposite_id"`
}
```

`StringInt` statt `int`, obwohl die Response Zahlen liefert: Kanboard ist bei der Typisierung
inkonsistent, und `StringInt` deckt beide Varianten ab — konsistent zum Rest der Bibliothek.

### R2 — Methode `GetAllLinks`

```go
func (c *Client) GetAllLinks(ctx context.Context) ([]Link, error)
```

RPC-Methode `getAllLinks`, keine Parameter. Fehlerbehandlung wie bei `GetAllTaskLinks`.

### R3 — optional: `GetOppositeLinkID`

```go
func (c *Client) GetOppositeLinkID(ctx context.Context, linkID int) (int, error)
```

RPC `getOppositeLinkId`, Param `link_id`. Nur ergänzen, wenn es ohne Mehraufwand mitgeht — hqcli
kommt mit `GetAllLinks` allein aus, weil `opposite_id` dort schon enthalten ist.

## Folgeschritt in hqcli (nach dem Lib-Release)

In `internal/kanboard/links.go`:

1. `AddLink` löst den Relationsnamen primär über `GetAllLinks` auf: Label des Kanboard-Typs
   normalisieren (kleinschreiben, „a"/„to" verwerfen, Leerzeichen → `_`) und gegen den übergebenen
   Namen matchen — `"is a milestone of"` → `is_milestone_of`, `"targets milestone"` →
   `targets_milestone`.
2. Ergebnis pro Prozess cachen (ein Call pro CLI-Aufruf genügt).
3. `linkRelations` bleibt als **Fallback**, wenn der Call fehlschlägt oder kein Label matcht — dann
   wie bisher verfahren, aber mit Hinweis auf stderr.
4. Hilfe- und Fehlertexte weiterhin aus `linkRelationNames()` erzeugen; die Namensliste bleibt also
   die stabile Nutzerschnittstelle, nur die ID-Zuordnung wird dynamisch.
