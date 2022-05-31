package model

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_getID(t *testing.T) {
	cm := NewCookieMap()
	cm.Add(Cookie{Name: "A", Date: "2018-12-09"})

	val, ok := cm.IDToName[1]
	require.True(t, ok)
	require.EqualValues(t, "A", val)

	id, ok := cm.NameToID["A"]
	require.True(t, ok)
	require.EqualValues(t, 1, id)

	// Add another cookie
	cm.Add(Cookie{Name: "B", Date: "2018-12-09"})

	val, ok = cm.IDToName[2]
	require.True(t, ok)
	require.EqualValues(t, "B", val)

	id, ok = cm.NameToID["B"]
	require.True(t, ok)
	require.EqualValues(t, 2, id)
}

func Test_CookieMap_Integration(t *testing.T) {
	cm := NewCookieMap()

	cookie1 := Cookie{Name: "A", Date: "2018-12-09"}
	cookie2 := Cookie{Name: "B", Date: "2018-12-09"}
	cookie3 := Cookie{Name: "C", Date: "2018-12-09"}

	// Same cookie with different date
	cookie4 := Cookie{Name: "C", Date: "2018-11-09"}

	addToCookieToMap(cookie1, cm, 10)
	addToCookieToMap(cookie2, cm, 2)
	addToCookieToMap(cookie3, cm, 7)

	// Add cookie4 for a different data
	addToCookieToMap(cookie4, cm, 3)

	freqMap := cm.Get("2018-12-09")

	freq, ok := freqMap["A"]
	require.True(t, ok)
	require.EqualValues(t, 10, freq)

	freq, ok = freqMap["B"]
	require.True(t, ok)
	require.EqualValues(t, 2, freq)

	freq, ok = freqMap["C"]
	require.True(t, ok)
	require.EqualValues(t, 7, freq)

	// Check for other date
	freqMap = cm.Get("2018-11-09")
	freq, ok = freqMap["C"]
	require.True(t, ok)
	require.EqualValues(t, 3, freq)
}

func addToCookieToMap(c Cookie, cookieMap *CookieMap, times int) {
	for i := 0; i < times; i++ {
		cookieMap.Add(c)
	}
}
