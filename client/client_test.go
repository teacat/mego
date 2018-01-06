package client

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var client *Client

type dataStruct struct {
	Username, Password string
}

func init() {
	/*go func() {
		e := mego.New()
		e.Register("Test", func(c *mego.Context) {
			var d dataStruct
			c.MustBind(&d)
			c.Respond(d)
		})
		e.Register("TestParams", func(c *mego.Context) {
			c.Respond(dataStruct{
				Username: c.Param(0).GetString(),
				Password: c.Param(1).GetString(),
			})
		})
		e.Register("TestFile", func(c *mego.Context) {
			p := c.MustGetFile("File1").Path
			f, err := os.Open(p)
			if err != nil {
				panic(err)
			}
			b, err := ioutil.ReadAll(f)
			if err != nil {
				panic(err)
			}
			c.Respond(b)
		})
		e.Register("TestFiles", func(c *mego.Context) {
			f := c.MustGetFiles("File1")
			f2 := c.MustGetFiles("File2")
			c.Respond(len(f) + len(f2))
		})
		e.Register("TestFileSlice", func(c *mego.Context) {
			f := c.MustGetFiles("File1")
			c.Respond(len(f))
		})
		e.Register("TestFileMix", func(c *mego.Context) {
			var d dataStruct
			c.MustBind(&d)

			f := c.MustGetFiles("File1")
			f2 := c.MustGetFiles("File2")
			c.Respond(H{
				"Files":    len(f) + len(f2),
				"Username": d.Username,
				"Password": d.Password,
			})
		})
		e.Register("TestFileChunks", func(c *mego.Context) {
			p := c.MustGetFile("File1").Path
			f, err := os.Open(p)
			if err != nil {
				panic(err)
			}
			b, err := ioutil.ReadAll(f)
			if err != nil {
				panic(err)
			}
			c.Respond(b)
		})
		e.Register("TestFileChunksMix", func(c *mego.Context) {
			var d dataStruct
			c.MustBind(&d)

			f := c.MustGetFiles("File1")
			c.Respond(H{
				"Files":    len(f),
				"Username": d.Username,
				"Password": d.Password,
			})
		})
		e.Run()
	}()*/
	// 等待一秒讓引擎初始化並運作。
	// <-time.After(time.Second * 1)
}

func TestClientMain(t *testing.T) {
	//assert := assert.New(t)
	client = New("ws://localhost:5000")
}

func TestClientConnect(t *testing.T) {
	assert := assert.New(t)
	err := client.Connect()
	assert.NoError(err)
}

func TestClientClose(t *testing.T) {
	//assert := assert.New(t)
	//err := client.Close()
	//assert.NoError(err)
}

func TestClientReconnect(t *testing.T) {
	assert := assert.New(t)
	err := client.Reconnect()
	assert.NoError(err)
}

func TestClientSendStruct(t *testing.T) {
	var resp dataStruct
	assert := assert.New(t)
	err := client.
		Call("Test").
		Send(dataStruct{
			Username: "YamiOdymel",
			Password: "yamiodymel12345",
		}).
		EndStruct(&resp)
	assert.NoError(err)
	assert.Equal(dataStruct{
		Username: "YamiOdymel",
		Password: "yamiodymel12345",
	}, resp)
}

func TestClientSendJSON(t *testing.T) {
	var resp dataStruct
	assert := assert.New(t)
	err := client.
		Call("Test").
		Send(`{"Username": "YamiOdymel", "Password": "yamiodymel12345"}`).
		EndStruct(&resp)
	assert.NoError(err)
	assert.Equal(dataStruct{
		Username: "YamiOdymel",
		Password: "yamiodymel12345",
	}, resp)
}

func TestClientSendMap(t *testing.T) {
	var resp dataStruct
	assert := assert.New(t)
	err := client.
		Call("Test").
		Send(H{
			"Username": "YamiOdymel",
			"Password": "yamiodymel12345",
		}).
		EndStruct(&resp)
	assert.NoError(err)
	assert.Equal(dataStruct{
		Username: "YamiOdymel",
		Password: "yamiodymel12345",
	}, resp)
}

func TestClientSendParamsSlice(t *testing.T) {
	var resp dataStruct
	assert := assert.New(t)
	err := client.
		Call("TestParams").
		Send([]string{"YamiOdymel", "yamiodymel12345"}).
		EndStruct(&resp)
	assert.NoError(err)
	assert.Equal(dataStruct{
		Username: "YamiOdymel",
		Password: "yamiodymel12345",
	}, resp)
}

func TestClientSendParamsJSON(t *testing.T) {
	var resp dataStruct
	assert := assert.New(t)
	err := client.
		Call("TestParams").
		Send(`["YamiOdymel", "yamiodymel12345"]`).
		EndStruct(&resp)
	assert.NoError(err)
	assert.Equal(dataStruct{
		Username: "YamiOdymel",
		Password: "yamiodymel12345",
	}, resp)
}

func TestClientSendFileReader(t *testing.T) {
	var resp []byte
	assert := assert.New(t)
	file, err := os.Open("./README.md")
	assert.NoError(err)
	content, err := ioutil.ReadAll(file)
	assert.NoError(err)
	err = client.
		Call("TestFile").
		SendFile(file).
		EndStruct(&resp)
	assert.NoError(err)
	assert.Equal(content, resp)
}

func TestClientSendFileString(t *testing.T) {
	var resp []byte
	assert := assert.New(t)
	file, err := os.Open("./README.md")
	assert.NoError(err)
	content, err := ioutil.ReadAll(file)
	assert.NoError(err)
	err = client.
		Call("TestFile").
		SendFile("./README.md").
		EndStruct(&resp)
	assert.NoError(err)
	assert.Equal(content, resp)
}

func TestClientSendFileBytes(t *testing.T) {
	var resp []byte
	assert := assert.New(t)
	file, err := os.Open("./README.md")
	assert.NoError(err)
	content, err := ioutil.ReadAll(file)
	assert.NoError(err)
	err = client.
		Call("TestFile").
		SendFile(content).
		EndStruct(&resp)
	assert.NoError(err)
	assert.Equal(content, resp)
}

func TestClientSendFiles(t *testing.T) {
	var resp int
	assert := assert.New(t)
	file, err := os.Open("./README.md")
	assert.NoError(err)
	file2, err := os.Open("./client.go")
	assert.NoError(err)
	err = client.
		Call("TestFiles").
		SendFile(file).
		SendFile(file2).
		EndStruct(&resp)
	assert.NoError(err)
	assert.Equal(2, resp)
}

func TestClientSendFileSlices(t *testing.T) {
	var resp int
	assert := assert.New(t)
	file, err := os.Open("./README.md")
	assert.NoError(err)
	file2, err := os.Open("./client.go")
	assert.NoError(err)
	err = client.
		Call("TestFileSlice").
		SendFile(file, "File1").
		SendFile(file2, "File1").
		EndStruct(&resp)
	assert.NoError(err)
	assert.Equal(2, resp)
}

func TestClientSendFileMix(t *testing.T) {
	var resp H
	assert := assert.New(t)
	file, err := os.Open("./README.md")
	assert.NoError(err)
	file2, err := os.Open("./client.go")
	assert.NoError(err)
	file3, err := os.Open("./error.go")
	assert.NoError(err)
	err = client.
		Call("TestFileMix").
		SendFile(file, "File1").
		SendFile(file2, "File1").
		SendFile(file3, "File2").
		Send(dataStruct{
			Username: "YamiOdymel",
			Password: "yamiodymel12345",
		}).
		EndStruct(&resp)
	assert.NoError(err)
	assert.Equal(H{
		"Files":    3,
		"Username": "YamiOdymel",
		"Password": "yamiodymel12345",
	}, resp)
}

func TestClientSendFileChunksReader(t *testing.T) {
	var resp []byte
	assert := assert.New(t)
	file, err := os.Open("./README.md")
	assert.NoError(err)
	content, err := ioutil.ReadAll(file)
	assert.NoError(err)
	err = client.
		Call("TestFileChunks").
		SendFileChunks(file).
		EndStruct(&resp)
	assert.NoError(err)
	assert.Equal(content, resp)
}

func TestClientSendFileChunksString(t *testing.T) {
	var resp []byte
	assert := assert.New(t)
	file, err := os.Open("./README.md")
	assert.NoError(err)
	content, err := ioutil.ReadAll(file)
	assert.NoError(err)
	err = client.
		Call("TestFileChunks").
		SendFileChunks("./README.md").
		EndStruct(&resp)
	assert.NoError(err)
	assert.Equal(content, resp)
}

func TestClientSendFileChunksBytes(t *testing.T) {
	var resp []byte
	assert := assert.New(t)
	file, err := os.Open("./README.md")
	assert.NoError(err)
	content, err := ioutil.ReadAll(file)
	assert.NoError(err)
	err = client.
		Call("TestFileChunks").
		SendFileChunks(content).
		EndStruct(&resp)
	assert.NoError(err)
	assert.Equal(content, resp)
}

func TestClientSendFileChunksMix(t *testing.T) {
	var resp []byte
	assert := assert.New(t)
	file, err := os.Open("./README.md")
	assert.NoError(err)
	content, err := ioutil.ReadAll(file)
	assert.NoError(err)
	err = client.
		Call("TestFileChunksMix").
		SendFileChunks(content).
		Send(dataStruct{
			Username: "YamiOdymel",
			Password: "yamiodymel12345",
		}).
		EndStruct(&resp)
	assert.NoError(err)
	assert.Equal(H{
		"Files":    1,
		"Username": "YamiOdymel",
		"Password": "yamiodymel12345",
	}, resp)
}

func TestClientSubscribe(t *testing.T) {
	assert := assert.New(t)
	err := client.Subscribe("TestEvent", "TestChannel")
	assert.NoError(err)
}

func TestClientSubscribeError(t *testing.T) {
	assert := assert.New(t)
	err := client.Subscribe("TestRefuseEvent", "TestChannel")
	assert.Error(err)
	assert.Equal(ErrSubscriptionRefused, err)
}

func TestClientOn(t *testing.T) {
	//assert := assert.New(t)
	client.On("TestEvent", func(e *Event) {})
}

func TestClientOff(t *testing.T) {
	//assert := assert.New(t)
	client.Off("TestEvent")
}

func TestClientUnsubscribe(t *testing.T) {
	assert := assert.New(t)
	err := client.Unsubscribe("TestEvent", "TestChannel")
	assert.NoError(err)
}
