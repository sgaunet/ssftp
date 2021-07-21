package pathh

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// setUp function, add a number to numbers slice
func setUp() {
	fmt.Printf("setUp tests\n")
}

// tearDown function, delete a number to numbers slice
func tearDown() {
	fmt.Printf("TEARDOWN\n")
}

func TestMain(m *testing.M) {

	setUp()
	code := m.Run()
	tearDown()

	os.Exit(code)
}

func TestIsRemote(t *testing.T) {

	p := New("localFile")
	//assert.Nil(t, err)
	assert.Equal(t, p.IsRemote(), false)

	p = New("user@server")
	assert.Equal(t, p.IsRemote(), true)

	p = New("user@127.0.0.1")
	assert.Equal(t, p.IsRemote(), true)

	p = New("use654r@server2:dir/df")
	assert.Equal(t, p.IsRemote(), true)

	p = New("user@server:.")
	assert.Equal(t, p.IsRemote(), true)
}

func TestGetUser(t *testing.T) {

	p := New("localFile")
	//assert.Nil(t, err)
	assert.Equal(t, p.GetUser(), "")

	p = New("user@server")
	assert.Equal(t, p.GetUser(), "user")

	p = New("use654r@server:dir/df")
	assert.Equal(t, p.GetUser(), "use654r")

}

func TestGetServer(t *testing.T) {
	p := New("localFile")
	//assert.Nil(t, err)
	assert.Equal(t, p.GetServer(), "")

	p = New("user@server")
	assert.Equal(t, p.GetServer(), "server")

	p = New("use654r@server2:dir/df")
	assert.Equal(t, p.GetServer(), "server2")

	p = New("use654r@192.168.0.5:dir/df")
	assert.Equal(t, p.GetServer(), "192.168.0.5")

	p = New("vagrant@192.168.0.5:.")
	assert.Equal(t, p.GetServer(), "192.168.0.5")
}

func TestGetFilePath(t *testing.T) {
	p := New("localFile")
	//assert.Nil(t, err)
	assert.Equal(t, p.GetFilePath(), "localFile")

	p = New("user@server")
	assert.Equal(t, "", p.GetFilePath())

	p = New("use654r@server2:dir/df")
	assert.Equal(t, "dir/df", p.GetFilePath())
}
