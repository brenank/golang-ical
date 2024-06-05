//go:build go1.16
// +build go1.16

package ics

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
)

func TestCalendar_Recurrence(t *testing.T) {
	testDir := "testdata/recurrence"
	expectedDir := filepath.Join(testDir, "expected")
	actualDir := filepath.Join(testDir, "actual")

	testFileNames := []string{
		"input1.ics",
	}

	for _, filename := range testFileNames {
		t.Run(fmt.Sprintf("test rrule for: %s", filename), func(t *testing.T) {
			//given
			originalSeriailizedCal, err := os.ReadFile(filepath.Join(testDir, filename))
			require.NoError(t, err)

			//when
			deserializedCal, err := ParseCalendar(bytes.NewReader(originalSeriailizedCal))
			require.NoError(t, err)
			until := time.Date(1999, 1, 1, 0, 0, 0, 0, time.UTC)
			events, err := deserializedCal.EventsWithRecurrence(until)
			require.NoError(t, err)
			newCal := Calendar{
				CalendarProperties: deserializedCal.CalendarProperties,
			}
			for _, e := range events {
				newCal.AddVEvent(e)
			}
			serializedCal := newCal.Serialize()

			//then
			expectedCal, err := os.ReadFile(filepath.Join(expectedDir, filename))
			require.NoError(t, err)
			if diff := cmp.Diff(string(expectedCal), serializedCal); diff != "" {
				err = os.MkdirAll(actualDir, 0755)
				if err != nil {
					t.Logf("failed to create actual dir: %v", err)
				}
				err = os.WriteFile(filepath.Join(actualDir, filename), []byte(serializedCal), 0644)
				if err != nil {
					t.Logf("failed to write actual file: %v", err)
				}
				t.Error(diff)
			}
		})

		t.Run(fmt.Sprintf("compare deserialized -> serialized -> deserialized: %s", filename), func(t *testing.T) {
			//given
			loadIcsContent, err := os.ReadFile(filepath.Join(testDir, filename))
			require.NoError(t, err)
			originalDeserializedCal, err := ParseCalendar(bytes.NewReader(loadIcsContent))
			require.NoError(t, err)

			//when
			serializedCal := originalDeserializedCal.Serialize()
			deserializedCal, err := ParseCalendar(strings.NewReader(serializedCal))
			require.NoError(t, err)

			//then
			if diff := cmp.Diff(originalDeserializedCal, deserializedCal); diff != "" {
				t.Error(diff)
			}
		})
	}
}
