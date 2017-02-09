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
	// "fmt"
	"testing"
)

func Test_MakeResortName(t *testing.T) {
	if val := MakeResortName("La Pierre St Martin"); val != "la-pierre-st-martin" {
		t.Fatalf("Invalid resort name: %s", val)
	}
}

func Test_MakeResortNameCustom(t *testing.T) {
	if val := MakeResortName("AILLONS MARGERIAZ"); val != "les-aillons-margeriaz" {
		t.Fatalf("Invalid custom resort name: %s", val)
	}
}
