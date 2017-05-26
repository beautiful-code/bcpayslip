// +build ignore

package main

import (
	"bytes"
	"encoding/base64"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/knq/chromedp/kb"
)

var (
	flagOut = flag.String("out", "keys.go", "out source")
	flagPkg = flag.String("pkg", "kb", "out package name")
)

const (
	// chromiumSrc is the base chromium source repo location
	chromiumSrc = "https://chromium.googlesource.com/chromium/src"

	// domUsLayoutDataH contains the {printable,non-printable} DomCode -> DomKey
	// also contains DomKey -> VKEY (not used)
	domUsLayoutDataH = chromiumSrc + "/+/master/ui/events/keycodes/dom_us_layout_data.h?format=TEXT"

	// keycodeConverterDataInc contains DomKey -> Key Name
	keycodeConverterDataInc = chromiumSrc + "/+/master/ui/events/keycodes/dom/keycode_converter_data.inc?format=TEXT"

	// domKeyDataInc contains DomKey -> Key Name + unicode (non-printable)
	domKeyDataInc = chromiumSrc + "/+/master/ui/events/keycodes/dom/dom_key_data.inc?format=TEXT"

	// keyboardCodesPosixH contains the scan code definitions for posix (ie native) keys.
	keyboardCodesPosixH = chromiumSrc + "/+/master/ui/events/keycodes/keyboard_codes_posix.h?format=TEXT"

	// keyboardCodesWinH contains the scan code definitions for windows keys.
	keyboardCodesWinH = chromiumSrc + "/+/master/ui/events/keycodes/keyboard_codes_win.h?format=TEXT"

	// windowsKeyboardCodesH contains the actual #defs for windows.
	windowsKeyboardCodesH = chromiumSrc + "/third_party/+/master/WebKit/Source/platform/WindowsKeyboardCodes.h?format=TEXT"
)

const (
	hdr = `package %s

// DOM keys.
const (
	%s)
	
// Keys is the map of unicode characters to their DOM key data.
var Keys = map[rune]*Key{
	%s}
`
)

func main() {
	var err error

	flag.Parse()

	// special characters
	keys := map[rune]kb.Key{
		'\b': {"Backspace", "Backspace", "", "", int64('\b'), int64('\b'), false, false},
		'\t': {"Tab", "Tab", "", "", int64('\t'), int64('\t'), false, false},
		'\r': {"Enter", "Enter", "\r", "\r", int64('\r'), int64('\r'), false, true},
	}

	// load keys
	err = loadKeys(keys)
	if err != nil {
		log.Fatal(err)
	}

	// process keys
	constBuf, mapBuf, err := processKeys(keys)
	if err != nil {
		log.Fatal(err)
	}

	// output
	err = ioutil.WriteFile(
		*flagOut,
		[]byte(fmt.Sprintf(hdr, *flagPkg, string(constBuf), string(mapBuf))),
		0644,
	)
	if err != nil {
		log.Fatal(err)
	}

	// format
	err = exec.Command("goimports", "-w", *flagOut).Run()
	if err != nil {
		log.Fatal(err)
	}

	// format
	err = exec.Command("gofmt", "-s", "-w", *flagOut).Run()
	if err != nil {
		log.Fatal(err)
	}
}

// loadKeys loads the dom key definitions from the chromium source tree.
func loadKeys(keys map[rune]kb.Key) error {
	var err error

	// load key converter data
	keycodeConverterMap, err := loadKeycodeConverterData()
	if err != nil {
		return err
	}

	// load dom code map
	domKeyMap, err := loadDomKeyData()
	if err != nil {
		return err
	}

	// load US layout data
	layoutBuf, err := grab(domUsLayoutDataH)
	if err != nil {
		return err
	}

	// load scan code map
	scanCodeMap, err := loadScanCodes(keycodeConverterMap, domKeyMap, layoutBuf)
	if err != nil {
		return err
	}

	// process printable
	err = loadPrintable(keys, keycodeConverterMap, domKeyMap, layoutBuf, scanCodeMap)
	if err != nil {
		return err
	}

	// process non-printable
	err = loadNonPrintable(keys, keycodeConverterMap, domKeyMap, layoutBuf, scanCodeMap)
	if err != nil {
		return err
	}

	return nil
}

var fixRE = regexp.MustCompile(`,\n\s{10,}`)
var usbKeyRE = regexp.MustCompile(`(?m)^\s*USB_KEYMAP\((.*?), (.*?), (.*?), (.*?), (.*?), (.*?), (.*?)\)`)

// loadKeycodeConverterData loads the key codes from the keycode_converter_data.inc.
func loadKeycodeConverterData() (map[string][]string, error) {
	buf, err := grab(keycodeConverterDataInc)
	if err != nil {
		return nil, err
	}
	buf = fixRE.ReplaceAllLiteral(buf, []byte(", "))

	domMap := make(map[string][]string)
	matches := usbKeyRE.FindAllStringSubmatch(string(buf), -1)
	for _, m := range matches {
		vkey := m[7]
		if _, ok := domMap[vkey]; ok {
			panic(fmt.Sprintf("vkey %s already defined", vkey))
		}
		domMap[vkey] = m[1:]
	}

	return domMap, nil
}

// decodeRune is a wrapper around parsing a printable c++ int/char definition to a unicode
// rune value.
func decodeRune(s string) rune {
	if strings.HasPrefix(s, "0x") {
		i, err := strconv.ParseInt(s, 0, 16)
		if err != nil {
			panic(err)
		}
		return rune(i)
	}

	if !strings.HasPrefix(s, "'") || !strings.HasSuffix(s, "'") {
		panic(fmt.Sprintf("expected character, got: %s", s))
	}

	if len(s) == 4 {
		if s[1] != '\\' {
			panic(fmt.Sprintf("expected escaped character, got: %s", s))
		}
		return rune(s[2])
	}

	if len(s) != 3 {
		panic(fmt.Sprintf("expected character, got: %s", s))
	}

	return rune(s[1])
}

// getCode is a simple wrapper around parsing the code definition.
func getCode(s string) string {
	if !strings.HasPrefix(s, `"`) || !strings.HasSuffix(s, `"`) {
		panic(fmt.Sprintf("expected string, got: %s", s))
	}

	return s[1 : len(s)-1]
}

// addKey is a simple map add wrapper to panic if the key is already defined,
// and to lookup the correct scan code.
func addKey(keys map[rune]kb.Key, r rune, key kb.Key, scanCodeMap map[string][]int64, shouldPanic bool) {
	if _, ok := keys[r]; ok {
		if shouldPanic {
			panic(fmt.Sprintf("rune %U (%s/%s) already defined in keys", r, key.Code, key.Key))
		}
		return
	}

	sc, ok := scanCodeMap[key.Code]
	if ok {
		key.Native = sc[0]
		key.Windows = sc[1]
	}

	keys[r] = key
}

var printableKeyRE = regexp.MustCompile(`\{DomCode::(.+?), \{(.+?), (.+?)\}\}`)

// loadPrintable loads the printable key definitions.
func loadPrintable(keys map[rune]kb.Key, keycodeConverterMap, domKeyMap map[string][]string, layoutBuf []byte, scanCodeMap map[string][]int64) error {
	buf := extract(layoutBuf, "kPrintableCodeMap")

	matches := printableKeyRE.FindAllStringSubmatch(string(buf), -1)
	for _, m := range matches {
		domCode := m[1]

		// ignore domCodes that are duplicates of other unicode characters
		if domCode == "INTL_BACKSLASH" || domCode == "INTL_HASH" || strings.HasPrefix(domCode, "NUMPAD") {
			continue
		}

		kc, ok := keycodeConverterMap[domCode]
		if !ok {
			panic(fmt.Sprintf("could not find key %s in keycode map", domCode))
		}

		code := getCode(kc[5])
		r1, r2 := decodeRune(m[2]), decodeRune(m[3])
		addKey(keys, r1, kb.Key{
			Code:       code,
			Key:        string(r1),
			Text:       string(r1),
			Unmodified: string(r1),
			Print:      true,
		}, scanCodeMap, true)

		// shifted value is same as non-shifted, so skip
		if r2 == r1 {
			continue
		}
		// skip for duplicate keys
		if r2 == '|' && domCode != "BACKSLASH" {
			continue
		}

		addKey(keys, r2, kb.Key{
			Code:       code,
			Key:        string(r2),
			Text:       string(r2),
			Unmodified: string(r1),
			Shift:      true,
			Print:      true,
		}, scanCodeMap, true)
	}

	return nil
}

var domKeyRE = regexp.MustCompile(`(?m)^\s+DOM_KEY_(?:UNI|MAP)\("(.+?)",\s*(.+?),\s*(0x[0-9A-F]{4})\)`)

// loadDomKeyData loads the dom key data definitions.
func loadDomKeyData() (map[string][]string, error) {
	buf, err := grab(domKeyDataInc)
	if err != nil {
		return nil, err
	}
	buf = fixRE.ReplaceAllLiteral(buf, []byte(", "))

	keyMap := make(map[string][]string)
	matches := domKeyRE.FindAllStringSubmatch(string(buf), -1)
	for _, m := range matches {
		keyMap[m[2]] = m[1:]
	}

	return keyMap, nil
}

var nonPrintableKeyRE = regexp.MustCompile(`\n\s{4}\{DomCode::(.+?), DomKey::(.+?)\}`)

// loadNonPrintable loads the not printable key definitions.
func loadNonPrintable(keys map[rune]kb.Key, keycodeConverterMap, domKeyMap map[string][]string, layoutBuf []byte, scanCodeMap map[string][]int64) error {
	buf := extract(layoutBuf, "kNonPrintableCodeMap")
	matches := nonPrintableKeyRE.FindAllStringSubmatch(string(buf), -1)
	for _, m := range matches {
		code, key := m[1], m[2]

		// get code, key definitions
		dc, ok := keycodeConverterMap[code]
		if !ok {
			panic(fmt.Sprintf("no dom code definition for %s", code))
		}
		dk, ok := domKeyMap[key]
		if !ok {
			panic(fmt.Sprintf("no dom key definition for %s", key))
		}

		// some scan codes do not have names defined, so use key name
		c := dk[0]
		if dc[5] != "NULL" {
			c = getCode(dc[5])
		}

		// convert rune
		r, err := strconv.ParseInt(dk[2], 0, 32)
		if err != nil {
			return err
		}

		addKey(keys, rune(r), kb.Key{
			Code: c,
			Key:  dk[0],
		}, scanCodeMap, false)
	}

	return nil
}

var nameRE = regexp.MustCompile(`[A-Z][a-z]+:`)

// processKeys processes the generated keys.
func processKeys(keys map[rune]kb.Key) ([]byte, []byte, error) {
	// order rune keys
	idx := make([]rune, len(keys))
	var i int
	for c := range keys {
		idx[i] = c
		i++
	}
	sort.Slice(idx, func(a, b int) bool {
		return idx[a] < idx[b]
	})

	// process
	var constBuf, mapBuf bytes.Buffer
	for _, c := range idx {
		key := keys[c]

		g, isGoCode := goCodes[c]
		s := fmt.Sprintf("\\u%04x", c)
		if isGoCode {
			s = g
		} else if key.Print {
			s = fmt.Sprintf("%c", c)
		}

		// add key definition
		v := strings.TrimPrefix(fmt.Sprintf("%#v", key), "kb.")
		v = nameRE.ReplaceAllString(v, "")
		mapBuf.WriteString(fmt.Sprintf("'%s': &%s,\n", s, v))

		// fix 'Quote' const
		if s == `\'` {
			s = `'`
		}

		// add const definition
		if (isGoCode && c != '\n') || !key.Print {
			n := strings.TrimPrefix(key.Key, ".")
			if n == `'` || n == `\` {
				n = key.Code
			}

			constBuf.WriteString(fmt.Sprintf("%s = \"%s\"\n", n, s))
		}
	}

	return constBuf.Bytes(), mapBuf.Bytes(), nil
}

var domCodeVkeyFixRE = regexp.MustCompile(`,\n\s{5,}`)
var domCodeVkeyRE = regexp.MustCompile(`(?m)^\s*\{DomCode::(.+?), (.+?)\}`)

// loadScanCodes loads the scan codes for the dom key definitions.
func loadScanCodes(keycodeConverterMap, domKeyMap map[string][]string, layoutBuf []byte) (map[string][]int64, error) {
	vkeyCodeMap, err := loadPosixWinKeyboardCodes()
	if err != nil {
		return nil, err
	}

	buf := extract(layoutBuf, "kDomCodeToKeyboardCodeMap")
	buf = domCodeVkeyFixRE.ReplaceAllLiteral(buf, []byte(", "))

	scanCodeMap := make(map[string][]int64)
	matches := domCodeVkeyRE.FindAllStringSubmatch(string(buf), -1)
	for _, m := range matches {
		domCode, vkey := m[1], m[2]

		kc, ok := keycodeConverterMap[domCode]
		if !ok {
			panic(fmt.Sprintf("dom code %s not defined in keycode map", domCode))
		}

		sc, ok := vkeyCodeMap[vkey]
		if !ok {
			panic(fmt.Sprintf("vkey %s is not defined in keyboardCodeMap", vkey))
		}

		scanCodeMap[getCode(kc[5])] = sc
	}

	return scanCodeMap, nil
}

var defineRE = regexp.MustCompile(`(?m)^#define\s+(.+?)\s+([0-9A-Fx]+)`)

// loadPosixWinKeyboardCodes loads the native and windows keyboard scan codes
// mapped to the DOM key.
func loadPosixWinKeyboardCodes() (map[string][]int64, error) {
	var err error

	lookup := map[string]string{
		// mac alias
		"VKEY_LWIN": "0x5B",

		// no idea where these are defined in chromium code base (assuming in
		// windows headers)
		//
		// manually added here as pulled from various online docs
		"VK_CANCEL":       "0x03",
		"VK_OEM_ATTN":     "0xF0",
		"VK_OEM_FINISH":   "0xF1",
		"VK_OEM_COPY":     "0xF2",
		"VK_DBE_SBCSCHAR": "0xF3",
		"VK_DBE_DBCSCHAR": "0xF4",
		"VK_OEM_BACKTAB":  "0xF5",
		"VK_OEM_AX":       "0xE1",
	}

	// load windows key lookups
	buf, err := grab(windowsKeyboardCodesH)
	if err != nil {
		return nil, err
	}

	matches := defineRE.FindAllStringSubmatch(string(buf), -1)
	for _, m := range matches {
		lookup[m[1]] = m[2]
	}

	// load posix and win keyboard codes
	keyboardCodeMap := make(map[string][]int64)
	err = loadKeyboardCodes(keyboardCodeMap, lookup, keyboardCodesPosixH, 0)
	if err != nil {
		return nil, err
	}
	err = loadKeyboardCodes(keyboardCodeMap, lookup, keyboardCodesWinH, 1)
	if err != nil {
		return nil, err
	}

	return keyboardCodeMap, nil
}

var keyboardCodeRE = regexp.MustCompile(`(?m)^\s+(VKEY_.+?)\s+=\s+(.+?),`)

// loadKeyboardCodes loads the enum definition from the specified path, saving
// the resolved symbol value to the specified position for the resulting dom
// key name in the vkeyCodeMap.
func loadKeyboardCodes(vkeyCodeMap map[string][]int64, lookup map[string]string, path string, pos int) error {
	buf, err := grab(path)
	if err != nil {
		return err
	}
	buf = extract(buf, "KeyboardCode")

	matches := keyboardCodeRE.FindAllStringSubmatch(string(buf), -1)
	for _, m := range matches {
		v := m[2]
		switch {
		case strings.HasPrefix(m[2], "'"):
			v = fmt.Sprintf("0x%04x", m[2][1])

		case !strings.HasPrefix(m[2], "0x") && m[2] != "0":
			z, ok := lookup[v]
			if !ok {
				panic(fmt.Sprintf("could not find %s in lookup", v))
			}
			v = z
		}

		// load the value
		i, err := strconv.ParseInt(v, 0, 32)
		if err != nil {
			panic(fmt.Sprintf("could not parse %s // %s // %s", m[1], m[2], v))
		}

		vkey, ok := vkeyCodeMap[m[1]]
		if !ok {
			vkey = make([]int64, 2)
		}
		vkey[pos] = i
		vkeyCodeMap[m[1]] = vkey
	}

	return nil
}

var endRE = regexp.MustCompile(`\n}`)

// extract extracts a block of next from a block of c++ code.
func extract(buf []byte, name string) []byte {
	extractRE := regexp.MustCompile(`\s+` + name + `.+?{`)
	buf = buf[extractRE.FindIndex(buf)[0]:]
	return buf[:endRE.FindIndex(buf)[1]]
}

// grab retrieves a file from the chromium source code.
func grab(path string) ([]byte, error) {
	res, err := http.Get(path)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	buf, err := base64.StdEncoding.DecodeString(string(body))
	if err != nil {
		return nil, err
	}

	return buf, nil
}

var goCodes = map[rune]string{
	'\a': `\a`,
	'\b': `\b`,
	'\f': `\f`,
	'\n': `\n`,
	'\r': `\r`,
	'\t': `\t`,
	'\v': `\v`,
	'\\': `\\`,
	'\'': `\'`,
}
