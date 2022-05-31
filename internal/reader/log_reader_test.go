package reader

import (
	"bufio"
	"io"
	"os"
	"sync"
	"testing"

	"github.com/sprakhar77/filereader/internal/model"
	"github.com/stretchr/testify/require"
)

func Test_MostActiveCookies(t *testing.T) {
	mp := model.NewCookieMap()
	mp.Add(model.Cookie{Name: "A", Date: "2018-12-09"})
	mp.Add(model.Cookie{Name: "A", Date: "2018-12-09"})
	mp.Add(model.Cookie{Name: "B", Date: "2018-12-09"})
	mp.Add(model.Cookie{Name: "B", Date: "2018-12-09"})
	mp.Add(model.Cookie{Name: "C", Date: "2018-12-09"})
	mp.Add(model.Cookie{Name: "D", Date: "2018-12-10"})
	mp.Add(model.Cookie{Name: "E", Date: "2018-12-10"})

	lr := &logReader{cookieMap: mp}

	cookies := lr.MostActiveCookies("2018-12-09")
	require.EqualValues(t, 2, len(cookies))
	require.ElementsMatch(t, []string{"A", "B"}, cookies)

	cookies = lr.MostActiveCookies("2018-12-10")
	require.EqualValues(t, 2, len(cookies))
	require.ElementsMatch(t, []string{"D", "E"}, cookies)
}

func Test_Read_Integration(t *testing.T) {
	lr := NewLogReader("./dummy.txt")

	err := lr.Read(GB, 1)
	require.NoError(t, err)
	require.ElementsMatch(t, []string{"AtY0laUfhglK3lC7"}, lr.MostActiveCookies("2018-12-09"))
}


func Test_readChunk(t *testing.T) {
	outputChan := make(chan model.Cookie)

	// dummy2 contains all items in a fuzzy manner i.e with extra newlines and
	// space characters
	lr := logReader{
		filePath:   "./dummy2.txt",
		outputChan: outputChan,
		cookieMap:  model.NewCookieMap(),
	}

	var linCount int
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for _ = range outputChan {
			linCount++
		}
		wg.Done()
	}()
	err := lr.readChunk(0, MB)
	close(outputChan)
	wg.Wait()
	require.NoError(t, err)
	require.Equal(t, 8, linCount)
}

func Test_readLine(t *testing.T) {
	file, err := os.Open("./dummy.txt")
	defer file.Close()

	require.NoError(t, err)
	reader := bufio.NewReader(file)
	lineCount := 0

	for {
		_, err := readLine(reader)
		if err == io.EOF {
			break
		}
		lineCount++
	}

	require.EqualValues(t, 8, lineCount)
}

func Test_fileSizeInBytes_Success(t *testing.T) {
	bytes, err := fileSizeInBytes("./dummy.txt")
	require.NoError(t, err)
	require.Equal(t, int64(344), bytes)
}

func Test_fileSizeInBytes_Fail(t *testing.T) {
	_, err := fileSizeInBytes("./nonexistantfile.txt")
	require.Error(t, err)
}
