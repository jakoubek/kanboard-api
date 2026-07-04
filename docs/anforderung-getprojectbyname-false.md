# Anforderung: `GetProjectByName`/`GetProjectByID` — Kanboards `false`-Antwort sauber behandeln

> Quelle: gefunden bei der hqcli-Konsolidierung (Kanboard API-only). Umsetzung hier im Lib-Projekt
> gemäß `AGENTS.md` (Plan-Mode → manueller Test → Commit/Release).

## Ziel

`GetProjectByName` (und ggf. `GetProjectByID`) sollen bei „kein Treffer" einen sauberen
`ErrProjectNotFound` zurückgeben — so wie `GetTask` es für sein „nicht gefunden" bereits tut — statt
mit einem JSON-Unmarshal-Fehler zu scheitern.

## Kontext / Warum

Kanboards JSON-RPC gibt bei `getProjectByName` (und teils `getProjectById`) im Fall „nicht gefunden"
den Wert **`false`** zurück (nicht `null`). Die aktuelle Implementierung (`projects.go`)
deserialisiert das Ergebnis direkt in `*Project`:

```go
var result *Project
if err := c.call(ctx, "getProjectByName", params, &result); err != nil {
    return nil, fmt.Errorf("getProjectByName: %w", err)
}
if result == nil { return nil, fmt.Errorf("%w: project %q", ErrProjectNotFound, name) }
```

`json.Unmarshal(false, *Project)` schlägt fehl → der Aufrufer bekommt:
`getProjectByName: failed to unmarshal result: json: cannot unmarshal bool into Go value of type kanboard.Project`.

Der `result == nil`-Check greift nur bei `null`, nicht bei `false`. Zum Vergleich funktioniert
`GetTask` sauber, weil `getTask` bei „nicht gefunden" `null` liefert (→ `result == nil`). Das
`false`-Muster kennt die Lib bereits von Create-Operationen (`IntOrFalse` in `types.go`) — nur für
Struct-Ergebnisse fehlt das Pendant.

## Anforderungen

### R1 — `false`/`null` als „nicht gefunden" behandeln
`GetProjectByName` so anpassen, dass eine `false`- **oder** `null`-Antwort zu `ErrProjectNotFound`
führt, nicht zu einem Unmarshal-Fehler. Umsetzungsvorschlag: Ergebnis zuerst in `json.RawMessage`
(bzw. `any`) aufnehmen, auf `false`/`null` prüfen, sonst in `*Project` deserialisieren. Alternativ ein
wiederverwendbarer Decode-Helper analog `IntOrFalse` („struct-or-false").

### R2 — `GetProjectByID` gleich mitprüfen
Prüfen, ob `getProjectById` dasselbe `false`-Verhalten zeigt; falls ja, identisch behandeln. Andere
`GetX`-Methoden, die ein Objekt-oder-`false` liefern können, ggf. mit demselben Helper absichern.

### R3 — Tests
Analog zu `projects_test.go`: je ein Fall, in dem der Server `false` bzw. `null` liefert → Methode
gibt `ErrProjectNotFound` zurück (kein Unmarshal-Fehler); der bestehende Erfolgsfall (Objekt) bleibt
grün.

## Abnahmekriterien
1. `mage Test` / `mage Lint` grün.
2. `GetProjectByName` mit `false`- und `null`-Response → `errors.Is(err, ErrProjectNotFound)` == true.
3. Erfolgsfall (echtes Projekt) unverändert.

## Release
- Conventional Commit, English: `fix: return ErrProjectNotFound on false/null from getProjectByName`.
- CHANGELOG unter **Fixed**; `fix:` ⇒ PATCH-Bump (→ v1.6.1).
- hqcli nutzt die Lib lokal via `replace` → kein `go get`-Bump nötig.
