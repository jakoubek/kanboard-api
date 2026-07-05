# Anforderung: `TaskLink`-Mapping an die reale `getAllTaskLinks`-Antwort anpassen

> Quelle: Kanboard-Ticket #4610 ("Bug: linked_task_id ist immer 0 bei Task-Links").
> Umsetzung hier im Lib-Projekt gemäß `AGENTS.md` (Plan-Mode → manueller Test →
> Commit/Release).

## Ziel

`GetAllTaskLinks` soll die ID des tatsächlich verlinkten (gegenüberliegenden) Tasks
korrekt liefern, statt immer `0`.

## Kontext / Warum

Die aktuelle Implementierung (`types.go`, `TaskLink`) geht von folgender JSON-Form aus
und wird von `links_test.go` auch nur gegen diese (erfundene) Form getestet:

```json
{"id": "1", "link_id": "2", "task_id": "42", "opposite_task_id": "100", "label": "...", "title": "..."}
```

Die echte Kanboard-JSON-RPC-Antwort von `getAllTaskLinks` sieht aber völlig anders aus.
Live-Verifikation (read-only, per `curl` direkt gegen den JSON-RPC-Endpunkt von
`kanboard.jakoubek.net`, für die Tickets #4590 und #4580) ergab z. B. für Ticket #4580:

```json
{
  "id": 2355,
  "task_id": 4579,
  "label": "relates to",
  "title": "(Barcode) Zählersuche über Cache + Logging entschlacken + Zähler unabhängig von Sendungsgröße",
  "is_active": 0,
  "project_id": 31,
  "column_id": 162,
  "color_id": "blue",
  "date_completed": 1783036803,
  "date_started": 1782295140,
  "date_due": 0,
  "task_time_spent": 0,
  "task_time_estimated": 0,
  "task_assignee_id": 1,
  "task_assignee_username": "oli",
  "task_assignee_name": "Oliver Jakoubek",
  "column_title": "Erledigt",
  "project_name": "BURG | Kassenprogramm"
}
```

Erkenntnisse:

- Es gibt **kein** Feld `opposite_task_id` → `TaskLink.OppositeTaskID` bleibt beim
  JSON-Unmarshal immer beim Zero-Value `0`. Das ist der im Ticket gemeldete Bug.
- Das Feld `task_id` in der Antwort ist bereits die **ID des verlinkten
  (gegenüberliegenden) Tasks** — nicht die ID des abgefragten Tasks (der wird als
  Request-Parameter mitgegeben, taucht in der Antwort nicht separat auf). Das erklärt
  auch, warum `Title` bisher korrekt ankam: er hängt am selben, nur falsch benannten
  Datensatz.
- Es gibt **kein** Feld `link_id` — nur `label` (Klartext-Relation, z. B.
  `"relates to"`, `"is blocked by"`). `TaskLink.LinkID` war damit ebenfalls immer `0`
  und im Ergebnis nutzlos.
- Kanboard liefert stattdessen viele denormalisierte Zusatzfelder zum verlinkten Task
  (`project_id`, `column_id`, `project_name`, `column_title`, `task_assignee_*`,
  `date_*`, `task_time_*`, `is_active`, `color_id`), die aktuell komplett ignoriert
  werden und außerhalb des Scopes dieser Anforderung bleiben.

Der Fehler sitzt, wie im Ticket vermutet, in dieser Bibliothek — `hqcli`
(`internal/kanboard/links.go`) greift lediglich korrekt auf `l.OppositeTaskID` zu,
bekommt von der Lib aber den falschen (immer `0`) Wert.

## Anforderungen

### R1 — `OppositeTaskID` auf `task_id` mappen
`TaskLink.OppositeTaskID` bekommt den JSON-Tag `task_id` statt `opposite_task_id`.

### R2 — `TaskID`-Feld entfernen
Das bisherige `TaskLink.TaskID`-Feld ergibt in der echten API keinen Sinn (kein
separates Feld für die ID des abgefragten Tasks in der Antwort) und wird entfernt.

### R3 — `LinkID`-Feld entfernen
`TaskLink.LinkID` (JSON-Tag `link_id`) entfernen — die echte Antwort enthält kein
solches Feld, es war nie befüllbar. Betrifft nur das Result-Struct; `CreateTaskLink`
und `RemoveTaskLink` nutzen `link_id`/`task_link_id` weiterhin unverändert als
Request-Parameter.

### R4 — Tests auf reale Response-Struktur umstellen
`links_test.go` (`TestClient_GetAllTaskLinks`, `TestClient_GetAllTaskLinks_Empty`,
`TestTaskScope_GetLinks`) und ggf. `types_test.go`: Mock-Responses im echten
Kanboard-Format aufbauen (inkl. zusätzlicher, im Go-Struct nicht abgebildeter Felder
wie `is_active`, `project_id`, um zu zeigen, dass unbekannte Felder toleriert werden).
Assertions auf `OppositeTaskID` prüfen gegen den Wert des `task_id`-Felds der
Mock-Response.

## Nicht im Scope

- Die zusätzlichen denormalisierten Felder (`project_id`, `column_id`, `project_name`,
  `column_title`, `task_assignee_*`, `date_*`, `task_time_*`) werden in diesem Zug
  **nicht** in `TaskLink` aufgenommen — nur der ID-Fix.

## Abnahmekriterien

1. `mage Test` / `mage Lint` grün.
2. `GetAllTaskLinks` liefert bei einer Mock-Response im echten Kanboard-Format die
   korrekte `OppositeTaskID` (= Wert des `task_id`-Felds der Antwort).
3. Kein Zugriff mehr auf `TaskLink.TaskID`/`TaskLink.LinkID` im Repo (entfernt).

## Release

- Conventional Commit, English: `fix: map linked task id from task_id field in getAllTaskLinks response`.
- CHANGELOG unter **Fixed**; PATCH-Bump.
- Breaking Change am `TaskLink`-Struct (Felder `TaskID`/`LinkID` entfernt, `OppositeTaskID`
  liest jetzt `task_id`) — da `hqcli` der einzige bekannte Konsument ist und die Lib
  lokal per `replace` eingebunden wird, reicht laut bisherigem Vorgehen in diesem Repo
  ein `fix:`-Commit ohne SemVer-MAJOR; im CHANGELOG kurz auf die Struct-Änderung
  hinweisen.
