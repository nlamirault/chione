// Copyright (C) 2017 Nicolas Lamirault <nicolas.lamirault@gmail.com>

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package skiinfo

import (
	"strings"
	"unicode"

	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

var (
	resortsParticularCase = map[string]string{
		"ABRIÈS EN QUEYRAS":    "abries",
		"AIGUILLES EN QUEYRAS": "aiguilles",
		"AILLONS MARGERIAZ":    "les-aillons-margeriaz",
		"ARVIEUX EN QUEYRAS":   "arvieux",
		"BEILLE":               "plateau-de-beille",
		"Grand Tourmalet (La Mongie / Barèges)": "la-mongie-bareges",
	}
)

func isMn(r rune) bool {
	return unicode.Is(unicode.Mn, r) // Mn: nonspacing marks
}

func MakeResortName(name string) string {
	if val, ok := resortsParticularCase[name]; ok {
		return val
	}
	str := strings.ToLower(name)
	str = strings.Replace(str, " ", "-", -1)
	t := transform.Chain(norm.NFD, transform.RemoveFunc(isMn), norm.NFC)
	str, _, _ = transform.String(t, str)
	return strings.TrimSpace(str)

}
