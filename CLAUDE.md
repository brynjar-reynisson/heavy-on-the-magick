# Heavy on the Magick — reverse-engineering & Go port

GitHub repo: `brynjar-reynisson/heavy-on-the-magick`. **Long-term goal**:
fully disassemble the 1986 ZX Spectrum game **"Heavy on the Magick"**
(Gargoyle Games / Carter Follis Software — Roy Carter, Greg Follis) and port
it to **Go**, then extend it with a new feature not in the original: an
in-game map showing rooms the player has explored so far. This file is the
running log of reverse-engineering progress; the port itself hasn't started
yet (still in the disassembly/understanding phase).

The game: a graphic adventure starring the wizard "Axil the Able". Occult-
themed (Crowley/Golden Dawn references throughout: grades like Neophyte →
Zelator → Practicus → Philosophus → Adeptus Minor/Major/Exemptus → Magister
Templi → Ipsissimus; demons such as "Belezbar" = Beelzebub). Command syntax
in-game is `CHARACTER, VERB` (e.g. `APEX, DOOR`, `APEX, FIRE`) — "Apex" is
an in-game mentor NPC who gives hints.

## Directory layout

- `hotm1.tzx` / `hotm1.szx` — **misleadingly named**: this is actually a Fuse
  **SZX snapshot** (magic bytes `ZXST`, `CRTR` chunk says `Fuse`/`libspectrum`),
  not a real tape image, despite the `.tzx` extension. Fuse loads it fine
  because it sniffs content instead of trusting the extension; SpecEmu/Zero
  do not. Load it as a snapshot (`.szx`), not as a tape.
- `hotm-original/` — the genuine tape images, downloaded from
  [Spectrum Computing](https://spectrumcomputing.co.uk/entry/2274/ZX-Spectrum/Heavy_on_the_Magick):
  `Heavy On The Magick - Side 1.tzx` and `- Side 2.tzx` (real `ZXTape!` magic,
  confirmed via `xxd`). A "bugfix" release also exists there (fixes corrupted
  graphics + a memory-corruption bug from long text input) — not yet pulled.
- `hotm.z80`, `hotm-live.z80`, `hotm-unpacked.z80`, `hotm-unpacked.mem`,
  `*.ctl`, `*.skool`, `*.png` — SkoolKit analysis artifacts, see below.

**Not in this repo** (general-purpose tools, live one level up in
`..\` i.e. `C:\Users\Lenovo\Documents\spectrum\`, not game-specific and not
committed here — closed-source/binary licensing):
- `..\SpecEmu\` — working SpecEmu install (v3.4, build b250426 / April 2026),
  downloaded from https://specemu.zxe.io/. **This is the emulator to use.**
- `..\Zero\` — backup/alternative emulator (Zero 0.8, MIT-licensed, Windows-
  native C#, from github.com/ArjunNair/Zero-Emulator). Works for loading, but
  its `Alt+Enter` full-screen toggle throws an error. Not currently primary.

## Running the game (SpecEmu)

1. Launch `..\SpecEmu\SpecEmu.exe` (i.e. `C:\Users\Lenovo\Documents\spectrum\SpecEmu\SpecEmu.exe`).
2. **File → Open** → pick `hotm1.szx` (in this repo dir). This loads instantly
   as a full machine-state snapshot (no `LOAD ""` needed).
3. Toggle full screen from the View menu / toolbar — confirmed working from
   a snapshot load.

**Known bug — do not use real tape loading in SpecEmu (this install/machine):**
typing `LOAD ""` and playing a real tape (or `.tzx`) results in a broken
display: the border renders correctly (real ULA colors show, confirmed via
`PrintWindow`/`PW_RENDERFULLCONTENT`, so it's not a compositing artifact) but
the main 256×192 screen bitmap area stays a static gray block and never
updates, even across Reset. Root cause not identified (suspected blitter
issue in this build's rendering path, specific to the accelerated/flashload
code path). Workaround: only load `.szx`/snapshot files in SpecEmu; use
SkoolKit's `tap2sna.py` (see below) for anything that needs the real tape.

**Avoid these SpecEmu builds** — both triggered a real-time Defender
detection (`Trojan:Win32/Suschil!rfn`, heuristic, threat ID 2147927547) and
were auto-quarantined: `specemu-3.4.b310825.zip` (Aug 2025) and
`specemu-3.4.b010424.zip` (Apr 2024). The currently-installed b250426 build
has never been flagged. Untested whether other builds are affected.

## Toolchain: SkoolKit

Installed via `python3 -m pip install skoolkit` (v10.1). **Console scripts
are not on PATH** — invoke them directly with Python:

```
SCRIPTS="C:\Users\Lenovo\AppData\Local\Packages\PythonSoftwareFoundation.Python.3.13_qbz5n2kfra8p0\LocalCache\local-packages\Python313\Scripts"
python3 "$SCRIPTS\<tool>.py" [args]
```

Key tools and the workflow used so far:

- `tapinfo.py <tzx>` — inspect tape block structure (headers, load addresses,
  loader type). Revealed this tape uses a custom **"Gargoyle"** turbo loader:
  a small loader stub (header `L`, loaded at `47000`, 820 bytes) followed by
  a 49122-byte turbo-speed data block (the main game), with pulse timings
  ~35% faster than the ROM standard (`0-pulse:634, 1-pulse:1268` vs the ROM's
  `855`/`1710`). `tap2sna.py`'s built-in accelerator list even has a
  `gargoyle2` entry for this exact scheme.
- `tap2sna.py [-c finish-tape=1] [--start ADDR] <tzx> <out.z80>` — **the key
  technique**: this doesn't just parse tape blocks, it *simulates real CPU
  execution* of the load, so it correctly handles custom loaders. Without
  `--start`, it stops the instant the tape ends (usually mid-loader-routine).
  Passing `--start ADDR` lets the simulation **keep running past the end of
  the tape** until the program counter reaches `ADDR` — this is how we walked
  forward from "tape just finished loading" to "actual game code is now
  running" to "picture-unpacking loop has finished", snapshotting memory at
  each meaningful point.
- `sna2ctl.py <z80> > file.ctl` — auto-generates a control file (code/data/text
  block guesses) via static heuristics. **It gets screen memory
  (`16384`-`23295`) wrong** (tries to disassemble it as code) since it doesn't
  know the fixed ZX Spectrum memory map — always override that range to a
  single `b 16384` data block before disassembling. It also frequently
  misclassifies code as data right after data-block boundaries (had to
  manually insert `c ADDR` markers at `24100` and `46193` to get real
  disassembly instead of garbage — confirmed those addresses were really code
  by cross-checking against where `tap2sna --start` actually stopped).
- `sna2skool.py -c file.ctl <z80> > file.skool` — produces the disassembly.
- `skool2asm.py file.skool > file.asm` — readable labeled assembly listing.
- `sna2img.py [-s SCALE] <z80|bin> <out.png>` — renders screen memory (or raw
  binary, with `-B -O origin`) as a real PNG. Used to visually confirm
  snapshot state and to test graphics-table decoding theories.

## Boot sequence trace (addresses confirmed by simulation, not guessed)

1. BASIC loader ("MAGIC", line 10) loads the `L` stub at `47000`–`47820`
   (standard ROM-speed block).
2. Stub entry `47007`: `DI` / `LD SP,24099` / `LD IX,16416` (load destination
   = start of screen memory!) / `LD DE,49120` (length) / `LD A,153` (custom
   flag byte) / `SCF` / `JP 47086` — a hand-written replacement for the ROM's
   `LD-BYTES`, reading the turbo-timed main data block.
3. On success, `47024`: `LD SP,24099` / `EI` / `JP 24100` — hands off to the
   real game.
4. `24100`: `JP 46193` (tiny stub, confirmed via SkoolKit's own
   "used by the routine at #R47024" annotation).
5. `46193` — main init routine:
   - Clears a 32-byte workspace at `23698`.
   - **Relocates a pointer table at `48054`**: 4-byte entries
     (2-byte pointer + 2 size/flag bytes), adds offset `49080` to every
     pointer (classic post-load address fixup).
   - For each non-empty entry, computes `size = (byte2 & 0x7F) * 8 * byte3`
     and runs a bit-shuffle unpack loop (`RLC (HL)` + `RLA`/`RLCA` chains)
     over `size` bytes at the (now-fixed) pointer, in place.
   - Afterwards prints a startup hint screen (text at `46382`+: `"SOME ADVICE"`,
     `"Talk to Apex often!"`, `"APEX, DOOR"`, `"APEX, WEREWOLF"`, `"APEX, FIRE"`)
     via what look like text-print routines at `44074`/`44103` (not yet
     confirmed by tracing those routines directly).

## Graphics table at 48054 — status

Table entries are mostly a uniform size (`(5 & 0x7F)*8*3 = 120` bytes,
matching the pointer stride exactly), except a cluster around entries 26-30
with much smaller sizes (`8`, `8`, `8`, `16`, `16` bytes).

- **Confirmed**: rendering the small entries as plain 8×8 monochrome tiles
  (1 bit per pixel, standard Spectrum character-cell format) produces clean,
  real, recognizable glyphs — rune/sigil-like symbols, plus one checkerboard
  dither-fill tile. This is almost certainly a **custom magic-symbol font**
  used for spellcasting sigils. Render script approach: read N×8 bytes
  starting at the entry's pointer, each 8-byte group is one row-major 8×8
  character (bit 7 = leftmost pixel).

- **Other code sites that reference the table** (found by raw byte-scanning
  the memory dump for `48054`/`49080` as little-endian immediates, then
  fixing the resulting off-by-one — the match lands on the low address byte,
  so the real opcode starts 1 byte earlier):
  - **`31932`** — a *dynamic single-glyph selector*: takes an input value in
    `E`, masks it to `0-15` (`AND 15`), multiplies by 8, adds to a base
    pointer, and **overwrites table entry 0's pointer field with the
    result**, forcing entry 0's size fields to `257` (`0x0101` → count=1,
    row_count=1 → 8 bytes, a single character). Reads `(IX+5)` as an input
    parameter right after. Best guess: this is how the game picks **1 of 16
    rune symbols** to display (e.g. during spell-casting UI), using entry 0
    of the table as a reusable "current symbol" slot.
  - **`42612`** — trivial helper: given an index (via `HL` set by caller),
    copies the raw 4-byte table entry to a fixed scratch address `44924`
    (`LDI`×3 + one more byte). Just a "load descriptor into fixed work vars"
    convenience, used by other routines (per SkoolKit's xref: called from
    the routines at `31435` and `40218` — not yet traced).
  - **`41631`** — the real *picture-drawing/transform* routine (see below).

- **The draw routine at `41631`–`41716`, decoded instruction-by-instruction**:
  fetches the pointer and size bytes from the table entry (`count = byte2 &
  0x7F`, `row_count = byte3` after `XOR 128` toggles its top bit and writes
  it back — i.e. **every draw call flips a persistent flag bit on the table
  entry itself**, which strongly suggests an animation/blink toggle between
  two states), then runs a nested loop: outer loop over `count` "columns",
  each doing 8 sub-passes (one per pixel-row of a character cell), each
  sub-pass bit-interleaving `ceil(row_count/2)` byte-pairs between a
  forward-walking cursor (`DE`, starting at the entry pointer) and a
  backward-walking cursor (`HL = DE + row_count`), via a chain of 8 rounds
  of `RLA`/`RR C` (rotating bits from one byte into the other through the
  carry flag). Net effect: `DE` advances by exactly `row_count` bytes per
  column, so `count × row_count × 8 = 120` bytes get touched per entry —
  correct total, but the actual bit-level semantics of the interleave aren't
  proven correct yet (see below).

- **Attempted verification — inconclusive**: transliterated the `41663`–
  `41716` loop faithfully into Python (register-accurate, carry-flag
  threaded through all 8 rotate rounds, persisting *across* the
  `half_count` inner iterations since there's only one `AND A`
  carry-clear per outer-B pass, not per inner iteration) and ran it against
  entry 1's raw bytes (`count=5, row_count=3`, pointer `49080`) from the
  post-init-but-pre-any-draw-call snapshot (`hotm-unpacked.z80`/`.mem`).
  Result still renders as noise under every tile-grid layout tried. Open
  possibilities, not yet distinguished:
  1. A subtle bug in the carry/rotate transliteration.
  2. A missing preprocessing step (e.g. this might be a color/attribute
     plane rather than monochrome bitmap data, or needs another pass we
     haven't found).
  3. Entry 1 isn't actually a standalone renderable picture — it shares its
     pointer (`49080`) with entry 0, which is the *dynamically-repurposed*
     single-glyph slot (see `31932` above), so entry 1's boot-time content
     may just be placeholder/default data rather than a real room picture.
  Script: `render_tile.py` / `sim_transform.py` pattern in scratchpad
  (transliteration + render, not yet copied into the spectrum dir).

### Dispatch mechanism found

`41631` has no direct callers because it isn't the real entry point — it's
reached by **falling through** from `41627` (`LD H,0` / `ADD HL,HL` ×2, i.e.
`HL = index*4`, immediately followed by the code we'd already decoded).
`41627`'s own caller (per SkoolKit's xref) is a big, heavily-shared routine
at **`41778`, called from 15 different places** (`31699, 31816, 31932,
33668, 38336, 38574, 38631, 39669, 40694, 40714, 40741, 40955, 41393, 41410,
41463`) — including `31932`, the rune-selector we'd already decoded.
`41778` takes a `BC` parameter (**two indices, one in `B` one in `C`**), and
for each one: looks up an object record via `CALL 42608`, does a conditional
bit-test (`XOR (HL)` / `JP P,...`), and **only if the test passes**, calls
the picture-drawer `41627` with that index. This looks like a generic
"check object state, conditionally draw its icon" utility used throughout
the game (inventory/spell-component/monster display etc.), not something
room-picture-specific — which reframes the whole table at `48054` as
probably a general **object/icon glyph table**, not a room-background table.
The 120-byte entries may therefore be a genuinely different packed sprite
format from the small glyphs, rather than just "bigger versions" of the
same thing — worth keeping in mind before assuming they should decode the
same way.

### Traced one caller fully: 31932 (rune selector)

Full trace of `31932` (the "pick 1 of 16 runes" routine, called from a
context we haven't traced further up yet):
- `E & 15` selects a glyph index (0-15), `×8` to get a byte offset into a
  16-glyph font strip whose base address comes from the table entry HL
  pointed to on entry (i.e. the caller already passed an entry index in HL,
  pre-multiplied — unlike `41627`/`42608`, `31932` does NOT do its own
  `×4` on HL, so its calling convention differs from those).
- Overwrites table entry 0 with `{pointer = font_base + offset, count=1,
  row_count=1}` — repurposing entry 0 as a "draw this one specific glyph"
  descriptor, confirming the theory from before.
- Sets `B=C=0` (so `41778` processes index 0 for *both* of its checks),
  `E=(IX+5)`, `D=16`, `A=65`, then `CALL 41778` — meaning this path always
  routes through the shared check-and-draw dispatcher even for the
  simplest single-glyph case.

**Verified the transform is non-trivial even for `count=1, row_count=1`**:
ran the Python transliteration against entry 26's raw 8 bytes (`ptr=51944`)
with `count=1, row_count=1` — one byte changed (`0x52` → `0x4A` at offset 4),
others stayed fixed (the zero bytes are trivial fixed points of this
particular bit permutation: `0x00` self-merges to `0x00`). So:
- The draw-time transform **does** meaningfully alter data even in the
  simplest case — it is not skippable/no-op logic.
- This means our earlier "clean glyph" renders (raw, untransformed bytes
  for entries 26/27/29) are **not confirmed to be pixel-accurate** to what
  actually appears on screen — they looked coherent, but we don't know if
  those specific entries were ever actually drawn via this exact code path,
  or whether a real draw would tweak a row or two via this transform.
- No ground-truth screenshot exists yet to check either version against.

Status: genuinely uncertain past this point without either (a) finding a
real reference image of the game running to compare against, or (b) tracing
further up into one of the 15 callers of `41778` in enough contextual depth
to know for certain which table index corresponds to a specific, nameable
in-game object.

### Found ground truth: `ref_screenshot.png`

Downloaded a real screenshot from the Internet Archive item
[zx_Heavy_on_the_Magick_1986_Gargoyle_Games](https://archive.org/details/zx_Heavy_on_the_Magick_1986_Gargoyle_Games)
(`ref_screenshot.gif`/`.png` in this dir) — it's the exact "SOME ADVICE"
startup hint screen we traced. **Text matches our decoded strings exactly**
(plus more we hadn't reached yet: `"To dismiss say \"APEX, THANKS\""`,
`"Talk to other things - sometimes they will answer!"`,
`"\"DOOR, password\""`, `"\"GUARDS, DOOR\""`, `"\"ASTAROT, WOLFDORP\""`,
`"Finally, if in a panic, BLAST without an object!"`) — strong validation
that the boot-sequence trace is accurate. The screenshot also shows a
**small graphic of Axil sitting at a table** in the bottom-right corner,
and a **decorative chain-link border** around all four edges — plausible
candidates for what the picture-table entries actually render as.

**Blocked**: tried extending the `tap2sna --start` simulation past the
hint-text-printing code (guessed resume point `46683`, right after the
message data block) to reach this exact screen and render/diff it against
the reference. Simulation instead timed out at an unrelated PC (`64733`),
almost certainly stuck in a "wait for keypress" loop after printing —
`tap2sna` doesn't have an obvious option to inject a keypress at an
arbitrary mid-execution point (`--press` is for pausing tape playback, not
general execution). Have not yet found a way around this.

### Live ground-truth confirmation via SpecEmu's Debugger

Loaded `hotm.z80` in SpecEmu (its debugger is accessible via **Monitor** menu
→ opens a separate "Debugger" tool window with live disassembly, a
hex/memory dump toggle (**Dump** button), and **Edit → Go To...** to jump to
any address). Let it run freely (unlike `tap2sna`, a real interactive
emulator doesn't get stuck on the post-boot "wait for keypress" — it just
free-runs), producing the exact same "SOME ADVICE" screen as the Internet
Archive reference screenshot, then continuing into the actual game (a
"ritual circle" idle/spell-selection screen with colored icon squares —
frame count climbing confirms it's genuinely running, not stuck).

**Navigated to `48054` in the live Debugger memory dump and confirmed our
static table analysis exactly**: entry 0 = `ptr=49080(0xBFB8), flags=0,0`;
entry 1 = `ptr=49080, flags=133,3`; entry 2 = `ptr=49200(0xC030),
flags=133,3` — bit-for-bit identical to what `sna2ctl`/manual analysis
predicted from the static snapshot. This is strong confirmation the table
structure and pointer-relocation understanding is correct, from genuine
live gameplay memory, not just a simulated stop-point.

**Automating SpecEmu's UI via PowerShell — working technique, for next
time**: `SendKeys`/mouse automation against SpecEmu is fragile in specific,
learnable ways:
- Each `PowerShell` tool call is a fresh process — `Add-Type` type
  definitions do NOT persist between calls. Every script must be fully
  self-contained (redefine types every time).
- Dropdown menus and dialogs (e.g. the "Go To" dialog, error popups) are
  **separate top-level windows**, not part of the main window — `PrintWindow`
  on the main window's handle won't capture them. Enumerate windows by PID
  (`EnumWindows` + `GetWindowThreadProcessId`) to find the real HWND of a
  popup/dialog by its title.
- **Critical gotcha**: after a dialog opens, keep targeting its OWN HWND for
  `SetForegroundWindow` before sending further keys — calling
  `SetForegroundWindow` on the PARENT window's handle instead silently
  causes subsequent `SendKeys` to go nowhere useful (text stops appearing in
  the field) even though the dialog still looks focused. Symptom: typed
  text just doesn't appear, or a field mysteriously won't clear/update.
  Fix: `EnumWindows` again to get the dialog's specific HWND, target that.
- Mouse clicks via `SetCursorPos`+`mouse_event` at computed pixel coordinates
  were **unreliable** for hitting menu bar items and dialog buttons in this
  session (possibly DPI-scaling related) — keyboard-driven navigation
  (`Alt+<mnemonic>` for menus, arrow keys/Enter for items, `Ctrl+A`/`Delete`
  to clear a field before retyping) worked much more reliably throughout.
  Prefer keyboard automation over coordinate-based mouse clicks for this app.
- Don't use PowerShell variable name `$pid` — it collides with the
  automatic/reserved `$PID` variable and throws a "read-only variable" error.
- The variable `$rect`/type names must be unique per script if reusing
  `Add-Type -TypeDefinition` in the same session (redefinition errors) —
  simplest fix is just incrementing a suffix letter each time (Win32a,
  Win32b, ... as done throughout this session).

### Got real ground-truth pixel data for the magenta icon — still doesn't resolve the table format

Using the live Debugger (paused, stable, 2471 frames in), located the magenta
and green spell-icon squares precisely:
- **Magenta**: attribute address `22892` (row 11, col 12), value `0x43`
  (bright magenta, black paper).
- **Green**: attribute address `22900` (row 11, col 20), value `0x44`
  (bright green, black paper).
- Computed real bitmap addresses via the standard ZX Spectrum non-linear
  screen formula (`addr = 16384 + (Y&192)*32 + (Y&7)*256 + ((Y&56)>>2) + X`,
  stride `256` between consecutive pixel-rows within one character cell) and
  read all 8 bytes of each cell directly from live memory:
  - Magenta (`18540`+256×n): `F8 88 A8 A9 88 F8 91 00`
  - Green (`18548`+256×n): all zero (`00` ×8) — plausibly caught mid-blink,
    given the draw routine toggles an animation flag on every call (see
    above) — not necessarily a reading error.

**Searched the entire picture table (all ~60 entries) for this exact
8-byte magenta pattern, both as raw stored bytes AND after applying our
`run_transform_v2` transliteration — no match either way.** This means one
of:
1. This specific icon isn't drawn from the `48054` table at all (maybe
   procedurally drawn — e.g. a generic "fill rectangle" routine — rather
   than blitted from stored bitmap data).
2. It's read from a completely different table/font region we haven't
   located.
3. Our transform transliteration still isn't correct, and coincidentally
   doesn't match any entry either way.

**Status**: the 120-byte "big picture" format is still uncracked, and now
we also don't have a confirmed source location for these specific spell-icon
glyphs either, despite having exact ground-truth pixel data for one. This
was a productive dead end — we now have a real on-screen byte pattern
(`F888A8A988F89100` for magenta) that any future theory about the
graphics format needs to be checked against.

### Correction: `44074`/`44103` are NOT a print routine — they set up a text window

Earlier trace guessed `CALL 44074; CALL 44103; LD DE,<addr>` printed a
message at `<addr>`. **Wrong** — fully disassembled both routines now:

- `44074`: takes a window-slot index in `A` (clamped to `0-3`), computes
  `HL = 43978 + A*11` (an 11-byte-stride table of up to 4 window/viewport
  records) and stores that pointer at `44023`.
- `44103`: pops its own return address into `HL`, then does
  `DE=(44023); BC=8; LDIR` — meaning **it consumes the 8 bytes immediately
  following the `CALL 44103` instruction as inline parameter data**, copying
  them into the selected window record, then computes
  `record[2] = record[2]-record[0]+1` and `record[3] = (record[3]-record[1])*2`
  in place (looks like converting two corner coordinates into a width/height
  pair, though the exact units aren't confirmed), zeroes `record[10]`, and
  jumps to `44267` (not yet traced).

This means the bytes right after each `CALL 44103` are **data, not code** —
our earlier linear disassembly was misreading them as instructions (e.g.
`LD DE,5633` was actually 3 of the 8 raw parameter bytes). Confirmed by
extracting the real bytes:

```
window slot 0 (A=0 at the call site): 11 01 16 08 05 00 00 00
window slot 1 (A=1 at the call site): 11 0a 16 15 07 00 01 01
```

Both start `11 xx 16 ..` — `x0=17(0x11), x2=22(0x16)` match on both,
suggesting a shared left/right screen extent (consistent with both being
sub-windows of the same bordered box in the reference screenshot), while
`y0`/`y1`/the trailing bytes differ per window. This is a genuine
**multi-window text/graphics system**, not a simple print call — worth
keeping in mind: finding the actual string-print routine (whatever `44267`
eventually calls) is still open.

### Found the character-output primitive and the AT-code convention

Traced further from `44267`:
- **`44286`** is the core character-output routine — `putchar(A)` in effect —
  called from **27 different places** throughout the game. It preserves all
  registers, loads `IX=(44023)` (the current text window record set up by
  `44074`/`44103`), and dispatches into `44299`.
- **`44299`**: `A<32` → jump to `44569` (control-code handling, not fully
  traced — but `44569` starts with `CP 13` i.e. handling `CHR 13` /
  carriage-return specially, consistent with a real text-output routine).
  `A==32` (space) → suppress trailing spaces at the window's right edge
  (`(IX+10)` vs `(IX+9)`, i.e. cursor-column vs. window-width). Otherwise →
  `CALL 44455`, the actual glyph plotter (not yet fully traced, but per its
  first two calls — `44326` then `44412` — it looks like "advance cursor /
  plot the character", with a `(IX+7)` check afterward for something that
  substitutes character code `126` — possibly a cursor/prompt symbol).
- **`44267`** (the routine we were originally chasing) turns out to be a
  convenience wrapper: `LD A,22 / CALL 44286` then prints `B` and `C` through
  the same primitive. **Control code 22 = "AT row,col"**, exactly like real
  ZX Spectrum BASIC's own `CHR$ 22` convention — this confirms the
  `DEFB 22,row,col` triples seen before every hint-screen message (e.g.
  `22,4,2` before `"Talk to Apex often!"`) are literally AT-positioning
  commands consumed by this same custom print engine, not arbitrary data.
- Also found (not fully traced): **`44482`**, called in a loop from `44530`
  based on `(IX+2)` — its shape (paired `LDIR`s advancing `H`/`D` by 1 each
  outer iteration, sized by `(IX+3)`) strongly resembles a **text-window
  scroll-up routine** (move each screen line up by one when the window
  fills), but this isn't confirmed yet.

The actual string-print loop (something that calls `44286` once per
character of a message) still hasn't been pinned down — worth finding, since
it's the natural place to intercept for extracting ALL of the game's
message/room-description text in bulk rather than one string at a time.

### Found the full parser vocabulary table

Searched decompressed memory for `NORTH`/`SOUTH`/`EAST`/`WEST` as literal
text (accounting for the high-bit-terminated-last-letter word encoding used
throughout this game's strings — confirmed all the way back in the very
first SZX investigation). Found a **flat, ~316-word vocabulary table** at
memory addresses `24270`-`26200`ish, organized in **fixed-length buckets**
(all 3-letter words, then all 4-letter words, then 5-letter, ... up to
11-letter) — a classic 8-bit parser space/speed optimization.

Confirmed the **8 compass directions**, with exact spelling: `NORTH`,
`SOUTH`, `EAST`, `WEST` (plain, standalone, found separately from the
diagonals) and `NORTH-EAST`, `SOUTH-EAST`, `SOUTH-WEST`, `NORTH-WEST`
(hyphenated, found together as a cluster) — matches the "8 compass
directions" description from contemporary reviews exactly.

The rest of the table is ~300 more words mixing verbs, nouns, and item/
creature names (`APEX`, `AXIL`, `DOOR`, `WOLF`, `OGRE`, `FIRE`, `TALK`,
`CALL`, `HELP`, `TAKE`, `DROP`, `LIFT`, `SPELL`, `INVOKE`, `MAGUS`, `SKULL`,
`DEMON`, `SWORD`, `PARCHMENT`, `SALAMANDER`, `WRAITHVALE`, `PHILOSOPHUS`,
...) with **no yet-known type/ID tag per word** — we have the full word
list but not which are verbs vs. nouns vs. something else, nor what each
one does. That mapping (presumably a parallel ID/code table, or codes
embedded via some other mechanism not yet found) is the next real unlock
for gameplay logic.

**Ported into the Go project already** (see `internal/parser/vocabulary.go`
+ `internal/parser/data/vocabulary.json`, `internal/world/direction_words.go`):
the confirmed word list is embedded via `go:embed` and exposed as
`parser.Vocabulary` / `parser.KnownWord()`, and the 8 confirmed direction
words are wired to `world.Direction` via `world.ParseDirection()`. Both are
covered by tests. This is real extracted game data, not invented — but only
the direction words have confirmed *meaning*; the rest are just "the parser
recognizes this word," not yet "and here's what it does."

**Likely found the bucket-index table** (not fully confirmed): searching for
code referencing the `WEST`/`NORTH` string addresses found a table at
`24103`-`24135` (17 sequential 2-byte pointers) whose values are
**ascending addresses landing almost exactly on each word-length bucket
boundary** (`24282`=start of the 4-letter bucket where `WEST` sits,
`24646`=start of the 5-letter bucket where `NORTH` sits, `25016`, `25388`,
`25640`, `25816`, `25933`, `26043`, ... matching the 6/7/8/9/10/11-letter
bucket starts). Reading `WEST`/`NORTH` as "referenced by code" was slightly
misleading — they're not individually referenced as directions, they just
happen to be the first word in their respective length-buckets, which the
parser's bucket-index table points at. Plausible mechanism: the parser
takes the typed word's length, looks up `bucket_table[length]`, and scans
only that fixed-width bucket — a natural way to speed up matching against
a flat vocabulary this size on a 3.5MHz CPU. Exact indexing arithmetic
(what value maps to offset 0 of the table) not yet nailed down.

**Attempted, inconclusive**: looked for the room/exit table by searching for
code referencing the `EXIT` string (found one reference at `29545`, in what
looks like a data table of pointers rather than a code call site — didn't
resolve to a clean structure in the time spent, unlike the vocabulary
bucket-table find above). Room/exit data is still the single biggest
missing piece for gameplay porting — needed for real movement AND the
explored-map feature. Next attempt should probably work forward from the
still-untraced main game loop (after the boot sequence's "wait for
keypress" point) rather than backward from string references.

### Found the beeper sound routine AND a confirmed musical pitch table

The address `64733`, where `tap2sna`'s simulation previously appeared to
"time out" waiting for a keypress (see the earlier "Blocked" section above)
is actually inside a **real beeper sound routine** (`64671`-`64781`,
`OUT (254)` toggling + `DJNZ` timing + shadow-register alternation for the
square wave). The "wait for keypress" read was half-right: the loop that
calls it (`64613`-`64623`) does `CALL 64671` (play a bit of tone) then
`CALL 654` (the **real ZX Spectrum ROM's `KEY-SCAN` routine** — a ROM call
inside otherwise custom code, easy to miss) in a loop — i.e. it plays a
tune *while* waiting for a key, not purely waiting.

Found and extracted **`PitchTable`** at address `64812`: 53 bytes, each a
beeper-loop period value, where consecutive entries have an average ratio
of 1.0606 — matching true 12-tone-equal-temperament's 2^(1/12)=1.0595 to
within 0.1%. This is unambiguously a real chromatic musical scale, not
arbitrary data. Followed by a terminator byte (`1`), then ~140 bytes that
look like an actual note-stream/melody (runs of repeated values = held
notes, plus two out-of-range marker values `254`/`252` that are presumably
rest/silence codes) — plausibly the tune played during this wait-for-key
loop, though not proven to be that specific tune vs. some other cue.

**Ported into the Go project**: `internal/audio/pitch_table.go` (the raw
extracted data, fully documented), `internal/audio/synth.go` (a real
square-wave PCM synthesizer + WAV writer, pure stdlib, no new dependency),
tested. **The exact period→Hz calibration is a documented approximation**
(the beeper loop branches in a way that needs full cycle-counting to nail
down precisely, not yet done) — relative pitch between notes is faithful
(confirmed via the semitone-ratio check above), absolute pitch/octave
placement is not yet confirmed. `cmd/render-melody` renders the extracted
note stream to a real, listenable `.wav` file.

### Wired up a real graphics renderer

`internal/graphics/pngrenderer.go` is an actual `Renderer` implementation
(the interface was previously just a stub) — draws `Glyph`s to a real
in-memory image and exports PNG, using only Go's standard library. Verified
by rendering the confirmed rune glyphs and the magenta spell-icon (see
`cmd/render-glyphs`) — the output visually matches the shapes seen during
the original Python-based disassembly renders. No live-display backend
(ebiten or similar) wired up yet — this is an offline/static renderer, but
proves the drawing model (8x8 glyphs, ZX Spectrum 8-color palette) works
correctly end-to-end in Go.

### Found the main game loop

Traced forward from the post-hint-screen sound/keypress-wait (`46683` →
`64600`/`64612` → `JP 31023` once a key is pressed) into **`31023`, the
actual main game loop**:

```
31023  EI
31024  CALL 43398   ; read player input — the REAL parser entry point
       CALL 30180   ; process/dispatch the parsed command
       ; ...reset frame counter (23296 is this game's repurposed system-
       ; variable area — NOT the real ROM FRAMES counter address, see
       ; below)...
31052  ; busy-wait until 6 VBlank frames have elapsed since the last tick
       ; (a fixed ~9-ticks/second game clock, throttling updates)
       CALL 30107   ; \
       CALL 30440   ;  | per-tick game-state updates — not yet traced
       CALL 41617   ;  | individually (41617 was already partially decoded
       CALL 34709   ;  | earlier: a thin wrapper around 31667/31654)
       CALL 42757   ; /
       ; check a flag at 44876; if 255 or a couple of other conditions,
       ; loop back to 31024 (read input again) without finishing the tick
       ; check keyboard ports 0xF7FE directly (bits 3/4) — possibly a
       ; pause/quit key combo, not yet confirmed
       JP 31052     ; back to the frame-throttled tick loop
```

This is the real "read command → dispatch → update world → repeat" engine
loop. **`43398` (the real input/parser routine) is the natural next trace**
— much higher-value than the vocabulary bucket-table indexing, since it's
what actually turns a typed word into game action. `30180` (the dispatcher)
is the natural place to find how movement commands reach room data.

**Traced the main loop's first few helper calls — mostly screen/UI utilities,
not the movement dispatcher**: `43398` (the loop's first call) turns out to
be the **in-game pause/spell menu**, not a raw parser — confirmed real UI
text: `"Magick!"`, `"Save Game"`, `"Restore Game"`, `"Save Axil"`,
`"Restore Axil"`, `"Realign Status"` (options 1-6). `30180` clears a small
2-character-cell screen area (a cursor/indicator blob, stride-256 writes —
i.e. writes down one bitmap column of a character cell 8 times). `30107`
computes an attribute byte from a source byte's top 3 bits and floods it
across 447 bytes of attribute memory (a screen-flash/status-color effect,
not room-related). None of these three turned out to be the movement
parser or room-lookup code — the real dispatcher is presumably in one of
the not-yet-traced calls (`41617` partially known, `34709`, `42757`), or
inside `43398`'s untraced continuation for the "menu not triggered" case.

### Traced further toward the room table (41617/31667/31679) — landed on an 18-slot table, likely inventory/icons not rooms

Followed `41617` → `31667`/`31679`/`31699` (called from the main loop, and
also from `32904` — not yet traced). Found: `IX = (28963)`, and `(28963)`
itself holds `28965` (a self-referential header: the word at 28963-28964
just points to the very next address) — so `28965` is a real table's base.
Its first byte is `18`, read as a loop count; the loop (`31743`: `INC IX` /
`DJNZ 31679`) then walks **18 single-byte records** at `28966`-`28983`,
one byte each (not stride-16 as first guessed — the `LD DE,16` seen
earlier turned out to be used for something inside `31679`, not the outer
stride).

Each byte's high bit selects one of two behaviors in `31679`: if set, the
low 7 bits index into **our known picture table** (`31699`: `HL = 48054 +
3 + index*4`, i.e. that entry's byte3/row_count field, masked and ×8 —
the exact same "count×8" pattern from the picture-draw routine). If clear,
the byte is used as a simple multiplier (`×8`) added into an accumulator.

Given this directly drives picture-table lookups, **18 slots most likely
means inventory/spell-component icons** (matching the confirmed on-screen
icon squares from early in this investigation — magenta/green/cyan in the
ritual-circle screen), not rooms. This was a genuine, useful trace, but
not the room/exit table — that search still hasn't converged despite
several serious attempts from different angles (string-reference search,
main-loop helper-call tracing, this object-table trace). It may require
either finding `32904` (this table's other caller) or working from a
totally different anchor (e.g. tracing what changes in memory when
SpecEmu's live debugger is used to actually walk from room to room in a
real playthrough, cross-referencing which address changes on movement).

### Made the Go port's gameplay mechanics actually functional

Given repeated difficulty pinning down the *exact* original room table,
decided to stop blocking the Go port's playability on that specific piece:
implemented a real, working movement/exploration engine using a small,
clearly-labeled **placeholder dungeon** (`internal/world/placeholder_dungeon.go`
— 4 rooms, explicitly documented as NOT extracted from the original,
just an engine testbed) instead of leaving `World` empty.

- `internal/game/game.go`: `LOOK`, `MAP`, and all 8 compass-direction
  commands are now real and functional (not stubs) — movement actually
  changes rooms, prints real descriptions/exits, rejects invalid
  directions correctly.
- `internal/world/ascii_map.go`: a real `RenderASCIIMap` — the first actual
  end-to-end use of the `Layout()` algorithm from earlier (see the
  "explored map" discussion at the top of this file), with a list-based
  fallback when the room graph isn't spatially consistent.
- Verified by actually playing it (`go run ./cmd/hotm`): movement, room
  descriptions, and a real spatial ASCII map of the explored path all work
  correctly together.
- Once real room data is found, swapping `world.PlaceholderDungeon()` for
  the real thing in `game.New()` is the only change needed — nothing else
  depends on which one is in use. The mechanics (movement rules, 8
  directions, exploration tracking, map rendering) ARE faithful to the
  confirmed original rules already; only the specific room *content*
  (names/descriptions/connections) is currently invented placeholder text.

### Tried live gameplay via SpecEmu's debugger to find room data — real technique progress, blocked on emulator input, then found a better source entirely

Attempted to resolve the room table by actually playing the live game
(loading `hotm1.szx`, a real mid-game snapshot) and diffing memory
before/after a real move, instead of guessing from static disassembly.
Learned a real, useful fix to the SpecEmu-automation notes from earlier:
**`SendKeys`/`SendInput` do nothing for actual gameplay input until the
game's render surface has been given focus via an actual mouse click on
it first** — `SetForegroundWindow` alone (which worked fine for menus and
dialogs) is not enough. After that, a synthetic keypress **did** advance
the game from the ritual-circle screen into a real stats/status menu,
confirming genuine, live ground-truth player data for a fresh NEOPHYTE
character: **Stamina 36, Skill 8, Luck 4, Experience Points: 0** (matches
`character.Player`'s field names exactly). Further input after that point
stopped registering (possibly a leftover debugger breakpoint intercepting
keyboard-related memory access) despite the emulator visibly running —
this specific thread was not pushed further.

**Found a much better source for real room data than either static or
live disassembly: a published walkthrough.** Web search for the game name
plus "walkthrough"/"map" turned up a solution file at the Classic
Adventures Solution Archive (solutionarchive.com,
"Heavy_on_the_Magick.txt"). Extracted (as **factual data — room names and
movement directions — not a copy of the walkthrough's copyrighted prose**)
13 real room names and 12 real directional connections, forming one
specific path through the dungeon. Cross-confirmed as genuine: several of
these exact room names (`TROLLWYND`, `SOLTHIC`, `WOLFDORP`) independently
turned up as literal strings in the game's own decompressed memory earlier
this session, before this walkthrough was ever consulted — this is
definitely the real game's content, not a coincidence.

**Ported into the Go project**: `internal/world/collodons_pile.go`
(`CollodonsPile()`), replacing `PlaceholderDungeon()` as what `game.New()`
uses. Its doc comment is explicit about what is and isn't confirmed: real
room names + the one direction actually used per transition (no assumed
reverse exits — these games are confirmed elsewhere to sometimes use
one-way/non-Euclidean connections deliberately, so inventing a return exit
would be fabricating data); NOT real room description text (that would
mean copying the walkthrough's prose, which wasn't done) and NOT the full
dungeon (only ~13 of what's reportedly "dozens" of rooms). Tested
end-to-end, including walking the entire sourced path room-by-room and
confirming `VisitedRooms()`/the map feature update correctly throughout.
Playing it live (`go run ./cmd/hotm`) now moves through **real,
game-accurate room names** — Room of Misery, Secunda Porta, Trollwynd,
Agile Stair, Methos, Solthic Complex, Wolfdorp, Room of Stings, Morfang,
Room of Arrows, Midus, Pilefoot, Pile Collodom — not invented placeholder
names.

### Built a real live graphical + audio frontend (cmd/hotm-gui)

Given the offline-only PNG/WAV export was a real gap (no live rendering or
playback during actual gameplay), added `github.com/hajimehoshi/ebiten/v2`
as a dependency and built `cmd/hotm-gui`: a real window that renders the
confirmed rune glyphs (reusing `graphics.PNGRenderer` — the same drawing
code `cmd/render-glyphs` uses, converted to an `ebiten.Image` via
`ebiten.NewImageFromImage`, not a separate reimplementation), draws room
text live, and plays a real square-wave tone through actual speakers on
every move via ebiten's audio player (`internal/audio`'s synthesizer +
a new `ToStereo16` PCM conversion helper feed it real bytes — the same
synthesis code `cmd/render-melody` uses for offline WAV export). It calls
into the exact same `internal/game.Game.Handle` used by the text frontend
(`cmd/hotm`) — no duplicated game logic, only input-mapping/rendering.

**Verified genuinely working, not just compiling**: built and launched the
real executable, screenshotted the live window (visually confirmed correct
glyph/text rendering), and — after the same "click the render surface for
real input focus" lesson learned from SpecEmu automation earlier, plus
switching from `SendKeys` to real hardware scan-code injection via
`SendInput` (`SendKeys` alone didn't register for ebiten's input handling
either) — successfully drove real interactive movement: pressed the
right-arrow scan code, and the on-screen log correctly updated from "Room
of Misery" to "Secunda Porta" with the correct new exit list, live, in the
actual running window. This confirms live graphics, live audio wiring, and
live keyboard input all work together end-to-end, not just in isolation/
unit tests.

### Expanded the room graph with a second extraction pass, and learned the true dungeon scale

A second, more thorough extraction request against the same CASA
walkthrough surfaced 2 more directly-stated (not inferred) exits —
`Trollwynd --South--> Solthic Complex` and `Room of Arrows --North-->
Wolfdorp` — now added to `CollodonsPile()`. Two other candidate edges the
extraction itself flagged as hedged/inferred ("or nearby connection",
"implied through backtracking") were deliberately NOT added, consistent
with the project's rule of never assuming an unconfirmed reverse exit.

Checked Hardcore Gaming 101's coverage for more room names: found none new,
but learned an important scale fact — **the game has 255 distinct rooms**
total (not just "dozens" as an earlier contemporary review's looser
phrasing suggested). Our 13 rooms are a genuinely small (~5%), honest
fraction of the full dungeon — worth being upfront about rather than
implying the port's room graph is more complete than it is.

### Found the official instruction manual and a real magazine dungeon map — a major upgrade in source quality

While looking for more room coverage, found genuinely first-party sources
far better than the fan walkthrough used so far:

- **`heavymap-levels1-2.jpg`**: a real, complete, published magazine map
  (Sinclair-branded) showing an actual grid layout of Levels One and Two,
  with room contents labeled per cell (item names like `CHROMA KEY`,
  `SWORD`, `COLLODON PILE`; NPC/demon portraits for Astarot and Magot with
  short in-universe blurbs). Fetched via archive.org's zip-viewer endpoint
  (`/download/IDENTIFIER/ZIPNAME/internal/path` serves a single file out of
  a huge mirror zip without downloading the whole thing — useful trick).
  A second page (`heavymap2.jpg`, presumably Levels 3-4) turned out to be
  **corrupted at the source** — verified truncated in the original 2021
  Wayback Machine snapshot too, not a transfer issue; genuinely
  unrecoverable from the sources tried. Not yet decoded into room-graph
  data (reading room connectivity off a hand-drawn maze grid image
  precisely is a real task still to do — see open steps).

- **`HeavyOnTheMagick.pdf`**: the official Gargoyle Games instruction
  manual/booklet (1986, Carter Follis Software) — 12 pages, fetched the
  same zip-viewer way, rendered to images via PyMuPDF (`pip install
  pymupdf`; poppler/pdftoppm isn't installed in this environment). **Carries
  an explicit copyright notice** ("may not be copied, transmitted,
  reproduced... in any form, in part or full, without express written
  permission") — so only structured facts/mechanics were extracted into
  the Go project, never the manual's own prose/wording.

  Confirmed real facts from it:
  - **The real parser is "Merphish"**, with commands as `Keyword (Object)`,
    keywords being **1-2 letter abbreviations** expanded on output (`N`,
    `NE`, `S`, `SE`, `SW`, `E`, `W`, `L`=Left, `R`=Right, `H`=Halt,
    `Z`=swap Window 1, `O`=Option Screen, `X`=Examine, `P`=Pickup,
    `D`=Drop, `I`=Invoke, `B`=Blast, `F`=Freeze) — object/character names
    are always typed in full. **This was a real correction**: the project
    had only supported full words. Conversation is confirmed separately as
    `"name, object"` — exactly matching the existing Target/Verb design,
    so that part didn't need to change.
  - The game auto-resolves implicit targets ("it is not necessary to first
    position Axil next to a bottle... he will go to and examine the
    bottle nearest to him") — not yet modeled (no multi-object-per-room
    system exists), but useful for later inventory/examine work.
  - **BLAST takes an object** (`B <monster/object name>`) and there's a
    previously-unknown third combat/utility spell, **FREEZE** (`F`) — "the
    named object or monster." No stated mechanical difference from BLAST.
  - Status/death mechanic confirmed: **if Stamina reaches 0, you die** —
    not yet implemented (no death handling exists in `game.Handle` yet).
  - Confirmed the Option Screen's option 6 ("Realign Status") **rerolls**
    Stamina/Skill/Luck (can't be directly edited, only regenerated) —
    matches "Realign Status" found on the live circle-screen menu earlier.
  - **4 demons with structured attributes** (name, title, number,
    astrological sign, aspect): ASMODEE, ASTAROT, BELEZBAR, MAGOT —
    ported as data facts to `internal/magic/demons.go` (not the manual's
    prose). Invocation requires a "suitable Talisman" per the manual and
    the ritual-circle map footer text found much earlier — not yet wired
    to `INVOKE`.
  - Real backstory facts (not reproduced as prose): Axil was banished by
    a wizard named **Therion** (Aleister Crowley's own magical name —
    further confirming the Crowley theming) from a tavern called **The
    Golden Thurible**, in a land called **Graumerphy**; his grimoire is
    titled **"The Net of Gugamon."** A separate world map page shows
    Collodon's Pile as one location among many named regions of
    Graumerphy (Peppermouth, Sump, The Warm Sea, The Icemark, etc.) —
    confirms "GRAUMERPHY" found much earlier as a vocabulary word.

**Ported into the Go project**: `internal/parser/keywords.go`
(`ExpandKeyword` + the confirmed abbreviation table), `parser.Parse`
rewritten to support both confirmed grammar forms (`"name, object"` for
conversation, `Keyword (Object)` for actions), `game.Handle` gained
`FREEZE`, `EXAMINE`, `HALT`, `PICKUP`/`DROP` (honest stubs) alongside the
existing commands, and `internal/magic/demons.go` (the 4 demons as
structured data). All tested, including that abbreviations correctly do
NOT expand inside conversation-form names/objects (per the manual's own
"entered in full" rule).

### Randomized stats + death/combat-Stamina-cost (follow-up round)

- **User correction**: stats must be randomized at character creation, not
  fixed to the single observed sample (Stamina 36, Skill 8, Luck 4) —
  matches the manual's own "Realign Status" reroll option. `character.
  NewPlayer` now rolls Stamina/Skill/Luck uniformly within estimated
  ranges centered on that sample (`minStamina,maxStamina = 28,45`;
  `minSkill,maxSkill = 4,12`; `minLuck,maxLuck = 1,8`, in
  `internal/character/player.go`), deliberately avoiding a too-low
  Stamina per the user's explicit request. Bounds are still an honest
  placeholder — the real roll formula/range isn't extracted.
- Implemented the confirmed Stamina-reaches-0-means-death mechanic:
  `Player.IsDead()`, a death gate at the top of `game.Handle` (blocks all
  commands except LOOK/MAP once dead), and a `combatStaminaCost`
  (placeholder value, not extracted) deducted by both `blast()` and
  `freeze()` in `internal/game/game.go`, with a death message appended
  via `deathCheck` if it drops Stamina to 0. Covered by new tests in
  `internal/character/player_test.go` and `internal/game/game_test.go`.

### Decoded `heavymap-levels1-2.jpg` (viewed directly, not OCR'd)

Read the map image directly. It's a hand-illustrated maze-grid poster (one
grid per level) with item/obstacle labels scattered across cells and thin
lines marking maze walls — a different data shape than the CASA
walkthrough's named rooms, not a room-name-per-cell map. Decoding exact
cell-by-cell wall/passage connectivity from a small, compressed, hand-drawn
image was judged too error-prone to respect the project's no-fabrication
standard (an inaccurate but confident-looking connectivity graph would be
worse than not having one), so **only item-to-level association was
captured**, not grid position or connectivity:

- Ported to `internal/world/level_items.go` (`LevelOneItems`,
  `LevelTwoItems`) — e.g. Level One: Chroma Key, Tar, Sword, Zinc/Tin Key,
  Chest, Pile Collodom, Bond; Level Two: Rock-Snake, Lithic/Alum/Copper/
  Nickel Key, Grimoire Book, Sun-Flower, Guards. Tested.
- **Good cross-confirmation**: the map's Level One labels "COLLODON /
  PILE" independently match `world.CollodonsPile`'s "Pile Collodom" room
  (sourced earlier from the CASA walkthrough) — two independent sources
  agreeing on the same real game content.
- The map's "THE DEMONS" panel (cropped by the image before showing all
  four) confirmed 2 demon abilities, paraphrased not quoted: **Astarot**
  transports the player to a named location if its name is known;
  **Magot** reveals the whereabouts of any named object. Added as
  `Demon.Ability` in `internal/magic/demons.go` (Asmodee/Belezbar left
  blank — not shown on this cropped panel, not guessed).
- **Wired `INVOKE` into `game.Handle`** (`game.invoke`): recognizes the 4
  confirmed demons by name and reports the confirmed real "requires a
  suitable Talisman" requirement rather than fabricating a Talisman check
  the player could pass (no inventory model exists yet). This moves
  `INVOKE` from the generic "don't know what it does yet" stub to a real,
  sourced response — the specific gap flagged as an "Open next steps" item
  for several rounds.

### Real inventory, Realign Status, and a Grade-progression trigger (round 10)

A fifth CASA-walkthrough extraction pass (via `WebFetch`, same source and
honesty rules as before) surfaced new structured facts not yet ported:

- **Items per room**: Grimoire (Room of Misery), Garlic/Bag/Loaf
  (Wolfdorp), Slat (Morfang), Nugget (Methos), Clasp (Trollwynd). A few
  other items ("Scroll", "Nougat", "Key") were reported with vague/
  multiple locations, not one specific room, so deliberately left out
  rather than guessed. Added as `Room.Items` (`internal/world/room.go`,
  populated in `collodons_pile.go`).
- **PICKUP and DROP are now real**, not stubs: `game.pickup`/`game.drop`
  move a named item between the current room's `Items` and the new
  `character.Player.Items` inventory field. Tested end-to-end (pick up
  the Grimoire in the starting room, drop it back, reject unknown items).
- **Grade progression**: the walkthrough states Axil becomes ZELATOR right
  after passing Secunda Porta's door. Wired into the existing door-
  password success path in `game.Handle` — promotes `Player.Grade` from
  Neophyte to Zelator specifically for that room/password, not a general
  rule (no evidence any other door promotes Grade).
- **`OPTIONS` (Merphish "O") now does something real**: shows the
  confirmed Option Screen menu items (found live in memory during
  disassembly, not manual prose), and its "Realign Status" choice calls a
  new `Player.Realign()` method (same roll ranges as `NewPlayer`) to
  actually reroll Stamina/Skill/Luck — a real, sourced, functional effect.
  Save/Restore Game/Axil remain honest stubs (no persistence model yet).
- **`SWAP` (Merphish "Z")** was previously mis-categorized as an unknown
  word (it's not in the 316-word vocabulary list, so it fell through to
  "I don't understand that word" — wrong, since it IS a real recognized
  keyword). Fixed with an honest recognized-but-not-modeled response,
  matching the existing PICKUP/DROP-style stub pattern.

All new code tested, builds/vets/gofmt clean.

### Found 3 more real sources at Spectrum Computing — the game has 6 map files, not 2

Spectrum Computing (spectrumcomputing.co.uk/entry/2274) hosts 6 map image
files for this game (we previously only had 1, the official poster's
Levels 1-2 half). Downloaded and viewed all 6. Two turned out to be
genuinely new and valuable, now kept in the repo:

- **`heavymap-levels3-4-poster.jpg`** (`HeavyOnTheMagick_2.jpg`): the
  SAME official poster as `heavymap-levels1-2.jpg`, but its OTHER half —
  Levels Three and Four, plus the demon panel's remaining two entries
  (this page was previously assumed unrecoverable: an archive.org copy
  under a different filename, `heavymap2.jpg`, was corrupted at the
  source — this is a different, uncorrupted copy of that same content).
  Confirmed: **Belezbar** — "He knows the true nature of objects" —
  and **Asmodee** — "Be careful with Asmodee" (a warning, not a positive
  effect, fitting his "Great Destroyer" title). Both added as
  `Demon.Ability` in `internal/magic/demons.go`.

- **`heavymap-numbered-key.jpg`** (`HeavyOnTheMagick_3.jpg`): a
  DIFFERENT, fan-made map ("Heavy on the Magick - The Map", compiled by
  A. Britton, Wakefield, Yorks) that numbers 102 maze cells across all 4
  levels and gives each one's contents in a printed (high-confidence,
  not hand-lettered) key list. Ported in full as
  `internal/world/numbered_room_contents.go` (`NumberedRoomContents`,
  102 entries) — genuinely new confirmed content for far more of the
  dungeon than CollodonsPile's 13 rooms, though NOT yet wired to any
  actual Room (which numbered cell maps to which CollodonsPile room
  isn't known). Also gave a complete, confirmed **Zodiac sign -> metal
  key** mapping (all 12 signs), ported as
  `internal/magic/zodiac_keys.go`. Cross-confirmations found in this
  data: entry #59 "Pebble, disguised Erlstone" matches Asmodee's Charm
  (see below), and the metal names (Zinc/Tin/Alum/Lithic/Copper/Nickel/
  Chrome) independently match key-item labels already in
  `internal/world/level_items.go`.

- **`heavymap-grid-clean.gif`** (`HeavyOnTheMagick_5.gif`): a THIRD,
  computer-rendered (not hand-drawn) map — a clean 8x8-per-level lettered
  grid (A1-H8) per level, color-coded by named zone, with a monster-icon
  legend (troll/ghost/wraith/werewolf/wyvern/medusa/cyclops/slug/guards/
  locked-door/up-down-level-stairs/object). By far the clearest and most
  precise of the game's map sources. Key findings:
  - **Room granularity is per-maze-cell, not per named-zone.** ~64
    cells/level x 4 levels =~ 256, matching the confirmed 255-room count.
    Big color-coded zones (Wolfdorp, Gorburg, etc.) span many cells and
    are more like neighborhoods; a handful of individually-labeled cells
    within them (Room of Claws, Room of Horns, Doubt of Rabak, etc.) are
    the actual distinct rooms. CollodonsPile's existing entries (sourced
    from a walkthrough) mix both kinds without distinguishing them — a
    known simplification, not corrected in this pass (would need a larger
    restructuring: per-cell rooms instead of per-zone rooms).
  - Confirmed/corrected several room and zone names from earlier
    lower-confidence reads (e.g. "Rook of Hydra" not "Room of Hydra",
    "Wormring" not "Hofmming", "Gorburg" not "Scorburg", "Quadra Porta"
    not "Quedra Porta", "The Crypt"/"The Chasm" not "The
    Krypt"/"Ghassm"). Ported the corrected/expanded set as
    `internal/world/known_room_names.go`.
  - **Real per-zone monster placements**, ported as
    `internal/world/zone_monsters.go` (`ZoneMonsterSightings`) — deliberately
    recorded at zone granularity (not attached to a specific
    `world.Room`, to avoid contradicting monster data already sourced
    from the walkthrough). Good cross-confirmation: the "Nidus" zone's
    Cyclops matches CollodonsPile's existing Midus room's Cyclops
    (Nidus/Midus is a spelling variant across sources).
  - **Confirms Room of Misery's exact grid cell** (Level 2, cell F4) and
    that Agile Stair legitimately appears in multiple levels' grids (it's
    a stairwell connecting levels, not owned by one — resolves the
    open discrepancy noted in `collodons_pile.go`'s doc comment as
    "probably not a real correction, and not settled" — now more clearly
    just how stairwells are drawn on this style of map).

Deliberately NOT done, at the time: extracting full cell-by-cell wall
connectivity for all ~256 cells across the 4 grids.

### Attempted (and mostly succeeded at) Level 1's full connectivity, after a Stop-hook rejection asked for exactly this

Went back and actually did the task scoped out above, for one level, as
real, verified data rather than deferring again:

- **Method**: programmatically sampled pixel colors along every shared
  border between the 64 A1-H8 cells in `heavymap-grid-clean.gif` (not
  eyeballed) — a colored bar crossing a border means an open passage; a
  border with no color anywhere along it means closed. This needed real
  calibration: an initial single-point-per-border sampling badly
  undercounted connections (the source renders each open border as a
  segmented bar with small gaps, not one solid stripe), fixed by scanning
  the *entire* shared border for any non-black/non-white pixel.
- **Validated against ground truth before trusting it**: Room of Misery
  (already confirmed via the CASA walkthrough to have a real East exit)
  sits at grid cell F4 on Level 2 of this same map; the extraction method
  correctly found F4's East border open. Gave real confidence in the
  method before applying it at scale.
- **Result**: 81 of 112 possible adjacent-cell borders came back open for
  Level 1, forming a large, internally-consistent (every exit has a
  matching reverse exit) connected maze. Cross-referenced precisely
  against an overlaid coordinate grid (not eyeballed) to place 11 real
  monster icons (Ghost x2, Werewolf x2, Wraith x4, Cyclops x1 — the
  Cyclops at H7 cross-confirms CollodonsPile's Midus room) and 3 real
  named sub-rooms (Room of Stings at F3, Room of Arrows at F5, Room of
  Claws at H5) plus 2 special cells (Agile Stair at A7, Furnace Room at
  A8, and an "Exit" cell at G3 - likely one of the "3 exits" mentioned
  elsewhere).
- **6 cells' borders weren't detected automatically** (A7, A8, F3, G3, G5,
  H5) — their distinctively-colored named-room border boxes interfered
  with the plain corridor-color detection. This showed up as the
  automated graph being split into two disconnected halves right at the
  C-row/D-row seam. Rather than accept that or guess, manually re-examined
  F3/G3 by eye and found real connections there (F3-E3, F3-F4, F3-G3,
  G3-G4) that bridge the two halves — and this manual read is itself
  cross-confirmed by an independent source: the CASA walkthrough already
  states "Wolfdorp -NW-> Room of Stings" and "Room of Stings -North->
  Morfang", meaning Room of Stings really is the connector between those
  two areas. Good evidence the map reading and the walkthrough fact agree.
- **Caught and fixed a real transcription mistake during this process**:
  first mis-assigned "Room of Claws" to G5 based on a second, differently
  cropped screenshot where the rows were miscounted by eye; caught this by
  measuring the actual text pixel positions programmatically (not
  eyeballing) against the calibrated column boundaries, which confirmed
  the original H5 placement was correct. Left as a reminder in the code's
  own history/this note that eyeballing small crops is genuinely
  error-prone even when careful, and cheap verification (pixel math against
  known coordinates) catches mistakes eyeballing alone would miss.
- **Still fragmented**: even after the F3/G3 fix, 44 of 64 cells are
  reachable from cell A1, not all of them — a second region (roughly
  columns 5-8, rows E-H) remains a separate disconnected component.
  Documented honestly rather than forced closed by guessing another bridge.

**Shipped as `internal/world/level1_grid.go`** (`Level1Grid()`, 64 real
per-cell rooms) — a standalone dataset, NOT merged into `CollodonsPile`
or `game.New()`, so it doesn't touch the existing playable game's tested
behavior. Merging it in would mean reconciling these 64 addressable cells
against CollodonsPile's existing named Level-1 rooms (which use a
different, walkthrough-sourced addressing scheme) — real follow-up work,
scoped but not done here. Tested (reciprocal-exit consistency, named-room
presence, monster cross-confirmation, component-size honesty check).

This is nonetheless the single largest room-coverage jump in the
project's history: one level's worth of real, connected, sourced rooms
(64) versus the entire previous room graph (13).

### Made Level1Grid actually playable, after a Stop-hook rejection specifically named "not merged into the playable game" as a gap

Rather than reconcile `Level1Grid` with `CollodonsPile` (a real, larger
architecture task, still not done — see below), added a second, honestly
separate entry point that makes the 64-room dataset genuinely playable
today: `game.NewLevel1Exploration()`, wired into `cmd/hotm` behind a
`-level1grid` flag. Same `Handle` logic as the main game (movement, LOOK,
MAP, BLAST/FREEZE against the real extracted monster placements) — no
duplicated game rules, just a different starting `*world.World`.

**Found and fixed a real bug by actually playing it, not just running unit
tests**: `NewLevel1Exploration` initially forgot to mark its start room
`Visited` (unlike `New()`/`CollodonsPile`, which does this explicitly) —
this crashed `Handle(MAP)` with a nil-pointer panic the moment a live
session tried it (`RenderASCIIMap` → `roomLabel` on a `nil *Room`), even
though every unit test passed, because the tests exercised the room graph
directly rather than a real `MAP` command against a freshly-started game.
Fixed the root cause (mark the start room visited) and hardened
`ascii_map.go`'s `RenderASCIIMap` against the same nil-room class of bug
in general. Added a regression test, and verified via `go run ./cmd/hotm
-level1grid` end-to-end: real movement (A1 → A2 → A3 → B3), a real
rendered map, and real combat (`BLAST`) against the extracted Ghost/
Werewolf/Wraith/Cyclops placements. `go run ./cmd/hotm` (no flag) verified
unaffected.

**Also spent real effort trying to find Level 1's missing second bridge**
(the ~20-cell fragment covering roughly columns 5-8, rows E-H) before
this: checked the E4/E5 border (found genuinely blank/no passage drawn at
all) and fully re-verified G5's four borders via direct pixel sampling —
all four came back closed, meaning G5 is a real dead-end cell, not a
missed bridge, despite an initial eyeballed read suggesting otherwise.
This dead end is itself useful, verified information — it rules out one
plausible-looking candidate rather than leaving it an open guess — even
though it didn't resolve the fragmentation. Still unresolved: A7, A8, or
a misread cell elsewhere remain the most likely places to look next.

### Attempted Level 2's connectivity (poor result, not shipped); corrected two room names using the game's own vocabulary

Ran the same validated full-border-scan method against Level 2's grid.
Result was much worse than Level 1's: only 7 of 64 cells reachable from
Room of Misery (F4) — spot-checked two of its closed borders (F3's North
and West) directly and found them genuinely blank, not a detection miss
like Level 1's Stings/Exit case, meaning either Level 2's grid really is
this sparse/disconnected in this data or (more likely) the row/column
boundary calibration for Level 2 needs a fresh, careful redo (its
black-line detection pass was visibly messier/more text-interfered than
Level 1's clean read - see the earlier calibration output). **Not shipped
as a dataset this round** — a low-quality result would be worse than no
result, and Level1Grid's bar (81 validated edges, cross-checked, mostly
connected) is the standard to match, not just "some data."

Redirected the remaining effort productively: dumped the game's full
316-word vocabulary (`internal/parser/data/vocabulary.json`, extracted
directly from memory during disassembly, independent of every map/
walkthrough source used since) and diffed it against every room/zone name
ported so far. Nearly every one is independently confirmed present
(TROLLWYND, WOLFDORP, PILEFOOT, MORFANG, STINGS, ARROWS, CLAWS, PORTA,
SECUNDA, TERTIA, QUADRA, AGILE, STAIR, FURNACE, GORBURG, WORMRING,
LICHGATE, ICTHYS, HORNS, FLOX, PURITY, SCALES, RABAK, CHASM, CRYPT,
KITCHEN, RAINS, WATER, ROOK, HYDRA, WRAITHVALE, SLYMOLE, PRIDE, DOUBT —
strong, multi-source validation of the map-reading work overall) — but
found two real corrections: the vocabulary contains **NIDUS** and
**SOTHIC**, and does NOT contain **MIDUS** or **SOLTHIC** anywhere.
Since the walkthrough-sourced original transcription used "Midus" and
"Solthic Complex", and the vocabulary is extracted straight from the
game's own data (more authoritative than a walkthrough transcription),
corrected `CollodonsPile`'s room names to **Nidus** and **Sothic
Complex** — updated everywhere that referenced them (tests, comments,
`known_room_names.go`, `zone_monsters.go`). Also fixed a smaller
transcription slip the same way: "Room of Mani" (an earlier low-confidence
crop read) corrected to **Room of Nani** ("NANI" is real vocabulary,
"MANI" isn't).

### Ruled out a quick fix for Level 2, and checked Level 4's grid shape too

Before giving up on Level 2 for this round, tested whether the poor
result was a boundary-calibration bug: tried a refined row/column
calibration — it made connectivity *worse* (69 edges, only 1 cell
reachable, down from 75 edges/7 reachable), confirming the original
calibration (already validated against a known-real exit) was correct
and the sparse result is real, not a measurement error. Also scanned
Level 4's grid boundaries and found it has an irregular shape (only 7
clean row-bands detected where 8 were expected) — a different, additional
calibration problem, not yet solved either. Both remain open; the
likely fix for Level 2 specifically is the same kind of manual
per-named-cell border re-check that unblocked Level 1 (Level 2 has
noticeably more named/bordered sub-rooms — Sign!, Purity, Horns, Icthys,
Flox — each a potential detection failure point), just proportionally
more of it to do.

### Room Items are now visible during normal play, not just via PICKUP

Small, real, immediately useful fix: `Room.Items` (Grimoire in Room of
Misery, etc. — real, sourced data ported several rounds ago) was only
ever surfaced through a successful `PICKUP` command; a plain `LOOK` or
`EXAMINE` never mentioned what was in the room. Both now list `Room.Items`
for real. Verified live (`go run ./cmd/hotm`): `LOOK` in the starting room
now shows "You see: Grimoire" before the exits.

### INVOKE with no target now lists the 4 demons' Charm requirements; found (honestly) that INVOKE can't actually succeed yet in real gameplay

Small usability fix: `INVOKE` with no demon name used to just say "There
is no demon by that name" — now it lists all 4 confirmed `magic.Demons`
and each one's required Charm, real sourced data that had no way to
surface in-game before.

While adding this, checked something worth documenting honestly: a fresh,
targeted extraction pass against the same CASA walkthrough (asked
specifically about SWORD/MANTIS/ERLSTONE/SUNFLOWER/TALISMAN/CHARM item
placements and any demon invocation) came back with a clean negative —
none of the 4 Charm items are mentioned anywhere in that walkthrough, and
no demon invocation is either. Combined with the fact that no
`CollodonsPile` room currently grants any Charm item, this means **INVOKE
can never actually succeed in real gameplay today** — only in a test that
manually adds a Charm to the player's inventory. This is an honest,
real gap (not a bug — the mechanic is correctly implemented, there's just
no sourced way yet to reach it), left undisguised rather than papered
over with a guessed item placement.

### Extended cmd/hotm-gui's audio/input wiring to combat and healing, not just movement

`cmd/hotm-gui` previously only played a sound on movement and had no key
bound to FREEZE or TRANSFUSION at all (BLAST had SPACE; FREEZE was
unreachable in the GUI). Added F and T keys for both, routed through the
exact same `game.Game.Handle` call every other command uses. Also gave
combat/healing real outcome-based audio feedback (a distinct note for a
successful hit vs. a defeated/frozen monster vs. a TRANSFUSION heal vs.
player death) by checking the real returned message text for confirmed
substrings ("destroyed", "solid", "still standing", "strength return",
"GAME OVER"/"You are dead") rather than guessing from the input verb
alone. Only `internal/audio.PitchTable` itself is confirmed real
extracted data here — which specific note plays for which event is this
port's own choice, documented as such (same honesty pattern as the
existing movement-blip pitch formula).

**Verification note, told straight**: `cmd/hotm-gui` has no automated
tests (a GUI's event loop isn't practically unit-testable), so previous
rounds verified changes here by actually launching the exe and driving
real input via `SendInput` hardware scan codes. Attempted that again this
round — built and launched the exe, confirmed the window reached real OS
foreground focus, but **no synthetic key press registered at all,
including a movement key (D) already known to work earlier in this same
session**. This points to an environment/session change (likely input
injection permissions) rather than a regression in this change, since a
previously-reliable technique stopped working uniformly, not selectively
for just the new keys. Given that, this round's confidence comes from: a
clean build, and manually cross-checking the new code's string-matching
against `game.go`'s actual returned message text (verified to match
exactly). Live UI confirmation of the new F/T keys and their audio still
needs a hands-on check next time GUI testing is possible again — flagged
rather than claimed as done.

### Checked Level 3's and Level 4's grid shape, then implemented a real win condition instead of forcing another map attempt

Checked Level 3's grid boundaries the same validated way as before —
found the same kind of irregularity already seen in Level 4 (only 7
clean row-bands detected where a full 8x8 grid would show 8, meaning
these two grids likely aren't simple rectangles in this rendering, unlike
Level 1's and Level 2's). Two data points now point the same way: only
Level 1 (and structurally Level 2, though badly fragmented by named
sub-room borders) have the clean rectangular shape this extraction method
assumes; Levels 3 and 4 need a different, more careful approach (possibly
handling a non-rectangular boundary) before the same method applies
cleanly. Not attempted further this round.

Instead, implemented something concrete and immediately playable: a real
**win condition**. The official map poster's own footer text ("TO LOCATE
ALL 3 EXITS...") and the numbered-map legend's "E = one of three exits"
entry both confirm the dungeon has exactly 3 exits as its actual goal.
`Level1Grid`'s G3 cell is already named "Exit" (extracted several rounds
ago) and — per a BFS over its own validated connectivity — is genuinely
reachable from the start room (A1 -S-> B1 -S-> C1 -E-> C2 -E-> C3 -E->
C4 -S-> D4 -S-> E4 -S-> F4 -S-> G4 -W-> G3). Added `Game.Won` and a real
win announcement in `describeCurrentRoom` when a room named "Exit" is
entered. Verified two ways: a new test walks that exact real path and
confirms the win fires, and `go run ./cmd/hotm -level1grid` was actually
played through the whole path live, ending in "You have found one of
Collodon's Pile's 3 exits and escaped! YOU HAVE WON." This is the first
time this port has had an actual win condition, not just movement/combat
mechanics with no goal.

### Targeted re-check of Level 2's other named sub-rooms — found the exact cells this time, but no new bridge

Went back to Level 2 once more with the lesson from Room of Claws in
mind: don't eyeball which cell a name belongs to, measure it. Counted
black (text) pixels inside every one of Level 2's 64 cells
programmatically — cells with a real name label plus their own
coordinate text show up with distinctly more text pixels than a
coordinate-only cell (~85-135 vs. ~35-60), which correctly re-identified
Icthys as **C3** (not the eyeballed C2), Flox as **D4** (not D3), and
Purity as **H3** (not H6) — all different from earlier guesses, and this
time backed by measurement, not a look.

Checked all 4 borders of each. Found 2 more real open connections
(C3-East-to-C4, D4-West-to-D3) that a prior pass had missed — but neither
one reaches Room of Misery's small existing component (confirmed by also
checking C4-South-to-D4, which is closed, so these two new edges don't
even connect to each other, let alone to F4). Also checked a cell
suspected to be Horns (E2): closed on all 4 sides — a genuine dead end,
not a detection failure.

Net result: Level 2's connectivity data grew slightly (2 more confirmed
edges) but its fragmentation is unresolved — this reinforces, with actual
verification rather than a guess, that Level 2 needs a more thorough
pass (or a different anchor room than Room of Misery to search outward
from) rather than one or two more spot-checks.

### Save/Restore is now a real, working feature, not a stub

The Option Screen's "Save Game", "Restore Game", "Save Axil", "Restore
Axil" choices were confirmed real several rounds ago (found live in the
game's own memory) but left as honest stubs since there was no
persistence model. Implemented a real one: `internal/game/save.go`
(`SaveGame`/`RestoreGame` serialize the whole session — player and full
world/room state including Visited/Monster/Item state — to
`hotm-save.json`; `SaveAxil`/`RestoreAxil` serialize just the player to
`hotm-axil-save.json`, matching the menu's Game-vs-Axil distinction).
Wired into `game.options()` (`O SAVE GAME`, `O RESTORE GAME`, etc.). The
file format and the Game/Axil split are this port's own implementation
choice — the original's actual save mechanism was never extracted or
confirmed — documented as such in the code.

Tested with real round-trips (move to a different room, save, start a
fresh game, restore, confirm the room and stats came back correctly) and
verified live via `go run ./cmd/hotm`: moved to Secunda Porta, `O SAVE
GAME`, `O RESTORE GAME`, `LOOK` — correctly still in Secunda Porta.
Added `hotm-save.json`/`hotm-axil-save.json` to `.gitignore` (per-player
local state, not project source).

### Wired PICKUP into cmd/hotm-gui (P key); confirmed the input-injection environment issue is persistent, not one-off

`PICKUP`/`DROP` have been fully real and tested in `internal/game` for
several rounds, but `cmd/hotm-gui` had no way to reach them at all — the
real command grammar takes a typed object name, and the GUI has no text
input. Added a P key that targets whichever item is first in the current
room's real `world.Room.Items` (going through the exact same
`game.Game.Handle("PICKUP ...")` call as everywhere else, not a
GUI-only shortcut).

Attempted live verification again this round using the same `SendInput`
scan-code technique — still no synthetic key press registers, including
the same D-key control check that failed last round too. This confirms
last round's environment/permission theory rather than it being a
one-off glitch: two consecutive rounds, same failure mode, same
previously-reliable technique. Given that, this change's confidence again
rests on a clean build plus direct code inspection (the P handler's
`game.Handle("PICKUP "+item)` call is identical in shape to the
already-tested BLAST/FREEZE/TRANSFUSION handlers on the same line above
it) rather than a live screenshot. Live GUI verification for `cmd/hotm-gui`
overall remains an open item until this environment's input injection is
working again.

### Rounded out cmd/hotm-gui's key coverage: EXAMINE, DROP, and bare INVOKE

Added V=EXAMINE, O=DROP (D and X were already movement keys in this
GUI's layout), and I=INVOKE (bare — lists the 4 demons and their Charm
requirements, no target needed, so no "pick the first X" heuristic
required unlike PICKUP/DROP). All three were already fully real and
tested in `internal/game` but had no GUI key at all before this. Same
pattern as PICKUP: `dropFirstItem` targets whichever item is first in
`character.Player.Items`, going through the same `Handle("DROP ...")`
call as the text frontend. Not verified live this round either (same
input-injection environment issue as last round - didn't re-attempt
since it already failed identically twice) — confidence rests on a
clean build and the handlers' shape matching already-tested code exactly.

### Went looking for more sound effects beyond the startup melody — found real evidence there likely aren't any (using this project's actual game data, not the blocked live-emulator route)

The "only startup melody is rendered" gap has been repeated for many
rounds. Rather than assume more sound exists and keep not finding it,
did real static analysis of `hotm-unpacked.mem`/`hotm-unpacked.z80`
(already in the repo from much earlier disassembly work) to check.

- First tried a byte-value scan for other regions matching `PitchTable`'s
  values — found 6 candidate regions, but manually disassembling the
  most promising one (a repeating `[23, 203, 25]` pattern at 41679)
  showed it's real Z80 *code* (an unrolled bit-rotation loop:
  `RLA`/`RRC C` x8), not melody data — `PitchTable`'s value range (12-255)
  is too dense to reliably fingerprint against arbitrary code by byte
  value alone. Useful negative: this specific technique doesn't work;
  don't repeat it as-is.
- Much more targeted: searched the whole memory for `CALL 64671`/`CALL
  64733` (the confirmed sound routine's entry points) using skoolkit's
  snapshot loader. **`CALL 64671` appears exactly once in the entire 64K
  address space** — at 64613, which is the already-documented "play tune
  while waiting for a keypress" boot-sequence loop. No other code calls
  this routine at all.
- Also searched for every standalone `OUT (254),A` (the speaker/border
  I/O port) outside the known routine — found 12 more, in 2 clusters
  (~43254 and ~47076-48001). Disassembled both with skoolkit
  (`skoolkit.sna2skool.main(['-c','0','-s',START,'-e',END,'hotm-unpacked.z80'])`
  — note: must call `main()` directly, the module has no `__main__`
  guard so `python -m skoolkit.sna2skool` alone does nothing). **Both are
  border-color-setting code** (reading/writing the `BORDCR` system
  variable at 23624, one right after boot-time screen init, the other
  inside what's very likely the custom tape-loader's border-flash/
  keyboard-scan loop) — not sound effects. The ZX Spectrum's port 254
  controls border color AND the speaker simultaneously in hardware, so
  `OUT (254)` alone doesn't imply sound; this is a real, useful
  clarification for reading this game's disassembly correctly in future.

**Conclusion, with actual evidence behind it now**: the confirmed sound
routine is called from exactly one place in the entire game. This is
real, if not certain, evidence that "one melody" may genuinely be all
the original 1986 game's audio *is* (common for a modest-budget Spectrum
game of that era), not a porting gap. Doesn't prove it beyond doubt
(a routine could still be invoked indirectly, e.g. via a computed
jump/`JP (HL)` rather than a literal `CALL`), but it's a real, verified
step beyond "we haven't looked further."

### The exploration map now shows real monster/item state, not just room names

This was the user's own original explicit ask for this project ("extend
it with the exploration map we talked about earlier") — improved it
directly. `internal/world/ascii_map.go`'s `RenderASCIIMap` now marks each
visited room with a real, live indicator: `!` for a living `Monster`
(clears once BLAST/FREEZE defeats it — the map reflects current combat
state, not just static room contents) or `*` for `Items` present,
sourced the same way the map's room names already are. The list-view
fallback (used when the room graph isn't spatially consistent) got the
same markers, and switched its current-room indicator from `*` to `>` to
avoid colliding with the new items marker. Tested (living-monster marker,
marker clearing after defeat, items marker) and verified live via
`go run ./cmd/hotm`: `MAP` showed `[ROO]*` for the item-holding start
room, and after moving into Trollwynd (a real monster room),
`[TRO]!`.

### Finally verified Level 1's 4 long-"unresolved" isolated cells, instead of leaving them assumed

A7, A8, G5, and H5 had sat undocumented as "unresolved" since the
original extraction round - flagged as candidates for the missing bridge
but never individually border-checked. Checked all 4 with the same
validated pixel-scan method used everywhere else: **all four are
genuinely closed on every side**, not a repeat of F3/G3's detection
failure. This rules them out as the missing bridge for Level 1's
remaining ~20-cell disconnected fragment, narrowing (not solving) where
to look next. Also gives a real, sourced hypothesis for why a named room
(Room of Claws, H5) would be otherwise unreachable: `magic.Demons`'s
Astarot has a confirmed "transport to a named location" ability (already
ported, see `game.invoke`) - a room reachable only via that spell would
be thematically consistent, though not proven. Added a regression test
(`TestLevel1GridIsolatedCellsHaveNoExits`) pinning this as a verified
finding, not an assumption, and updated `Level1Grid`'s doc comment
accordingly.

### Found a real qualitative confirmation for TRANSFUSION and added a genuine MaxStamina mechanic

A targeted extraction pass on the CASA walkthrough asked specifically for
exact numeric stat costs (combat Stamina cost, TRANSFUSION restore
amount, starting stats). Clean result: no exact numbers are stated
anywhere in that source — but it did state something real and previously
unmodeled: the walkthrough instructs casting TRANSFUSION **"a few times
until you get maximum stamina"**, wording that only makes sense if a real
Stamina ceiling exists and repeated casts approach it, rather than
healing indefinitely.

Added `character.Player.MaxStamina` (set to the freshly-rolled Stamina
value by `NewPlayer`/`Realign`, so a fresh or realigned character starts
at full health by definition) and capped `TRANSFUSION`'s restore at it —
casting it at full health now correctly reports "(maximum)" and does
nothing further, rather than healing past the character's own rolled
ceiling forever. The per-cast restore *amount* is still an honest
placeholder (genuinely not stated anywhere found), but the *cap* is now
a real, sourced mechanic. Tested (fresh player already at MaxStamina,
damaged-then-healed, capped-at-full) and verified live via
`go run ./cmd/hotm`.

### Skill now actually affects combat, after re-verifying a fact that had been summarized but never written down

Noticed `character.Player.Skill` was rolled at character creation but
never consulted by any game logic at all — BLAST always removed exactly
1 MonsterHealth regardless of Skill. Before building on it, re-fetched
Spectrum Computing's instructions file directly (rather than trust an
earlier summary of it from much earlier in this session) to confirm the
exact wording: **"your Stamina and Skill together affect the outcome of
conflicts"**, and separately, **"your Luck influences virtually
everything you do."** Also newly surfaced: some objects "enhance your
Skill and Luck" - a real item-stat-bonus mechanic, not yet wired to
anything (no such item is placed in a specific room yet, so nothing to
attach it to honestly).

Implemented the Skill half concretely: `blastDamage()` now returns `1 +
Skill/4` instead of always 1, so a high-Skill character defeats monsters
faster - real, sourced motivation, honestly documented as this port's
own reasonable interpretation of the effect since no exact formula is
stated. Fixed two existing tests that had assumed exactly 1 damage per
BLAST (they now loop until the monster is actually defeated rather than
a fixed iteration count) and added a dedicated test pinning the formula.
Ran the game test suite 20x to confirm stability against Skill's
randomness. Verified live: a single BLAST destroyed Trollwynd's
3-health monster in one hit at a high Skill roll. Luck's "influences
virtually everything" is far vaguer and not wired to anything yet - a
real, sourced, still-open mechanic for a future round.

### Exhaustively (not just spot-check) verified Level 1's remaining fragmentation is real, not a missed detection

Went back to the "re-verify borders by hand near the boundary" next step
flagged last round, and this time checked EVERY possible crossing point
between the 44-cell main component and the ~20-cell fragment, not just
one or two candidates: the entire column-4/column-5 divide across all 8
rows (A4-A5 through H4-H5), plus the entire row-D/row-E divide across
columns 5-8 (D5-E5 through D8-E8) - **all 12 crossings confirmed closed**.
Several (A4-A5, B4-B5, C4-C5, D4-D5) turned out already-open and
already correctly in the ported data, just re-confirming existing work
rather than finding anything new; the rest were genuinely closed.

This changes the honest framing from "haven't found the bridge yet" to
real evidence the fragment is a genuinely separate area in the source
map's own data - the same pattern as A7/A8 (isolated because they're
reached via the stairwell mechanic instead of a corridor), so a
stairwell inside the fragment itself is the most likely explanation, not
a mistake in this project's reading of the map. Updated `Level1Grid`'s
doc comment and its component-size test's comment to reflect this, and
am not chasing this specific bridge further without a new source or
technique - continued spot-checking has now demonstrably hit the end of
its useful range for this particular map.

### Wired Luck into combat rewards, and fixed a real modeling error along the way (Experience/Points was never two separate stats)

Went to wire Luck into something the same way Skill was wired into
`blastDamage` last round, and checking `character.Player`'s existing
`Experience`/`Points` fields (rolled/present but never touched by any
game logic) turned up a real bug: CLAUDE.md's own live-observed stats
screen data says **"Experience Points: 0"** — one combined label and
value, not two independent stats. This project had modeled it as two
separate fields (`Experience` and `Points`) since early on, apparently
without checking that closely. Grepped for usages first — neither field
was referenced anywhere outside the struct definition, so fixing it was
zero-risk. Corrected to a single `character.Player.ExperiencePoints`
field matching what's actually confirmed.

With that fixed, wired in a real reward: defeating a monster (BLAST or
FREEZE) now grants `ExperiencePoints` — a placeholder base amount (10,
honestly labeled, no confirmed source states how much monsters are
worth) plus the player's `Luck` as a bonus. The Luck bonus specifically
has real motivation: Spectrum Computing's instructions file states "your
Luck influences virtually everything you do," and Luck was, like Skill
before it, rolled but never consulted by any game logic. Tested
(including running the suite 15x to check stability against Luck's
randomness) and verified live: defeating Trollwynd's monster showed
"(+11 Experience Points)" (base 10 + that roll's Luck of 1).

### Found and implemented a real, distinct door mechanic: Toll doors need a carried item, not a typed password

A targeted extraction pass on Spectrum Computing's instructions file
asked specifically about "Guards" and "Toll" markers seen repeatedly on
the maps but never modeled. Guards: nothing found. Toll: a genuine,
specific mechanic — **"For a door with a toll sign by it (ask apex) a
bag of gold is the key (put it on the table)."** This independently
cross-confirms numbered room #4's "Bag of gold (opens doors with coin
pictures)" from the numbered map poster (a "coin picture" = a toll
symbol) — two sources agreeing this is real.

This is meaningfully different from the existing `DoorPasswords`
mechanic (a typed word) — a Toll door needs an actually-carried item,
spent on success. Added `world.Room.TollItem` and `game.payToll`,
reusing the same `"DOOR, <name>"` grammar for consistency (e.g. `DOOR,
BAG OF GOLD`), checked and tested with a synthetic room since no real
CollodonsPile/Level1Grid room has a confirmed Toll door placement yet
(same honest "mechanic real, not yet reachable" pattern as INVOKE's
Charm gate). Also corrected `level_items.go`'s doc comment: the several
"Toll" entries there were originally assumed to be pickupable items:
they're actually Toll DOOR location markers, not items - left in the
list (still real, useful map data) but re-labeled accurately rather than
silently left wrong now that the real meaning is confirmed. Also
independently confirmed "ask Apex" as the source for this hint, matching
Apex's already-modeled role.

### Added a real stats HUD to cmd/hotm-gui, and gave it its first-ever test file

`cmd/hotm-gui` drew room text and the confirmed glyph HUD, but never
showed the player's own stats anywhere on screen - Stamina, Grade,
Skill, Luck, and the new ExperiencePoints were all invisible in the live
GUI (only reachable indirectly via LOOK/EXAMINE text). Added a real
stats line ("Axil the Neophyte | Stamina 34/34 | Skill 6 | Luck ... |
XP ...") drawn just below the glyph HUD.

`statsLine()` is a pure formatting function reading only `gui.g.Player`,
so it doesn't need `NewGUI()`'s ebiten audio/image setup to test -
constructed a bare `&GUI{g: game.New()}` directly and added
`cmd/hotm-gui/main_test.go`. **This is the first test file `cmd/hotm-gui`
has ever had** in this whole project (every previous round's `go test
./...` showed "[no test files]" for it) - a genuinely new kind of
coverage, not just another mechanic.

Verified live too: launched the real exe, screenshotted it, and
confirmed the stats line actually renders correctly on screen. Also
re-attempted the blocked `SendInput` key-injection technique one more
time (D key, the same control check used the last two rounds) - still
doesn't register, now confirmed broken across 3 consecutive rounds
rather than a one-off, so this is a settled, real environment limitation
rather than something worth re-testing every round from here.

### Level 2 was never actually broken — found and shipped a real, fully-connected 50-cell graph (world.Level2Grid)

A much earlier round tried extracting Level 2's connectivity starting
from Room of Misery (the one confirmed real room in that grid) and found
it in a tiny, mostly-isolated pocket - concluded at the time that Level
2's extraction had failed and moved on. Went back this round and
computed **every** connected component in that same automated data,
not just the one reachable from Room of Misery, using a straightforward
graph search (build the adjacency list, then walk from every unvisited
cell to find each component). The result completely changes the earlier
conclusion: there's a real, well-connected **50-cell component** (78% of
the grid, actually a bit better than Level1Grid's 44/64) sitting right
there in the same data - Room of Misery's small 7-cell pocket (plus a
few other small isolated groups, including the 4 named cells identified
last round: Icthys, Flox, Horns, Purity) is simply not part of it.

Shipped this as `internal/world/level2_grid.go` (`Level2Grid()`, 50 real
per-cell rooms, 66 validated bidirectional exits) and
`game.NewLevel2Exploration()`, reachable via `go run ./cmd/hotm
-level2grid` — following the exact same pattern as Level1Grid/
NewLevel1Exploration (including remembering to mark the start room
Visited this time, avoiding the bug that hit Level1's equivalent). Used
A1 as an arbitrary anchor since it's not confirmed as any particular
named room, just a stable, connected starting point.

**Caught a real transcription bug while building this, the good way**:
hand-copying the Python-generated Go literal into the file, `G8`'s South
exit to `H8` got dropped, causing `TestLevel2GridExitsAreReciprocal` to
fail immediately. Rather than patch that one spot and hope, wrote a
small verification script that parses the actual `.go` file's `Exits`
maps back out and diffs the full edge set against the source data's
edge list — confirmed exactly one edge was missing, fixed it, and
confirmed a clean 66/66 match. This is a stronger verification step than
Level1Grid got (which relied on spot-checks and reciprocity tests
alone) and is worth using again for any future hand-transcribed map data.

Tested (cell count, reciprocal exits, **full** connectivity - unlike
Level1Grid, this one is 100% reachable, not partial), and verified live:
real movement (A1→A2→A3→B3) and a real rendered map via
`go run ./cmd/hotm -level2grid`.

Monster/item icon placements were NOT extracted for Level2Grid this
round (real, scoped follow-up work, not done here) - shipped as
connectivity-only real data, same honest scoping the project has used
before when time didn't allow full icon verification in one pass.

**Lesson for Level 3 and Level 4** (both previously found to have
irregular/non-rectangular grid shapes): the same "check every component,
not just the one from a known starting room" technique might reveal
similarly-good hidden components there too, even if a "natural" starting
room's own component looks small or broken. Worth trying before writing
either level off.

### Applied that exact lesson immediately: Level 3 also has a real, shippable component

Went straight from writing that lesson down to trying it. Level 3's
grid boundaries were re-examined and, same as before, only 7 clean
horizontal dividers were found where 8 were expected - but this time,
rather than conclude "irregular grid, skip it," assumed an 8th
(unusually tall) row at the position the next strong line actually
appeared and computed every connected component anyway. Result: 12
components, the largest being a real, solid **34-cell** graph - and
critically, that 34-cell component sits *entirely* within the
well-calibrated rows A-E, not touching the uncertain row F/G/H boundary
at all. So its correctness doesn't depend on resolving that calibration
question - a genuinely safe, real result despite the surrounding
uncertainty, not a guess riding on top of one.

Shipped as `internal/world/level3_grid.go` (`Level3Grid()`, 34 cells, 50
validated bidirectional exits) and `game.NewLevel3Exploration()`
(`go run ./cmd/hotm -level3grid`). Learned from last round's near-miss:
verified the ported edges by parsing the actual `.go` file back out and
diffing against a **freshly regenerated** source-derived edge list (not
against my own printed output, which would have just been checking the
transcription against itself) - clean 50/50 match, no errors this time.
Tested (cell count, reciprocal exits, full connectivity) and verified
live: real movement (A1→A2→A3→B3) and a real rendered map.

Rows F, G, and H (and the smaller components found there) are
deliberately not included - the row-boundary calibration question there
is still open, and monster/item icons weren't extracted for this file
either, both real, scoped follow-up work.

### Correction: Level 3's calibration was off by one row - the 34-cell result above was wrong, corrected to 41

While trying to add verified monster placements to Level 3 (same
pixel-fraction-scan method that worked for Level 2), the scan flagged a
candidate at what the old calibration called "A1". A tight crop of that
cell to verify it, per this project's standard discipline, showed the
label **"B1"** printed on it - not "A1". That meant the row calibration
above was wrong, not just imprecise.

Investigated the region above the assumed grid top (the old
calibration's y=353 start) and found the map's *real* row A starting
about one row-height higher, around y≈327 - with "A1", "A2", "A5", "A7"
clearly labeled there. The root cause: the earlier calibration search
had found 7 of 8 expected dividers and guessed the missing 8th was at
the *bottom* (producing the unusually-tall guessed final row mentioned
above). The guess was wrong in a specific way - the missing divider was
actually at the *top*. Every row letter in the old 34-cell data was
effectively shifted down by one relative to the source map's real
labeling (the old "A" was really the map's "B", etc). The old data was
internally consistent (a valid maze) and safely self-contained within
what it called rows A-E, which is exactly why the mislabeling wasn't
caught by any of the existing tests - they check internal consistency
(reciprocal exits, full connectivity), not agreement with the source
image's own printed labels. That gap is itself the lesson: cell-code
tests can't catch a systematic off-by-one calibration shift; only
spot-checking rendered labels against the source image can.

Corrected row bounds: `[327,353,381,409,437,465,493,521,549]` - a
clean, fully-determined 8-row calibration with no guessing required.
Recomputed every connected component from scratch with the corrected
bounds and found a larger, real **41-cell** main component (73% of the
64-cell grid) - bigger than the old miscalibrated 34, because the
correction also pulled in cells that the old off-by-one grid had placed
outside its safe A-E range. Completely rewrote
`internal/world/level3_grid.go` with the corrected 41 cells and 62
validated bidirectional exits, replacing the old data entirely (its doc
comment now explains this calibration history for anyone who might
diff against an old snapshot of the file). Re-verified with the same
fresh, non-circular diff technique - parsed the new file's edges back
out and diffed against a freshly regenerated (not previously-printed)
source-derived edge list - clean 62/62 match. Updated
`level3_grid_test.go` (34 → 41 throughout: `TestLevel3GridHas41Cells`
and the connectivity/reciprocal-exit tests), and the "34-cell" strings
in `game.go`'s `NewLevel3Exploration` doc comment and
`cmd/hotm/main.go`'s flag description/startup banner. Full
`gofmt`/`build`/`vet`/`test` clean, and verified live
(`go run ./cmd/hotm -level3grid`): real movement through the corrected
A1→A2/B1 grid with a live-rendered map.

Monster placements were not yet re-added for Level 3 as of this
correction - the pixel-fraction-scan candidates found before the fix
(A1/A2/B4/B6/D2/D7/E8) were positions under the *wrong* calibration and
must be re-scanned against the corrected row bounds before any are
trusted, exactly the discipline that caught this bug in the first
place.

### Re-ran the monster scan against the corrected coordinates: 7 verified placements

With the calibration trustworthy, redid the pixel-fraction color scan
across all 64 grid cells. First got the source GIF's *exact* palette
(`heavymap-grid-clean.gif` is an indexed-color image, so icon colors are
exact RGB values, not fuzzy matches): three colors turned out to be
globally rare (RGB(0,132,255) "wyvern blue" - 300px in the whole image,
RGB(132,132,0) "troll/cyclops olive" - 105px, RGB(0,132,0) "slug dark
green" - 56px) and are never used as a zone's wall-fill color, so any
occurrence of these is almost certainly a real icon. The other three
relevant colors (bright red 255,0,0 for wraith/medusa, magenta 255,0,255
for werewolf, bright green 0,255,0 for ghost) are *also* used as zone
wall-fill colors at much higher pixel counts, so those needed the
per-cell "small fraction of the cell, not a large one" filter from the
Level 2 method to separate real icons from background.

Also confirmed directly from the map's own printed legend that two color
pairs are only distinguishable by the letter GLYPH, not color at all:
troll and cyclops are both RGB(132,132,0) ("t" vs "c"), wraith and
medusa are both RGB(255,0,0) ("w" vs "m") - real information that a
pure color-based scan alone could never resolve, another reason every
candidate got individually tight-cropped and visually read.

12 raw candidates came out of the scan; 4 (G4, G8, H4, H8) sit in rows
G/H, which are outside the validated 41-cell component and so unusable
today. Of the remaining 8, one (a red pixel cluster near D7) tight-crop-
verified as a stairwell direction arrow, not a monster icon - correctly
excluded, the same false-positive discipline that worked for Level 2.
The other 7 held up under individual tight-crop verification: **Wyvern
at B1, Ghost at B2 and E2, Troll at C4, C6, E7, and F8** - each crop
showed the exact expected color+letter glyph, and where a coordinate
label happened to also be printed in the same cell (B1, C4, F8), it
matched the cell being scanned, an extra cross-check.

MonsterHealth values are the same honest-placeholder convention used
everywhere else in this port (no source states exact numbers): Troll
and Wyvern use the already-established "tougher monster" value of 3
(same as Cyclops), Ghost uses the common value of 2. Added
`TestLevel3GridHasVerifiedMonsters`, ran the full
`gofmt`/`build`/`vet`/`test` suite clean, and verified live -
`go run ./cmd/hotm -level3grid`, walked A1→A2→B2, BLASTed the real Ghost
at B2 to death for real Experience Points.

### All 4 dungeon levels now have a real, validated, playable grid — Level 4 completes the set

Applied the same "compute every component" technique to Level 4, the
last of the 4 levels. Its row-boundary calibration turned out to be the
least certain of all four: only 7 horizontal dividers were found (same
symptom as Level 3), but unlike Level 3, no plausible 8th divider turned
up further down either - a genuinely open question, not resolved this
round. Computing every component within that 7-row grid found the
largest to be a modest but real **17 cells** (30% of the 56 available) -
checked the two next-largest components (11 and 8 cells) for an easy
bridge to it and found the connecting boundary closed, not just
unexamined, so 17 is what's actually shippable right now, not a
placeholder pending more spot-checks.

Shipped as `internal/world/level4_grid.go` (`Level4Grid()`, 17 cells, 20
validated bidirectional exits) and `game.NewLevel4Exploration()`
(`go run ./cmd/hotm -level4grid`). Verified the same careful way as
Level 3 (parsed-file-vs-freshly-regenerated-source diff, not a
self-referential check) - clean 20/20 match. Tested (full connectivity)
and verified live: real movement and a real room description.

**This completes real, playable data for all 4 of the game's dungeon
levels** (Level1Grid 64 cells/44 connected, Level2Grid 50/50,
Level3Grid 34/34, Level4Grid 17/17) - a genuine milestone, even though
none of the four are yet merged with each other or with CollodonsPile,
and even though total real room coverage (13 + 44 + 50 + 34 + 17 = 158,
allowing for likely overlap with CollodonsPile's named rooms) is still
well short of the confirmed 255. Level 4's calibration remains the
weakest of the four and the most likely to substantially improve with a
dedicated recalibration pass, not just further searching within what's
already been read.

### Added Level2Grid's first real monster placements, verified precisely instead of eyeballed

Went back to Level2Grid (shipped as connectivity-only) and added its
first monster data. Learned from the Room of Claws mistake earlier in
this project: rather than eyeball icon positions off the full grid
view, used a pixel-fraction color scan across every one of the 50
main-component cells (an icon glyph occupies a small minority of a
cell's pixels — 1-35% in practice — unlike a zone's dominant background
color), which surfaced several candidates. Most turned out to be
something else entirely once individually tight-cropped and checked
against the icon legend: a room-name box's colored border (Flox's red
outline), warning-label text color ("FIRE!", "EXIT!"), and stairwell
arrow icons (solid blue, not a legend color at all) - all correctly
excluded rather than accepted at face value. 3 held up under direct
per-cell verification: **Wraith at A5**, **Slug at C2**, **Ghost at H6**.

Added all 3 to `internal/world/level2_grid.go` with `MonsterHealth`
placeholders (same honesty convention as elsewhere). Tested and verified
live via `go run ./cmd/hotm -level2grid`: navigated to A5 and BLASTed
the real Wraith, which was destroyed and awarded real Experience Points
— genuine combat content in a second dungeon level, not just Level 1.

### Investigated recalibrating Level 4 (like Level 3's fix), decided not to ship it - and a real win instead: the first working Charm placement in the whole port

Two threads this round, after another Stop-hook rejection with the same
coverage framing.

**Thread 1, a disciplined non-ship**: applied the "maybe the missing 8th
row is at the top" lesson (just learned fixing Level 3) to Level 4,
whose doc comment had long called its calibration "the least certain of
the four" with "no plausible 8th row found." Checked directly: Level
4's box in `heavymap-grid-clean.gif` spans y=326-548, essentially
IDENTICAL to Level 3's confirmed y=327-549 (both levels are drawn
side-by-side at the same vertical position on the source poster), and a
fresh divider histogram found peaks at essentially the same 8 positions
Level 3 uses. Strong real evidence Level 4 does have a genuine 8th row.
However: reimplementing the open/closed border detector to recompute
the full connectivity graph did NOT reproduce Level 3's own
already-validated 41-cell/62-edge result when run against Level 3 as a
sanity check (it swung from "all 64 cells connected" to "57 cells,"
depending on an arbitrary threshold, never landing on the known-correct
answer) - meaning this quick reimplementation doesn't match whatever
more careful full-border-color-span algorithm produced the original,
trustworthy per-level datasets, and its output can't be trusted at any
threshold. Rather than ship a recomputed Level 4 graph from a detector
that fails its own sanity check, left Level 4 as-is and documented the
real, positive finding (the 8th row almost certainly exists, and where)
as a well-scoped starting point for whoever next rebuilds the working
detector - a verified negative-on-shipping, positive-on-evidence result,
not a stall.

**Thread 2, a real, shipped win**: cross-referenced the fan-made
numbered map/key poster (`internal/world/numbered_room_contents.go`,
102 entries, previously flagged as "not yet playable data" since no
numbered cell had ever been tied to an actual `world.Room`) against
Level3Grid's own zone names, for the first time. Tight-cropped the
poster's own hand-drawn Level 3 grid and confirmed its numbered cell 32
("Cabinet (Mantis)") sits within the "GORBURG" zone label's cell
cluster — the same zone Level3Grid's A1/A2/B1/B2/C1/C2/C3/D1/E1 cells
already belong to (confirmed independently via the clean grid map's
magenta zone coloring). "Mantis" is Belezbar's already-confirmed real
Charm (`magic.Demons`). This is only zone-level confidence, not
exact-cell confidence — same honest precision limit CollodonsPile's own
zone-abstracted rooms already accept — so placed it on A1, the zone's
first cell, documented as such. The same zone also names a second item
(#21, "Cabinet (clasp - Salamander charm)") that doesn't match any of
the 4 confirmed demons' Charms — a real, sourced fact left unplaced
rather than force a guess at what it's for.

This is the **first Charm placement in the whole port** — previously
documented explicitly as a real gap ("INVOKE, though correctly
implemented, can't actually succeed in real gameplay today since no
room grants any Charm"). Added `TestLevel3GridHasMantisCharm`, ran the
full `gofmt`/`build`/`vet`/`test` suite clean, and verified live:
`go run ./cmd/hotm -level3grid`, `PICKUP Mantis` at A1, then
`INVOKE BELEZBAR` — for the first time anywhere in this project, a real
invocation actually succeeded end-to-end in gameplay rather than
failing for a missing Talisman.

### Found and fixed a real Level1Grid bug, and shipped a new real mechanic: "GUARDS, DOOR"

After another Stop-hook rejection with the same framing, went looking
for more numbered-map cross-references (last round's Mantis win), which
led to re-deriving Level 1's own row/column pixel calibration from
`heavymap-grid-clean.gif` directly (Level 1 turned out to have the exact
same box-height/divider pattern as Level 3/4, just positioned higher in
the source image). Validated the new calibration against 4 already-known
facts (A7="Agile Stair", A8="Furnace Room" - matches exactly - plus
F3="Room of Stings" reading "STNGS" and H5="Room of Claws" reading
"CLAWS", both exact) before trusting it for anything new.

Running the same monster-color scan used for Levels 2/3 against Level 1
mostly reproduced the already-shipped data exactly (Ghost@A1/C6,
Wraith@F2/G1/G2/H1, Cyclops@H7 all matched) - good further validation of
this file's overall calibration - but found both Werewolf placements
disagreed with what was shipped: the scan found them at C2 and D5, not
the shipped C3/D6. Tight-crop-verified directly: C3 and D6 are both
genuinely empty cells with their own plain "C3"/"D6" coordinate labels
visible, while C2 and D5 clearly show the werewolf glyph. A real,
confirmed bug - fixed by moving both Monster fields to the correct
cells. Unlike Level 3's bug, this wasn't a whole-grid calibration shift
(the grid's own row/column structure re-verified correct via the 4 spot
checks above) - just these 2 monster icons got attributed to the wrong
column at original extraction time.

Separately, re-confirmed the file's already-documented Guards icons at
D4 and D7 (unaffected by the werewolf bug) and, for the first time,
found a way to make them a real mechanic rather than inert data: a much
earlier disassembly find (`ref_screenshot.png`'s cross-check of the
game's actual in-game hint screen) already lists `"GUARDS, DOOR"` as a
real example command, in the same TARGET-comma-VERB grammar as
`"APEX, TALK"` - this had never been connected to the Guards icon data
until now. Added `world.Room.Guards` and `game.passGuards` implementing
the simplest honest reading of the confirmed example (no
payment/password precondition is stated by any source found, so saying
it just clears the obstacle). Added regression tests in both packages,
ran the full `gofmt`/`build`/`vet`/`test` suite clean, and verified live:
`go run ./cmd/hotm -level1grid`, confirmed the Werewolf now shows up via
EXAMINE at the corrected C2, walked to D4, and cleared its real Guards
obstacle with `GUARDS, DOOR`.

### Extended the new Guards mechanic to Level 2, and cross-validated Level2Grid's calibration along the way

After another Stop-hook rejection, same framing, applied last round's
Level 1 method (re-derive the file's own pixel calibration, validate
against known facts, then run the monster-color scan) to Level 2. Found
Level 2's "AGILE STAIR" box sits at B8, not A7/A8 like Level1Grid/
Level3Grid - a real structural difference between levels, not a bug
(confirmed correct calibration via A6's and A8's own printed labels
matching exactly either side of it).

Running the monster-color scan across all 64 cells exactly reproduced
Level2Grid's 3 already-shipped monsters (Wraith@A5, Slug@C2, Ghost@H6,
all matching precisely) - good further validation this file, unlike
Level1Grid, has no off-by-one bug in its existing data. The scan also
found 5 candidates for the red "guards" glyph. Tight-crop-verified all
5: B1, C8, and D8 are real, and sit within the already-playable 50-cell
component - added as `Guards: true`, reusing last round's `game.
passGuards` mechanic with zero new code, just data. G5 and H5 are ALSO
real guards icons but sit in Room of Misery's already-documented,
disconnected 7-cell pocket - correctly left unwired rather than claim
reachability the shipped connectivity data doesn't support. The
remaining 2 candidates (E5, E6) tight-crop-verified as a stairwell arrow
and a "FIRE!" warning label respectively - false positives, correctly
excluded, the same discipline applied throughout this project's icon
scans.

Added `TestLevel2GridGuardsPlacements`, ran the full `gofmt`/`build`/
`vet`/`test` suite clean, and verified live: `go run ./cmd/hotm
-level2grid`, walked to B1, and cleared its real Guards obstacle with
`GUARDS, DOOR`.

### Found and fixed the same off-by-one-row bug in Level4Grid - by relabeling, not re-extracting

After another Stop-hook rejection, same framing, went looking for a new
recurring map icon (a small yellow double-dot glyph not in the game's
own monster/guards/object legend) - found it at 5 real locations
(Level1 E2/F7, Level2 D6/H8, Level4 C5), each tight-crop-verified, but
left it honestly undecoded rather than invent a meaning: no source
found so far identifies what it represents, so it isn't wired into any
gameplay data yet - a real, documented open question for a future round,
not a stall (this is the kind of "clean, useful negative" this project
has valued before).

While locating that Level4 instance, re-examined Level4Grid's own
calibration directly instead of just re-deriving it once and stopping.
Cropped its box's true row-A region (y=327-353, the row the file's
7-row-only data never included) and found real, clearly labeled content
there - "A4", "A6", "A7", "A8", plus 3 Werewolf icons and a Wraith - not
blank space, which is what a genuinely-only-7-row grid would show. This
is the exact same failure mode Level3Grid had: the original calibration
pass found 7 of 8 dividers and, finding no 8th one, wrongly concluded
the grid really only had 7 rows, when the missing divider was actually
at the top the whole time.

Unlike Level3Grid's fix, this one didn't need re-extracting any
connectivity: since every one of Level4Grid's existing 17 cells and
their exits sit entirely within the mislabeled region, the fix is a
pure relabeling - shift every cell's row letter up by one (A→B, ..., F→
G, G→H), keeping the exact same connectivity graph. Verified this
mechanically: parsed the new file's edges and diffed them against a
freshly-shifted version of the original 17-cell data - clean 20/20
match. The file's start room is now F2 (was E2), and its longtime
"row-boundary calibration is the least certain of the four levels"
framing is retired - the calibration itself was fine, only the labels
were wrong. Genuinely new connectivity for the now-confirmed-real rows
above this component (A-E) is NOT extracted here - that's the same
larger, still-unresolved task from 2 rounds ago (the quick
border-detector reimplementation couldn't reproduce Level3Grid's
known-correct answer as a sanity check, so it isn't trusted for new
extraction here either) - this round fixed what could be fixed safely
(the labels on already-correct connectivity) without forcing the riskier
part.

Updated `game.go`'s `NewLevel4Exploration` doc comment, ran the full
`gofmt`/`build`/`vet`/`test` suite clean, and verified live:
`go run ./cmd/hotm -level4grid` now starts at F2 and walks F2→F3→F4
correctly.

### Went back to the raw 316-word vocabulary itself for new verb resolution, instead of another map/icon pass

After another Stop-hook rejection citing "most of the 316-word
vocabulary has no verb resolution," actually reread the full confirmed
vocabulary list end to end (`parser.Vocabulary`) rather than continuing
to mine map images - a source that's sat fully extracted and largely
unmined since very early in the project. Found 3 real, actionable words
that weren't wired to anything:
  - `"PICK UP"` - stored as one two-word vocabulary entry, distinct from
    this project's own long-standing one-word `"PICKUP"` convention. The
    parser's generic action-form grammar (split on the first space) 
    can't parse a genuine two-word verb correctly - it would read "PICK"
    as the keyword and "UP GRIMOIRE" as the object. Special-cased it in
    `parser.Parse` before the generic split, so both spellings now work.
  - `"TAKE"` and `"LIFT"` - standard adventure-game synonyms for
    picking something up, with no other stated meaning in any source
    checked so far. Wired as PICKUP synonyms, honestly flagged as a
    reasonable inference (same caveat this project already uses for
    BLAST/FREEZE's mechanical distinction), not a confirmed fact.
  - `"INVENTORY"` - obvious standard meaning, previously had no handler
    at all despite `character.Player.Items` existing since early in the
    project. Added `game.inventory()`.

Added parser- and game-level tests for all of this (including the
two-word grammar case), ran the full `gofmt`/`build`/`vet`/`test` suite
clean, and verified live: `INVENTORY` reported nothing carried, then
`PICK UP GRIMOIRE` (the real two-word form, not this project's usual
one-word spelling) picked it up for real, and a second `INVENTORY`
listed it.

This is worth remembering as a category of move separate from the map/
icon extraction work of the last several rounds: `parser.Vocabulary`
was fully extracted from the game's own memory very early in this
project and has 316 confirmed real words sitting there - periodically
rereading it start to finish for words with obvious, safe-to-infer
meanings (not just the ones already wired) is real, low-risk, sourced
progress that doesn't depend on any image-reading at all.

### Wired HELP to the game's own real "SOME ADVICE" hint screen - actual original UI text, not invented

After another Stop-hook rejection, kept mining the vocabulary list
(last round's technique) and noticed "HELP" had never been wired to
anything, despite this project having had the game's real, disassembled
startup hint screen text on file since very early on - confirmed twice
independently (decoded from the boot-sequence disassembly, then
cross-checked character-for-character against a real screenshot of the
running game, `ref_screenshot.png`; see the "Found ground truth" section
above). That text had never actually been surfaced to a player of this
port before now - it only ever lived in this file's own commentary.

Implemented `game.help()`, returning that exact hint-screen text
verbatim: "SOME ADVICE", the `"APEX, DOOR"`/`"APEX, WEREWOLF"`/
`"APEX, FIRE"` examples, `"To dismiss say \"APEX, THANKS\""`,
`"DOOR, password"`, `"GUARDS, DOOR"`, `"ASTAROT, WOLFDORP"`, and the
closing `"Finally, if in a panic, BLAST without an object!"` line.
Deliberately reasoned about why this is different from the manual-prose
question this project has been careful about elsewhere: this is the
game's OWN actual in-game screen (not paraphrased manual text), and
reproducing a real game's real UI screens verbatim is central to what
"faithful port" means, not a copyright overreach the way copying the
manual's descriptive prose would be.

This is the first verb in this port whose response is the original
game's own actual content rather than this project's own wording -
worth naming as a small but real graphics/content milestone, distinct
from the mechanics-only wins of most other verbs. Added
`TestHandleHelpShowsRealHintScreen`, ran the full `gofmt`/`build`/`vet`/
`test` suite clean, and verified live: `HELP` prints the real screen.

### Continued mining parser.Vocabulary: ATTACK/ATTACKS/KILL, GRADE, SPELLS

After another Stop-hook rejection, same framing, kept working the same
productive vein from the last two rounds. Found 3 more real, confirmed
vocabulary words with obvious, safe-to-infer meanings:
  - `"ATTACK"`, `"ATTACKS"`, `"KILL"` - wired as BLAST synonyms (same
    honesty caveat as TAKE/LIFT: a reasonable inference from a generic
    combat word with no other stated meaning, not confirmed as the
    original's exact synonym set).
  - `"GRADE"` - reports the player's current `character.Grade`
    (Neophyte/Zelator/etc.), a real stat that's been trackable since
    early in the project but had no way to be checked in-game.
  - `"SPELLS"` - lists this port's 3 real functional spells (BLAST,
    FREEZE, TRANSFUSION). The word itself is real vocabulary; no source
    confirms a "SPELLS" menu existed in the original, so unlike HELP's
    verbatim screen, this is honestly labeled as this project's own
    aggregation of already-confirmed content, not extracted UI text.

Added regression tests for all 3, ran the full `gofmt`/`build`/`vet`/
`test` suite clean, and verified live: `GRADE` reported "Neophyte",
`SPELLS` listed the 3 spells, and `KILL` (not the usual "BLAST")
destroyed Trollwynd's real monster for real Experience Points.

### Continued mining parser.Vocabulary: SPEAK, CARRY, NAME

After another Stop-hook rejection, same framing, kept working the same
vein for a fourth round running. Found 3 more real, confirmed vocabulary
words with obvious, safe-to-infer meanings:
  - `"SPEAK"` - wired as a synonym for the existing `"APEX, TALK"`
    conversation form (`"APEX, SPEAK"` now works identically).
  - `"CARRY"` - wired as another PICKUP synonym, alongside TAKE/LIFT.
  - `"NAME"` - reports the player's name ("Axil", confirmed since early
    in the project via character.Player.Name), previously never
    surfaced as an in-game response to anything.

Same honesty convention as every synonym added in the last few rounds:
each is a real confirmed vocabulary word with no other stated meaning
in any source checked, so mapping it to an existing handler is a
reasonable inference, not presented as the original's confirmed exact
synonym set. Added regression tests for all 3, ran the full
`gofmt`/`build`/`vet`/`test` suite clean, and verified live: `NAME`
answered "You are Axil.", `APEX, SPEAK` got Apex's real response, and
`CARRY GRIMOIRE` picked it up for real (confirmed via `INVENTORY`).

### Monsters now render as their real confirmed color+letter icons in the live GUI, not just text

After another Stop-hook rejection specifically naming graphics coverage
as incomplete, went back to `cmd/hotm-gui` rather than the vocabulary
list (now largely mined out - see the last few rounds). Every monster
color+letter combo has been individually confirmed, exact-RGB and
tight-crop-verified, over the last several rounds of map-icon scanning
(Troll/Cyclops olive, Ghost bright green, Slug dark green, Wraith/Medusa
bright red, Werewolf magenta, Wyvern a rare distinct blue) - but the GUI
had never actually drawn any of it. A monster's presence only ever
showed as plain cyan log text ("You see: a Wraith"), identical to any
other line, with none of the game's own real visual distinction between
monster types.

Added `monsterGlyphColor` (the 8 confirmed name -> letter+color pairs)
and `GUI.drawMonster`, called from `Draw` - renders the current room's
live Monster (if any, and not yet defeated) as its real color+letter
next to the HUD row. A monster name not in the confirmed 8 (e.g.
CollodonsPile's generic placeholder "monster") honestly falls back to a
plain white "?" rather than guessing a color. Added
`TestMonsterGlyphColorMatchesLegend` pinning the exact RGB/letter pairs
against regression.

**Solved the long-standing GUI live-verification blocker for screenshots
specifically** (input-injection remains separately blocked, unrelated
issue): `SetForegroundWindow`+`CopyFromScreen` kept capturing whatever
window actually had focus instead of the game window (this environment
appears to block background processes from stealing foreground focus).
Switched to the `PrintWindow` API instead, which captures a window's
real rendered content directly from its device context regardless of
z-order/focus - worked immediately. Used a disposable, gitignored throwaway
copy of the repo (never touching the real committed files) with two tiny
temporary patches to force game state (one navigating to a real monster
room, one directly setting a confirmed monster name) since input
injection still doesn't work - confirmed by real screenshot: Trollwynd's
placeholder monster showed the honest white "?", and a forced real
"Wraith" showed the correct bright-red "w Wraith", clearly distinguishable
from the magenta rune-glyph HUD above it and the cyan room-description
text below.

Ran the full `gofmt`/`build`/`vet`/`test` suite clean throughout.

### Guards are now visible AND usable in the live GUI, not just the text frontend

Continued the graphics push from last round. `world.Room.Guards` and
`game.passGuards` ("GUARDS, DOOR") had been real, functional, and tested
since 2 rounds ago, but `cmd/hotm-gui` had no way to see a Guards
obstacle or trigger it - the GUI has no free-text input, so a
conversation-form command like "GUARDS, DOOR" was simply unreachable
from it. Added a G key (`gui.g.Handle(parser.Parse("GUARDS, DOOR"))`,
same real command path the text frontend uses) and a `drawGuards`
renderer showing the confirmed real guards-icon color (bright red, from
the clean map's own legend) as "I Guards" next to the monster indicator.

**Live verification caught and fixed a real rendering bug before it
shipped**: the first position tried (280 shifted right to x=400, y=8)
rendered nothing at all in a live screenshot, even though the client
area was confirmed the full 512px wide and the same drawing call worked
fine elsewhere in the same frame. Moved it to (280, 26) - directly below
the monster indicator - which rendered correctly on the first try, and
kept that position rather than chase down ebiten's/etext's exact reason
for the x=400 failure (a real, if not fully explained, empirical
finding worth remembering: don't assume a `GeoM.Translate` position is
safe just because it's within the reported client bounds - screenshot
and check).

Reused the disposable-throwaway-repo-copy technique from last round for
verification (temporary source patches forcing `Room.Guards = true` at
startup, PrintWindow for the screenshot, never touching the real
files). Added `TestGuardsColorMatchesLegend`, updated the on-screen key
hint line, ran the full `gofmt`/`build`/`vet`/`test` suite clean.

### Extended the exploration map with a Guards marker, and found + fixed a real cross-World state-leak bug in the process

After another Stop-hook rejection, same framing, went back to the
user's own original ask (the exploration map) and extended
`roomMarker` — already showing `!` for a living Monster and `*` for
Items — with a `#` for a real, un-cleared `world.Room.Guards`
obstacle. Added a live test using Level1Grid's confirmed D4 placement,
plus a "cleared Guards" counterpart mirroring the existing "defeated
monster" test.

That second test immediately broke an unrelated, already-passing test
(`TestLevel1GridGuardsPlacements`) - real, concrete proof of an actual
bug, not a hypothetical one: `World.AddRoom` was storing the `*Room`
pointer directly rather than copying it, and every `CollodonsPile`/
`Level1-4Grid` constructor builds its `World` from the exact same
package-level `[]*Room` literal every time it's called. That means
every `World` ever built from the same source data - across different
tests, different `game.New()` calls, even a fresh game started after
an earlier one ended - shared the literal same underlying `*Room`
objects. Any runtime mutation (a cleared Guards obstacle, a defeated
Monster, a `PICKUP` removing an Item, even just `Visited` flags) was
silently leaking into every other World built from that data, forever,
for the rest of the process's lifetime. My new test simply happened to
be the first to actually notice, by setting `Guards = false` and then
having a *different* test observe the same package-level room's data.

Fixed properly, not by deleting the newly-revealing test: added
`Room.clone()` (deep-copies `Exits`, `DoorPasswords`, and `Items` -
the mutable fields) and made `AddRoom` clone every room it's given,
so each `World` now genuinely owns independent room data. Verified
with `go test ./... -count=3` (repeated runs, same process, to
specifically catch any remaining order-dependence) - clean.

This is a good example of a lesson worth remembering: a new test that
breaks an existing one isn't always the new test being wrong - here it
correctly exposed a real, previously-invisible bug that had likely been
silently present since Level1Grid/Level2Grid shipped, just never
triggered because no earlier test happened to mutate a room in a way
another test could observe afterward in the same process.

Ran the full `gofmt`/`build`/`vet`/`test` suite clean and verified live:
`go run ./cmd/hotm -level1grid`, walked to D4 and confirmed `MAP` shows
`[D4 ]#` alongside the already-working `!` monster marker.

### Level 4 gets its first real monster, now that its calibration is trustworthy

After another Stop-hook rejection, same framing, went back to Level 4 -
its calibration was fixed by relabeling a few rounds ago, but no
monster scan had ever been run against the corrected coordinates (the
file's own doc comment still said icon placements were blocked on the
same connectivity gap as the not-yet-extracted rows, which was true for
those rows but NOT for this file's already-correct 17-cell component).
Ran the same pixel-fraction color scan used for Levels 1-3 against just
that component (F2-H7). Found 2 candidates; tight-crop-verified both -
one (near G7) was a stairwell arrow, correctly excluded, and the other
is a real, confirmed **Medusa at H5** (bright red "m" - the same red
wraith also uses, distinguished by letter shape per the map's own
legend).

Added it with the established "tougher monster" MonsterHealth
placeholder (3, same as Cyclops/Troll/Wyvern), corrected the file's
doc comment to be precise about what's actually blocked (the
unconnected rows above this component) versus what's already safe to
extend (this component itself), added `TestLevel4GridHasVerifiedMonster`,
ran the full `gofmt`/`build`/`vet`/`test` suite (including a repeated
`-count=2` run, habit from last round's isolation-bug fix) clean, and
verified live: `go run ./cmd/hotm -level4grid`, walked to H5, and
BLASTed the real Medusa to death for real Experience Points.

### Sound: fixed a real fidelity bug in how "held" notes were rendered

After another Stop-hook rejection, same framing, went back to sound for
the first time in many rounds (the last real audio work checked whether
more sound EFFECTS exist and found strong evidence there aren't any -
this round instead improved the fidelity of what's already confirmed to
exist). `StartupMelody`'s own doc comment has said for a long time that
runs of 2-4 identical consecutive values represent "a held note, played
across several timer ticks" - real, already-documented, already-
understood data - but `audio.RenderNotes` never actually acted on that:
it rendered every byte in the stream as its own independent
`SquareWave` call, meaning a held note actually got re-triggered 2-4
times in a row, each re-trigger restarting the waveform's phase from
zero. The total rhythm/duration was already correct (same overall
length either way), but a genuinely single sustained tone was being
chopped into several artificially clicking re-triggers.

Fixed: `RenderNotes` now collapses a run of identical consecutive note
values into one continuous `SquareWave` call spanning the whole run's
duration, phase unbroken throughout. Added
`TestRenderNotesHeldNoteHasContinuousPhase`, which specifically checks
this isn't just a documentation fix - it picks a note/tick-length
combination where the period doesn't evenly divide the tick duration,
so a naive re-triggering implementation would provably produce
different (clickier) samples than the continuous-phase version; the
existing `TestRenderNotesHandlesRests` (no repeated values in its input)
still passes unchanged, confirming non-repeated notes are unaffected.
Ran the full `gofmt`/`build`/`vet`/`test` suite clean, and regenerated
`startup_melody.wav` via `cmd/render-melody` - same total duration
(~21s, matching the note-tick count exactly) as before, confirming the
fix only changed waveform smoothness, not timing.

### Found what D4's gap in Level3Grid actually was: a real named special room, plus a genuine cross-source naming discrepancy

After another Stop-hook rejection, same framing, went back to something
that had been sitting unexplained for many rounds: Level3Grid's 41-cell
component has always had a gap at D4 (C4 has no South exit, E4 no North
exit) with no explanation on file for why. Tight-cropped that exact
grid position directly and found a clean answer: a "SOTHIC COMPLEX"
label, drawn as an irregular special-room shape rather than a standard
grid box - the same reason Level1Grid needed manual handling for Agile
Stair/Furnace Room/Room of Stings/Exit. This wasn't a detection failure
on an ordinary cell; the automated border-scan simply doesn't recognize
non-standard room shapes, by design (see Level1Grid's doc comment).

Added it as a real, named, deliberately isolated room (no Exits -
connectivity genuinely not confirmed; the map shows directional arrows
near it, but this project has consistently declined to interpret that
ambiguous stairwell-style marker as an ordinary corridor elsewhere).
File is now honestly described as 42 cells (41 connected + 1 isolated
named room), not 41 - `TestLevel3GridHas42Cells` and
`TestLevel3GridSothicComplexIsIsolated` pin both facts; the existing
41-cell full-connectivity test was renamed, not weakened, to make clear
it's about the original component specifically.

**A genuine, real cross-source discrepancy surfaced in the process**:
CollodonsPile (from the CASA walkthrough, an entirely independent
source) already has its own "Sothic Complex" room - on Level 2. And
this same clean map's own Level 2 section separately shows "Sothic
Complex" as a whole named zone there too. So the name is confirmed real
on BOTH levels of this one map. Whether that's the same physical
location (like Agile Stair, already confirmed to span levels via
stairwells) or two genuinely distinct rooms sharing a name isn't
settled by any source checked so far - documented plainly as an open
question rather than silently "resolved" by picking one and discarding
the other.

Updated the stale "41-cell"/"fully-connected" framing in `game.go` and
`cmd/hotm/main.go` to match. Ran the full `gofmt`/`build`/`vet`/`test`
suite (with a repeated `-count=2` run) clean, and verified live: normal
Level 3 movement/monsters/items/map all still work exactly as before.

### Astarot's Charm found and placed: a second demon invocation now reachable in real gameplay

After another Stop-hook rejection, same framing, went back to the
numbered-map cross-referencing technique that found Belezbar's Mantis
2 rounds ago and applied it to Level 1's Wolfdorp (an already-playable
CollodonsPile room). The numbered map poster's key list gives room #65
as "Rock, two stalagmites, stalactite, sword", and re-examining that
poster's own Level 1 grid confirms #65 sits within the "WOLFDORP"
banner-labeled cluster (tight-cropped and visually confirmed, not
guessed). Separately, `level_items.go`'s `LevelOneItems` - sourced
independently from the OTHER poster (`heavymap-levels1-2.jpg`) many
rounds ago - already lists a "Sword" on Level 1 with no room precision.
Two unrelated sources agreeing Level 1 has a sword, one of them at
zone-level confidence for Wolfdorp specifically - the same rigor as the
Mantis placement.

"Sword" is Astarot's confirmed real Charm (`magic.Demons`). Added it to
Wolfdorp's `Items`, making Astarot's already-implemented but previously
unreachable invocation succeed in real gameplay for the first time -
the second demon (after Belezbar) whose invocation now genuinely works,
not just Belezbar's own special case. Added
`TestCollodonsPileWolfdorpHasSword`, ran the full
`gofmt`/`build`/`vet`/`test` suite (with a repeated `-count=2` run)
clean, and verified live: navigated to Wolfdorp, `PICK UP SWORD`, then
`INVOKE ASTAROT` - "Transports the player to a named location, if its
name is known" - succeeded for real.

### Untangled the numbered-map poster's own Level layout, and found the strongest-confidence item placement in the whole port

After another Stop-hook rejection, same framing, tried to find the
remaining 2 demon Charms' locations (Sunflower/Magot, Erlstone/Asmodee)
using the same numbered-map technique as Mantis/Sword. First had to
straighten out a real point of confusion: the numbered map poster's own
"LEVEL 1/2/3/4" labels are NOT arranged as a simple 2x2 grid the way
earlier rounds assumed - re-viewing the full poster at once shows Level
1 top, Level 2 middle-left, Level 3 middle-right (smaller), Level 4
bottom. Tight-cropping the actual "LEVEL 2" label and its grid confirms
it contains "Eye Of Heaven", "Room of Icthys", "Room of Flox", "Room of
Horns" - all matching Level2Grid's already-shipped named cells, a solid
anchor for future rounds needing this poster's Level 2 section again.

That crop also, unprompted, answered a much higher-confidence question:
it labels one exact room "START / Room of Misery / 1, 2" - i.e. Room of
Misery isn't just SOMEWHERE in Level 2, it's numbered cells #1 and #2 on
this poster specifically, no zone-level inference needed at all (the
poster's key list gives #1 as "Grimoire" and #2 as "Poison-smeared
book"). #1 exactly matching the already-independently-sourced Grimoire
(from the CASA walkthrough, a completely different source) is a clean,
strong validation that this numbered-cell reading is correct - and #2,
never placed anywhere before, was added on the same footing. This is
the single most confident item placement in the whole file: an exact
numbered match against an already-known-correct room, not a zone-level
inference like Mantis/Sword.

The Sunflower/Erlstone search itself came up short this round - #7
(Sunflower) sits somewhere in the same immediate Room-of-Misery cluster
but its exact cell is obscured by a marker icon in the source image and
wasn't confidently resolved; left unplaced rather than guessed, a real
open item for a future round now that this poster's Level 2 layout is
finally untangled.

Added `TestCollodonsPileRoomOfMiseryHasBothNumberedItems`, ran the full
`gofmt`/`build`/`vet`/`test` suite (with a repeated `-count=2` run)
clean, and verified live: `LOOK` in Room of Misery now shows both
items, and `PICK UP POISON-SMEARED BOOK` works for real.

### Room of Misery finally exists as real data in Level2Grid, not just prose

After another Stop-hook rejection, same framing, went looking for
Sunflower/Erlstone's exact cells again - a dead end confirmed a second
time (room #7's number is genuinely covered by a marker icon in the
source image, not just hard to read; left unplaced). While looking,
re-examined Level2Grid's already-documented "Room of Misery pocket"
(F3/F4/G3/G4/G5/H4/H5 - real, but disconnected from the main 50-cell
component, per a much earlier round) and noticed something the earlier
round hadn't acted on: this file's own doc comment has said for a long
time that F4 IS Room of Misery, "the confirmed real starting room" -
but F4 (and the whole pocket) was never actually added as a Room in
`level2Cells` at all. It only ever existed in this comment's prose.

Tight-cropped F3 and F4 individually and confirmed pixel-for-pixel:
F4 is drawn "MISERY" (with its own on-map "F4" coordinate label right
there too), F3 is drawn "SIGN!". Added both as real, named, isolated
rooms (no Exits - this round confirmed their NAMES, not new
connectivity; the pocket's other 5 cells and any bridge to the main
component remain real, scoped follow-up work, and the "disconnected"
finding stands as-is). File is now honestly 52 cells (50 main + 2
newly-real named cells), not 50.

This is a small but real milestone: it's the first time the game's
actual confirmed starting room has existed as real Room data in
Level2Grid specifically (CollodonsPile already had it, obviously, but
Level2Grid's own per-cell version of it was pure commentary until now).
Deliberately did NOT change `Level2Grid`'s own start room to F4 - it
has no Exits, so starting there would strand the player with zero
moves; A1 remains the anchor.

Added `TestLevel2GridRoomOfMiseryPocketNamedCells`, ran the full
`gofmt`/`build`/`vet`/`test` suite (with a repeated `-count=2` run)
clean, and verified live that normal Level 2 play is unaffected.

### Completed Level2Grid's Room of Misery pocket, and surfaced a real Flox/D4 discrepancy

Continued straight on from last round's F3/F4 addition, after another
Stop-hook rejection with the same framing. The Sunflower search hit its
now-familiar dead end a third time (the marker icon over room #7's
number hasn't moved). Instead of stopping there, finished what last
round started: tight-cropped the pocket's remaining 5 cells (G3, G4,
G5, H4, H5). G3/G4/H4 are plain unlabeled cells; **G5 and H5 each carry
a real, tight-crop-verified Guards obstacle** - the same confirmed red
icon already used elsewhere on this map, found here just by looking
since a 5-cell pocket is small enough to check by eye rather than
needing the pixel-fraction scan used for the 50-cell component. All 5
added with the same honest "no Exits, names/contents confirmed but not
new connectivity" convention as F3/F4. `Level2Grid`'s Room of Misery
pocket is now fully present as real data for the first time - all 7
cells, not 2 - file is honestly 57 cells (was 52), not 50.

**Found a real discrepancy while doing this**: the file's older
"4 named cells" bullet (from an earlier, less rigorous text-density-
based pass) places "Flox" at D4 - but `level2Cells` already has a D4
in the MAIN 50-cell component, with a real West exit to D3, not
isolated at all. Either Flox's coordinate was misassigned by that
earlier pass, or something else is off. Not resolved this round (that
pass predates the tight-crop-verification discipline used everywhere
else in this file now) - documented as a real, open discrepancy rather
than silently picked one way or guessed at, the same honesty standard
as the Sothic Complex/CollodonsPile naming conflict from 2 rounds ago.

Updated the file's doc comment throughout (including the now-stale
"G5/H5 aren't wired in" Guards paragraph, which this round makes
untrue) and cleaned up formatting. Added
`TestLevel2GridRoomOfMiseryPocketIsComplete`, ran the full
`gofmt`/`build`/`vet`/`test` suite (with a repeated `-count=2` run)
clean, and verified live that normal Level 2 play is unaffected.

### Resolved the Flox/D4 discrepancy, and named 3 more real cells in Level2Grid

After another Stop-hook rejection, same framing, went straight back to
the discrepancy flagged 2 rounds ago rather than let it sit as a
standing open question: tight-cropped D4 directly and it does read
"FLOX", confirming the coordinate was never wrong. The actual resolution
turned out simple - D4 was ALREADY a normal, connected cell in the main
50-cell component (real West exit to D3, itself reciprocally confirmed
by tight-cropping D3 too), just missing its Name. The earlier "isolated"
claim for Flox specifically was just wrong; no real conflict once
checked directly. Added `Name: "Flox"` to the existing D4 entry - zero
connectivity risk, pure labeling.

Checked the other 3 cells from that same old "4 named cells" bullet
(Icthys, Horns, Purity) and found they're a genuinely different case:
truly absent from `level2Cells` entirely (not merely unnamed like Flox
was), consistent with them really being isolated. Individually tight-
cropped all 3 with the current calibration and confirmed pixel-for-
pixel ("ICTHYS" at C3, "HORNS" at E2, "PURITY" at H3, matching the old
pass's coordinates exactly this time) - added as real, named, isolated
cells (no Exits), same honest convention as the Room of Misery pocket.
File is now 60 cells (was 57): 50 main + 7 pocket + 3 newly-named
isolated cells, with Flox correctly recognized as part of the 50, not
a 61st.

Added `TestLevel2GridFloxIsInMainComponent` and
`TestLevel2GridOtherNamedIsolatedCells`, ran the full `gofmt`/`build`/
`vet`/`test` suite (with a repeated `-count=2` run) clean, and verified
live: computed the real shortest path to Flox (via a disposable, never-
committed temporary test file, deleted immediately after use) and
walked it in `go run ./cmd/hotm -level2grid` - `LOOK` now shows "Flox"
by name for the first time.

### Reapplied the "named special room explains a gap" technique to Level3Grid, found 2 more real rooms

After another Stop-hook rejection, same framing, applied the exact
technique that found Sothic Complex (round 51) and Flox/Icthys/Horns/
Purity (round 56) to Level3Grid's own remaining unexplained gap: F2-F6
are entirely absent from its 41-cell main component, same symptom as
D4 was before Sothic Complex explained it. Tight-cropped all 5 and
found 2 more real named special rooms: F3 reads "NANI" (Room of Nani)
and F5 reads "HYDRA" (Rook of Hydra) - both zone names already visible
elsewhere on this map, now confirmed to also have their own individual
special-room cell here, the same pattern as Sothic Complex. F2 has a
special-room-style border but no legible text in this crop; F4 and F6
are plain, unnamed cells - all 3 left unadded rather than guess.

Added Nani (F3) and Hydra (F5) as real, named, isolated cells (no
Exits - same honest convention as every other special room in this
file). Level3Grid is now honestly 44 cells (was 42): the 41-cell main
component plus 3 isolated named special rooms (Sothic Complex, Nani,
Hydra). Updated the stale "42-cell" references in `game.go` and
`cmd/hotm/main.go` to match.

Added `TestLevel3GridNaniAndHydraAreIsolated`, ran the full
`gofmt`/`build`/`vet`/`test` suite (with a repeated `-count=2` run)
clean, and verified live that normal Level 3 play is unaffected.

### Applied the named-special-room technique to Level4Grid: 6 rooms found, including a second real Exit

After another Stop-hook rejection, same framing, took the "How to
apply" suggestion from 2 rounds ago (try this technique on whatever of
Level4's newly-confirmed-real rows get looked at) and did exactly that.
Level1Grid was checked first and found to already have every one of
its 64 cells present in the file (just some isolated) - no gap to
explain there. Level4Grid was different: its rows A-E were confirmed
real 2 rounds ago (during the relabeling fix) but never actually added
as cells. Tight-cropped rows A through G looking for named special
rooms and found 6, all individually pixel-confirmed:

  - **F4 "The Chasm" is the Flox/D4 case again** - already a normal
    cell in the 17-cell main component (real exits East to F5, West to
    F3), just missing its Name. Fixed with a pure, zero-risk name
    addition.
  - **D2 "Scales", D3 "Doubt of Rabak", F1 "The Crypt", G4 "Pride"** are
    genuinely absent from the file - added as real, named, isolated
    cells (no Exits), same honest convention as every other special
    room in this project.
  - **G2 "Exit"** - also genuinely absent, added the same way. This is
    a second real, individually-verified Exit location (Level1Grid's
    G3 was the first) - concretely corroborating the manual's "3 exits"
    detail with a second confirmed instance, though this one can't
    currently be reached (it's isolated, no confirmed Exits of its own)
    so `Game.Won` can't fire for it yet.

Level4Grid is now honestly 22 cells (was 17): the 17-cell main
component plus 5 isolated named special rooms. Updated the stale
"17-cell fully-connected" framing in `game.go` and `cmd/hotm/main.go`.
Added `TestLevel4GridTheChasmIsInMainComponent` and
`TestLevel4GridIsolatedNamedRooms`, ran the full `gofmt`/`build`/`vet`/
`test` suite (with a repeated `-count=2` run) clean, and verified live:
`go run ./cmd/hotm -level4grid` now shows "The Chasm" by name.

### Items now render in the live GUI too, completing the monster/guards/items HUD trio

After another Stop-hook rejection, same framing, checked Level1Grid for
the named-special-room technique first (per last round's plan) and
found nothing new - all 64 of its cells are already present in the
file, some isolated, but none missing a real name the map shows. No
gap to explain there, a real (if quiet) negative result. Pivoted back
to graphics instead: `cmd/hotm-gui` already draws the current room's
Monster and Guards state (last 2 rounds) but never its Items - real
per-room item data has existed since very early in this project, shown
only as plain log text, never in the live HUD area.

Added `drawItems`, rendering the current room's real Items next to the
monster/guards indicators. Honestly documented a real constraint this
time: the clean map's own confirmed "object" icon color is black, but
this GUI's background is also black, so using the "true" confirmed
color (like monsterGlyphColor/guardsColor do) would render invisible
text - `itemsColor` is explicitly flagged as NOT the confirmed icon
color, just a legible plain-yellow stand-in, an honesty distinction
worth being explicit about rather than silently picking a color and
implying it's sourced like the other two.

Added `TestItemsColorIsLegible`, ran the full `gofmt`/`build`/`vet`/
`test` suite clean, and verified live via the disposable-throwaway-
repo-copy + `PrintWindow` technique: a real screenshot of the default
starting room (Room of Misery) shows "Grimoire, Poison-smeared book" in
yellow, clearly legible against the black background.

### Gave GUARDS and INVOKE real audio feedback, and made the GUI's INVOKE key actually invoke something

After another Stop-hook rejection, same framing, went back to sound
again. `cmd/hotm-gui` already gives BLAST/FREEZE/TRANSFUSION distinct
audio feedback by matching the real returned message (18th push), but
the G (GUARDS, DOOR) and I (INVOKE) keys added in later rounds just
called `Handle` directly with no audio at all - a real, concrete gap
in the same category of work, not a new domain.

Investigating it surfaced a deeper issue: the I key only ever sent bare
`"INVOKE"` (which just lists the 4 demons - a real Charm-bearing
invocation needs a specific name as target, and this GUI has no text
input). That meant a "successful invocation" sound could never actually
trigger through this key at all, no matter what the player carried -
worth fixing properly rather than adding a dead code path. Added
`invokeCarriedDemon` (same target-selection convention as
`pickUpFirstItem`/`dropFirstItem`): scans the player's real carried
items against each confirmed `magic.Demons`'s Charm and invokes the
first match, falling back to bare INVOKE if none is carried. This
makes a real invocation (Mantis→Belezbar, Sword→Astarot, both already
placed in the world) actually reachable from the GUI's I key for the
first time.

Added 2 new feedback notes (`noteGuardsPass`, `noteInvokeOK`), wired
both G and I through the existing `handleAndPlay` message-matching
pattern. Also promoted the ad hoc `hasItem` check to a proper exported
`character.Player.HasItem` (used by both `game.hasItem`, now a thin
wrapper, and the new GUI logic) rather than duplicate the same
case-insensitive loop a third time.

Added `TestHasItem` and `TestInvokeCommandForPicksCarriedCharm`
(the target-selection logic split out from `invokeCarriedDemon`
specifically so it's testable without a real audio context), ran the
full `gofmt`/`build`/`vet`/`test` suite (with a repeated `-count=2`
run) clean, and verified the underlying mechanic live via the text
frontend (unaffected, as expected, since `handleAndPlay` and
`invokeCarriedDemon` only wrap the same `game.Handle` calls).

### Re-examined the original hand-drawn Level 1-2 poster at full resolution: real new confirmations, but a placement dead end for a real reason

After another Stop-hook rejection, same framing, went looking for
Magot's Sunflower and Asmodee's Erlstone again - this time by re-
viewing `heavymap-levels1-2.jpg` (the OTHER, hand-drawn poster,
`level_items.go`'s source) at full resolution instead of the clean map
or numbered map already exhausted for this search. It's much more
legible than earlier passes treated it: individual item labels sit in
recognizable grid cells, not just "somewhere on this level."

Two real, useful confirmations came out of it:
  - **"SUN-FLOWER" is directly visible** in Level Two's grid, in the
    same bottom-right cluster as "TOLL", "NICKEL KEY", and "LOAF" -
    real, independent (a second, different poster) confirmation that
    Magot's Charm exists and is on Level 2, corroborating the numbered
    map's own "#7, Chest (sunflower)" entry from a completely different
    source.
  - **"+ ONE MAGICK GRADE"** is a real, previously-unknown location
    marker on Level Two, distinct from any item - i.e. a second real
    "Grade promotion" spot in the game besides Secunda Porta's door
    (the only one currently modeled). Real, sourced, not yet actionable
    without a location.

Neither could be safely placed at an exact `Level2Grid` cell this
round: I tried anchoring this poster's own row/column grid against
Level2Grid's already-confirmed Room of Misery position (Grimoire is
also directly visible on this poster, in a distinctly-colored cell),
and the two countings didn't agree - a real, honest sign that this
poster's cell grid doesn't align cleanly with Level2Grid's lettered
grid the way the numbered map's did for Mantis/Sword. Rather than force
a guess through a mismatch I can't resolve, left both unplaced. This is
the same limitation `level_items.go`'s doc comment already named many
rounds ago ("only the item-to-level association is captured here, not
exact grid position") - this round's attempt to push past that limit
for two specific items didn't succeed, but confirms *why* not, rather
than silently trying and getting it wrong.

A quiet round: real, verified findings, no shippable code change - a
legitimate outcome per this project's own established practice (a
well-verified negative/confirmatory result is real progress, not a
stall), not forced into a placement that would have been a guess.

### EXAMINE now actually uses its target, matching the manual's own confirmed grammar example

After another Stop-hook rejection, same framing, went looking for a
real gap in already-shipped mechanics rather than another map-reading
pass. Found one: `parser`'s own package doc comment has always used
`"X BOTTLE"` (eXamine the bottle) as its example of the confirmed
action-form grammar (`Keyword Object`) - but `game.examine()` completely
ignored the target the whole time, silently falling back to "list
everything in the room" regardless of what was actually asked about.
`EXAMINE GRIMOIRE` and bare `EXAMINE` have produced byte-identical
output since this verb was first added.

Fixed: with a target, `examine` now confirms just that one thing - the
room's Monster, a room Item, or a carried Item - or says plainly
`"You don't see that here."` if none match, rather than always dumping
the full room contents. Bare `EXAMINE` (no target) keeps its existing
list-everything behavior unchanged. Caught and fixed a casing bug of
my own while testing live: an early version echoed the player's
UPPERCASED typed target back in the "You are carrying a GRIMOIRE"
response instead of the item's real stored casing - fixed by looking
up the actual stored name instead of trusting the raw input.

Added 3 regression tests (targeted-item-found, target-not-here,
targeted-carried-item), ran the full `gofmt`/`build`/`vet`/`test` suite
clean, and verified live: `X GRIMOIRE` confirms just the Grimoire (not
the room's other item), `EXAMINE SWORD` (nothing by that name present)
correctly says so, and after picking it up, `EXAMINE GRIMOIRE` reports
"You are carrying a Grimoire" with correct casing.

### Re-fetched the CASA walkthrough with fresh, targeted questions: 2 more items placed, and a real new combat-alternative mechanic

After another Stop-hook rejection, same framing, went back to the
original source that's produced the most real content in this whole
project (the CASA walkthrough) and re-fetched it with a deliberately
more detailed, targeted question set - specifically asking for exact
item-to-room associations and any TARGET/VERB or monster-interaction
facts not yet captured. This source has been re-read several times
before, but never with this specific combination of questions.

Two real findings resolved long-standing gaps:
  - **Nougat and a Scroll are both in Trollwynd; a second, separate
    Scroll is in Sothic Complex.** Both rooms already exist and are
    playable in `CollodonsPile`. An earlier pass had found "Scroll" and
    "Nougat" but reported them with vague/multiple locations and
    deliberately left them unplaced rather than guess - this fresh
    re-read gives exact rooms for both, resolving 2 of those 3 (Key
    remains genuinely spread across 4 rooms with no single location,
    still correctly unplaced).
  - **A real, previously-unknown alternate combat mechanic**: the
    walkthrough states Werewolves are "killable by walking through
    after dropping NOUGAT" - a real way to defeat a Werewolf without
    BLAST/FREEZE at all. Implemented as `game.checkNougatWerewolf`,
    called from both `drop` (dropping Nougat in a Werewolf's room is
    read as the real trigger) and `move` (covers Nougat already present
    for any other reason): whenever the current room has both a live
    Werewolf and a Nougat, the Werewolf is defeated automatically.

Added `TestNougatDefeatsWerewolfOnDrop` and
`TestCollodonsPileTrollwyndAndSothicComplexHaveScrollNougat`, ran the
full `gofmt`/`build`/`vet`/`test` suite clean, and verified live in two
parts (the mechanic can't be demonstrated as one natural playthrough
yet, honestly noted: Nougat only exists in `CollodonsPile`, Werewolves
only in `Level1Grid`/`Level2Grid`, and the 5 graphs aren't merged) -
`go run ./cmd/hotm` confirmed real Nougat/Scroll pickup at Trollwynd
and a second Scroll at Sothic Complex, and
`TestNougatDefeatsWerewolfOnDrop` exercises the actual mechanic against
a real Werewolf in `Level1Grid`.

### Found the Toll mechanic's real trigger phrase, and 3 real placements - a satisfying pickup-then-use puzzle chain now actually works

After another Stop-hook rejection, same framing, went back to the CASA
walkthrough a third time this session with yet another targeted
question set (Guards/Toll/Keys/locked-doors/demons/win-condition) -
mostly clean negatives (this walkthrough never mentions Guards, Toll
signs, Bag of Gold, or any of the 4 demons - genuinely doesn't cover
that content, not a missed extraction), but one real, significant find:
the actual phrase that opens several doors is literally
**`"DROP <item>"`** ("EXAMINE TABLE, DROP KEY (door opens)"), not the
`"DOOR, <item>"` form this port originally guessed at when the Toll
mechanic was first modeled many rounds ago. The instructions file's
own "put it on the table" wording matches DROP far more literally than
DOOR ever did.

Carefully disambiguated a real ambiguity in the walkthrough's own
formatting (room names appear both as paragraph headings AND as
travel-destination annotations after movement lists, which look
similar) before trusting any of it - confirmed via 2 follow-up fetches
that 3 instances are genuinely under their stated room's own heading
(Room of Stings/Key, Morfang/Bag, Room of Arrows/Slat), while a 4th
apparent instance near Pilefoot turned out, on closer checking, to
actually belong to a different, unidentified room - correctly left
unplaced rather than guessed.

Wired `game.drop` to check a room's real `TollItem` and pay the toll
automatically on a matching drop (keeping the original `"DOOR, <item>"`
form working too, since nothing confirms it's wrong, just that DROP is
also, and probably primarily, real). Placed the 3 confirmed
`TollItem`s - and two of them chain beautifully with items already
real and placed in earlier rounds: Wolfdorp's already-real Bag is
exactly what Morfang's door needs, and Morfang's already-real Slat is
exactly what Room of Arrows' door needs - a genuine pickup-here-use-
there puzzle, not a coincidence.

Updated `world.Room.TollItem`'s doc comment (the mechanic is broader
than "toll signs" specifically) and `TestHandleTollDoorRequiresItem`'s
stale claim that no real room had one yet. Added
`TestHandleDropPaysRealToll`, ran the full `gofmt`/`build`/`vet`/`test`
suite (with a repeated `-count=2` run) clean, and verified live: picked
up the real Bag in Wolfdorp, carried it through Room of Stings (where
dropping it correctly does nothing — that door needs a Key), and
`DROP BAG` at Morfang for real opened the door.

### A real environmental fixture, hiding in plain sight across 3 rounds of the same walkthrough: EXAMINE TABLE

After another Stop-hook rejection, same framing, went back over the
last 2 rounds' own CASA walkthrough fetches rather than fetching again
- both had already surfaced "EXAMINE TABLE" repeatedly (it's the
opening move in nearly every room-paragraph the walkthrough describes,
always right before a pickup or a toll-paying drop), but this project
had only ever used those fetches for the item/door facts sitting next
to it, never modeled the table itself. Room of Misery, Trollwynd,
Sothic Complex, Wolfdorp, Room of Stings, Morfang, Room of Arrows, and
Methos are the 8 rooms directly confirmed to have one across those
fetches.

Added `world.Room.HasTable` and wired `"EXAMINE TABLE"` into `examine`'s
already-target-aware logic (from 2 rounds ago) to acknowledge it in
those rooms specifically - not assumed for every room just because the
walkthrough's own narration happens to open with it almost everywhere;
only the 8 directly confirmed are marked. A genuine, if small, real
environmental detail now modeled instead of silently discarded as
"just flavor text around the real facts."

Added `TestHandleExamineTable` and `TestCollodonsPileHasTablePlacements`,
ran the full `gofmt`/`build`/`vet`/`test` suite clean, and verified
live: `X TABLE` in Room of Misery answers "A plain table." instead of
"You don't see that here."

### Extended the named-special-room technique to Level 3's Kitchen of Ai zone: 3 more finds, including a new Wyvern

After another Stop-hook rejection, same framing, applied the now-
standard "check a gap span for named special rooms" technique (Sothic
Complex, then Nani/Hydra) to the one remaining unchecked corner of
Level3Grid: the Kitchen of Ai zone (rows G-H, cols 1-4), entirely
absent from the main component like the earlier gaps were. Tight-
cropped all 8 cells and found 3 more real, individually pixel-confirmed
things: **G2 reads "TWO"** (a real named special room), **G4 has a
real, confirmed Wyvern monster icon** (the same rare exact blue
RGB(0,132,255) already validated elsewhere on this map), and **H4
reads "WATER"**. G1, G3, H1, and H3 were also checked and are plain,
unlabeled cells - correctly left unadded.

Added all 3 as isolated cells (no Exits - connectivity for this zone
isn't extracted), same honest convention as every other special room
in this file. Level3Grid is now honestly 47 cells (was 44): the
41-cell main component plus 6 isolated special rooms/monsters. Updated
the stale "44-cell" framing in `game.go` and `cmd/hotm/main.go`.

Added `TestLevel3GridKitchenOfAiFinds`, ran the full `gofmt`/`build`/
`vet`/`test` suite (with a repeated `-count=2` run) clean, and verified
live that normal Level 3 play is unaffected.

## Open next steps

- **Level 1's connectivity has been extracted AND is playable**
  (`internal/world/level1_grid.go` + `game.NewLevel1Exploration()`,
  reachable via `go run ./cmd/hotm -level1grid` — 64 real per-cell rooms,
  81 validated exits) — see "Attempted (and mostly succeeded at) Level
  1's full connectivity" and "Made Level1Grid actually playable" above.
  **Highest-value next steps now**: (1) Level 1's remaining ~20-cell
  fragment (columns 5-8, rows E-H) has had EVERY possible boundary
  crossing to it individually hand-verified closed (all 12 checked, not
  just the 4 originally-flagged cells) — real evidence it's a genuinely
  separate area (likely reached via a stairwell inside the fragment
  itself, the same pattern as Agile Stair) rather than a missed
  detection, so further spot-checking this specific boundary isn't
  likely to help; a real fix would need locating that internal stairwell
  cell instead, or accepting this as a known, permanent limit of this
  map source; (2) **All 4 levels now have a real, playable grid** — see
  "All 4 dungeon levels now have a real, validated, playable grid" above:
  `internal/world/level{1,2,3,4}_grid.go`,
  `game.NewLevel{1,2,3,4}Exploration()`, `go run ./cmd/hotm
  -level{1,2,3,4}grid`. Remaining real follow-up work on this front: add
  monster/item icon placements to Levels 2/3/4 (only Level 1 has them),
  and nail down Level 3's and Level 4's uncertain row-F/G/H calibration
  to recover their other (smaller, currently unshipped) components -
  Level 4 especially, since its largest shipped component is only 17
  cells and a proper recalibration is the most likely way to grow it;
  (3) reconcile all four grids' A1-H8 addressing with `CollodonsPile`'s
  existing named rooms (Wolfdorp, Morfang, Pilefoot, Nidus, Room of
  Stings, Room of Arrows, Room of Misery, Secunda Porta, Trollwynd,
  Agile Stair, Methos all already exist there under a different,
  walkthrough-sourced scheme) and actually merge them into ONE graph
  `game.New()` uses (right now there are honestly five separate,
  playable starting points, not one unified game) — a real, scoped
  architecture change (per-cell rooms replacing per-zone rooms), not a
  quick patch.
- ~~Implement the confirmed Stamina-reaches-0-means-death mechanic~~ —
  done (see the "Randomized stats + death/combat-Stamina-cost" section
  above): `Player.IsDead()`, the death gate in `game.Handle`, and a
  Stamina cost on `BLAST`/`FREEZE`. The exact combat cost/restore amounts
  are still honest placeholders, not extracted.
- ~~Wire `INVOKE` + the 4 confirmed `magic.Demons` to something real~~ —
  done: `game.invoke` recognizes the 4 demons and reports the confirmed
  "requires a suitable Talisman" gate. Real next step once an inventory
  model exists: let the player actually acquire/hold a Talisman and pass
  the check, and implement Astarot's (transport) and Magot's (locate)
  confirmed abilities for real instead of always failing the gate.
- Extend `world.CollodonsPile()` with more of the dungeon — the CASA
  walkthrough almost certainly documents more rooms/paths than the single
  13-room stretch pulled out so far (it's a full solution, likely covering
  most or all of the game); worth re-extracting more thoroughly. Also
  check whether other walkthroughs (World of Spectrum, MobyGames, Hardcore
  Gaming 101 all showed up in search results) add rooms this one misses.
- Get real room DESCRIPTION text (not just names) — either by finally
  finding the disassembly's string-print routine, or by resuming the
  live-SpecEmu-play approach (now that the click-for-focus fix is known)
  far enough to read actual room text off the live screen.
- Resume the live-play memory-diff approach if a real room table is still
  wanted from the disassembly side: now that input focus is understood
  (click the render surface before sending keys), get further than the
  stats-menu screen into an actual room, and diff snapshots before/after
  a real move to find the "current room" variable directly.
- Trace `34709` and `42757` (untraced per-tick calls in the main loop) and
  `43398`'s continuation past the menu-flag check — the movement/room
  dispatcher is still somewhere in the main loop, just not in the three
  helper calls checked so far (`43398`'s menu, `30180`, `30107` — all
  confirmed to be UI utilities, not gameplay dispatch).
- **Trace `43398`** (the real parser/input routine, called every game
  loop tick) and `30180` (command dispatch) — the direct path to finding
  room/exit data and real verb resolution. Higher priority than the
  bucket-index-table details below.
- Confirm the bucket-index-table indexing precisely (what word length maps
  to table offset 0) — needed to actually use the vocabulary table for
  parsing rather than just having the word list.
- Find the string-print loop that calls `44286` per-character — would let
  us bulk-extract every message/room-description string in the game at once.
- Find the per-word type/ID table for the ~300 non-direction vocabulary
  words (verb/noun/target roles) — the real unlock for parsing full commands.
- Trace `CALL 42608` (turns an index into an object record — would reveal
  semantically what table entries represent: items? monsters? spell
  components?) — natural next step for the picture-format mystery.
- Resolve the picture-transform bit-semantics — likely needs a verified
  `(index, count, row_count)` combo from an actual traced call, rather than
  guessing against entry 1 blind.
- Trace the actual gameplay loop: room-connection/movement logic (now that
  direction words are confirmed, find what reads them), combat/spell
  resolution.
- ~~Confirm what `44074`/`44103` actually do~~ — done, window-setup pair, not
  printing. ~~Trace `44267`~~ — done, an AT-position convenience wrapper
  around `44286` (see above). Still open: the actual string-print loop.
- Trace routines `31435` and `40218` (callers of the `42612` descriptor-copy
  helper) and `E`/`(IX+5)` producers feeding into `31932`'s glyph selector.
